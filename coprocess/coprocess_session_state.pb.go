// Code generated by protoc-gen-go. DO NOT EDIT.
// source: coprocess_session_state.proto

package coprocess

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type AccessSpec struct {
	Url     string   `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Methods []string `protobuf:"bytes,2,rep,name=methods" json:"methods,omitempty"`
}

func (m *AccessSpec) Reset()                    { *m = AccessSpec{} }
func (m *AccessSpec) String() string            { return proto.CompactTextString(m) }
func (*AccessSpec) ProtoMessage()               {}
func (*AccessSpec) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *AccessSpec) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *AccessSpec) GetMethods() []string {
	if m != nil {
		return m.Methods
	}
	return nil
}

type AccessDefinition struct {
	ApiName     string        `protobuf:"bytes,1,opt,name=api_name,json=apiName" json:"api_name,omitempty"`
	ApiId       string        `protobuf:"bytes,2,opt,name=api_id,json=apiId" json:"api_id,omitempty"`
	Versions    []string      `protobuf:"bytes,3,rep,name=versions" json:"versions,omitempty"`
	AllowedUrls []*AccessSpec `protobuf:"bytes,4,rep,name=allowed_urls,json=allowedUrls" json:"allowed_urls,omitempty"`
}

func (m *AccessDefinition) Reset()                    { *m = AccessDefinition{} }
func (m *AccessDefinition) String() string            { return proto.CompactTextString(m) }
func (*AccessDefinition) ProtoMessage()               {}
func (*AccessDefinition) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{1} }

func (m *AccessDefinition) GetApiName() string {
	if m != nil {
		return m.ApiName
	}
	return ""
}

func (m *AccessDefinition) GetApiId() string {
	if m != nil {
		return m.ApiId
	}
	return ""
}

func (m *AccessDefinition) GetVersions() []string {
	if m != nil {
		return m.Versions
	}
	return nil
}

func (m *AccessDefinition) GetAllowedUrls() []*AccessSpec {
	if m != nil {
		return m.AllowedUrls
	}
	return nil
}

type BasicAuthData struct {
	Password string `protobuf:"bytes,1,opt,name=password" json:"password,omitempty"`
	Hash     string `protobuf:"bytes,2,opt,name=hash" json:"hash,omitempty"`
}

func (m *BasicAuthData) Reset()                    { *m = BasicAuthData{} }
func (m *BasicAuthData) String() string            { return proto.CompactTextString(m) }
func (*BasicAuthData) ProtoMessage()               {}
func (*BasicAuthData) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{2} }

func (m *BasicAuthData) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *BasicAuthData) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type JWTData struct {
	Secret string `protobuf:"bytes,1,opt,name=secret" json:"secret,omitempty"`
}

func (m *JWTData) Reset()                    { *m = JWTData{} }
func (m *JWTData) String() string            { return proto.CompactTextString(m) }
func (*JWTData) ProtoMessage()               {}
func (*JWTData) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{3} }

func (m *JWTData) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

type Monitor struct {
	TriggerLimits []float64 `protobuf:"fixed64,1,rep,packed,name=trigger_limits,json=triggerLimits" json:"trigger_limits,omitempty"`
}

func (m *Monitor) Reset()                    { *m = Monitor{} }
func (m *Monitor) String() string            { return proto.CompactTextString(m) }
func (*Monitor) ProtoMessage()               {}
func (*Monitor) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{4} }

func (m *Monitor) GetTriggerLimits() []float64 {
	if m != nil {
		return m.TriggerLimits
	}
	return nil
}

type SessionState struct {
	LastCheck               int64                        `protobuf:"varint,1,opt,name=last_check,json=lastCheck" json:"last_check,omitempty"`
	Allowance               float64                      `protobuf:"fixed64,2,opt,name=allowance" json:"allowance,omitempty"`
	Rate                    float64                      `protobuf:"fixed64,3,opt,name=rate" json:"rate,omitempty"`
	Per                     float64                      `protobuf:"fixed64,4,opt,name=per" json:"per,omitempty"`
	Expires                 int64                        `protobuf:"varint,5,opt,name=expires" json:"expires,omitempty"`
	QuotaMax                int64                        `protobuf:"varint,6,opt,name=quota_max,json=quotaMax" json:"quota_max,omitempty"`
	QuotaRenews             int64                        `protobuf:"varint,7,opt,name=quota_renews,json=quotaRenews" json:"quota_renews,omitempty"`
	QuotaRemaining          int64                        `protobuf:"varint,8,opt,name=quota_remaining,json=quotaRemaining" json:"quota_remaining,omitempty"`
	QuotaRenewalRate        int64                        `protobuf:"varint,9,opt,name=quota_renewal_rate,json=quotaRenewalRate" json:"quota_renewal_rate,omitempty"`
	AccessRights            map[string]*AccessDefinition `protobuf:"bytes,10,rep,name=access_rights,json=accessRights" json:"access_rights,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	OrgId                   string                       `protobuf:"bytes,11,opt,name=org_id,json=orgId" json:"org_id,omitempty"`
	OauthClientId           string                       `protobuf:"bytes,12,opt,name=oauth_client_id,json=oauthClientId" json:"oauth_client_id,omitempty"`
	OauthKeys               map[string]string            `protobuf:"bytes,13,rep,name=oauth_keys,json=oauthKeys" json:"oauth_keys,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	BasicAuthData           *BasicAuthData               `protobuf:"bytes,14,opt,name=basic_auth_data,json=basicAuthData" json:"basic_auth_data,omitempty"`
	JwtData                 *JWTData                     `protobuf:"bytes,15,opt,name=jwt_data,json=jwtData" json:"jwt_data,omitempty"`
	HmacEnabled             bool                         `protobuf:"varint,16,opt,name=hmac_enabled,json=hmacEnabled" json:"hmac_enabled,omitempty"`
	HmacSecret              string                       `protobuf:"bytes,17,opt,name=hmac_secret,json=hmacSecret" json:"hmac_secret,omitempty"`
	IsInactive              bool                         `protobuf:"varint,18,opt,name=is_inactive,json=isInactive" json:"is_inactive,omitempty"`
	ApplyPolicyId           string                       `protobuf:"bytes,19,opt,name=apply_policy_id,json=applyPolicyId" json:"apply_policy_id,omitempty"`
	DataExpires             int64                        `protobuf:"varint,20,opt,name=data_expires,json=dataExpires" json:"data_expires,omitempty"`
	Monitor                 *Monitor                     `protobuf:"bytes,21,opt,name=monitor" json:"monitor,omitempty"`
	EnableDetailedRecording bool                         `protobuf:"varint,22,opt,name=enable_detailed_recording,json=enableDetailedRecording" json:"enable_detailed_recording,omitempty"`
	Metadata                map[string]string            `protobuf:"bytes,23,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Tags                    []string                     `protobuf:"bytes,24,rep,name=tags" json:"tags,omitempty"`
	Alias                   string                       `protobuf:"bytes,25,opt,name=alias" json:"alias,omitempty"`
	LastUpdated             string                       `protobuf:"bytes,26,opt,name=last_updated,json=lastUpdated" json:"last_updated,omitempty"`
	IdExtractorDeadline     int64                        `protobuf:"varint,27,opt,name=id_extractor_deadline,json=idExtractorDeadline" json:"id_extractor_deadline,omitempty"`
	SessionLifetime         int64                        `protobuf:"varint,28,opt,name=session_lifetime,json=sessionLifetime" json:"session_lifetime,omitempty"`
	ApplyPolicies           []string                     `protobuf:"bytes,29,rep,name=apply_policies,json=applyPolicies" json:"apply_policies,omitempty"`
	Certificate             string                       `protobuf:"bytes,30,opt,name=certificate" json:"certificate,omitempty"`
}

