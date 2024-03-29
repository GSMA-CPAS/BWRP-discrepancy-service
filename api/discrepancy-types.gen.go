// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

// DataService defines model for DataService.
type DataService struct {
	Name  *string       `json:"name,omitempty"`
	Value *TelcoService `json:"value,omitempty"`
}

// Error defines model for Error.
type Error struct {

	// Error code
	Code int32 `json:"code"`

	// Error message
	Message string `json:"message"`
}

// GeneralInfoData defines model for GeneralInfoData.
type GeneralInfoData struct {
	InboundDiscrepancy   float64 `json:"inbound_discrepancy"`
	InboundOwnUsage      float64 `json:"inbound_own_usage"`
	InboundPartnerUsage  float64 `json:"inbound_partner_usage"`
	OutboundDiscrepancy  float64 `json:"outbound_discrepancy"`
	OutboundOwnUsage     float64 `json:"outbound_own_usage"`
	OutboundPartnerUsage float64 `json:"outbound_partner_usage"`
	Service              string  `json:"service"`
	Unit                 string  `json:"unit"`
}

// MOC defines model for MOC.
type MOC struct {
	EEA                 *TelcoService `json:"EEA,omitempty"`
	EU                  *TelcoService `json:"EU,omitempty"`
	ROW                 *TelcoService `json:"ROW,omitempty"`
	BackHome            *TelcoService `json:"backHome,omitempty"`
	International       *TelcoService `json:"international,omitempty"`
	Local               *TelcoService `json:"local,omitempty"`
	Premium             *TelcoService `json:"premium,omitempty"`
	Satellite           *TelcoService `json:"satellite,omitempty"`
	SpecialDestinations *TelcoService `json:"specialDestinations,omitempty"`
	VideoTelephony      *TelcoService `json:"videoTelephony,omitempty"`
}

// Settlement defines model for Settlement.
type Settlement struct {

	// Settlement body
	Body struct {
		Fromdate *string            `json:"fromdate,omitempty"`
		Inbound  SettlementServices `json:"inbound"`
		Outbound SettlementServices `json:"outbound"`
		Todate   *string            `json:"todate,omitempty"`
	} `json:"body"`

	// Settlement header
	Header struct {

		// Context
		Context string `json:"context"`

		// MSP owner
		MspOwner string `json:"mspOwner"`

		// Type of the document
		Type string `json:"type"`

		// Version of the document type
		Version string `json:"version"`
	} `json:"header"`
}

// SettlementDiscrepancyRecord defines model for SettlementDiscrepancyRecord.
type SettlementDiscrepancyRecord struct {
	DeltaCalculationPercent float64 `json:"delta_calculation_percent"`
	DeltaUsageAbs           float64 `json:"delta_usage_abs"`
	DeltaUsagePercent       float64 `json:"delta_usage_percent"`
	OwnCalculation          float64 `json:"own_calculation"`
	OwnUsage                float64 `json:"own_usage"`
	PartnerCalculation      float64 `json:"partner_calculation"`
	PartnerUsage            float64 `json:"partner_usage"`
	Service                 string  `json:"service"`
	Unit                    string  `json:"unit"`
}

// SettlementDiscrepancyReport defines model for SettlementDiscrepancyReport.
type SettlementDiscrepancyReport struct {
	HomePerspective *struct {
		Details            []SettlementDiscrepancyRecord `json:"details"`
		GeneralInformation []SettlementDiscrepancyRecord `json:"general_information"`
	} `json:"homePerspective,omitempty"`
	PartnerPerspective *struct {
		Details            []SettlementDiscrepancyRecord `json:"details"`
		GeneralInformation []SettlementDiscrepancyRecord `json:"general_information"`
	} `json:"partnerPerspective,omitempty"`
	SettlementReport *struct {
		HomeCharges            float64 `json:"homeCharges"`
		HomeDeltaCommitment    float64 `json:"homeDeltaCommitment"`
		HomeRevenue            float64 `json:"homeRevenue"`
		PartnerCharges         float64 `json:"partnerCharges"`
		PartnerDeltaCommitment float64 `json:"partnerDeltaCommitment"`
		PartnerRevenue         float64 `json:"partnerRevenue"`
	} `json:"settlementReport,omitempty"`
}

