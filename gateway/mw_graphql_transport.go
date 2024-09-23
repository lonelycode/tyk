package gateway

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/TykTechnologies/tyk/apidef"
)

type GraphQLEngineTransportType int

const (
	GraphQLEngineTransportTypeProxyOnly GraphQLEngineTransportType = iota
	GraphQLEngineTransportTypeMultiUpstream
)

func DetermineGraphQLEngineTransportType(apiSpec *APISpec) GraphQLEngineTransportType {
	switch apiSpec.GraphQL.ExecutionMode {
	case apidef.GraphQLExecutionModeSubgraph:
		fallthrough
	case apidef.GraphQLExecutionModeProxyOnly:
		return GraphQLEngineTransportTypeProxyOnly
	}

	return GraphQLEngineTransportTypeMultiUpstream
}

type contextKey struct{}

var graphqlProxyContextInfo = contextKey{}

type GraphQLProxyOnlyContextValues struct {
	forwardedRequest       *http.Request
	upstreamResponse       *http.Response
	ignoreForwardedHeaders map[string]bool
}

func SetProxyOnlyContextValue(ctx context.Context, req *http.Request) context.Context {
	value := &GraphQLProxyOnlyContextValues{
		forwardedRequest: req,
		ignoreForwardedHeaders: map[string]bool{
			http.CanonicalHeaderKey("date"):           true,
			http.CanonicalHeaderKey("content-type"):   true,
			http.CanonicalHeaderKey("content-length"): true,
		},
	}

	return context.WithValue(ctx, graphqlProxyContextInfo, value)
}

func GetProxyOnlyContextValue(ctx context.Context) *GraphQLProxyOnlyContextValues {
	val, ok := ctx.Value(graphqlProxyContextInfo).(*GraphQLProxyOnlyContextValues)
	if !ok {
		return nil
	}
	return val
}

type GraphQLProxyOnlyContext struct {
	context.Context
	forwardedRequest       *http.Request
	upstreamResponse       *http.Response
	ignoreForwardedHeaders map[string]bool
}

func NewGraphQLProxyOnlyContext(ctx context.Context, forwardedRequest *http.Request) *GraphQLProxyOnlyContext {
	return &GraphQLProxyOnlyContext{
		Context:          ctx,
		forwardedRequest: forwardedRequest,
		ignoreForwardedHeaders: map[string]bool{
			http.CanonicalHeaderKey("date"):           true,
			http.CanonicalHeaderKey("content-type"):   true,
			http.CanonicalHeaderKey("content-length"): true,
		},
	}
}

func (g *GraphQLProxyOnlyContext) Response() *http.Response {
	return g.upstreamResponse
}

type GraphQLEngineTransport struct {
	originalTransport http.RoundTripper
	transportType     GraphQLEngineTransportType
}

func NewGraphQLEngineTransport(transportType GraphQLEngineTransportType, originalTransport http.RoundTripper) *GraphQLEngineTransport {
	return &GraphQLEngineTransport{
		originalTransport: originalTransport,
		transportType:     transportType,
	}
}

func (g *GraphQLEngineTransport) RoundTrip(request *http.Request) (res *http.Response, err error) {
	switch g.transportType {
	case GraphQLEngineTransportTypeProxyOnly:
		val := GetProxyOnlyContextValue(request.Context())
		if val != nil {
			return g.handleProxyOnly(val, request)
		}
	}

	return g.originalTransport.RoundTrip(request)
}

func (g *GraphQLEngineTransport) handleProxyOnly(proxyOnlyCtx *GraphQLProxyOnlyContextValues, request *http.Request) (*http.Response, error) {
	request.Method = proxyOnlyCtx.forwardedRequest.Method
	g.setProxyOnlyHeaders(proxyOnlyCtx, request)

	response, err := g.originalTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		// In proxy-only mode, we keep the upstream error message to
		// insert into the library's error message.
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = response.Body.Close()
		}()
		// graphql-go-tools uses response.body to resolve the upstream response.
		// It's not possible to re-use io.ReadCloser. Because of that, we keep the
		// original error message for later use.
		// See TT-7808
		reusableBody, err := newNopCloserBuffer(io.NopCloser(bytes.NewReader(body)))
		if err != nil {
			return nil, err
		}
		response.Body = reusableBody
	}
	proxyOnlyCtx.upstreamResponse = response
	return response, err
}

func (g *GraphQLEngineTransport) setProxyOnlyHeaders(proxyOnlyCtx *GraphQLProxyOnlyContextValues, r *http.Request) {
	for forwardedHeaderKey, forwardedHeaderValues := range proxyOnlyCtx.forwardedRequest.Header {
		if proxyOnlyCtx.ignoreForwardedHeaders[forwardedHeaderKey] {
			continue
		}

		for _, forwardedHeaderValue := range forwardedHeaderValues {
			r.Header.Add(forwardedHeaderKey, forwardedHeaderValue)
		}
	}
}