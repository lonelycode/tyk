package main

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/jeffail/tunny"
	"github.com/oschwald/maxminddb-golang"
	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/TykTechnologies/tyk/config"
)

// AnalyticsRecord encodes the details of a request
type AnalyticsRecord struct {
	Method        string
	Path          string
	RawPath       string
	ContentLength int64
	UserAgent     string
	Day           int
	Month         time.Month
	Year          int
	Hour          int
	ResponseCode  int
	APIKey        string
	TimeStamp     time.Time
	APIVersion    string
	APIName       string
	APIID         string
	OrgID         string
	OauthID       string
	RequestTime   int64
	RawRequest    string
	RawResponse   string
	IPAddress     string
	Geo           GeoData
	Tags          []string
	Alias         string
	TrackPath     bool
	ExpireAt      time.Time `bson:"expireAt" json:"expireAt"`
}

type GeoData struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`

	City struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`

	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
		TimeZone  string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

const analyticsKeyName = "tyk-system-analytics"

func (a *AnalyticsRecord) GetGeo(ipStr string) {
	if !globalConf.AnalyticsConfig.EnableGeoIP {
		return
	}
	// Not great, tightly coupled
	if analytics.GeoIPDB == nil {
		return
	}

	record, err := geoIPLookup(ipStr)
	if err != nil {
		log.Error("GeoIP Failure (not recorded): ", err)
		return
	}
	if record == nil {
		return
	}

	log.Debug("ISO Code: ", record.Country.ISOCode)
	log.Debug("City: ", record.City.Names["en"])
	log.Debug("Lat: ", record.Location.Latitude)
	log.Debug("Lon: ", record.Location.Longitude)
	log.Debug("TZ: ", record.Location.TimeZone)

	a.Geo = *record
}

func geoIPLookup(ipStr string) (*GeoData, error) {
	if ipStr == "" {
		return nil, nil
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address %q", ipStr)
	}
	record := new(GeoData)
	if err := analytics.GeoIPDB.Lookup(ip, record); err != nil {
		return nil, fmt.Errorf("geoIPDB lookup of %q failed: %v", ipStr, err)
	}
	return record, nil
}

func initNormalisationPatterns() (pats config.NormaliseURLPatterns) {
	pats.UUIDs = regexp.MustCompile(`[0-9a-fA-F]{8}(-)?[0-9a-fA-F]{4}(-)?[0-9a-fA-F]{4}(-)?[0-9a-fA-F]{4}(-)?[0-9a-fA-F]{12}`)
	pats.IDs = regexp.MustCompile(`\/(\d+)`)

	for _, pattern := range globalConf.AnalyticsConfig.NormaliseUrls.Custom {
		if patRe, err := regexp.Compile(pattern); err != nil {
			log.Error("failed to compile custom pattern: ", err)
		} else {
			pats.Custom = append(pats.Custom, patRe)
		}
	}
	return
}

func (a *AnalyticsRecord) NormalisePath() {
	if globalConf.AnalyticsConfig.NormaliseUrls.NormaliseUUIDs {
		a.Path = globalConf.AnalyticsConfig.NormaliseUrls.CompiledPatternSet.UUIDs.ReplaceAllString(a.Path, "{uuid}")
	}
	if globalConf.AnalyticsConfig.NormaliseUrls.NormaliseNumbers {
		a.Path = globalConf.AnalyticsConfig.NormaliseUrls.CompiledPatternSet.IDs.ReplaceAllString(a.Path, "/{id}")
	}
	for _, r := range globalConf.AnalyticsConfig.NormaliseUrls.CompiledPatternSet.Custom {
		a.Path = r.ReplaceAllString(a.Path, "{var}")
	}
}

func (a *AnalyticsRecord) SetExpiry(expiresInSeconds int64) {
	expiry := time.Duration(expiresInSeconds) * time.Second
	if expiresInSeconds == 0 {
		// Expiry is set to 100 years
		expiry = (24 * time.Hour) * (365 * 100)
	}

	t := time.Now()
	t2 := t.Add(expiry)
	a.ExpireAt = t2
}

var AnalyticsPool *tunny.WorkPool

// RedisAnalyticsHandler will record analytics data to a redis back end
// as defined in the Config object
type RedisAnalyticsHandler struct {
	Store   StorageHandler
	Clean   Purger
	GeoIPDB *maxminddb.Reader
}

func (r *RedisAnalyticsHandler) Init() {
	var err error
	if globalConf.AnalyticsConfig.EnableGeoIP {
		db, err := maxminddb.Open(globalConf.AnalyticsConfig.GeoIPDBLocation)
		if err != nil {
			log.Error("Failed to init GeoIP Database: ", err)
		} else {
			r.GeoIPDB = db
		}
	}

	analytics.Store.Connect()

	ps := globalConf.AnalyticsConfig.PoolSize
	if ps == 0 {
		ps = 50
	}

	AnalyticsPool, err = tunny.CreatePoolGeneric(ps).Open()
	if err != nil {
		log.Error("Failed to init analytics pool")
	}
}

// RecordHit will store an AnalyticsRecord in Redis
func (r *RedisAnalyticsHandler) RecordHit(record AnalyticsRecord) error {

	AnalyticsPool.SendWork(func() {
		// If we are obfuscating API Keys, store the hashed representation (config check handled in hashing function)
		record.APIKey = publicHash(record.APIKey)

		if globalConf.SlaveOptions.UseRPC {
			// Extend tag list to include this data so wecan segment by node if necessary
			record.Tags = append(record.Tags, "tyk-hybrid-rpc")
		}

		if globalConf.DBAppConfOptions.NodeIsSegmented {
			// Extend tag list to include this data so wecan segment by node if necessary
			record.Tags = append(record.Tags, globalConf.DBAppConfOptions.Tags...)
		}

		// Lets add some metadata
		if record.APIKey != "" {
			record.Tags = append(record.Tags, "key-"+record.APIKey)
		}

		if record.OrgID != "" {
			record.Tags = append(record.Tags, "org-"+record.OrgID)
		}

		record.Tags = append(record.Tags, "api-"+record.APIID)

		encoded, err := msgpack.Marshal(record)

		if err != nil {
			log.Error("Error encoding analytics data: ", err)
		}

		r.Store.AppendToSet(analyticsKeyName, string(encoded))
	})

	return nil

}
