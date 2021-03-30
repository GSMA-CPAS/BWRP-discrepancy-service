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
	report := UsageDiscrepancyReport{}

	// general information
	aggregatedSubServicesMap := make(map[string]*GeneralInfoData, 0)

	// general information - inbound own usage
	for _, usageDataRecord := range ownUsage.Body.Inbound {
		value, ok := aggregatedSubServicesMap[*usageDataRecord.Service]
		if !ok {
			generalInfoData := GeneralInfoData{}
			generalInfoData.Service = *usageDataRecord.Service
			generalInfoData.Unit = *usageDataRecord.Unit
			generalInfoData.InboundOwnUsage = *usageDataRecord.Usage
			aggregatedSubServicesMap[*usageDataRecord.Service] = &generalInfoData

		} else {
			summary := value.InboundOwnUsage + *usageDataRecord.Usage
			value.InboundOwnUsage = summary
		}
	}

	// general information - inbound partner usage
	for _, usageDataRecord := range partnerUsage.Body.Outbound {
		value, ok := aggregatedSubServicesMap[*usageDataRecord.Service]
		if !ok {
			generalInfoData := GeneralInfoData{}
			generalInfoData.Service = *usageDataRecord.Service
			generalInfoData.Unit = *usageDataRecord.Unit
			generalInfoData.InboundPartnerUsage = *usageDataRecord.Usage
			aggregatedSubServicesMap[*usageDataRecord.Service] = &generalInfoData

		} else {
			summary := value.InboundPartnerUsage + *(usageDataRecord.Usage)
			value.InboundPartnerUsage = summary

		}
	}

	// general information - outbound own usage
	for _, usageDataRecord := range ownUsage.Body.Outbound {
		value, ok := aggregatedSubServicesMap[*usageDataRecord.Service]
		if !ok {
			generalInfoData := GeneralInfoData{}
			generalInfoData.Service = *usageDataRecord.Service
			generalInfoData.Unit = *usageDataRecord.Unit
			generalInfoData.OutboundOwnUsage = *usageDataRecord.Usage
			aggregatedSubServicesMap[*usageDataRecord.Service] = &generalInfoData

		} else {
			summary := value.OutboundOwnUsage + *(usageDataRecord.Usage)
			value.OutboundOwnUsage = summary
		}
	}

	// general information - outbound partner usage
	for _, usageDataRecord := range partnerUsage.Body.Inbound {
		value, ok := aggregatedSubServicesMap[*usageDataRecord.Service]
		if !ok {
			generalInfoData := GeneralInfoData{}
			generalInfoData.Service = *usageDataRecord.Service
			generalInfoData.Unit = *usageDataRecord.Unit
			generalInfoData.OutboundPartnerUsage = *usageDataRecord.Usage
			aggregatedSubServicesMap[*usageDataRecord.Service] = &generalInfoData

		} else {

			summary := value.OutboundPartnerUsage + *(usageDataRecord.Usage)
			value.OutboundPartnerUsage = summary

		}
	}

	// create general information array for sub-services
	generalInformationSubServiceArray := make([]GeneralInfoData, 0, len(aggregatedSubServicesMap))

	for _, value := range aggregatedSubServicesMap {
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

	// retrieve two settlements from the request body
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.NoContent(http.StatusNotAcceptable)
	}

	var req []Settlement // the non-struct body
	if body != nil {
		err := json.Unmarshal(body, &req)
		if err != nil {
			return ctx.NoContent(http.StatusNotAcceptable)
		}
	} else {
		return ctx.NoContent(http.StatusNotAcceptable)
	}

	homeSettlement := req[0] // assumption: first settlement is a home one
	partnerSettlement := req[1]

	fmt.Println(homeSettlement.Header.Context)
	fmt.Println(partnerSettlement.Header.Context)

	// home inbound
	homeInboundVoiceServicesMap := createVoiceServicesMap(homeSettlement.Body.Inbound)
	homeInboundSmsServicesMap := createSMSServicesMap(homeSettlement.Body.Inbound)
	homeInboundDataServicesMap := createDataServicesMap(homeSettlement.Body.Inbound)
	// home outbound
	homeOutboundVoiceServicesMap := createVoiceServicesMap(homeSettlement.Body.Outbound)
	homeOutboundSmsServicesMap := createSMSServicesMap(homeSettlement.Body.Outbound)
	homeOutboundDataServicesMap := createDataServicesMap(homeSettlement.Body.Outbound)

	// partner outbound
	partnerOutboundVoiceServicesMap := createVoiceServicesMap(partnerSettlement.Body.Outbound)
	partnerOutboundSmsServicesMap := createSMSServicesMap(partnerSettlement.Body.Outbound)
	partnerOutboundDataServicesMap := createDataServicesMap(partnerSettlement.Body.Outbound)
	// partner inbound
	partnerInboundVoiceServicesMap := createVoiceServicesMap(partnerSettlement.Body.Inbound)
	partnerInboundSmsServicesMap := createSMSServicesMap(partnerSettlement.Body.Inbound)
	partnerInboundDataServicesMap := createDataServicesMap(partnerSettlement.Body.Inbound)

	// HOME PERSPECTIVE
	// Home Perspective details: home inbound & partner outbound
	homePerspectiveDetails := make([]SettlementDiscrepancyRecord, 0)

	// voice sub-services details
	createSubServicesDetails(homeInboundVoiceServicesMap, partnerOutboundVoiceServicesMap, "min", &homePerspectiveDetails)

	// SMS sub-services details
	createSubServicesDetails(homeInboundSmsServicesMap, partnerOutboundSmsServicesMap, "SMS", &homePerspectiveDetails)

	// data sub-services details
	createSubServicesDetails(homeInboundDataServicesMap, partnerOutboundDataServicesMap, "MB", &homePerspectiveDetails)

	// Home Perspective general information: home inbound & partner outbound
	homePerspectiveGeneralInfo := make([]SettlementDiscrepancyRecord, 0)

	// voice general information
	createGeneralInformation(homeInboundVoiceServicesMap, partnerOutboundVoiceServicesMap, "Voice", "min", &homePerspectiveGeneralInfo)

	// SMS general information
	createGeneralInformation(homeInboundSmsServicesMap, partnerOutboundSmsServicesMap, "SMS", "SMS", &homePerspectiveGeneralInfo)

	// data general information
	createGeneralInformation(homeInboundDataServicesMap, partnerOutboundDataServicesMap, "Data", "MB", &homePerspectiveGeneralInfo)

	// PARTNER PERSPECTIVE
	// Partner Perspective details: partner inbound & home outbound
	partnerPerspectiveDetails := make([]SettlementDiscrepancyRecord, 0)

	// voice sub-services details
	createSubServicesDetails(partnerInboundVoiceServicesMap, homeOutboundVoiceServicesMap, "min", &partnerPerspectiveDetails)

	// SMS sub-services details
	createSubServicesDetails(partnerInboundSmsServicesMap, homeOutboundSmsServicesMap, "SMS", &partnerPerspectiveDetails)

	// data sub-services details
	createSubServicesDetails(partnerInboundDataServicesMap, homeOutboundDataServicesMap, "MB", &partnerPerspectiveDetails)

	// Partner Perspective general information: partner inbound & home outbound
	partnerPerspectiveGeneralInfo := make([]SettlementDiscrepancyRecord, 0)

	// voice general information
	createGeneralInformation(partnerInboundVoiceServicesMap, homeOutboundVoiceServicesMap, "Voice", "min", &partnerPerspectiveGeneralInfo)

	// SMS general information
	createGeneralInformation(partnerInboundSmsServicesMap, homeOutboundSmsServicesMap, "SMS", "SMS", &partnerPerspectiveGeneralInfo)

	// data general information
	createGeneralInformation(partnerInboundDataServicesMap, homeOutboundDataServicesMap, "Data", "MB", &partnerPerspectiveGeneralInfo)

	// create final report
	report := SettlementDiscrepancyReport{}

	report.HomePerspective = &(struct {
		Details            []SettlementDiscrepancyRecord `json:"details"`
		GeneralInformation []SettlementDiscrepancyRecord `json:"general_information"`
	}{homePerspectiveDetails, homePerspectiveGeneralInfo})

	report.PartnerPerspective = &(struct {
		Details            []SettlementDiscrepancyRecord `json:"details"`
		GeneralInformation []SettlementDiscrepancyRecord `json:"general_information"`
	}{partnerPerspectiveDetails, partnerPerspectiveGeneralInfo})

	return ctx.JSON(http.StatusOK, report)
}

