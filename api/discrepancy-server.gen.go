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
	CalculateUsageDiscrepancy(ctx echo.Context, usageId string, params CalculateUsageDiscrepancyParams) error
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
	var usageId string

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

	"H4sIAAAAAAAC/+xZS2/bOBf9KwK/bynEbgrMwqvpa6ZeeBLEaWfRCQxaurbZSCTDhxM38H8f8KGHJcqW",
	"HWfQRVetRd4Hz7338Eh5RgnLOaNAlUSjZySTFeTY/vcjVngKYk0SMD+5YByEImAXKc7tU7XhgEZIKkHo",
	"Em1jtMaZtisLJnKs0AgtMoYVioutVOdzEGi7LZ+w+XdIlDH+JAQT7WAJS63LFGQiCFeEUTRymyO7FlfR",
	"CFVvL6tohCpYmnAxykFKvOx0VCzHzTNtYyTgQRMBKRp9Qz5gsf1uG6M/gYLA2ZgumAGtfQBC50zTdJYS",
	"mQjgmCabPhDFpSF7pDNdZN/fjGOhKIhjTJlWJ+ZaWh6ZbGl3Qray6s9WJ2pKVGChUc7Cg98fgrwLzzhY",
	"1yAQnafsANx01eTqQ7uTbq7+7ofMHCf3n1neu2EUCIrNQOCsn0nGkr5buYCc6PxUVpiCUhnkQFUbjzlL",
	"N+2JriwiuyFuWC0Ey1Oswn3ji2rW/i9ggUbof4OKJAeeIQdVDE+Sst7Mp1kr1pFVo2mLFGvx7gLArQCn",
	"IPbC47fELcalCp5U2/SDX4jbuOWSXz3SULzJ9Dpidilg5h40TW43HCK2iNQKopQl2lY/YL0GIa1F08FX",
	"t9D0EVkPhzi+8Or31c4Wl9C0AW84KaG1LXi308gfq1m/gYSJtN3ZKWQKzxKcJTqzgznjIBI/BIdnzplb",
	"kpnhudwxSpmeZxDk4kdaD9mTwIN83x2kYMCjA3VfEN3BXuGGqNN6k82bsLchDZ8/3lPvPa3DmQiQ4orl",
	"cA1CckgUWUOotxQmmZMmCnLZn63ajVsxNhYCb8zvpdNDM0JdgXyBzx+qUalQ3Lg8bIghfSl+gXUYrP03",
	"c3mLtcW7FgK8iGxNoOw0K4R0LxzqryoBiKeTaTvA5Kof6UxuT9UtaxZ8dfK6bt95zBYb+sNpsZvXGXP8",
	"ZXCIHbIHb6+ybLUihSboS8HFfXTZe5ZuihvZGtrrOPLe4vBbU+8usB7t6QI9UBdmL/V2XjXmgPglxHoJ",
	"sReor6qgwQvzFqdkuY+kOiWEDK/0f4tdE0kUE90ZbACLCaNqFVIr5dGqKyB8ys+3k7+C7r9edyzsE5G9",
	"lOdRkvXI7wZn/lxwLMZd4usl93nzQ1KAx04jxUZrnI8fDzneWuwMFgWb4cSCBjkmGRqhTN9j+eMi1ekP",
	"QuU9+Z2tGbvgmeukHSZaERkRaVnk5tP0Nnp3PY6MbCMLklicowUTdrmWVjQt5bsiKjOpBVaNL1QjMTS8",
	"GF68sbhwoJgTNEJv7SPTd2pl0RnIUv7IwXP1Y5xubVvoAG3fwIMGqaJqd9T4eMRB2MOMU8Py/nUAgkrP",
	"vUvgHBQIiUbfmsEMrdUjjT/ab1ZoZM+AYv8VF9VTR3UaVUJD7L8J70xY11fWbdzMwU9pOI0HDWJT5eH3",
	"Ts+Xzp0zB6nee0FiLxDHSJjzzHfO4Lt0Y1o5P1KA24sXP42d1WWMckKrH82haPb2Z1Opf/RwePlb1EZM",
	"tlCwl53kjEpHOpfD4VGHO+GlwrJdIPX6NAm7K3okahXZ5sM0Lc/DQcjEv2VJN94LrDN1tsTdXw8CKWoK",
	"TyYwpBH4PTGSOs+x2ByeSbN5YC8ZOXAEb/PcHdQ/CE0tI7pavVJpnNg+VxF+rhooLaiMcJZFHJTcQf3Z",
	"/tuHVu3GfozavL96kanz7xHuolOf7V7qasneLt7siBhmzi8nRH5tivQ9e252rOPy3/JjhxA8Yiq9hrOD",
	"WeiuqPjy8hMSY3uqrC+jpotZ0cKIuQHmZLB+M0CmqxRehuboXZKAlJFiwWGtjxBqT0Vl3CmgWpoGbe+2",
	"/wYAAP//wMs+ruceAAA=",
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
