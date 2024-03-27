package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/TykTechnologies/tyk/swagger"
)

var licence = "https://github.com/TykTechnologies/tyk/blob/master/LICENSE.md"

func main() {
	r := openapi3.Reflector{
		Reflector: jsonschema.Reflector{},
	}
	r.DefaultOptions = append(r.DefaultOptions, jsonschema.StripDefinitionNamePrefix("Apidef", "Swagger", "Oas", "Gateway", "User"))

	r.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	r.Spec.WithServers(openapi3.Server{
		URL: "http://localhost:8080",
	})
	r.Spec.WithSecurity(map[string][]string{"api_key": {}})
	r.Spec.SetAPIKeySecurity("api_key", "X-Tyk-Authorization", openapi.InHeader, "Api key")
	r.Spec.Info.
		WithTitle("Tyk Gateway API").
		WithVersion("5.2.3").
		WithLicense(openapi3.License{Name: "Mozilla Public License Version 2.0", URL: &licence}).
		WithDescription(" The Tyk Gateway API is the primary means for integrating your application with the Tyk API Gateway")

	err := swagger.APIS(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.OasAPIS(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.Keys(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.OrgsApi(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.InvalidateCache(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.Certs(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.ReloadApi(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.SchemaAPi(&r)
	if err != nil {
		log.Fatal(err)
	}

	//err = swagger.DebugApi(&r)
	//if err != nil {
	//	return
	//}

	err = swagger.HealthEndpoint(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.PoliciesApis(&r)
	if err != nil {
		log.Fatal(err)
	}

	err = swagger.OAuthApi(&r)
	if err != nil {
		log.Fatal(err)
	}
	err = swagger.DebugApi(&r)
	if err != nil {
		log.Fatal(err)
	}

	schema, err := r.Spec.MarshalYAML()
	if err != nil {
		log.Fatal(err)
	}

	err = writeSchema(schema, "swagger.yaml")
	if err != nil {
		log.Fatal(err)
	}
}

func writeSchema(schema []byte, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(schema)
	return err
}