func createSubServicesDetails(ownMap, partnerMap map[string]float32, units string, details *[]SettlementDiscrepancyRecord) {
	for key, ownCalculation := range ownMap {
		partnerCalculation := partnerMap[key]
		var discrepancyRecord = SettlementDiscrepancyRecord{}
		discrepancyRecord.Service = key
		discrepancyRecord.Unit = units
		discrepancyRecord.OwnCalculation = ownCalculation
		discrepancyRecord.PartnerCalculation = partnerCalculation
		discrepancyRecord.DeltaCalculationPercent = calculateRelativeDelta(ownCalculation, partnerCalculation)
		*details = append(*details, discrepancyRecord)
	}
}

func createGeneralInformation(ownMap, partnerMap map[string]float32, service, units string, generalInfoArr *[]SettlementDiscrepancyRecord) {
	ownCalculationTotalAmount := float32(0)
	for _, value := range ownMap {
		ownCalculationTotalAmount += value
	}
	partnerCalculationTotalAmount := float32(0)
	for _, value := range partnerMap {
		partnerCalculationTotalAmount += value
	}
	discrepancyRecord := SettlementDiscrepancyRecord{}
	discrepancyRecord.Service = service
	discrepancyRecord.Unit = units
	discrepancyRecord.OwnCalculation = ownCalculationTotalAmount
	discrepancyRecord.PartnerCalculation = partnerCalculationTotalAmount
	discrepancyRecord.DeltaCalculationPercent = calculateRelativeDelta(ownCalculationTotalAmount, partnerCalculationTotalAmount)
	*generalInfoArr = append(*generalInfoArr, discrepancyRecord)
}

