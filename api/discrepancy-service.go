package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	fmt.Println("Start")

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

	fmt.Println(*(req[0].Header.Context))
	fmt.Println(*(req[1].Header.Context))
	fmt.Println((*(req[1].Body.Inbound))[0])

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
	generalInfo.Service = &bearerService

	units := "min"
	generalInfo.Unit = &units

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

	generalInfo.InboundOwnUsage = &inbound_own_usage
	generalInfo.InboundPartnerUsage = &inbound_partner_usage
	generalInfo.InboundDiscrepancy = &inbound_discrepancy
	generalInfo.OutboundOwnUsage = &outbound_own_usage
	generalInfo.OutboundPartnerUsage = &outbound_partner_usage
	generalInfo.OutboundDiscrepancy = &outbound_discrepancy

	generalInfoArray := make([]GeneralInfoData, 0)
	generalInfoArray = append(generalInfoArray, generalInfo)
	report.GeneralInformation = &generalInfoArray

	inbound := make([]UsageDiscrepancyData, 0)
	inbound = append(inbound, usageDiscrepancyData)
	report.Inbound = &inbound

	outbound := make([]UsageDiscrepancyData, 0)
	outbound = append(outbound, usageDiscrepancyData)
	report.Outbound = &outbound

	return ctx.JSON(http.StatusOK, report)
}

func (p *DiscrepancyServer) FindUsages(ctx echo.Context) error {
	fmt.Println("Start")

	var usage Usage
	dtag := "DTAG"
	version := "1.0"
	usage.Header.MspOwner = &dtag
	usage.Header.Version = &version

	return ctx.JSON(http.StatusOK, usage)
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
