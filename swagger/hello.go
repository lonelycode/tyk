package swagger

import (
	"net/http"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/TykTechnologies/tyk/apidef"
)

const (
	helloTag     = "Health Checking"
	helloTagDesc = `Check health status of the Gateway and loaded APIs
`
)

// Done
func HealthEndpoint(r *openapi3.Reflector) error {
	addTag(helloTag, helloTagDesc)
	oc, err := r.NewOperationContext(http.MethodGet, "/tyk/hello")
	if err != nil {
		return err
	}
	oc.SetTags(helloTag)
	oc.SetID("hello")
	oc.SetSummary("Check the Health of the Tyk Gateway")
	oc.SetDescription("From v2.7.5 you can now rename the `/hello`  endpoint by using the `health_check_endpoint_name` option\n        \n        Returns 200 response in case of success")
	oc.AddRespStructure(new(apidef.HealthCheckResponse), openapi.WithHTTPStatus(http.StatusOK), func(cu *openapi.ContentUnit) {
		cu.Description = "Success"
	})
	oc.AddRespStructure(new(apiStatusMessage), openapi.WithHTTPStatus(http.StatusMethodNotAllowed))
	return r.AddOperation(oc)
}