func (m *SessionState) Reset()                    { *m = SessionState{} }
func (m *SessionState) String() string            { return proto.CompactTextString(m) }
func (*SessionState) ProtoMessage()               {}
func (*SessionState) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{5} }

func (m *SessionState) GetLastCheck() int64 {
	if m != nil {
		return m.LastCheck
	}
	return 0
}

func (m *SessionState) GetAllowance() float64 {
	if m != nil {
		return m.Allowance
	}
	return 0
}

func (m *SessionState) GetRate() float64 {
	if m != nil {
		return m.Rate
	}
	return 0
}

func (m *SessionState) GetPer() float64 {
	if m != nil {
		return m.Per
	}
	return 0
}

func (m *SessionState) GetExpires() int64 {
	if m != nil {
		return m.Expires
	}
	return 0
}

func (m *SessionState) GetQuotaMax() int64 {
	if m != nil {
		return m.QuotaMax
	}
	return 0
}

func (m *SessionState) GetQuotaRenews() int64 {
	if m != nil {
		return m.QuotaRenews
	}
	return 0
}

func (m *SessionState) GetQuotaRemaining() int64 {
	if m != nil {
		return m.QuotaRemaining
	}
	return 0
}

func (m *SessionState) GetQuotaRenewalRate() int64 {
	if m != nil {
		return m.QuotaRenewalRate
	}
	return 0
}