func calculateRelativeDelta(A, B float32) float32 {
	// relative delta
	// [ (A-B) / A] x 100
	C := ((A - B) / A) * 100.0
	return C
}

func createVoiceServicesMap(input SettlementServices) map[string]float32 {
	fmt.Println("Voice services values:")

	fmt.Println(input.Services.Voice.MOC.BackHome)
	fmt.Println(input.Services.Voice.MOC.International)
	fmt.Println(input.Services.Voice.MOC.Local)
	fmt.Println(input.Services.Voice.MOC.Premium)
	fmt.Println(input.Services.Voice.MOC.ROW)

	voiceServicesMap := make(map[string]float32, 0)

	backHome := input.Services.Voice.MOC.BackHome
	local := input.Services.Voice.MOC.Local
	premium := input.Services.Voice.MOC.Premium
	international := input.Services.Voice.MOC.International
	ROW := input.Services.Voice.MOC.ROW

	if backHome != nil {
		fmt.Printf("backHome: %f\n", *backHome)
		voiceServicesMap["MOC Back Home"] = *backHome
	}
	if local != nil {
		fmt.Printf("local: %f\n", *local)
		voiceServicesMap["MOC Local"] = *local
	}
	if premium != nil {
		fmt.Printf("premium: %f\n", *premium)
		voiceServicesMap["MOC Premium"] = *premium
	}
	if international != nil {
		fmt.Printf("international: %f\n", *international)
		voiceServicesMap["MOC International"] = *international
	}
	if ROW != nil {
		fmt.Printf("ROW: %f\n", *ROW)
	}

	// TODO: Add support for MTC

	return voiceServicesMap
}

func createSMSServicesMap(input SettlementServices) map[string]float32 {
	smsMO := input.Services.SMS.MO
	smsMT := input.Services.SMS.MT

	smsServicesMap := make(map[string]float32, 0)

	if smsMO != nil {
		smsServicesMap["SMSMO"] = *smsMO
	}
	if smsMT != nil {
		smsServicesMap["SMSMT"] = *smsMT
	}

	return smsServicesMap
}

func createDataServicesMap(input SettlementServices) map[string]float32 {
	dataServicesMap := make(map[string]float32, 0)

	for _, element := range input.Services.Data {
		dataServicesMap[*element.Name] = *element.Value
	}

	return dataServicesMap
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