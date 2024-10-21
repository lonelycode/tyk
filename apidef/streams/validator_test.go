package streams

import (
	"context"
	"embed"
	"fmt"
	"github.com/TykTechnologies/gojsonschema"
	"github.com/TykTechnologies/tyk/apidef/oas"
	"github.com/buger/jsonparser"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/warpstreamlabs/bento/public/components/io"
	_ "github.com/warpstreamlabs/bento/public/components/kafka"
	"github.com/warpstreamlabs/bento/public/service"
	"strings"
	"testing"
)

//go:embed testdata/*-oas-template.json
var oasTemplateFS embed.FS

func TestValidateTykStreamsOASObject(t *testing.T) {
	t.Parallel()
	validOASObject := oas.OAS{
		T: openapi3.T{
			OpenAPI:    "3.0.3",
			Info:       &openapi3.Info{},
			Paths:      map[string]*openapi3.PathItem{},
			Extensions: map[string]interface{}{},
		},
	}

	validXTykAPIStreaming := XTykStreaming{
		Info: oas.Info{
			Name: "test-streams",
			State: oas.State{
				Active: true,
			},
		},
		Server: oas.Server{
			ListenPath: oas.ListenPath{
				Value: "/test-streams",
			},
		},
		Streams: map[string]interface{}{},
	}

	validOASObject.Extensions[ExtensionTykStreaming] = &validXTykAPIStreaming

	validOAS3Definition, _ := validOASObject.MarshalJSON()

	t.Run("valid Tyk Streaming Extension", func(t *testing.T) {
		t.Parallel()
		err := ValidateOASObject(validOAS3Definition, "3.0.3")
		assert.Nil(t, err)
	})

	invalidOASObject := validOASObject
	invalidXTykAPIGateway := validXTykAPIStreaming
	invalidXTykAPIGateway.Info = oas.Info{}
	invalidXTykAPIGateway.Server.GatewayTags = &oas.GatewayTags{Enabled: true, Tags: []string{}}
	invalidOASObject.Extensions[ExtensionTykStreaming] = &invalidXTykAPIGateway
	invalidOAS3Definition, _ := invalidOASObject.MarshalJSON()

	t.Run("invalid OAS object", func(t *testing.T) {
		t.Parallel()
		err := ValidateOASObject(invalidOAS3Definition, "3.0.3")
		expectedErrs := []string{
			`x-tyk-streaming.info.name: Does not match pattern '\S+'`,
		}
		actualErrs := strings.Split(err.Error(), "\n")
		assert.ElementsMatch(t, expectedErrs, actualErrs)
	})

	var wrongTypedOASDefinition = []byte(`{
    "openapi": "3.0.0",
    "info": {
        "version": "1.0.0",
        "title": "Tyk Streams Example",
        "license": {
            "name": "MIT"
        }
    },
    "servers": [
        {
            "url": "http://tyk-gateway:8081/test-streams"
        }
    ],
    "paths": {},
    "x-tyk-streaming": {
        
    }
}`)

	t.Run("wrong typed OAS object", func(t *testing.T) {
		t.Parallel()
		err := ValidateOASObject(wrongTypedOASDefinition, "3.0.3")
		expectedErr := fmt.Sprintf("%s\n%s\n%s",
			"x-tyk-streaming: info is required",
			"x-tyk-streaming: streams is required",
			"x-tyk-streaming: server is required")
		assert.Equal(t, expectedErr, err.Error())
	})

	t.Run("should error when requested oas schema not found", func(t *testing.T) {
		t.Parallel()
		reqOASVersion := "4.0.3"
		err := ValidateOASObject(validOAS3Definition, reqOASVersion)
		expectedErr := fmt.Errorf(oasSchemaVersionNotFoundFmt, reqOASVersion)
		assert.Equal(t, expectedErr, err)
	})
}

func TestValidateOASTemplate(t *testing.T) {
	t.Run("empty x-tyk ext", func(t *testing.T) {
		body, err := oasTemplateFS.ReadFile("testdata/empty-x-tyk-ext-oas-template.json")
		require.NoError(t, err)
		err = ValidateOASTemplate(body, "")
		assert.NoError(t, err)
	})

	t.Run("non-empty x-tyk-ext", func(t *testing.T) {
		body, err := oasTemplateFS.ReadFile("testdata/non-empty-x-tyk-ext-oas-template.json")
		require.NoError(t, err)
		err = ValidateOASTemplate(body, "")
		assert.NoError(t, err)
	})
}

func Test_loadOASSchema(t *testing.T) {
	t.Parallel()
	t.Run("load Tyk Streams OAS", func(t *testing.T) {
		t.Parallel()
		err := loadOASSchema()
		assert.Nil(t, err)
		assert.NotNil(t, oasJSONSchemas)
		for oasVersion := range oasJSONSchemas {
			var xTykStreaming, xTykStreams []byte
			xTykStreaming, _, _, err = jsonparser.Get(oasJSONSchemas[oasVersion], keyProperties, ExtensionTykStreaming)
			assert.NoError(t, err)
			assert.NotNil(t, xTykStreaming)

			xTykStreams, _, _, err = jsonparser.Get(oasJSONSchemas[oasVersion], keyDefinitions, "X-Tyk-Streams")
			assert.NoError(t, err)
			assert.NotNil(t, xTykStreams)
		}
	})
}

