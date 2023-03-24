package apidef

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	schema "github.com/xeipuuv/gojsonschema"
)

func TestSchema(t *testing.T) {
	schemaLoader := schema.NewBytesLoader([]byte(Schema))

	spec := DummyAPI()
	goLoader := schema.NewGoLoader(spec)
	result, err := schema.Validate(schemaLoader, goLoader)
	if err != nil {
		t.Error(err)
	}

	if !result.Valid() {
		for _, err := range result.Errors() {
			t.Error(err)
		}
	}
}

func TestEncodeForDB(t *testing.T) {
	t.Run("EncodeForDB persist schema objects from extended path", func(t *testing.T) {
		spec := DummyAPI()
		spec.EncodeForDB()
		var schemaNotEmpty bool
		for _, version := range spec.VersionData.Versions {
			for _, validateObj := range version.ExtendedPaths.ValidateJSON {
				schemaNotEmpty = schemaNotEmpty || (validateObj.Schema != nil)
			}
		}
		assert.True(t, schemaNotEmpty, "expected EncodeForDB to persist schema objects")
	})
}

func TestDecodeFromDB(t *testing.T) {
	t.Run("json schema validation middleware", func(t *testing.T) {
		apiDef := DummyAPI()
		var (
			bodySchema map[string]interface{}
			v1         = "v1"
			v1B64      = base64.StdEncoding.EncodeToString([]byte(v1))
		)
		err := json.Unmarshal([]byte(`{"$schema":"http://json-schema.org/draft-04/schema#","properties":{"id":{"type":"integer"}},"required":["id"],"type":"object"}`),
			&bodySchema)
		assert.NoError(t, err)
		apiDef.VersionData.Versions[v1] = VersionInfo{
			ExtendedPaths: ExtendedPathsSet{
				ValidateJSON: []ValidatePathMeta{
					{
						Path:   "/",
						Method: http.MethodPost,
						Schema: bodySchema,
					},
				},
			},
		}
		apiDef.EncodeForDB()
		copyAPIDef := apiDef
		copyAPIDef.DecodeFromDB()

		assert.Equal(t, apiDef.VersionData.Versions[v1B64].ExtendedPaths.ValidateJSON[0].Schema,
			copyAPIDef.VersionData.Versions[v1].ExtendedPaths.ValidateJSON[0].Schema)
		assert.Empty(t, copyAPIDef.VersionData.Versions[v1].ExtendedPaths.ValidateJSON[0].SchemaB64)
	})
}

func TestSchemaGraphqlConfig(t *testing.T) {
	schemaLoader := schema.NewBytesLoader([]byte(Schema))

	spec := DummyAPI()
	spec.GraphQL.ExecutionMode = ""

	goLoader := schema.NewGoLoader(spec)

	result, err := schema.Validate(schemaLoader, goLoader)
	if err != nil {
		t.Error(err)
	}

	if !result.Valid() {
		for _, err := range result.Errors() {
			t.Error(err)
		}
	}
}

func TestAPIDefinition_DecodeFromDB_AuthDeprecation(t *testing.T) {
	const authHeader = "authorization"

	spec := DummyAPI()
	spec.Auth = AuthConfig{AuthHeaderName: authHeader}
	spec.UseStandardAuth = true
	spec.DecodeFromDB()

	assert.Equal(t, spec.AuthConfigs, map[string]AuthConfig{
		"authToken": spec.Auth,
	})

	spec.EnableJWT = true
	spec.DecodeFromDB()

	assert.Equal(t, spec.AuthConfigs, map[string]AuthConfig{
		"authToken": spec.Auth,
		"jwt":       spec.Auth,
	})

}

func TestAPIDefinition_GenerateAPIID(t *testing.T) {
	a := APIDefinition{}
	a.GenerateAPIID()
	assert.NotEmpty(t, a.APIID)
}

func TestAPIDefinition_GetScopeClaimName(t *testing.T) {
	var (
		scopeName     = "scope"
		oidcScopeName = "oidc_scope"
	)

	getAPIDef := func(deprecatedScopeName, jwtScopeName, oidcScopeName string, useOIDC bool) APIDefinition {
		return APIDefinition{
			UseOpenID:         useOIDC,
			JWTScopeClaimName: deprecatedScopeName,
		}
	}

	testCases := []struct {
		name                string
		deprecatedScopeName string
		jwtScopeName        string
		oidcScopeName       string
		useOIDC             bool
		expectedScopeName   string
	}{
		{
			name:                "jwt",
			deprecatedScopeName: scopeName,
			expectedScopeName:   scopeName,
		},

		{
			name:                "oidc",
			deprecatedScopeName: oidcScopeName,
			expectedScopeName:   oidcScopeName,
			useOIDC:             true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiDef := getAPIDef(tc.deprecatedScopeName, tc.jwtScopeName, tc.oidcScopeName, tc.useOIDC)
			assert.Equal(t, tc.expectedScopeName, apiDef.GetScopeClaimName())
		})
	}
}

func TestAPIDefinition_GetScopeToPolicyMapping(t *testing.T) {
	var (
		scopeToPolicyMapping     = map[string]string{"jwtClaim": "pol1"}
		oidcScopeToPolicyMapping = map[string]string{"oidcClaim": "pol1"}
	)

	getAPIDef := func(deprecatedScopeToPolicy, jwtScopeToPolicy, oidcScopeToPolicy map[string]string, useOIDC bool) APIDefinition {
		return APIDefinition{
			UseOpenID:               useOIDC,
			JWTScopeToPolicyMapping: deprecatedScopeToPolicy,
		}
	}

	testCases := []struct {
		name                    string
		deprecatedScopeToPolicy map[string]string
		jwtScopeToPolicy        map[string]string
		oidcScopeToPolicy       map[string]string
		useOIDC                 bool
		expectedScopeToPolicy   map[string]string
	}{
		{
			name:                    "jwt",
			deprecatedScopeToPolicy: scopeToPolicyMapping,
			expectedScopeToPolicy:   scopeToPolicyMapping,
		},

		{
			name:                    "oidc",
			deprecatedScopeToPolicy: oidcScopeToPolicyMapping,
			expectedScopeToPolicy:   oidcScopeToPolicyMapping,
			useOIDC:                 true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiDef := getAPIDef(tc.deprecatedScopeToPolicy, tc.jwtScopeToPolicy, tc.oidcScopeToPolicy, tc.useOIDC)
			assert.Equal(t, tc.expectedScopeToPolicy, apiDef.GetScopeToPolicyMapping())
		})
	}

}