func (m *SessionState) GetAccessRights() map[string]*AccessDefinition {
	if m != nil {
		return m.AccessRights
	}
	return nil
}

func (m *SessionState) GetOrgId() string {
	if m != nil {
		return m.OrgId
	}
	return ""
}

func (m *SessionState) GetOauthClientId() string {
	if m != nil {
		return m.OauthClientId
	}
	return ""
}

func (m *SessionState) GetOauthKeys() map[string]string {
	if m != nil {
		return m.OauthKeys
	}
	return nil
}

func (m *SessionState) GetBasicAuthData() *BasicAuthData {
	if m != nil {
		return m.BasicAuthData
	}
	return nil
}

func (m *SessionState) GetJwtData() *JWTData {
	if m != nil {
		return m.JwtData
	}
	return nil
}

func (m *SessionState) GetHmacEnabled() bool {
	if m != nil {
		return m.HmacEnabled
	}
	return false
}

func (m *SessionState) GetHmacSecret() string {
	if m != nil {
		return m.HmacSecret
	}
	return ""
}

func (m *SessionState) GetIsInactive() bool {
	if m != nil {
		return m.IsInactive
	}
	return false
}

func (m *SessionState) GetApplyPolicyId() string {
	if m != nil {
		return m.ApplyPolicyId
	}
	return ""
}

func (m *SessionState) GetDataExpires() int64 {
	if m != nil {
		return m.DataExpires
	}
	return 0
}

func (m *SessionState) GetMonitor() *Monitor {
	if m != nil {
		return m.Monitor
	}
	return nil
}

func (m *SessionState) GetEnableDetailedRecording() bool {
	if m != nil {
		return m.EnableDetailedRecording
	}
	return false
}

func (m *SessionState) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *SessionState) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *SessionState) GetAlias() string {
	if m != nil {
		return m.Alias
	}
	return ""
}

func (m *SessionState) GetLastUpdated() string {
	if m != nil {
		return m.LastUpdated
	}
	return ""
}

func (m *SessionState) GetIdExtractorDeadline() int64 {
	if m != nil {
		return m.IdExtractorDeadline
	}
	return 0
}

func (m *SessionState) GetSessionLifetime() int64 {
	if m != nil {
		return m.SessionLifetime
	}
	return 0
}

func (m *SessionState) GetApplyPolicies() []string {
	if m != nil {
		return m.ApplyPolicies
	}
	return nil
}

func (m *SessionState) GetCertificate() string {
	if m != nil {
		return m.Certificate
	}
	return ""
}

func init() {
	proto.RegisterType((*AccessSpec)(nil), "coprocess.AccessSpec")
	proto.RegisterType((*AccessDefinition)(nil), "coprocess.AccessDefinition")
	proto.RegisterType((*BasicAuthData)(nil), "coprocess.BasicAuthData")
	proto.RegisterType((*JWTData)(nil), "coprocess.JWTData")
	proto.RegisterType((*Monitor)(nil), "coprocess.Monitor")
	proto.RegisterType((*SessionState)(nil), "coprocess.SessionState")
}

func init() { proto.RegisterFile("coprocess_session_state.proto", fileDescriptor4) }