func Test_findDefaultVersion(t *testing.T) {
	t.Parallel()
	t.Run("single version", func(t *testing.T) {
		rawVersions := []string{"3.0"}

		assert.Equal(t, "3.0", findDefaultVersion(rawVersions))
	})

	t.Run("multiple versions", func(t *testing.T) {
		rawVersions := []string{"3.0", "2.0", "3.1.0"}

		assert.Equal(t, "3.1", findDefaultVersion(rawVersions))
	})
}

func Test_setDefaultVersion(t *testing.T) {
	err := loadOASSchema()
	assert.NoError(t, err)

	setDefaultVersion()
	assert.Equal(t, "3.0", defaultVersion)
}

func TestGetOASSchema(t *testing.T) {
	err := loadOASSchema()
	assert.NoError(t, err)

	t.Run("return default version when req version is empty", func(t *testing.T) {
		_, err = GetOASSchema("")
		assert.NoError(t, err)
		assert.NotEmpty(t, oasJSONSchemas["3.0"])
	})

	t.Run("return minor version schema when req version is including patch version", func(t *testing.T) {
		_, err = GetOASSchema("3.0.8")
		assert.NoError(t, err)
		assert.NotEmpty(t, oasJSONSchemas["3.0"])
	})

	t.Run("return minor version 0 when only major version is requested", func(t *testing.T) {
		_, err = GetOASSchema("3")
		assert.NoError(t, err)
		assert.NotEmpty(t, oasJSONSchemas["3.0"])
	})

	t.Run("return error when non existing oas schema is requested", func(t *testing.T) {
		reqOASVersion := "4.0.3"
		_, err = GetOASSchema(reqOASVersion)
		expectedErr := fmt.Errorf(oasSchemaVersionNotFoundFmt, reqOASVersion)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("return error when requested version is not of semver", func(t *testing.T) {
		reqOASVersion := "a.0.3"
		_, err = GetOASSchema(reqOASVersion)
		expectedErr := fmt.Errorf("Malformed version: %s", reqOASVersion)
		assert.Equal(t, expectedErr, err)
	})
}

func Test_Bento_Config_Validation2(t *testing.T) {
	builder := service.NewStreamBuilder()
	data := `# Common config fields, showing default values
input:
  label: "foo"
  kafka:
    addresses: [] # No default (required)
    topics: [] # No default (required)
    target_version: 2.1.0 # No default (optional)
    consumer_group: ""
    checkpoint_limit: 1024
    auto_replay_nacks: sdfsd`

	err := builder.SetYAML(data)
	require.NoError(t, err)

	stream, err := builder.Build()
	require.NoError(t, err)

	// stream.Run returns the following error because `topics` array contains an invalid input.
	// failed to init input <no label> path root.input: it is not currently possible to include balanced and explicit partition topics in the same kafka input
	require.NoError(t, stream.Run(context.Background()))
}

func Test_Bento_Config_Validation(t *testing.T) {
	input := `{
	    "input": {
	        "label": "",
	        "kafka": {
	            "addresses": [],
	            "topics": [],
	            "target_version": "2.1.0",
	            "consumer_group": "",
	            "checkpoint_limit": 1024,
	            "auto_replay_nacks": true
	        },
			"foobar": {
				"barfoo": true
			}
	    },
"output": {
        "label": "",
    "drop_on": {
        "error": false,
        "error_patterns": [],
        "back_pressure": "30s",
        "output": null
    },
        "aws_sns": {
            "topic_arn": "",
            "message_group_id": "",
            "message_deduplication_id": "",
            "max_in_flight": 64,
            "metadata": {
                "exclude_prefixes": []
            }
        }
    }
	}`

	//sschema, err := service.GlobalEnvironment().FullConfigSchema("some-version", "date-built").MarshalJSONSchema()
	//require.NoError(t, err)
	//fmt.Println(string(sschema))

	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader([]byte(input))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	require.NoError(t, err)
	fmt.Println(result.Valid())
	fmt.Println(result.Errors())
}

var schema = []byte(`{
  "definitions": {
    "input": {
      "type": "object",
      "properties": {
        "label": {
          "type": "string"
        },
        "plugin": {
          "additionalProperties": false,
          "properties": {},
          "type": "object"
        },
        "processors": {
          "items": {
            "$ref": "#/definitions/processor"
          },
          "type": "array"
        },
        "type": {
          "type": "string"
        },
        "broker": {
          "additionalProperties": false,
          "properties": {
            "batching": {
              "additionalProperties": false,
              "properties": {
                "byte_size": {
                  "type": "number"
                },
                "check": {
                  "type": "string"
                },
                "count": {
                  "type": "number"
                },
                "period": {
                  "type": "string"
                },
                "processors": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "copies": {
              "type": "number"
            },
            "inputs": {
              "items": {
                "$ref": "#/definitions/input"
              },
              "type": "array"
            }
          },
          "required": [
            "inputs"
          ],
          "type": "object"
        },
        "http_client": {
          "additionalProperties": false,
          "properties": {
            "auto_replay_nacks": {
              "type": "boolean"
            },
            "backoff_on": {
              "items": {
                "type": "number"
              },
              "type": "array"
            },
            "basic_auth": {
              "additionalProperties": false,
              "properties": {
                "enabled": {
                  "type": "boolean"
                },
                "password": {
                  "type": "string"
                },
                "username": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "drop_empty_bodies": {
              "type": "boolean"
            },
            "drop_on": {
              "items": {
                "type": "number"
              },
              "type": "array"
            },
            "dump_request_log_level": {
              "type": "string"
            },
            "extract_headers": {
              "additionalProperties": false,
              "properties": {
                "include_patterns": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                },
                "include_prefixes": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "headers": {
              "patternProperties": {
                ".": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "jwt": {
              "additionalProperties": false,
              "properties": {
                "claims": {
                  "patternProperties": {
                    ".": {}
                  },
                  "type": "object"
                },
                "enabled": {
                  "type": "boolean"
                },
                "headers": {
                  "patternProperties": {
                    ".": {}
                  },
                  "type": "object"
                },
                "private_key_file": {
                  "type": "string"
                },
                "signing_method": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "max_retry_backoff": {
              "type": "string"
            },
            "metadata": {
              "additionalProperties": false,
              "properties": {
                "include_patterns": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                },
                "include_prefixes": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "oauth": {
              "additionalProperties": false,
              "properties": {
                "access_token": {
                  "type": "string"
                },
                "access_token_secret": {
                  "type": "string"
                },
                "consumer_key": {
                  "type": "string"
                },
                "consumer_secret": {
                  "type": "string"
                },
                "enabled": {
                  "type": "boolean"
                }
              },
              "type": "object"
            },
            "oauth2": {
              "additionalProperties": false,
              "properties": {
                "client_key": {
                  "type": "string"
                },
                "client_secret": {
                  "type": "string"
                },
                "enabled": {
                  "type": "boolean"
                },
                "endpoint_params": {
                  "patternProperties": {
                    ".": {}
                  },
                  "type": "object"
                },
                "scopes": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                },
                "token_url": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "payload": {
              "type": "string"
            },
            "proxy_url": {
              "type": "string"
            },
            "rate_limit": {
              "type": "string"
            },
            "retries": {
              "type": "number"
            },
            "retry_period": {
              "type": "string"
            },
            "stream": {
              "additionalProperties": false,
              "properties": {
                "codec": {
                  "type": "string"
                },
                "enabled": {
                  "type": "boolean"
                },
                "max_buffer": {
                  "type": "number"
                },
                "reconnect": {
                  "type": "boolean"
                },
                "scanner": {
                  "$ref": "#/definitions/scanner"
                }
              },
              "type": "object"
            },
            "successful_on": {
              "items": {
                "type": "number"
              },
              "type": "array"
            },
            "timeout": {
              "type": "string"
            },
            "tls": {
              "additionalProperties": false,
              "properties": {
                "client_certs": {
                  "items": {
                    "additionalProperties": false,
                    "properties": {
                      "cert": {
                        "type": "string"
                      },
                      "cert_file": {
                        "type": "string"
                      },
                      "key": {
                        "type": "string"
                      },
                      "key_file": {
                        "type": "string"
                      },
                      "password": {
                        "type": "string"
                      }
                    },
                    "type": "object"
                  },
                  "type": "array"
                },
                "enable_renegotiation": {
                  "type": "boolean"
                },
                "enabled": {
                  "type": "boolean"
                },
                "root_cas": {
                  "type": "string"
                },
                "root_cas_file": {
                  "type": "string"
                },
                "skip_cert_verify": {
                  "type": "boolean"
                }
              },
              "type": "object"
            },
            "url": {
              "type": "string"
            },
            "verb": {
              "type": "string"
            }
          },
          "required": [
            "url"
          ],
          "type": "object"
        },
        "http_server": {
          "additionalProperties": false,
          "properties": {
            "address": {
              "type": "string"
            },
            "allowed_verbs": {
              "items": {
                "type": "string"
              },
              "type": "array"
            },
            "cert_file": {
              "type": "string"
            },
            "cors": {
              "additionalProperties": false,
              "properties": {
                "allowed_origins": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                },
                "enabled": {
                  "type": "boolean"
                }
              },
              "type": "object"
            },
            "key_file": {
              "type": "string"
            },
            "path": {
              "type": "string"
            },
            "rate_limit": {
              "type": "string"
            },
            "sync_response": {
              "additionalProperties": false,
              "properties": {
                "headers": {
                  "patternProperties": {
                    ".": {
                      "type": "string"
                    }
                  },
                  "type": "object"
                },
                "metadata_headers": {
                  "additionalProperties": false,
                  "properties": {
                    "include_patterns": {
                      "items": {
                        "type": "string"
                      },
                      "type": "array"
                    },
                    "include_prefixes": {
                      "items": {
                        "type": "string"
                      },
                      "type": "array"
                    }
                  },
                  "type": "object"
                },
                "status": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "timeout": {
              "type": "string"
            },
            "ws_path": {
              "type": "string"
            },
            "ws_rate_limit_message": {
              "type": "string"
            },
            "ws_welcome_message": {
              "type": "string"
            }
          },
          "type": "object"
        },
        "kafka": {
          "additionalProperties": false,
          "properties": {
            "addresses": {
              "items": {
                "type": "string"
              },
              "type": "array"
            },
            "auto_replay_nacks": {
              "type": "boolean"
            },
            "batching": {
              "additionalProperties": false,
              "properties": {
                "byte_size": {
                  "type": "number"
                },
                "check": {
                  "type": "string"
                },
                "count": {
                  "type": "number"
                },
                "period": {
                  "type": "string"
                },
                "processors": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "checkpoint_limit": {
              "type": "number"
            },
            "client_id": {
              "type": "string"
            },
            "commit_period": {
              "type": "string"
            },
            "consumer_group": {
              "type": "string"
            },
            "extract_tracing_map": {
              "type": "string"
            },
            "fetch_buffer_cap": {
              "type": "number"
            },
            "group": {
              "additionalProperties": false,
              "properties": {
                "heartbeat_interval": {
                  "type": "string"
                },
                "rebalance_timeout": {
                  "type": "string"
                },
                "session_timeout": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "max_processing_period": {
              "type": "string"
            },
            "multi_header": {
              "type": "boolean"
            },
            "rack_id": {
              "type": "string"
            },
            "sasl": {
              "additionalProperties": false,
              "properties": {
                "access_token": {
                  "type": "string"
                },
                "mechanism": {
                  "type": "string"
                },
                "password": {
                  "type": "string"
                },
                "token_cache": {
                  "type": "string"
                },
                "token_key": {
                  "type": "string"
                },
                "user": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "start_from_oldest": {
              "type": "boolean"
            },
            "target_version": {
              "type": "string"
            },
            "tls": {
              "additionalProperties": false,
              "properties": {
                "client_certs": {
                  "items": {
                    "additionalProperties": false,
                    "properties": {
                      "cert": {
                        "type": "string"
                      },
                      "cert_file": {
                        "type": "string"
                      },
                      "key": {
                        "type": "string"
                      },
                      "key_file": {
                        "type": "string"
                      },
                      "password": {
                        "type": "string"
                      }
                    },
                    "type": "object"
                  },
                  "type": "array"
                },
                "enable_renegotiation": {
                  "type": "boolean"
                },
                "enabled": {
                  "type": "boolean"
                },
                "root_cas": {
                  "type": "string"
                },
                "root_cas_file": {
                  "type": "string"
                },
                "skip_cert_verify": {
                  "type": "boolean"
                }
              },
              "type": "object"
            },
            "topics": {
              "items": {
                "type": "string"
              },
              "type": "array"
            }
          },
          "required": [
            "addresses",
            "topics"
          ],
          "type": "object"
        }
      }
    },
    "output": {
      "type": "object",
      "properties": {
        "label": {
          "type": "string"
        },
        "plugin": {
          "additionalProperties": false,
          "properties": {},
          "type": "object"
        },
        "processors": {
          "items": {
            "$ref": "#/definitions/processor"
          },
          "type": "array"
        },
        "type": {
          "type": "string"
        },
        "broker": {
          "additionalProperties": false,
          "properties": {
            "batching": {
              "additionalProperties": false,
              "properties": {
                "byte_size": {
                  "type": "number"
                },
                "check": {
                  "type": "string"
                },
                "count": {
                  "type": "number"
                },
                "period": {
                  "type": "string"
                },
                "processors": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "copies": {
              "type": "number"
            },
            "outputs": {
              "items": {
                "$ref": "#/definitions/output"
              },
              "type": "array"
            },
            "pattern": {
              "type": "string"
            }
          },
          "required": [
            "outputs"
          ],
          "type": "object"
        },
        "http_client": {
          "additionalProperties": false,
          "properties": {
            "backoff_on": {
              "items": {
                "type": "number"
              },
              "type": "array"
            },
            "basic_auth": {
              "additionalProperties": false,
              "properties": {
                "enabled": {
                  "type": "boolean"
                },
                "password": {
                  "type": "string"
                },
                "username": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "batch_as_multipart": {
              "type": "boolean"
            },
            "batching": {
              "additionalProperties": false,
              "properties": {
                "byte_size": {
                  "type": "number"
                },
                "check": {
                  "type": "string"
                },
                "count": {
                  "type": "number"
                },
                "period": {
                  "type": "string"
                },
                "processors": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "drop_on": {
              "items": {
                "type": "number"
              },
              "type": "array"
            },
            "dump_request_log_level": {
              "type": "string"
            },
            "extract_headers": {
              "additionalProperties": false,
              "properties": {
                "include_patterns": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                },
                "include_prefixes": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "headers": {
              "patternProperties": {
                ".": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "jwt": {
              "additionalProperties": false,
              "properties": {
                "claims": {
                  "patternProperties": {
                    ".": {}
                  },
                  "type": "object"
                },
                "enabled": {
                  "type": "boolean"
                },
                "headers": {
                  "patternProperties": {
                    ".": {}
                  },
                  "type": "object"
                },
                "private_key_file": {
                  "type": "string"
                },
                "signing_method": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "max_in_flight": {
              "type": "number"
            },
            "max_retry_backoff": {
              "type": "string"
            },
            "metadata": {
              "additionalProperties": false,
              "properties": {
                "include_patterns": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                },
                "include_prefixes": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "multipart": {
              "items": {
                "additionalProperties": false,
                "properties": {
                  "body": {
                    "type": "string"
                  },
                  "content_disposition": {
                    "type": "string"
                  },
                  "content_type": {
                    "type": "string"
                  }
                },
                "type": "object"
              },
              "type": "array"
            },
            "oauth": {
              "additionalProperties": false,
              "properties": {
                "access_token": {
                  "type": "string"
                },
                "access_token_secret": {
                  "type": "string"
                },
                "consumer_key": {
                  "type": "string"
                },
                "consumer_secret": {
                  "type": "string"
                },
                "enabled": {
                  "type": "boolean"
                }
              },
              "type": "object"
            },
            "oauth2": {
              "additionalProperties": false,
              "properties": {
                "client_key": {
                  "type": "string"
                },
                "client_secret": {
                  "type": "string"
                },
                "enabled": {
                  "type": "boolean"
                },
                "endpoint_params": {
                  "patternProperties": {
                    ".": {}
                  },
                  "type": "object"
                },
                "scopes": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                },
                "token_url": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "propagate_response": {
              "type": "boolean"
            },
            "proxy_url": {
              "type": "string"
            },
            "rate_limit": {
              "type": "string"
            },
            "retries": {
              "type": "number"
            },
            "retry_period": {
              "type": "string"
            },
            "successful_on": {
              "items": {
                "type": "number"
              },
              "type": "array"
            },
            "timeout": {
              "type": "string"
            },
            "tls": {
              "additionalProperties": false,
              "properties": {
                "client_certs": {
                  "items": {
                    "additionalProperties": false,
                    "properties": {
                      "cert": {
                        "type": "string"
                      },
                      "cert_file": {
                        "type": "string"
                      },
                      "key": {
                        "type": "string"
                      },
                      "key_file": {
                        "type": "string"
                      },
                      "password": {
                        "type": "string"
                      }
                    },
                    "type": "object"
                  },
                  "type": "array"
                },
                "enable_renegotiation": {
                  "type": "boolean"
                },
                "enabled": {
                  "type": "boolean"
                },
                "root_cas": {
                  "type": "string"
                },
                "root_cas_file": {
                  "type": "string"
                },
                "skip_cert_verify": {
                  "type": "boolean"
                }
              },
              "type": "object"
            },
            "url": {
              "type": "string"
            },
            "verb": {
              "type": "string"
            }
          },
          "required": [
            "url"
          ],
          "type": "object"
        },
        "http_server": {
          "additionalProperties": false,
          "properties": {
            "address": {
              "type": "string"
            },
            "allowed_verbs": {
              "items": {
                "type": "string"
              },
              "type": "array"
            },
            "cert_file": {
              "type": "string"
            },
            "cors": {
              "additionalProperties": false,
              "properties": {
                "allowed_origins": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                },
                "enabled": {
                  "type": "boolean"
                }
              },
              "type": "object"
            },
            "key_file": {
              "type": "string"
            },
            "path": {
              "type": "string"
            },
            "stream_path": {
              "type": "string"
            },
            "timeout": {
              "type": "string"
            },
            "ws_path": {
              "type": "string"
            }
          },
          "type": "object"
        },
        "kafka": {
          "additionalProperties": false,
          "properties": {
            "ack_replicas": {
              "type": "boolean"
            },
            "addresses": {
              "items": {
                "type": "string"
              },
              "type": "array"
            },
            "backoff": {
              "additionalProperties": false,
              "properties": {
                "initial_interval": {
                  "type": "string"
                },
                "max_elapsed_time": {
                  "type": "string"
                },
                "max_interval": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "batching": {
              "additionalProperties": false,
              "properties": {
                "byte_size": {
                  "type": "number"
                },
                "check": {
                  "type": "string"
                },
                "count": {
                  "type": "number"
                },
                "period": {
                  "type": "string"
                },
                "processors": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "client_id": {
              "type": "string"
            },
            "compression": {
              "type": "string"
            },
            "custom_topic_creation": {
              "additionalProperties": false,
              "properties": {
                "enabled": {
                  "type": "boolean"
                },
                "partitions": {
                  "type": "number"
                },
                "replication_factor": {
                  "type": "number"
                }
              },
              "type": "object"
            },
            "idempotent_write": {
              "type": "boolean"
            },
            "inject_tracing_map": {
              "type": "string"
            },
            "key": {
              "type": "string"
            },
            "max_in_flight": {
              "type": "number"
            },
            "max_msg_bytes": {
              "type": "number"
            },
            "max_retries": {
              "type": "number"
            },
            "metadata": {
              "additionalProperties": false,
              "properties": {
                "exclude_prefixes": {
                  "items": {
                    "type": "string"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            "partition": {
              "type": "string"
            },
            "partitioner": {
              "type": "string"
            },
            "rack_id": {
              "type": "string"
            },
            "retry_as_batch": {
              "type": "boolean"
            },
            "sasl": {
              "additionalProperties": false,
              "properties": {
                "access_token": {
                  "type": "string"
                },
                "mechanism": {
                  "type": "string"
                },
                "password": {
                  "type": "string"
                },
                "token_cache": {
                  "type": "string"
                },
                "token_key": {
                  "type": "string"
                },
                "user": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "static_headers": {
              "patternProperties": {
                ".": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            "target_version": {
              "type": "string"
            },
            "timeout": {
              "type": "string"
            },
            "tls": {
              "additionalProperties": false,
              "properties": {
                "client_certs": {
                  "items": {
                    "additionalProperties": false,
                    "properties": {
                      "cert": {
                        "type": "string"
                      },
                      "cert_file": {
                        "type": "string"
                      },
                      "key": {
                        "type": "string"
                      },
                      "key_file": {
                        "type": "string"
                      },
                      "password": {
                        "type": "string"
                      }
                    },
                    "type": "object"
                  },
                  "type": "array"
                },
                "enable_renegotiation": {
                  "type": "boolean"
                },
                "enabled": {
                  "type": "boolean"
                },
                "root_cas": {
                  "type": "string"
                },
                "root_cas_file": {
                  "type": "string"
                },
                "skip_cert_verify": {
                  "type": "boolean"
                }
              },
              "type": "object"
            },
            "topic": {
              "type": "string"
            }
          },
          "required": [
            "addresses",
            "topic"
          ],
          "type": "object"
        }
      }
    },
    "processor": {
      "allOf": [
        {
          "anyOf": [
            {
              "properties": {
                "archive": {
                  "additionalProperties": false,
                  "properties": {
                    "format": {
                      "type": "string"
                    },
                    "path": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "format"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "bloblang": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "bounds_check": {
                  "additionalProperties": false,
                  "properties": {
                    "max_part_size": {
                      "type": "number"
                    },
                    "max_parts": {
                      "type": "number"
                    },
                    "min_part_size": {
                      "type": "number"
                    },
                    "min_parts": {
                      "type": "number"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "branch": {
                  "additionalProperties": false,
                  "properties": {
                    "processors": {
                      "items": {
                        "$ref": "#/definitions/processor"
                      },
                      "type": "array"
                    },
                    "request_map": {
                      "type": "string"
                    },
                    "result_map": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "processors"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "cache": {
                  "additionalProperties": false,
                  "properties": {
                    "key": {
                      "type": "string"
                    },
                    "operator": {
                      "type": "string"
                    },
                    "resource": {
                      "type": "string"
                    },
                    "ttl": {
                      "type": "string"
                    },
                    "value": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "resource",
                    "operator",
                    "key"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "cached": {
                  "additionalProperties": false,
                  "properties": {
                    "cache": {
                      "type": "string"
                    },
                    "key": {
                      "type": "string"
                    },
                    "processors": {
                      "items": {
                        "$ref": "#/definitions/processor"
                      },
                      "type": "array"
                    },
                    "skip_on": {
                      "type": "string"
                    },
                    "ttl": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "cache",
                    "key",
                    "processors"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "catch": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "command": {
                  "additionalProperties": false,
                  "properties": {
                    "args_mapping": {
                      "type": "string"
                    },
                    "name": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "name"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "compress": {
                  "additionalProperties": false,
                  "properties": {
                    "algorithm": {
                      "type": "string"
                    },
                    "level": {
                      "type": "number"
                    }
                  },
                  "required": [
                    "algorithm"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "decompress": {
                  "additionalProperties": false,
                  "properties": {
                    "algorithm": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "algorithm"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "dedupe": {
                  "additionalProperties": false,
                  "properties": {
                    "cache": {
                      "type": "string"
                    },
                    "drop_on_err": {
                      "type": "boolean"
                    },
                    "key": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "cache",
                    "key"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "for_each": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "grok": {
                  "additionalProperties": false,
                  "properties": {
                    "expressions": {
                      "items": {
                        "type": "string"
                      },
                      "type": "array"
                    },
                    "named_captures_only": {
                      "type": "boolean"
                    },
                    "pattern_definitions": {
                      "patternProperties": {
                        ".": {
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "pattern_paths": {
                      "items": {
                        "type": "string"
                      },
                      "type": "array"
                    },
                    "remove_empty_values": {
                      "type": "boolean"
                    },
                    "use_default_patterns": {
                      "type": "boolean"
                    }
                  },
                  "required": [
                    "expressions"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "group_by": {
                  "items": {
                    "additionalProperties": false,
                    "properties": {
                      "check": {
                        "type": "string"
                      },
                      "processors": {
                        "items": {
                          "$ref": "#/definitions/processor"
                        },
                        "type": "array"
                      }
                    },
                    "required": [
                      "check"
                    ],
                    "type": "object"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "group_by_value": {
                  "additionalProperties": false,
                  "properties": {
                    "value": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "value"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "http": {
                  "additionalProperties": false,
                  "properties": {
                    "backoff_on": {
                      "items": {
                        "type": "number"
                      },
                      "type": "array"
                    },
                    "basic_auth": {
                      "additionalProperties": false,
                      "properties": {
                        "enabled": {
                          "type": "boolean"
                        },
                        "password": {
                          "type": "string"
                        },
                        "username": {
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "batch_as_multipart": {
                      "type": "boolean"
                    },
                    "drop_on": {
                      "items": {
                        "type": "number"
                      },
                      "type": "array"
                    },
                    "dump_request_log_level": {
                      "type": "string"
                    },
                    "extract_headers": {
                      "additionalProperties": false,
                      "properties": {
                        "include_patterns": {
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        },
                        "include_prefixes": {
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        }
                      },
                      "type": "object"
                    },
                    "headers": {
                      "patternProperties": {
                        ".": {
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "jwt": {
                      "additionalProperties": false,
                      "properties": {
                        "claims": {
                          "patternProperties": {
                            ".": {}
                          },
                          "type": "object"
                        },
                        "enabled": {
                          "type": "boolean"
                        },
                        "headers": {
                          "patternProperties": {
                            ".": {}
                          },
                          "type": "object"
                        },
                        "private_key_file": {
                          "type": "string"
                        },
                        "signing_method": {
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "max_retry_backoff": {
                      "type": "string"
                    },
                    "metadata": {
                      "additionalProperties": false,
                      "properties": {
                        "include_patterns": {
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        },
                        "include_prefixes": {
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        }
                      },
                      "type": "object"
                    },
                    "oauth": {
                      "additionalProperties": false,
                      "properties": {
                        "access_token": {
                          "type": "string"
                        },
                        "access_token_secret": {
                          "type": "string"
                        },
                        "consumer_key": {
                          "type": "string"
                        },
                        "consumer_secret": {
                          "type": "string"
                        },
                        "enabled": {
                          "type": "boolean"
                        }
                      },
                      "type": "object"
                    },
                    "oauth2": {
                      "additionalProperties": false,
                      "properties": {
                        "client_key": {
                          "type": "string"
                        },
                        "client_secret": {
                          "type": "string"
                        },
                        "enabled": {
                          "type": "boolean"
                        },
                        "endpoint_params": {
                          "patternProperties": {
                            ".": {}
                          },
                          "type": "object"
                        },
                        "scopes": {
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        },
                        "token_url": {
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "parallel": {
                      "type": "boolean"
                    },
                    "proxy_url": {
                      "type": "string"
                    },
                    "rate_limit": {
                      "type": "string"
                    },
                    "retries": {
                      "type": "number"
                    },
                    "retry_period": {
                      "type": "string"
                    },
                    "successful_on": {
                      "items": {
                        "type": "number"
                      },
                      "type": "array"
                    },
                    "timeout": {
                      "type": "string"
                    },
                    "tls": {
                      "additionalProperties": false,
                      "properties": {
                        "client_certs": {
                          "items": {
                            "additionalProperties": false,
                            "properties": {
                              "cert": {
                                "type": "string"
                              },
                              "cert_file": {
                                "type": "string"
                              },
                              "key": {
                                "type": "string"
                              },
                              "key_file": {
                                "type": "string"
                              },
                              "password": {
                                "type": "string"
                              }
                            },
                            "type": "object"
                          },
                          "type": "array"
                        },
                        "enable_renegotiation": {
                          "type": "boolean"
                        },
                        "enabled": {
                          "type": "boolean"
                        },
                        "root_cas": {
                          "type": "string"
                        },
                        "root_cas_file": {
                          "type": "string"
                        },
                        "skip_cert_verify": {
                          "type": "boolean"
                        }
                      },
                      "type": "object"
                    },
                    "url": {
                      "type": "string"
                    },
                    "verb": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "url"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "insert_part": {
                  "additionalProperties": false,
                  "properties": {
                    "content": {
                      "type": "string"
                    },
                    "index": {
                      "type": "number"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "jmespath": {
                  "additionalProperties": false,
                  "properties": {
                    "query": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "query"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "jq": {
                  "additionalProperties": false,
                  "properties": {
                    "output_raw": {
                      "type": "boolean"
                    },
                    "query": {
                      "type": "string"
                    },
                    "raw": {
                      "type": "boolean"
                    }
                  },
                  "required": [
                    "query"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "json_schema": {
                  "additionalProperties": false,
                  "properties": {
                    "schema": {
                      "type": "string"
                    },
                    "schema_path": {
                      "type": "string"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "log": {
                  "additionalProperties": false,
                  "properties": {
                    "fields": {
                      "patternProperties": {
                        ".": {
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "fields_mapping": {
                      "type": "string"
                    },
                    "level": {
                      "type": "string"
                    },
                    "message": {
                      "type": "string"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "mapping": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "metric": {
                  "additionalProperties": false,
                  "properties": {
                    "labels": {
                      "patternProperties": {
                        ".": {
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "name": {
                      "type": "string"
                    },
                    "type": {
                      "type": "string"
                    },
                    "value": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "type",
                    "name"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "mutation": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "noop": {
                  "additionalProperties": false,
                  "properties": {},
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "parallel": {
                  "additionalProperties": false,
                  "properties": {
                    "cap": {
                      "type": "number"
                    },
                    "processors": {
                      "items": {
                        "$ref": "#/definitions/processor"
                      },
                      "type": "array"
                    }
                  },
                  "required": [
                    "processors"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "parse_log": {
                  "additionalProperties": false,
                  "properties": {
                    "allow_rfc3339": {
                      "type": "boolean"
                    },
                    "best_effort": {
                      "type": "boolean"
                    },
                    "codec": {
                      "type": "string"
                    },
                    "default_timezone": {
                      "type": "string"
                    },
                    "default_year": {
                      "type": "string"
                    },
                    "format": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "format",
                    "codec"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "processors": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "rate_limit": {
                  "additionalProperties": false,
                  "properties": {
                    "resource": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "resource"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "resource": {
                  "type": "string"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "retry": {
                  "additionalProperties": false,
                  "properties": {
                    "backoff": {
                      "additionalProperties": false,
                      "properties": {
                        "initial_interval": {
                          "type": "string"
                        },
                        "max_elapsed_time": {
                          "type": "string"
                        },
                        "max_interval": {
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "max_retries": {
                      "type": "number"
                    },
                    "parallel": {
                      "type": "boolean"
                    },
                    "processors": {
                      "items": {
                        "$ref": "#/definitions/processor"
                      },
                      "type": "array"
                    }
                  },
                  "required": [
                    "processors"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "select_parts": {
                  "additionalProperties": false,
                  "properties": {
                    "parts": {
                      "items": {
                        "type": "number"
                      },
                      "type": "array"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "sleep": {
                  "additionalProperties": false,
                  "properties": {
                    "duration": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "duration"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "split": {
                  "additionalProperties": false,
                  "properties": {
                    "byte_size": {
                      "type": "number"
                    },
                    "size": {
                      "type": "number"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "subprocess": {
                  "additionalProperties": false,
                  "properties": {
                    "args": {
                      "items": {
                        "type": "string"
                      },
                      "type": "array"
                    },
                    "codec_recv": {
                      "type": "string"
                    },
                    "codec_send": {
                      "type": "string"
                    },
                    "max_buffer": {
                      "type": "number"
                    },
                    "name": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "name"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "switch": {
                  "items": {
                    "additionalProperties": false,
                    "properties": {
                      "check": {
                        "type": "string"
                      },
                      "fallthrough": {
                        "type": "boolean"
                      },
                      "processors": {
                        "items": {
                          "$ref": "#/definitions/processor"
                        },
                        "type": "array"
                      }
                    },
                    "type": "object"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "sync_response": {
                  "additionalProperties": false,
                  "properties": {},
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "try": {
                  "items": {
                    "$ref": "#/definitions/processor"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "unarchive": {
                  "additionalProperties": false,
                  "properties": {
                    "format": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "format"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "while": {
                  "additionalProperties": false,
                  "properties": {
                    "at_least_once": {
                      "type": "boolean"
                    },
                    "check": {
                      "type": "string"
                    },
                    "max_loops": {
                      "type": "number"
                    },
                    "processors": {
                      "items": {
                        "$ref": "#/definitions/processor"
                      },
                      "type": "array"
                    }
                  },
                  "required": [
                    "processors"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "workflow": {
                  "additionalProperties": false,
                  "properties": {
                    "branch_resources": {
                      "items": {
                        "type": "string"
                      },
                      "type": "array"
                    },
                    "branches": {
                      "patternProperties": {
                        ".": {
                          "additionalProperties": false,
                          "properties": {
                            "processors": {
                              "items": {
                                "$ref": "#/definitions/processor"
                              },
                              "type": "array"
                            },
                            "request_map": {
                              "type": "string"
                            },
                            "result_map": {
                              "type": "string"
                            }
                          },
                          "required": [
                            "processors"
                          ],
                          "type": "object"
                        }
                      },
                      "type": "object"
                    },
                    "meta_path": {
                      "type": "string"
                    },
                    "order": {
                      "items": {
                        "items": {
                          "type": "string"
                        },
                        "type": "array"
                      },
                      "type": "array"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            }
          ]
        },
        {
          "properties": {
            "label": {
              "type": "string"
            },
            "plugin": {
              "additionalProperties": false,
              "properties": {},
              "type": "object"
            },
            "type": {
              "type": "string"
            }
          },
          "type": "object"
        }
      ]
    },
    "scanner": {
      "allOf": [
        {
          "anyOf": [
            {
              "properties": {
                "chunker": {
                  "additionalProperties": false,
                  "properties": {
                    "size": {
                      "type": "number"
                    }
                  },
                  "required": [
                    "size"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "csv": {
                  "additionalProperties": false,
                  "properties": {
                    "continue_on_error": {
                      "type": "boolean"
                    },
                    "custom_delimiter": {
                      "type": "string"
                    },
                    "lazy_quotes": {
                      "type": "boolean"
                    },
                    "parse_header_row": {
                      "type": "boolean"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "decompress": {
                  "additionalProperties": false,
                  "properties": {
                    "algorithm": {
                      "type": "string"
                    },
                    "into": {
                      "$ref": "#/definitions/scanner"
                    }
                  },
                  "required": [
                    "algorithm"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "json_documents": {
                  "additionalProperties": false,
                  "properties": {},
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "lines": {
                  "additionalProperties": false,
                  "properties": {
                    "custom_delimiter": {
                      "type": "string"
                    },
                    "max_buffer_size": {
                      "type": "number"
                    },
                    "omit_empty": {
                      "type": "boolean"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "re_match": {
                  "additionalProperties": false,
                  "properties": {
                    "max_buffer_size": {
                      "type": "number"
                    },
                    "pattern": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "pattern"
                  ],
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "skip_bom": {
                  "additionalProperties": false,
                  "properties": {
                    "into": {
                      "$ref": "#/definitions/scanner"
                    }
                  },
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "switch": {
                  "items": {
                    "additionalProperties": false,
                    "properties": {
                      "re_match_name": {
                        "type": "string"
                      },
                      "scanner": {
                        "$ref": "#/definitions/scanner"
                      }
                    },
                    "required": [
                      "scanner"
                    ],
                    "type": "object"
                  },
                  "type": "array"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "tar": {
                  "additionalProperties": false,
                  "properties": {},
                  "type": "object"
                }
              },
              "type": "object"
            },
            {
              "properties": {
                "to_the_end": {
                  "additionalProperties": false,
                  "properties": {},
                  "type": "object"
                }
              },
              "type": "object"
            }
          ]
        },
        {
          "properties": {
            "plugin": {
              "additionalProperties": false,
              "properties": {},
              "type": "object"
            },
            "type": {
              "type": "string"
            }
          },
          "type": "object"
        }
      ]
    }
  },
  "properties": {
    "http": {
      "additionalProperties": false,
      "properties": {
        "address": {
          "type": "string"
        },
        "basic_auth": {
          "additionalProperties": false,
          "properties": {
            "algorithm": {
              "type": "string"
            },
            "enabled": {
              "type": "boolean"
            },
            "password_hash": {
              "type": "string"
            },
            "realm": {
              "type": "string"
            },
            "salt": {
              "type": "string"
            },
            "username": {
              "type": "string"
            }
          },
          "type": "object"
        },
        "cert_file": {
          "type": "string"
        },
        "cors": {
          "additionalProperties": false,
          "properties": {
            "allowed_origins": {
              "items": {
                "type": "string"
              },
              "type": "array"
            },
            "enabled": {
              "type": "boolean"
            }
          },
          "type": "object"
        },
        "debug_endpoints": {
          "type": "boolean"
        },
        "enabled": {
          "type": "boolean"
        },
        "key_file": {
          "type": "string"
        },
        "root_path": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "input": {
      "$ref": "#/definitions/input"
    },
    "input_resources": {
      "items": {
        "$ref": "#/definitions/input"
      },
      "type": "array"
    },
    "output": {
      "$ref": "#/definitions/output"
    },
    "output_resources": {
      "items": {
        "$ref": "#/definitions/output"
      },
      "type": "array"
    },
    "processor_resources": {
      "items": {
        "$ref": "#/definitions/processor"
      },
      "type": "array"
    },
    "shutdown_delay": {
      "type": "string"
    },
    "shutdown_timeout": {
      "type": "string"
    }
  }
}`)
