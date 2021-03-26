package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type DiscrepancyServer struct {
	Pets   map[int64]Usage
	NextId int64
	Lock   sync.Mutex
}

func NewDiscrepancyServer() *DiscrepancyServer {
	return &DiscrepancyServer{
		Pets:   make(map[int64]Usage),
		NextId: 1000,
	}
}

func (p *DiscrepancyServer) CalculateUsageDiscrepancy(ctx echo.Context, usageId int32, params CalculateUsageDiscrepancyParams) error {
	fmt.Println("Start: CalculateUsageDiscrepancy")

	// retrieve two usage reports from the request body
	b, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.NoContent(http.StatusNotAcceptable)
	}

	var req []Usage // the non-struct body
	if b != nil {
		err := json.Unmarshal(b, &req)
		if err != nil {
			return ctx.NoContent(http.StatusNotAcceptable)
		}
	} else {
		return ctx.NoContent(http.StatusNotAcceptable)
	}

	ownUsage := req[0] // assumption: first usage is a home one
	partnerUsage := req[1]

	fmt.Println(ownUsage.Header.Context)
	fmt.Println(partnerUsage.Header.Context)

	// create output usage discrepancy report
	var report UsageDiscrepancyReport
	report = UsageDiscrepancyReport{}

	// general information
	generalInformationMap := make(map[string]*GeneralInfoData, 0)

	// general information - inbound own usage
	for _, usageDataRecord := range ownUsage.Body.Inbound {
		value, ok := generalInformationMap[*usageDataRecord.Service]
		if !ok {
			generalInfoData := GeneralInfoData{}
			generalInfoData.Service = *usageDataRecord.Service
			generalInfoData.Unit = *usageDataRecord.Unit
			generalInfoData.InboundOwnUsage = *usageDataRecord.Usage
			generalInformationMap[*usageDataRecord.Service] = &generalInfoData

		} else {
			summary := value.InboundOwnUsage + *usageDataRecord.Usage
			value.InboundOwnUsage = summary
		}
	}

	// general information - outbound partner usage
	for _, usageDataRecord := range partnerUsage.Body.Outbound {
		generalInfoRecord, ok := generalInformationMap[*usageDataRecord.Service]
		if !ok {
			generalInfoData := GeneralInfoData{}
			generalInfoData.Service = *usageDataRecord.Service
			generalInfoData.Unit = *usageDataRecord.Unit
			generalInfoData.InboundPartnerUsage = *usageDataRecord.Usage
			generalInformationMap[*usageDataRecord.Service] = &generalInfoData

		} else {
			summary := generalInfoRecord.InboundPartnerUsage + *(usageDataRecord.Usage)
			generalInfoRecord.InboundPartnerUsage = summary

		}
	}

	// general information - outbound own usage
	for _, usageDataRecord := range ownUsage.Body.Outbound {
		value, ok := generalInformationMap[*usageDataRecord.Service]
		if !ok {
			generalInfoData := GeneralInfoData{}
			generalInfoData.Service = *usageDataRecord.Service
			generalInfoData.Unit = *usageDataRecord.Unit
			generalInfoData.OutboundOwnUsage = *usageDataRecord.Usage
			generalInformationMap[*usageDataRecord.Service] = &generalInfoData

		} else {
			summary := value.OutboundOwnUsage + *(usageDataRecord.Usage)
			value.OutboundOwnUsage = summary
		}
	}

	// general information - outbound partner usage
	for _, usageDataRecord := range partnerUsage.Body.Inbound {
		value, ok := generalInformationMap[*usageDataRecord.Service]
		if !ok {
			generalInfoData := GeneralInfoData{}
			generalInfoData.Service = *usageDataRecord.Service
			generalInfoData.Unit = *usageDataRecord.Unit
			generalInfoData.OutboundPartnerUsage = *usageDataRecord.Usage
			generalInformationMap[*usageDataRecord.Service] = &generalInfoData

		} else {

			summary := value.OutboundPartnerUsage + *(usageDataRecord.Usage)
			value.OutboundPartnerUsage = summary

		}
	}

	// create general information array for sub-services
	generalInformationSubServiceArray := make([]GeneralInfoData, 0, len(generalInformationMap))

	for _, value := range generalInformationMap {
		generalInformationSubServiceArray = append(generalInformationSubServiceArray, *value)
	}

	// VOICE general information
	voiceGeneralInformation := GeneralInfoData{}
	moc := "MOC"
	voiceGeneralInformation.Service = moc
	min := "min"
	voiceGeneralInformation.Unit = min

	voiceGeneralInformation.InboundOwnUsage = 0
	voiceGeneralInformation.InboundPartnerUsage = 0
	voiceGeneralInformation.OutboundOwnUsage = 0
	voiceGeneralInformation.OutboundPartnerUsage = 0

	// SMS general information
	smsGeneralInformation := GeneralInfoData{}
	sms := "SMS"
	smsGeneralInformation.Service = sms
	smsUnit := "#"

	smsGeneralInformation.Unit = smsUnit
	smsGeneralInformation.InboundOwnUsage = 0
	smsGeneralInformation.InboundPartnerUsage = 0
	smsGeneralInformation.OutboundOwnUsage = 0
	smsGeneralInformation.OutboundPartnerUsage = 0

	// DATA general information
	dataGeneralInformation := GeneralInfoData{}
	dataServices := "Data"
	dataGeneralInformation.Service = dataServices
	dataUnit := "MB"
	dataGeneralInformation.Unit = dataUnit
	dataGeneralInformation.InboundOwnUsage = 0
	dataGeneralInformation.InboundPartnerUsage = 0
	dataGeneralInformation.OutboundOwnUsage = 0
	dataGeneralInformation.OutboundPartnerUsage = 0

	for _, element := range generalInformationSubServiceArray {
		if element.Unit == "min" {
			voiceGeneralInformation.InboundOwnUsage += element.InboundOwnUsage
			voiceGeneralInformation.InboundPartnerUsage += element.InboundPartnerUsage
			voiceGeneralInformation.OutboundOwnUsage += element.OutboundOwnUsage
			voiceGeneralInformation.OutboundPartnerUsage += element.OutboundPartnerUsage

		} else if element.Unit == "SMS" {
			smsGeneralInformation.InboundOwnUsage += element.InboundOwnUsage
			smsGeneralInformation.InboundPartnerUsage += element.InboundPartnerUsage
			smsGeneralInformation.OutboundOwnUsage += element.OutboundOwnUsage
			smsGeneralInformation.OutboundPartnerUsage += element.OutboundPartnerUsage

		} else if element.Unit == "MB" {
			dataGeneralInformation.InboundOwnUsage += element.InboundOwnUsage
			dataGeneralInformation.InboundPartnerUsage += element.InboundPartnerUsage
			dataGeneralInformation.OutboundOwnUsage += element.OutboundOwnUsage
			dataGeneralInformation.OutboundPartnerUsage += element.OutboundPartnerUsage
		}
	}

	generalInformationSubServiceArray = nil

	generalInformationBearerServiceArray := make([]GeneralInfoData, 3, 3)
	generalInformationBearerServiceArray[0] = calculateInOutDiscrepancies(&voiceGeneralInformation)
	generalInformationBearerServiceArray[1] = calculateInOutDiscrepancies(&smsGeneralInformation)
	generalInformationBearerServiceArray[2] = calculateInOutDiscrepancies(&dataGeneralInformation)

	report.GeneralInformation = &generalInformationBearerServiceArray

	// inbound details
	homeInboundMap := p.convertUsageDataArrayToMap(ownUsage.Body.Inbound)
	partnerOutboundMap := p.convertUsageDataArrayToMap(partnerUsage.Body.Outbound)

	inbound := make([]UsageDiscrepancyData, 0)

	for key, inUsage := range homeInboundMap {
		fmt.Println("Key:", key)
		outUsage, ok := partnerOutboundMap[key]
		if ok {
			inboundUsageDiscrepancyData := createInOutDetailsRecord(inUsage, outUsage)
			inbound = append(inbound, inboundUsageDiscrepancyData)
		}
	}

	report.Inbound = &inbound

	// outbound details
	homeOutboundMap := p.convertUsageDataArrayToMap(ownUsage.Body.Outbound)
	partnerInboundMap := p.convertUsageDataArrayToMap(partnerUsage.Body.Inbound)

	outbound := make([]UsageDiscrepancyData, 0)

	for key, outUsage := range homeOutboundMap {
		fmt.Println("Key:", key)
		inUsage, ok := partnerInboundMap[key]
		if ok {
			outboundUsageDiscrepancyData := createInOutDetailsRecord(outUsage, inUsage)
			outbound = append(outbound, outboundUsageDiscrepancyData)
		}
	}

	report.Outbound = &outbound

	return ctx.JSON(http.StatusOK, report)
}