var fileDescriptor4 = []byte{
	// 933 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0x4f, 0x6f, 0x1b, 0xb7,
	0x13, 0x85, 0x2c, 0xdb, 0x92, 0x66, 0x25, 0x5b, 0x61, 0xec, 0x84, 0xb6, 0xe3, 0xdf, 0x4f, 0x16,
	0x90, 0x54, 0x01, 0x52, 0xa3, 0x75, 0x2f, 0x46, 0x5a, 0xa0, 0x75, 0x63, 0x1d, 0xd4, 0xc6, 0x69,
	0xb1, 0x6e, 0xd0, 0x4b, 0x01, 0x82, 0xde, 0x1d, 0x4b, 0x8c, 0xf7, 0x5f, 0x49, 0xca, 0xb2, 0xbe,
	0x47, 0x4f, 0xfd, 0xb4, 0x05, 0x67, 0xb9, 0xf2, 0x1a, 0x6e, 0x0e, 0xbd, 0x71, 0xde, 0x7b, 0x33,
	0xfb, 0xc8, 0x19, 0x72, 0xe1, 0x30, 0xca, 0x0b, 0x9d, 0x47, 0x68, 0x8c, 0x30, 0x68, 0x8c, 0xca,
	0x33, 0x61, 0xac, 0xb4, 0x78, 0x5c, 0xe8, 0xdc, 0xe6, 0xac, 0xb3, 0xa2, 0x87, 0xa7, 0x00, 0x67,
	0x91, 0x5b, 0x5d, 0x16, 0x18, 0xb1, 0x3e, 0x34, 0xe7, 0x3a, 0xe1, 0x8d, 0x41, 0x63, 0xd4, 0x09,
	0xdd, 0x92, 0x71, 0x68, 0xa5, 0x68, 0x67, 0x79, 0x6c, 0xf8, 0xda, 0xa0, 0x39, 0xea, 0x84, 0x55,
	0x38, 0xfc, 0xbb, 0x01, 0xfd, 0x32, 0xf5, 0x1c, 0xaf, 0x55, 0xa6, 0xac, 0xca, 0x33, 0xb6, 0x07,
	0x6d, 0x59, 0x28, 0x91, 0xc9, 0x14, 0x7d, 0x95, 0x96, 0x2c, 0xd4, 0x07, 0x99, 0x22, 0xdb, 0x85,
	0x4d, 0x47, 0xa9, 0x98, 0xaf, 0x11, 0xb1, 0x21, 0x0b, 0x35, 0x89, 0xd9, 0x3e, 0xb4, 0x6f, 0x51,
	0x3b, 0x8b, 0x86, 0x37, 0xe9, 0x0b, 0xab, 0x98, 0x9d, 0x42, 0x57, 0x26, 0x49, 0xbe, 0xc0, 0x58,
	0xcc, 0x75, 0x62, 0xf8, 0xfa, 0xa0, 0x39, 0x0a, 0x4e, 0x76, 0x8f, 0x57, 0xf6, 0x8f, 0xef, 0xbd,
	0x87, 0x81, 0x97, 0x7e, 0xd4, 0x89, 0x19, 0x7e, 0x0f, 0xbd, 0x1f, 0xa5, 0x51, 0xd1, 0xd9, 0xdc,
	0xce, 0xce, 0xa5, 0x95, 0xee, 0x33, 0x85, 0x34, 0x66, 0x91, 0xeb, 0xd8, 0x1b, 0x5b, 0xc5, 0x8c,
	0xc1, 0xfa, 0x4c, 0x9a, 0x99, 0xf7, 0x45, 0xeb, 0xe1, 0x11, 0xb4, 0x7e, 0xfa, 0xfd, 0x37, 0x4a,
	0x7d, 0x06, 0x9b, 0x06, 0x23, 0x8d, 0xd6, 0x27, 0xfa, 0x68, 0xf8, 0x15, 0xb4, 0x2e, 0xf2, 0x4c,
	0xd9, 0x5c, 0xb3, 0x97, 0xb0, 0x65, 0xb5, 0x9a, 0x4e, 0x51, 0x8b, 0x44, 0xa5, 0xca, 0x1a, 0xde,
	0x18, 0x34, 0x47, 0x8d, 0xb0, 0xe7, 0xd1, 0xf7, 0x04, 0x0e, 0xff, 0x0a, 0xa0, 0x7b, 0x59, 0xf6,
	0xe3, 0xd2, 0xb5, 0x83, 0x1d, 0x02, 0x24, 0xd2, 0x58, 0x11, 0xcd, 0x30, 0xba, 0xa1, 0xf2, 0xcd,
	0xb0, 0xe3, 0x90, 0x77, 0x0e, 0x60, 0x2f, 0xa0, 0x43, 0x9b, 0x92, 0x59, 0x84, 0xe4, 0xae, 0x11,
	0xde, 0x03, 0xce, 0xb6, 0x96, 0x16, 0x79, 0x93, 0x08, 0x5a, 0xbb, 0x06, 0x16, 0xa8, 0xf9, 0x3a,
	0x41, 0x6e, 0xe9, 0x1a, 0x88, 0x77, 0x85, 0xd2, 0x68, 0xf8, 0x06, 0xd5, 0xaf, 0x42, 0x76, 0x00,
	0x9d, 0x3f, 0xe7, 0xb9, 0x95, 0x22, 0x95, 0x77, 0x7c, 0x93, 0xb8, 0x36, 0x01, 0x17, 0xf2, 0x8e,
	0x1d, 0x41, 0xb7, 0x24, 0x35, 0x66, 0xb8, 0x30, 0xbc, 0x45, 0x7c, 0x40, 0x58, 0x48, 0x10, 0xfb,
	0x02, 0xb6, 0x2b, 0x49, 0x2a, 0x55, 0xa6, 0xb2, 0x29, 0x6f, 0x93, 0x6a, 0xcb, 0xab, 0x3c, 0xca,
	0xde, 0x00, 0xab, 0xd5, 0x92, 0x89, 0x20, 0xdb, 0x1d, 0xd2, 0xf6, 0xef, 0x2b, 0xca, 0x24, 0x74,
	0x5b, 0xf8, 0x00, 0x3d, 0x49, 0x5d, 0x15, 0x5a, 0x4d, 0x67, 0xd6, 0x70, 0xa0, 0xae, 0xbf, 0xae,
	0x75, 0xbd, 0x7e, 0x86, 0x7e, 0x04, 0x42, 0xd2, 0x8e, 0x33, 0xab, 0x97, 0x61, 0x57, 0xd6, 0x20,
	0x37, 0x77, 0xb9, 0x9e, 0xba, 0xb9, 0x0b, 0xca, 0xb9, 0xcb, 0xf5, 0x74, 0x12, 0xb3, 0x57, 0xb0,
	0x9d, 0xcb, 0xb9, 0x9d, 0x89, 0x28, 0x51, 0x98, 0x59, 0xc7, 0x77, 0x89, 0xef, 0x11, 0xfc, 0x8e,
	0xd0, 0x49, 0xcc, 0xc6, 0x00, 0xa5, 0xee, 0x06, 0x97, 0x86, 0xf7, 0xc8, 0xcb, 0xab, 0xcf, 0x79,
	0xf9, 0xc5, 0x29, 0x7f, 0xc6, 0xa5, 0x37, 0xd2, 0xc9, 0xab, 0x98, 0xfd, 0x00, 0xdb, 0x57, 0x6e,
	0x20, 0x05, 0xd5, 0x8a, 0xa5, 0x95, 0x7c, 0x6b, 0xd0, 0x18, 0x05, 0x27, 0xbc, 0x56, 0xeb, 0xc1,
	0xc8, 0x86, 0xbd, 0xab, 0x07, 0x13, 0xfc, 0x25, 0xb4, 0x3f, 0x2d, 0x6c, 0x99, 0xba, 0x4d, 0xa9,
	0xac, 0x96, 0xea, 0x87, 0x35, 0x6c, 0x7d, 0x5a, 0x58, 0x92, 0x1f, 0x41, 0x77, 0x96, 0xca, 0x48,
	0x60, 0x26, 0xaf, 0x12, 0x8c, 0x79, 0x7f, 0xd0, 0x18, 0xb5, 0xc3, 0xc0, 0x61, 0xe3, 0x12, 0x62,
	0xff, 0x07, 0x0a, 0x85, 0x9f, 0xee, 0x27, 0xb4, 0x7d, 0x70, 0xd0, 0x25, 0x21, 0x4e, 0xa0, 0x8c,
	0x50, 0x99, 0x8c, 0xac, 0xba, 0x45, 0xce, 0xa8, 0x04, 0x28, 0x33, 0xf1, 0x88, 0x3b, 0x44, 0x59,
	0x14, 0xc9, 0x52, 0x14, 0x79, 0xa2, 0xa2, 0xa5, 0x3b, 0xc4, 0xa7, 0xe5, 0x21, 0x12, 0xfc, 0x2b,
	0xa1, 0x93, 0xd8, 0x99, 0x71, 0xbe, 0x45, 0x35, 0x89, 0x3b, 0xe5, 0x34, 0x39, 0x6c, 0xec, 0xa7,
	0xf1, 0x0d, 0xb4, 0xd2, 0xf2, 0x36, 0xf1, 0xdd, 0x47, 0xbb, 0xf3, 0xf7, 0x2c, 0xac, 0x24, 0xec,
	0x2d, 0xec, 0x95, 0x1b, 0x13, 0x31, 0x5a, 0xa9, 0x12, 0x8c, 0x85, 0xc6, 0x28, 0xd7, 0xb1, 0x9b,
	0xc2, 0x67, 0xe4, 0xf3, 0x79, 0x29, 0x38, 0xf7, 0x7c, 0x58, 0xd1, 0xec, 0x0c, 0xda, 0x29, 0x5a,
	0x49, 0x07, 0xf9, 0x9c, 0xfa, 0xf9, 0xf2, 0x73, 0xfd, 0xbc, 0xf0, 0xba, 0xb2, 0x9d, 0xab, 0x34,
	0x77, 0xf5, 0xac, 0x9c, 0x1a, 0xce, 0xe9, 0xc1, 0xa2, 0x35, 0xdb, 0x81, 0x0d, 0x99, 0x28, 0x69,
	0xf8, 0x9e, 0x7f, 0xde, 0x5c, 0xe0, 0x76, 0x4e, 0x37, 0x7c, 0x5e, 0xc4, 0xd2, 0x62, 0xcc, 0xf7,
	0x89, 0x0c, 0x1c, 0xf6, 0xb1, 0x84, 0xd8, 0x09, 0xec, 0xaa, 0x58, 0xe0, 0x9d, 0xd5, 0x32, 0xb2,
	0xb9, 0x16, 0x31, 0xca, 0x38, 0x51, 0x19, 0xf2, 0x03, 0x3a, 0xa5, 0xa7, 0x2a, 0x1e, 0x57, 0xdc,
	0xb9, 0xa7, 0xd8, 0x6b, 0xe8, 0x57, 0x0f, 0x7b, 0xa2, 0xae, 0xd1, 0xaa, 0x14, 0xf9, 0x0b, 0x92,
	0x6f, 0x7b, 0xfc, 0xbd, 0x87, 0xdd, 0xdb, 0x54, 0xeb, 0x91, 0x42, 0xc3, 0x0f, 0xc9, 0x75, 0xad,
	0x45, 0x0a, 0x0d, 0x1b, 0x40, 0x10, 0xa1, 0xb6, 0xea, 0x5a, 0x45, 0xee, 0x76, 0xfe, 0xaf, 0xf4,
	0x59, 0x83, 0xf6, 0xff, 0x80, 0x27, 0x8f, 0xee, 0x9a, 0x7b, 0x70, 0x6e, 0x70, 0x59, 0xfd, 0x31,
	0x6e, 0x70, 0xc9, 0xbe, 0x86, 0x8d, 0x5b, 0x99, 0xcc, 0xcb, 0x07, 0x2b, 0x38, 0x39, 0x78, 0xf4,
	0x5a, 0xdf, 0xff, 0x2e, 0xc2, 0x52, 0xf9, 0x76, 0xed, 0xb4, 0xb1, 0xff, 0x1d, 0x6c, 0x3d, 0xbc,
	0x3d, 0xff, 0x52, 0x7a, 0xa7, 0x5e, 0xba, 0x53, 0xcf, 0xfe, 0x16, 0x7a, 0x0f, 0x7a, 0xf5, 0x5f,
	0x92, 0xaf, 0x36, 0xe9, 0xaf, 0xf8, 0xcd, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xdc, 0xad, 0x0c,
	0x63, 0x36, 0x07, 0x00, 0x00,
}
