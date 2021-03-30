// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Request settlement discrepancy
	// (PUT /settlements/{settlementId})
	CalculateSettlementDiscrepancy(ctx echo.Context, settlementId int32, params CalculateSettlementDiscrepancyParams) error
	// Returns all pets
	// (GET /usages/)
	FindUsages(ctx echo.Context) error
	// Request usage discrepancy
	// (PUT /usages/{usageId})
	CalculateUsageDiscrepancy(ctx echo.Context, usageId int32, params CalculateUsageDiscrepancyParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CalculateSettlementDiscrepancy converts echo context to params.
func (w *ServerInterfaceWrapper) CalculateSettlementDiscrepancy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "settlementId" -------------
	var settlementId int32

	err = runtime.BindStyledParameter("simple", false, "settlementId", ctx.Param("settlementId"), &settlementId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter settlementId: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params CalculateSettlementDiscrepancyParams
	// ------------- Required query parameter "partnerSettlementId" -------------

	err = runtime.BindQueryParameter("form", true, true, "partnerSettlementId", ctx.QueryParams(), &params.PartnerSettlementId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter partnerSettlementId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CalculateSettlementDiscrepancy(ctx, settlementId, params)
	return err
}

// FindUsages converts echo context to params.
func (w *ServerInterfaceWrapper) FindUsages(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindUsages(ctx)
	return err
}

// CalculateUsageDiscrepancy converts echo context to params.
func (w *ServerInterfaceWrapper) CalculateUsageDiscrepancy(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "usageId" -------------
	var usageId int32

	err = runtime.BindStyledParameter("simple", false, "usageId", ctx.Param("usageId"), &usageId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter usageId: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params CalculateUsageDiscrepancyParams
	// ------------- Required query parameter "partnerUsageId" -------------

	err = runtime.BindQueryParameter("form", true, true, "partnerUsageId", ctx.QueryParams(), &params.PartnerUsageId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter partnerUsageId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CalculateUsageDiscrepancy(ctx, usageId, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.PUT(baseURL+"/settlements/:settlementId", wrapper.CalculateSettlementDiscrepancy)
	router.GET(baseURL+"/usages/", wrapper.FindUsages)
	router.PUT(baseURL+"/usages/:usageId", wrapper.CalculateUsageDiscrepancy)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xYy3LbNhT9FQ7aJcdSnJkutGpip7UWqj2Wky5SjwYiryTEJADjIVvx6N87ePBhEpQp",
	"Wc5kkZUtAvd1cO/BIZ9QwnLOKFAl0egJyWQFObb/nmOFpyDWJAHzkwvGQSgCdpHi3D5VGw5ohKQShC7R",
	"NkZrnGm7smAixwqN0CJjWKG42Ep1PgeBttvyCZt/g0QZ409CMNEOlrDUukxBJoJwRRhFI7c5smtxFY1Q",
	"9f60ikaogqUJF6McpMTLTkfFctysaRsjAfeaCEjR6CvyAYvtt9sY/Q0UBM7GdMEMaO0CCJ0zTdNZSmQi",
	"gGOabPpAFJeG7IHOdJF9fzOOhaIg9jFlWh2Ya2m5Z7Kl3QHZyqo/W52oKVGBhcZxFh78/hDkXXjGwXMN",
	"AtFZZQfgpqsml2ftTrq+/LcfMnOc3F2wvHfDKBAUm4HAWT+TjCV9t3IBOdH5oawwBaUyyIGqNh5zlm7a",
	"E11ZRHZD3LBaCJanWIX7xh+qWftdwAKN0G+DiiQHniEHVQxPkrLezIdZK9aRVaNpixRr8W4DwK0ApyB2",
	"wuO3xC3GpQoeVdv0zC/EbdxyyS8faCjeZHoVMbsUMHMPmiY3Gw4RW0RqBVHKEm1PP2C9BiGtRdPBF7fQ",
	"9BFZDy9xfOHV76vVFpfQtAFvOCmhtS14+6yRz6tZv4aEibTd2SlkCs8SnCU6s4M54yASPwQ9WPWB1o17",
	"Dqrnpr0N34CHmwWEs4t34LQDcs5EgExWLIcrEJJDosgaQmeiMMncla4gl/2nvH3gFdNhIfDG/F46HTEj",
	"1CHu4T9+qAb2obhxWWyIWfxR/ALrZbB232gl+7dFrxYCvPhqzZTsNCsEaC8c6hI/APF0Mm0HmFz2o4TJ",
	"zaH3/ZoFXzm8HtpVj9liQ58dFrt5DTDHSAaH2CH7IuuXx1Y7pNAEfS5Ebh8985Glm+Ims4b2Gou8tzj8",
	"ttG7C6xHW12gB+qC5rXejqtiHBC/BEwvAfMK1VIdaPDCvMEpWe4iqX1EQYz2ePlbE0kUE90JbACLCaNq",
	"FZIfZWXVDRAu8uJm8k/Q/ZerjgUnSWwlMzyX/aqpG+2t9PZA7chv2fti3KW9XnOdN7+/BGjsME5stMbx",
	"6PElx1uLncGiIDOcWNAgxyRDI5TpOyy/n6Q6/U6ovCN/sjVjJzxznfSMiFZERkRaErn+NL2JPlyNI6Pa",
	"yIIkFudowYRdrqUVTUs9rojKTGqBVeML1TgMDU+GJ+8sLhwo5gSN0Hv7yPSdWll0BrJUP3LwVP0Yp1vb",
	"FjrA2tdwr0GqqNodNb65cBC2mHFqSN6/DUBQ6LlXCZyDAiHR6GszmGG1eqTxuf3Ug0a2BhT7j5+onjqq",
	"s6gSGmL/KfXZhHV9nNzGzRz8lIbTuNcgNlUefu/0eOncOnOQ6qPXI/b+cIyEOc985wy+STemlfM99be9",
	"d/Hj2FmdxigntPrRHIpmb1+Yk/pPD4enf0RtxGQLBXvXSc6odKRzOhzuVdwB7xSW7QKp16dJ2F3RA1Gr",
	"yDYfpmlZDwchE/+SJd14L7DO1NESdx/dAylqCo8mMKQR+D0xkjrPsdi8PJNm88BeMnLgCN7m+XxQ/yI0",
	"tYzozuqNjsZp7WMdws91BkoLKiOcZREHJZ+h/mT/9qFVu7Efozbvr15k6vx7hLvo1Gf7VkzakUOYSz8f",
	"JZe3plHf18dm0DpSP5ZDO8TiHpPrdZ4d3kKbRcXHmZ+QPNuTZ30ZxV3MkxZG8A0wJ4P1uwEyXaXwMjRr",
	"H5IEpIwUCw50fcxQe04q406R1dI9aHu7/T8AAP//EKY+6kIeAAA=",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}