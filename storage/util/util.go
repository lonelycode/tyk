package util

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"github.com/TykTechnologies/murmur3"
	"github.com/TykTechnologies/tyk/internal/uuid"
	"github.com/buger/jsonparser"
)

const defaultHashAlgorithm = "murmur64"
const MongoBsonIdLength = 24

// If hashing algorithm is empty, use legacy key generation
func GenerateToken(orgID, keyID, hashAlgorithm string) (string, error) {
	if keyID == "" {
		keyID = uuid.NewHex()
	}

	if hashAlgorithm != "" {
		_, err := hashFunction(hashAlgorithm)
		if err != nil {
			hashAlgorithm = defaultHashAlgorithm
		}

		jsonToken := fmt.Sprintf(`{"org":"%s","id":"%s","h":"%s"}`, orgID, keyID, hashAlgorithm)
		return base64.StdEncoding.EncodeToString([]byte(jsonToken)), err
	}

	// Legacy keys
	return orgID + keyID, nil
}

// `{"` in base64
const B64JSONPrefix = "ey"

func TokenHashAlgo(token string) string {
	// Legacy tokens not b64 and not JSON records
	if strings.HasPrefix(token, B64JSONPrefix) {
		if jsonToken, err := base64.StdEncoding.DecodeString(token); err == nil {
			hashAlgo, _ := jsonparser.GetString(jsonToken, "h")
			return hashAlgo
		}
	}

	return ""
}

// TODO: add checks
func TokenID(token string) (id string, err error) {
	jsonToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	return jsonparser.GetString(jsonToken, "id")
}

func TokenOrg(token string) string {
	if strings.HasPrefix(token, B64JSONPrefix) {
		if jsonToken, err := base64.StdEncoding.DecodeString(token); err == nil {
			// Checking error in case if it is a legacy tooken which just by accided has the same b64JSON prefix
			if org, err := jsonparser.GetString(jsonToken, "org"); err == nil {
				return org
			}
		}
	}

	// 24 is mongo bson id length
	if len(token) > MongoBsonIdLength {
		newToken := token[:MongoBsonIdLength]
		_, err := hex.DecodeString(newToken)
		if err == nil {
			return newToken
		}
	}

	return ""
}

var (
	HashSha256    = "sha256"
	HashMurmur32  = "murmur32"
	HashMurmur64  = "murmur64"
	HashMurmur128 = "murmur128"
)

func hashFunction(algorithm string) (hash.Hash, error) {
	switch algorithm {
	case HashSha256:
		return sha256.New(), nil
	case HashMurmur64:
		return murmur3.New64(), nil
	case HashMurmur128:
		return murmur3.New128(), nil
	case "", HashMurmur32:
		return murmur3.New32(), nil
	default:
		return murmur3.New32(), fmt.Errorf("Unknown key hash function: %s. Falling back to murmur32.", algorithm)
	}
}

func HashStr(in string, withAlg ...string) string {
	var algo string
	if len(withAlg) > 0 && withAlg[0] != "" {
		algo = withAlg[0]
	} else {
		algo = TokenHashAlgo(in)
	}

	h, _ := hashFunction(algo)
	h.Write([]byte(in))
	return hex.EncodeToString(h.Sum(nil))
}

func HashKey(in string, hashKey bool) string {
	if !hashKey {
		// Not hashing? Return the raw key
		return in
	}
	return HashStr(in)
}
