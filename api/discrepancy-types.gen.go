// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

// Error defines model for Error.
type Error struct {

	// Error code
	Code int32 `json:"code"`

	// Error message
	Message string `json:"message"`
}

// GeneralInfoData defines model for GeneralInfoData.
type GeneralInfoData struct {
	InboundDiscrepancy   *float32 `json:"inbound_discrepancy,omitempty"`
	InboundOwnUsage      *float32 `json:"inbound_own_usage,omitempty"`
	InboundPartnerUsage  *float32 `json:"inbound_partner_usage,omitempty"`
	OutboundDiscrepancy  *float32 `json:"outbound_discrepancy,omitempty"`
	OutboundOwnUsage     *float32 `json:"outbound_own_usage,omitempty"`
	OutboundPartnerUsage *float32 `json:"outbound_partner_usage,omitempty"`
	Service              *string  `json:"service,omitempty"`
	Unit                 *string  `json:"unit,omitempty"`
}

// Usage defines model for Usage.
type Usage struct {

	// Body of the Usage type object
	Body *struct {
		Inbound  *[]UsageData `json:"inbound,omitempty"`
		Outbound *[]UsageData `json:"outbound,omitempty"`
	} `json:"body,omitempty"`

	// Usage header
	Header struct {

		// Context
		Context *string `json:"context,omitempty"`

		// MSP owner
		MspOwner *string `json:"mspOwner,omitempty"`

		// Type of the document
		Type *string `json:"type,omitempty"`

		// Version of the document type
		Version *string `json:"version,omitempty"`
	} `json:"header"`
}

// UsageData defines model for UsageData.
type UsageData struct {
	HomeTadig    *string `json:"homeTadig,omitempty"`
	Service      *string `json:"service,omitempty"`
	Usage        *int32  `json:"usage,omitempty"`
	VisitorTadig *string `json:"visitorTadig,omitempty"`
	YearMonth    *string `json:"yearMonth,omitempty"`
}

// UsageDiscrepancyData defines model for UsageDiscrepancyData.
type UsageDiscrepancyData struct {
	HTMN              *string  `json:"HTMN,omitempty"`
	VPMN              *string  `json:"VPMN,omitempty"`
	DeltaUsageAbs     *float32 `json:"delta_usage_abs,omitempty"`
	DeltaUsagePercent *float32 `json:"delta_usage_percent,omitempty"`
	OwnUsage          *float32 `json:"own_usage,omitempty"`
	PartnerUsage      *float32 `json:"partner_usage,omitempty"`
	Service           *string  `json:"service,omitempty"`
	YearMonth         *string  `json:"yearMonth,omitempty"`
}

// UsageDiscrepancyReport defines model for UsageDiscrepancyReport.
type UsageDiscrepancyReport struct {
	GeneralInformation *[]GeneralInfoData      `json:"general_information,omitempty"`
	Inbound            *[]UsageDiscrepancyData `json:"inbound,omitempty"`
	Outbound           *[]UsageDiscrepancyData `json:"outbound,omitempty"`
}

// CalculateUsageDiscrepancyJSONBody defines parameters for CalculateUsageDiscrepancy.
type CalculateUsageDiscrepancyJSONBody []Usage

// CalculateUsageDiscrepancyParams defines parameters for CalculateUsageDiscrepancy.
type CalculateUsageDiscrepancyParams struct {

	// partner usage report ID
	PartnerUsageId int32 `json:"partnerUsageId"`
}

// CalculateUsageDiscrepancyJSONRequestBody defines body for CalculateUsageDiscrepancy for application/json ContentType.
type CalculateUsageDiscrepancyJSONRequestBody CalculateUsageDiscrepancyJSONBody