func calculateInOutDiscrepancies(value *GeneralInfoData) GeneralInfoData {
	delta64 := float64(value.InboundOwnUsage) - float64(value.InboundPartnerUsage)
	absDelta64 := math.Abs(delta64)
	absDelta32 := float32(absDelta64)
	value.InboundDiscrepancy = absDelta32

	delta64 = float64(value.OutboundOwnUsage) - float64(value.OutboundPartnerUsage)
	absDelta64 = math.Abs(delta64)
	absDelta32 = float32(absDelta64)
	value.OutboundDiscrepancy = absDelta32

	return *value
}

func createInOutDetailsRecord(ownUsage UsageData, partnerUsage UsageData) UsageDiscrepancyData {

	var record UsageDiscrepancyData
	record = UsageDiscrepancyData{}

	record.HTMN = ownUsage.HomeTadig
	record.VPMN = ownUsage.VisitorTadig
	record.YearMonth = ownUsage.YearMonth
	record.Service = ownUsage.Service
	record.OwnUsage = ownUsage.Usage
	record.PartnerUsage = partnerUsage.Usage
	// absolute delta
	delta64 := float64(*ownUsage.Usage) - float64(*partnerUsage.Usage)
	absDelta64 := math.Abs(delta64)
	absDelta32 := float32(absDelta64)
	record.DeltaUsageAbs = &absDelta32
	// relative delta
	// [ (A-B) / A] x 100
	A := *ownUsage.Usage
	B := *partnerUsage.Usage
	C := ((A - B) / A) * 100

	record.DeltaUsagePercent = &C

	return record
}