// SettlementServices defines model for SettlementServices.
type SettlementServices struct {
	Currency string `json:"currency"`
	Services struct {
		SMS struct {
			MO *TelcoService `json:"MO,omitempty"`
			MT *TelcoService `json:"MT,omitempty"`
		} `json:"SMS"`
		Data  []DataService `json:"data"`
		Voice struct {
			MOC *MOC          `json:"MOC,omitempty"`
			MTC *TelcoService `json:"MTC,omitempty"`
		} `json:"voice"`
	} `json:"services"`
}

// TelcoService defines model for TelcoService.
type TelcoService struct {
	DealValue         float64 `json:"dealValue"`
	ShortOfCommitment float64 `json:"shortOfCommitment"`
	Usage             float64 `json:"usage"`
}

// Usage defines model for Usage.
type Usage struct {

	// Body of the Usage type object
	Body struct {
		Inbound  []UsageData `json:"inbound"`
		Outbound []UsageData `json:"outbound"`
	} `json:"body"`

	// Usage header
	Header struct {

		// Context
		Context string `json:"context"`

		// MSP owner
		MspOwner *string `json:"mspOwner,omitempty"`

		// Type of the document
		Type string `json:"type"`

		// Version of the document type
		Version string `json:"version"`
	} `json:"header"`
}

// UsageData defines model for UsageData.
type UsageData struct {
	HomeTadig    *string  `json:"homeTadig,omitempty"`
	Service      *string  `json:"service,omitempty"`
	Units        *string  `json:"units,omitempty"`
	Usage        *float64 `json:"usage,omitempty"`
	VisitorTadig *string  `json:"visitorTadig,omitempty"`
	YearMonth    *string  `json:"yearMonth,omitempty"`
}

// UsageDiscrepancyData defines model for UsageDiscrepancyData.
type UsageDiscrepancyData struct {
	HTMN              *string  `json:"HTMN,omitempty"`
	VPMN              *string  `json:"VPMN,omitempty"`
	DeltaUsageAbs     *float64 `json:"delta_usage_abs,omitempty"`
	DeltaUsagePercent *float64 `json:"delta_usage_percent,omitempty"`
	OwnUsage          *float64 `json:"own_usage,omitempty"`
	PartnerUsage      *float64 `json:"partner_usage,omitempty"`
	Service           *string  `json:"service,omitempty"`
	YearMonth         *string  `json:"yearMonth,omitempty"`
}

// UsageDiscrepancyReport defines model for UsageDiscrepancyReport.
type UsageDiscrepancyReport struct {
	GeneralInformation *[]GeneralInfoData      `json:"general_information,omitempty"`
	Inbound            *[]UsageDiscrepancyData `json:"inbound,omitempty"`
	Outbound           *[]UsageDiscrepancyData `json:"outbound,omitempty"`
}

// CalculateSettlementDiscrepancyJSONBody defines parameters for CalculateSettlementDiscrepancy.
type CalculateSettlementDiscrepancyJSONBody []Settlement

// CalculateSettlementDiscrepancyParams defines parameters for CalculateSettlementDiscrepancy.
type CalculateSettlementDiscrepancyParams struct {

	// partner settlement ID
	PartnerSettlementId string `json:"partnerSettlementId"`
}

// CalculateUsageDiscrepancyJSONBody defines parameters for CalculateUsageDiscrepancy.
type CalculateUsageDiscrepancyJSONBody []Usage

// CalculateUsageDiscrepancyParams defines parameters for CalculateUsageDiscrepancy.
type CalculateUsageDiscrepancyParams struct {

	// partner usage report ID
	PartnerUsageId string `json:"partnerUsageId"`
}

// CalculateSettlementDiscrepancyJSONRequestBody defines body for CalculateSettlementDiscrepancy for application/json ContentType.
type CalculateSettlementDiscrepancyJSONRequestBody CalculateSettlementDiscrepancyJSONBody

// CalculateUsageDiscrepancyJSONRequestBody defines body for CalculateUsageDiscrepancy for application/json ContentType.
type CalculateUsageDiscrepancyJSONRequestBody CalculateUsageDiscrepancyJSONBody
