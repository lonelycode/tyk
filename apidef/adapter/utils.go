package adapter

import (
	"net/http"
)

func convertApiDefinitionHeadersToHttpHeaders(apiDefHeaders map[string]string) http.Header {
	if len(apiDefHeaders) == 0 {
		return nil
	}

	engineV2Headers := make(http.Header)
	for apiDefHeaderKey, apiDefHeaderValue := range apiDefHeaders {
		engineV2Headers.Add(apiDefHeaderKey, apiDefHeaderValue)
	}

	return engineV2Headers
}

func removeDuplicateApiDefinitionHeaders(headers ...map[string]string) map[string]string {
	hdr := make(map[string]string)
	// headers priority depends on the order of arguments
	for _, header := range headers {
		for k, v := range header {
			keyCanonical := http.CanonicalHeaderKey(k)
			if _, ok := hdr[keyCanonical]; ok {
				// skip because header is present
				continue
			}
			hdr[keyCanonical] = v
		}
	}
	return hdr
}