func (p *DiscrepancyServer) convertUsageDataArrayToMap(arr []UsageData) map[string]UsageData {
	fmt.Println("Start: convertUsageDataArrayToMap")

	// create output map
	m := make(map[string]UsageData)

	for index, element := range arr {
		fmt.Println("At index", index, "value is", toString(element))

		compositeUsageId := makeUsageIdentifier(element)
		fmt.Println("compositeUsageId", compositeUsageId)

		var data = []byte(compositeUsageId)
		var dataBase64 = base64.StdEncoding.EncodeToString(data)
		sha256 := sha256.Sum256([]byte(dataBase64))
		hashKey := hex.EncodeToString(sha256[:])

		// sets the hash based key to the given element
		m[hashKey] = element
		fmt.Println("Hash key: ", hashKey)
	}

	return m
}

func toString(usageData UsageData) string {
	return (*usageData.HomeTadig + ", " + *usageData.VisitorTadig + ", " + *usageData.Service + ", " + *usageData.YearMonth)
}

func makeUsageIdentifier(usageData UsageData) string {
	return (*usageData.HomeTadig + *usageData.VisitorTadig + *usageData.Service + *usageData.YearMonth)
}

func (p *DiscrepancyServer) CreateUsageDiscrepancyReport() UsageDiscrepancyReport {

	// create usage discrepancy report
	usageDiscrepancyData := UsageDiscrepancyData{}
	hpmn := "DTAG"
	usageDiscrepancyData.HTMN = &hpmn
	vpmn := "DTAG"
	usageDiscrepancyData.VPMN = &vpmn
	service := "MOC Back Home"
	usageDiscrepancyData.Service = &service
	year_month := "202003"
	usageDiscrepancyData.YearMonth = &year_month

	// own usage
	var ownUsage float32
	ownUsage = 100.0
	usageDiscrepancyData.OwnUsage = &ownUsage
	// partner usage
	var partnerUsage float32
	partnerUsage = 110.0
	usageDiscrepancyData.PartnerUsage = &partnerUsage
	// delta usage abs
	var deltaUsageAbs float32
	deltaUsageAbs = 45.7
	usageDiscrepancyData.DeltaUsageAbs = &deltaUsageAbs
	// delta usage percent
	var deltaUsagePercent float32
	deltaUsagePercent = 10.5
	usageDiscrepancyData.DeltaUsagePercent = &deltaUsagePercent

	var report UsageDiscrepancyReport
	report = UsageDiscrepancyReport{}

	var generalInfo GeneralInfoData
	generalInfo = GeneralInfoData{}

	bearerService := "MOC"
	generalInfo.Service = bearerService

	units := "min"
	generalInfo.Unit = units

	var inbound_own_usage float32
	inbound_own_usage = 300.00
	var inbound_partner_usage float32
	inbound_partner_usage = 330.00
	var inbound_discrepancy float32
	inbound_discrepancy = 30.00
	var outbound_own_usage float32
	outbound_own_usage = 600.00
	var outbound_partner_usage float32
	outbound_partner_usage = 630.00
	var outbound_discrepancy float32
	outbound_discrepancy = 30.00

	generalInfo.InboundOwnUsage = inbound_own_usage
	generalInfo.InboundPartnerUsage = inbound_partner_usage
	generalInfo.InboundDiscrepancy = inbound_discrepancy
	generalInfo.OutboundOwnUsage = outbound_own_usage
	generalInfo.OutboundPartnerUsage = outbound_partner_usage
	generalInfo.OutboundDiscrepancy = outbound_discrepancy

	generalInfoArray := make([]GeneralInfoData, 0)
	generalInfoArray = append(generalInfoArray, generalInfo)
	report.GeneralInformation = &generalInfoArray

	inbound := make([]UsageDiscrepancyData, 0)
	inbound = append(inbound, usageDiscrepancyData)
	report.Inbound = &inbound

	outbound := make([]UsageDiscrepancyData, 0)
	outbound = append(outbound, usageDiscrepancyData)
	report.Outbound = &outbound

	return report

}

func (p *DiscrepancyServer) FindUsages(ctx echo.Context) error {
	fmt.Println("Start: FindUsages")

	var usage Usage
	dtag := "DTAG"
	version := "1.0"
	usage.Header.MspOwner = &dtag
	usage.Header.Version = version

	return ctx.JSON(http.StatusOK, usage)
}

func (p *DiscrepancyServer) CalculateSettlementDiscrepancy(ctx echo.Context, settlementId int32, params CalculateSettlementDiscrepancyParams) error {
	fmt.Println("Start: CalculateSettlementDiscrepancy")

	var bearerServiceData SettlementDiscrepancyData
	bearerServiceData = SettlementDiscrepancyData{}

	voice := "Voice"
	bearerServiceData.Service = &voice

	unit := "min"
	bearerServiceData.Unit = &unit

	// own calculation
	var ownCalculation float32
	ownCalculation = 4390.83
	bearerServiceData.OwnCalculation = &ownCalculation

	// partner calculation
	var partnerCalculation float32
	partnerCalculation = 4390.83
	bearerServiceData.PartnerCalculation = &partnerCalculation

	// delta calculation percent
	var deltaCalculationPercent float32
	deltaCalculationPercent = 4390.83
	bearerServiceData.DeltaCalculationPercent = &deltaCalculationPercent

	report := SettlementDiscrepancyReport{}

	generalInfoArray := make([]SettlementDiscrepancyData, 0)
	generalInfoArray = append(generalInfoArray, bearerServiceData)

	report.HomePerspective = &(struct {
		Details            *[]SettlementDiscrepancyData `json:"details,omitempty"`
		GeneralInformation *[]SettlementDiscrepancyData `json:"general_information,omitempty"`
	}{&generalInfoArray, &generalInfoArray})

	fmt.Println(report.HomePerspective)
	fmt.Println(*(report.HomePerspective))

	// (*(settlementDiscrepancyReport.HomePerspective.GeneralInformation))[0] = bearerServiceData

	return ctx.JSON(http.StatusOK, report)
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendDiscrepancyError(ctx echo.Context, code int, message string) error {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, petErr)
	return err
}
