package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/tkanos/gonfig"
	"gopkg.in/yaml.v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configuration struct {
	Server struct {
		Connection_String string `yaml:"connection_string" envconfig:"MONGO_CONN_URL"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"user" envconfig:"DB_USER"`
		Password string `yaml:"pass" envconfig:"DB_PASSWD"`
	} `yaml:"database"`
}

type ServiceUsage struct {
	ID    string  `bson:"_id,omitempty"`
	Total float64 `bson:"total,omitempty"`
}

type DiscrepancyServer struct {
	NextId     int64
	Lock       sync.Mutex
	config     Configuration
	credential options.Credential
}

func NewDiscrepancyServer() *DiscrepancyServer {
	fmt.Println("Starting service...")

	var config Configuration

	// db connection string and credentials provisioned in config file but they can be overriden from env vars
	readFile(&config)
	readEnv(&config)

	fmt.Printf("DB connection string: %s\n", config.Server.Connection_String)
	fmt.Printf("DB username:: %s\n", config.Database.Username)
	fmt.Printf("DB password: %s\n", config.Database.Password)

	dbAccessCredentials := options.Credential{
		Username: config.Database.Username,
		Password: config.Database.Password,
	}

	return &DiscrepancyServer{
		config:     config,
		credential: dbAccessCredentials,
	}
}

func readFile(cfg *Configuration) {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Configuration) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func (p *DiscrepancyServer) CalculateUsageDiscrepancy(ctx echo.Context, usageId string, params CalculateUsageDiscrepancyParams) error {
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

	// later on we can get usage aggregations for the settlement discrepancy purpose
	p.saveUsageReportsToLocalDB(ownUsage, partnerUsage)

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
			generalInfoData.Unit = *usageDataRecord.Units
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
			generalInfoData.Unit = *usageDataRecord.Units
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
			generalInfoData.Unit = *usageDataRecord.Units
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
			generalInfoData.Unit = *usageDataRecord.Units
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
	// MOC
	voiceGeneralInformation := GeneralInfoData{}
	moc := "MOC"
	voiceGeneralInformation.Service = moc
	min := "min"
	voiceGeneralInformation.Unit = min

	voiceGeneralInformation.InboundOwnUsage = 0
	voiceGeneralInformation.InboundPartnerUsage = 0
	voiceGeneralInformation.OutboundOwnUsage = 0
	voiceGeneralInformation.OutboundPartnerUsage = 0

	// MTC
	voiceMTCGeneralInformation := GeneralInfoData{}
	mtc := "MTC"
	voiceMTCGeneralInformation.Service = mtc
	voiceMTCGeneralInformation.Unit = min

	voiceMTCGeneralInformation.InboundOwnUsage = 0
	voiceMTCGeneralInformation.InboundPartnerUsage = 0
	voiceMTCGeneralInformation.OutboundOwnUsage = 0
	voiceMTCGeneralInformation.OutboundPartnerUsage = 0

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
			if element.Service == "MTC" {
				voiceMTCGeneralInformation.InboundOwnUsage += element.InboundOwnUsage
				voiceMTCGeneralInformation.InboundPartnerUsage += element.InboundPartnerUsage
				voiceMTCGeneralInformation.OutboundOwnUsage += element.OutboundOwnUsage
				voiceMTCGeneralInformation.OutboundPartnerUsage += element.OutboundPartnerUsage
			} else {

				voiceGeneralInformation.InboundOwnUsage += element.InboundOwnUsage
				voiceGeneralInformation.InboundPartnerUsage += element.InboundPartnerUsage
				voiceGeneralInformation.OutboundOwnUsage += element.OutboundOwnUsage
				voiceGeneralInformation.OutboundPartnerUsage += element.OutboundPartnerUsage
			}

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

	generalInformationBearerServiceArray := make([]GeneralInfoData, 4, 4)
	generalInformationBearerServiceArray[0] = calculateInOutDiscrepancies(&voiceGeneralInformation)
	generalInformationBearerServiceArray[1] = calculateInOutDiscrepancies(&voiceMTCGeneralInformation)
	generalInformationBearerServiceArray[2] = calculateInOutDiscrepancies(&smsGeneralInformation)
	generalInformationBearerServiceArray[3] = calculateInOutDiscrepancies(&dataGeneralInformation)

	report.GeneralInformation = &generalInformationBearerServiceArray

	// inbound details
	homeInboundMap := p.convertUsageDataArrayToMap(ownUsage.Body.Inbound)
	partnerOutboundMap := p.convertUsageDataArrayToMap(partnerUsage.Body.Outbound)

	inbound := make([]UsageDiscrepancyData, 0)

	for key, inUsage := range homeInboundMap {
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
		inUsage, ok := partnerInboundMap[key]
		if ok {
			outboundUsageDiscrepancyData := createInOutDetailsRecord(outUsage, inUsage)
			outbound = append(outbound, outboundUsageDiscrepancyData)
		}
	}

	report.Outbound = &outbound

	return ctx.JSON(http.StatusOK, report)
}

func (p *DiscrepancyServer) saveUsageReportsToLocalDB(home, partner Usage) {

	client, err := mongo.NewClient(options.Client().ApplyURI(p.config.Server.Connection_String).SetAuth(p.credential))
	if err != nil {
		log.Fatal(err)
		return
	}

	dbCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(dbCtx)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("nomad").Collection("usages")

	// only one timeline is supported
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the usages collection\n", deleteResult.DeletedCount)

	// insert home usage
	insertResult, err := collection.InsertOne(context.TODO(), home)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a Single Document: ", insertResult.InsertedID)

	// insert partner usage
	insertResult, err = collection.InsertOne(context.TODO(), partner)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a Single Document: ", insertResult.InsertedID)

	defer client.Disconnect(dbCtx)

}

func calculateInOutDiscrepancies(value *GeneralInfoData) GeneralInfoData {
	delta64 := float64(value.InboundOwnUsage) - float64(value.InboundPartnerUsage)
	absDelta64 := math.Abs(delta64)
	value.InboundDiscrepancy = absDelta64

	delta64 = float64(value.OutboundOwnUsage) - float64(value.OutboundPartnerUsage)
	absDelta64 = math.Abs(delta64)
	value.OutboundDiscrepancy = absDelta64

	return *value
}

func createInOutDetailsRecord(ownUsage, partnerUsage UsageData) UsageDiscrepancyData {

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
	record.DeltaUsageAbs = &absDelta64
	// relative delta
	C := calculateRelativeDelta64(*ownUsage.Usage, *partnerUsage.Usage)
	record.DeltaUsagePercent = &C

	return record
}

func (p *DiscrepancyServer) convertUsageDataArrayToMap(arr []UsageData) map[string]UsageData {
	fmt.Println("Start: convertUsageDataArrayToMap")

	// create output map
	m := make(map[string]UsageData)

	for _, element := range arr {
		// fmt.Println("At index", index, "value is", toString(element))

		compositeUsageId := makeUsageIdentifier(element)
		// fmt.Println("compositeUsageId", compositeUsageId)

		var data = []byte(compositeUsageId)
		var dataBase64 = base64.StdEncoding.EncodeToString(data)
		sha256 := sha256.Sum256([]byte(dataBase64))
		hashKey := hex.EncodeToString(sha256[:])

		// sets the hash based key to the given element
		m[hashKey] = element

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

	configuration := Configuration{}
	err := gonfig.GetConf("config/config.json", &configuration)
	if err != nil {

	}

	fmt.Printf("Connection string: %s\n", configuration.Server.Connection_String)

	var usage Usage
	dtag := "DTAG"
	version := "1.0"
	usage.Header.MspOwner = &dtag
	usage.Header.Version = version

	return ctx.JSON(http.StatusOK, usage)
}

func (p *DiscrepancyServer) createSubServicesWithUsagesMap(perspective, direction string) map[string]float64 {
	fmt.Println("createServicesWithUsagesMap")
	fmt.Println(perspective)
	fmt.Println(direction)

	client, err := mongo.NewClient(options.Client().ApplyURI(p.config.Server.Connection_String).SetAuth(p.credential))
	if err != nil {
		log.Fatal(err)
	}
	dbCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(dbCtx)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// db.usages.aggregate( [ { $unwind: "$header" }, { $match: { "header.context": "home" } }, { $unwind: "$body.inbound" },
	// { $group: { _id: "$body.inbound.service", total: { $sum: "$body.inbound.usage" } } } ] )

	usagesCollection := client.Database("nomad").Collection("usages")
	unwindStage1 := bson.D{{"$unwind", "$header"}}
	matchStage := bson.D{{"$match", bson.D{{"header.context", perspective}}}}

	bodyDirection := "$body." + direction
	fmt.Println(bodyDirection)
	bodyDirectionWithUsages := "$body." + direction + ".usage"
	bodyDirectionWithService := "$body." + direction + ".service"

	unwindStage2 := bson.D{{"$unwind", bodyDirection}}

	// groupStage := bson.D{{"$group", bson.D{{"_id", "$body.inbound.service"}, {"total", bson.D{{"$sum", bodyDirectionWithUsages}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", bodyDirectionWithService}, {"total", bson.D{{"$sum", bodyDirectionWithUsages}}}}}}

	serviceUsageCursor, err := usagesCollection.Aggregate(dbCtx, mongo.Pipeline{unwindStage1, matchStage, unwindStage2, groupStage})

	if err != nil {
		panic(err)
	}

	var serviceUsages []ServiceUsage
	if serviceUsageCursor.TryNext(dbCtx) {
		if err = serviceUsageCursor.All(dbCtx, &serviceUsages); err != nil {
			panic(err)
		}
	}

	servicesMap := make(map[string]float64, len(serviceUsages))

	for _, element := range serviceUsages {
		// fmt.Println(element.ID)
		// fmt.Println(element.Total)
		servicesMap[element.ID] = element.Total
	}

	fmt.Println(servicesMap)

	defer client.Disconnect(dbCtx)

	return servicesMap
}

func (p *DiscrepancyServer) createBearerServicesWithUsagesMap(perspective, direction string) map[string]float64 {
	fmt.Println("createBearerServicesWithUsagesMap")
	fmt.Println(perspective)
	fmt.Println(direction)

	client, err := mongo.NewClient(options.Client().ApplyURI(p.config.Server.Connection_String).SetAuth(p.credential))
	if err != nil {
		log.Fatal(err)
	}
	dbCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(dbCtx)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// db.usages.aggregate( [ { $unwind: "$header" }, { $match: { "header.context": "home" } }, { $unwind: "$body.inbound" },
	// { $group: { _id: "$body.inbound.units", total: { $sum: "$body.inbound.usage" } } } ] )

	usagesCollection := client.Database("nomad").Collection("usages")
	unwindStage1 := bson.D{{"$unwind", "$header"}}
	matchStage := bson.D{{"$match", bson.D{{"header.context", perspective}}}}

	bodyDirection := "$body." + direction
	fmt.Println(bodyDirection)
	bodyDirectionWithUsages := "$body." + direction + ".usage"
	bodyDirectionWithService := "$body." + direction + ".units"

	unwindStage2 := bson.D{{"$unwind", bodyDirection}}

	// groupStage := bson.D{{"$group", bson.D{{"_id", "$body.inbound.service"}, {"total", bson.D{{"$sum", bodyDirectionWithUsages}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", bodyDirectionWithService}, {"total", bson.D{{"$sum", bodyDirectionWithUsages}}}}}}

	serviceUsageCursor, err := usagesCollection.Aggregate(dbCtx, mongo.Pipeline{unwindStage1, matchStage, unwindStage2, groupStage})

	if err != nil {
		panic(err)
	}

	var serviceUsages []ServiceUsage
	if serviceUsageCursor.TryNext(dbCtx) {
		if err = serviceUsageCursor.All(dbCtx, &serviceUsages); err != nil {
			panic(err)
		}
	}

	fmt.Println(serviceUsages)

	servicesMap := make(map[string]float64, len(serviceUsages))

	for _, element := range serviceUsages {
		fmt.Println(element.ID)
		fmt.Println(element.Total)
		servicesMap[element.ID] = element.Total
	}

	defer client.Disconnect(dbCtx)

	return servicesMap
}

func mergeMaps(bearserServiceUsageMap, subServiceUsageMap map[string]float64) map[string]float64 {
	fmt.Println("mergeMaps")

	MTC := "MTC"
	MOC := "MOC"
	mtcUsage := subServiceUsageMap[MTC]

	if mtcUsage != 0 {
		bearserServiceUsageMap[MTC] = mtcUsage
		bearserServiceUsageMap[MOC] = bearserServiceUsageMap["min"] - mtcUsage
	} else {
		bearserServiceUsageMap[MTC] = 0
		bearserServiceUsageMap[MOC] = bearserServiceUsageMap["min"]
	}

	for key, element := range bearserServiceUsageMap {
		fmt.Println("Service:", key, "=>", "Usage:", element)
	}

	return bearserServiceUsageMap
}

func (p *DiscrepancyServer) CalculateSettlementDiscrepancy(ctx echo.Context, settlementId string, params CalculateSettlementDiscrepancyParams) error {
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

	// prettyJSON, err := json.MarshalIndent(homeSettlement, "", "    ")
	// if err != nil {
	// 	log.Fatal("Failed to generate json", err)
	// }
	// fmt.Printf("%s\n", string(prettyJSON))

	fmt.Println(homeSettlement.Header.Context)
	fmt.Println(partnerSettlement.Header.Context)

	// SERVICES MAPS:
	// home inbound
	homeInboundMOCServicesMap := createMOCServicesMap(homeSettlement.Body.Inbound)
	homeInboundMTCServicesMap := createMTCServicesMap(homeSettlement.Body.Inbound)
	homeInboundSmsServicesMap := createSMSServicesMap(homeSettlement.Body.Inbound)
	homeInboundDataServicesMap := createDataServicesMap(homeSettlement.Body.Inbound)
	// home outbound
	homeOutboundMOCServicesMap := createMOCServicesMap(homeSettlement.Body.Outbound)
	homeOutboundMTCServicesMap := createMTCServicesMap(homeSettlement.Body.Outbound)
	homeOutboundSmsServicesMap := createSMSServicesMap(homeSettlement.Body.Outbound)
	homeOutboundDataServicesMap := createDataServicesMap(homeSettlement.Body.Outbound)

	// partner outbound
	partnerOutboundMOCServicesMap := createMOCServicesMap(partnerSettlement.Body.Outbound)
	partnerOutboundMTCServicesMap := createMTCServicesMap(partnerSettlement.Body.Outbound)
	partnerOutboundSmsServicesMap := createSMSServicesMap(partnerSettlement.Body.Outbound)
	partnerOutboundDataServicesMap := createDataServicesMap(partnerSettlement.Body.Outbound)
	// partner inbound
	partnerInboundMOCServicesMap := createMOCServicesMap(partnerSettlement.Body.Inbound)
	partnerInboundMTCServicesMap := createMTCServicesMap(partnerSettlement.Body.Inbound)
	partnerInboundSmsServicesMap := createSMSServicesMap(partnerSettlement.Body.Inbound)
	partnerInboundDataServicesMap := createDataServicesMap(partnerSettlement.Body.Inbound)

	// PRECOMMITMENT VALUES BLOCK

	// PRECOMMITMENT VALUES BLOCK - DELTA
	homeMOCDeltaCommitment := calculateDelta(homeInboundMOCServicesMap)
	homeMTCDeltaCommitment := calculateDelta(homeInboundMTCServicesMap)
	homeSMSDeltaCommitment := calculateDelta(homeInboundSmsServicesMap)
	homeDataDeltaCommitment := calculateDelta(homeInboundDataServicesMap)

	homeDeltaCommitment := homeMOCDeltaCommitment + homeMTCDeltaCommitment + homeSMSDeltaCommitment + homeDataDeltaCommitment

	// fmt.Printf("homeMOCDeltaCommitment : %f \n", homeMOCDeltaCommitment)
	// fmt.Printf("homeMTCDeltaCommitment : %f \n", homeMTCDeltaCommitment)
	// fmt.Printf("homeSMSDeltaCommitment : %f \n", homeSMSDeltaCommitment)
	// fmt.Printf("homeDataDeltaCommitment : %f \n", homeDataDeltaCommitment)

	partnerMOCDeltaCommitment := calculateDelta(partnerInboundMOCServicesMap)
	partnerMTCDeltaCommitment := calculateDelta(partnerInboundMTCServicesMap)
	partnerSMSDeltaCommitment := calculateDelta(partnerInboundSmsServicesMap)
	partnerDataDeltaCommitment := calculateDelta(partnerInboundDataServicesMap)

	partnerDeltaCommitment := partnerMOCDeltaCommitment + partnerMTCDeltaCommitment + partnerSMSDeltaCommitment + partnerDataDeltaCommitment

	// PRECOMMITMENT VALUES BLOCK - RECALCULATE TO POSTCOMMITMENT VALUES

	recalculateDealValues(&homeInboundMOCServicesMap)
	recalculateDealValues(&homeInboundMTCServicesMap)
	recalculateDealValues(&homeInboundSmsServicesMap)
	recalculateDealValues(&homeInboundDataServicesMap)

	recalculateDealValues(&homeOutboundMOCServicesMap)
	recalculateDealValues(&homeOutboundMTCServicesMap)
	recalculateDealValues(&homeOutboundSmsServicesMap)
	recalculateDealValues(&homeOutboundDataServicesMap)

	recalculateDealValues(&partnerInboundMOCServicesMap)
	recalculateDealValues(&partnerInboundMTCServicesMap)
	recalculateDealValues(&partnerInboundSmsServicesMap)
	recalculateDealValues(&partnerInboundDataServicesMap)

	recalculateDealValues(&partnerOutboundMOCServicesMap)
	recalculateDealValues(&partnerOutboundMTCServicesMap)
	recalculateDealValues(&partnerOutboundSmsServicesMap)
	recalculateDealValues(&partnerOutboundDataServicesMap)

	// USAGES:

	// HOME PERSPECTIVE
	// sub-services with usages maps
	homeInboundServiceUsageMap := p.createSubServicesWithUsagesMap("home", "inbound")
	partnerOutboundServiceUsageMap := p.createSubServicesWithUsagesMap("partner", "outbound")
	// bearer services with usages maps
	homeInboundBearerServiceUsageMap := p.createBearerServicesWithUsagesMap("home", "inbound")
	homeInboundBearerServiceUsageMap = mergeMaps(homeInboundBearerServiceUsageMap, homeInboundServiceUsageMap)

	partnerOutboundBearerServiceUsageMap := p.createBearerServicesWithUsagesMap("partner", "outbound")
	partnerOutboundBearerServiceUsageMap = mergeMaps(partnerOutboundBearerServiceUsageMap, partnerOutboundServiceUsageMap)

	// PARTNER PERSPECTIVE
	// sub-services with usages maps
	partnerInboundServiceUsageMap := p.createSubServicesWithUsagesMap("partner", "inbound")
	homeOutboundServiceUsageMap := p.createSubServicesWithUsagesMap("home", "outbound")
	// bearer services with usages maps
	partnerInboundBearerServiceUsageMap := p.createBearerServicesWithUsagesMap("partner", "inbound")
	partnerInboundBearerServiceUsageMap = mergeMaps(partnerInboundBearerServiceUsageMap, partnerInboundServiceUsageMap)

	homeOutboundBearerServiceUsageMap := p.createBearerServicesWithUsagesMap("home", "outbound")
	homeOutboundBearerServiceUsageMap = mergeMaps(homeOutboundBearerServiceUsageMap, homeOutboundServiceUsageMap)

	// DISCREPANCY REPORT:

	// HOME PERSPECTIVE
	// Home Perspective details: home inbound & partner outbound
	homePerspectiveDetails := make([]SettlementDiscrepancyRecord, 0)

	// MOC sub-services details
	createSubServicesDetails(homeInboundMOCServicesMap, partnerOutboundMOCServicesMap, "min", &homePerspectiveDetails,
		homeInboundServiceUsageMap, partnerOutboundServiceUsageMap)

	// MTC sub-services details
	createSubServicesDetails(homeInboundMTCServicesMap, partnerOutboundMTCServicesMap, "min", &homePerspectiveDetails,
		homeInboundServiceUsageMap, partnerOutboundServiceUsageMap)

	// SMS sub-services details
	createSubServicesDetails(homeInboundSmsServicesMap, partnerOutboundSmsServicesMap, "SMS", &homePerspectiveDetails,
		homeInboundServiceUsageMap, partnerOutboundServiceUsageMap)

	// data sub-services details
	createSubServicesDetails(homeInboundDataServicesMap, partnerOutboundDataServicesMap, "MB", &homePerspectiveDetails,
		homeInboundServiceUsageMap, partnerOutboundServiceUsageMap)

	// Home Perspective general information: home inbound & partner outbound
	homePerspectiveGeneralInfo := make([]SettlementDiscrepancyRecord, 0)

	// MOC general information
	homeInboundTotalMOCCalculation, partnerOutboundTotalMOCCalculation := createGeneralInformation(homeInboundMOCServicesMap, partnerOutboundMOCServicesMap, "MOC", "MOC", &homePerspectiveGeneralInfo,
		homeInboundBearerServiceUsageMap, partnerOutboundBearerServiceUsageMap)

	// MTC general information
	homeInboundTotalMTCCalculation, partnerOutboundTotalMTCCalculation := createGeneralInformation(homeInboundMTCServicesMap, partnerOutboundMTCServicesMap, "MTC", "MTC", &homePerspectiveGeneralInfo,
		homeInboundBearerServiceUsageMap, partnerOutboundBearerServiceUsageMap)

	// SMS general information
	homeInboundTotalSMSCalculation, partnerOutboundTotalSMSCalculation := createGeneralInformation(homeInboundSmsServicesMap, partnerOutboundSmsServicesMap, "SMS", "SMS", &homePerspectiveGeneralInfo,
		homeInboundBearerServiceUsageMap, partnerOutboundBearerServiceUsageMap)

	// data general information
	homeInboundTotalDataCalculation, partnerOutboundTotalDataCalculation := createGeneralInformation(homeInboundDataServicesMap, partnerOutboundDataServicesMap, "Data", "MB", &homePerspectiveGeneralInfo,
		homeInboundBearerServiceUsageMap, partnerOutboundBearerServiceUsageMap)

	// PARTNER PERSPECTIVE
	// Partner Perspective details: partner inbound & home outbound
	partnerPerspectiveDetails := make([]SettlementDiscrepancyRecord, 0)

	// MOC sub-services details
	createSubServicesDetails(partnerInboundMOCServicesMap, homeOutboundMOCServicesMap, "min", &partnerPerspectiveDetails,
		partnerInboundServiceUsageMap, homeOutboundServiceUsageMap)

	// MTC sub-services details
	createSubServicesDetails(partnerInboundMTCServicesMap, homeOutboundMTCServicesMap, "min", &partnerPerspectiveDetails,
		partnerInboundServiceUsageMap, homeOutboundServiceUsageMap)

	// SMS sub-services details
	createSubServicesDetails(partnerInboundSmsServicesMap, homeOutboundSmsServicesMap, "SMS", &partnerPerspectiveDetails,
		partnerInboundServiceUsageMap, homeOutboundServiceUsageMap)

	// data sub-services details
	createSubServicesDetails(partnerInboundDataServicesMap, homeOutboundDataServicesMap, "MB", &partnerPerspectiveDetails,
		partnerInboundServiceUsageMap, homeOutboundServiceUsageMap)

	// Partner Perspective general information: partner inbound & home outbound
	partnerPerspectiveGeneralInfo := make([]SettlementDiscrepancyRecord, 0)

	// MOC general information
	partnerInboundTotalMOCCalculation, homeOutboundTotalMOCCalculation := createGeneralInformation(partnerInboundMOCServicesMap, homeOutboundMOCServicesMap, "MOC", "MOC", &partnerPerspectiveGeneralInfo,
		partnerInboundBearerServiceUsageMap, homeOutboundBearerServiceUsageMap)

	// MTC general information
	partnerInboundTotalMTCCalculation, homeOutboundTotalMTCCalculation := createGeneralInformation(partnerInboundMTCServicesMap, homeOutboundMTCServicesMap, "MTC", "MTC", &partnerPerspectiveGeneralInfo,
		partnerInboundBearerServiceUsageMap, homeOutboundBearerServiceUsageMap)

	// SMS general information
	partnerInboundTotalSMSCalculation, homeOutboundTotalSMSCalculation := createGeneralInformation(partnerInboundSmsServicesMap, homeOutboundSmsServicesMap, "SMS", "SMS", &partnerPerspectiveGeneralInfo,
		partnerInboundBearerServiceUsageMap, homeOutboundBearerServiceUsageMap)

	// data general information
	partnerInboundTotalDataCalculation, homeOutboundTotalDataCalculation := createGeneralInformation(partnerInboundDataServicesMap, homeOutboundDataServicesMap, "Data", "MB", &partnerPerspectiveGeneralInfo,
		partnerInboundBearerServiceUsageMap, homeOutboundBearerServiceUsageMap)

	// SETTLEMENT REPORT CALCULATIONS

	homeRevenue := homeInboundTotalMOCCalculation + homeInboundTotalMTCCalculation + homeInboundTotalSMSCalculation + homeInboundTotalDataCalculation
	partnerCharges := partnerOutboundTotalMOCCalculation + partnerOutboundTotalMTCCalculation + partnerOutboundTotalSMSCalculation + partnerOutboundTotalDataCalculation

	partnerRevenue := partnerInboundTotalMOCCalculation + partnerInboundTotalMTCCalculation + partnerInboundTotalSMSCalculation + partnerInboundTotalDataCalculation
	homeCharges := homeOutboundTotalMOCCalculation + homeOutboundTotalMTCCalculation + homeOutboundTotalSMSCalculation + homeOutboundTotalDataCalculation

	// create discrepancy report
	report := SettlementDiscrepancyReport{}

	report.SettlementReport = &(struct {
		HomeCharges            float64 `json:"homeCharges"`
		HomeDeltaCommitment    float64 `json:"homeDeltaCommitment"`
		HomeRevenue            float64 `json:"homeRevenue"`
		PartnerCharges         float64 `json:"partnerCharges"`
		PartnerDeltaCommitment float64 `json:"partnerDeltaCommitment"`
		PartnerRevenue         float64 `json:"partnerRevenue"`
	}{homeCharges, homeDeltaCommitment, homeRevenue, partnerCharges, partnerDeltaCommitment, partnerRevenue})

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

func createSubServicesDetails(ownMap, partnerMap map[string]TelcoService, units string, details *[]SettlementDiscrepancyRecord,
	ownUsageMap, partnerUsageMap map[string]float64) {

	fmt.Println("createSubServicesDetails")

	for key, ownTelcoService := range ownMap {
		partnerTelcoService := partnerMap[key]

		if !(ownTelcoService.DealValue == 0 && partnerTelcoService.DealValue == 0) {
			var discrepancyRecord = SettlementDiscrepancyRecord{}
			discrepancyRecord.Service = key
			discrepancyRecord.Unit = units
			////
			fmt.Printf("key: %s and associoated usages: own = %f, partner = %f\n", key, ownUsageMap[key], partnerUsageMap[key])
			////
			discrepancyRecord.OwnUsage = ownUsageMap[key]
			discrepancyRecord.PartnerUsage = partnerUsageMap[key]
			discrepancyRecord.DeltaUsageAbs = math.Abs(discrepancyRecord.OwnUsage - discrepancyRecord.PartnerUsage)
			discrepancyRecord.DeltaUsagePercent = calculateRelativeDelta64(discrepancyRecord.OwnUsage, discrepancyRecord.PartnerUsage)
			////
			fmt.Printf("DeltaUsageAbs : %f DeltaUsagePercent %f\n", discrepancyRecord.DeltaUsageAbs, discrepancyRecord.DeltaUsagePercent)
			///
			discrepancyRecord.OwnCalculation = ownTelcoService.DealValue
			discrepancyRecord.PartnerCalculation = partnerTelcoService.DealValue
			////
			fmt.Printf("Own calculation : %f partner calculation %f\n", discrepancyRecord.OwnCalculation, discrepancyRecord.PartnerCalculation)
			////
			discrepancyRecord.DeltaCalculationPercent = calculateRelativeDelta64(ownTelcoService.DealValue, partnerTelcoService.DealValue)
			////
			fmt.Printf("DeltaCalculationPercent %f\n", discrepancyRecord.DeltaCalculationPercent)
			////
			*details = append(*details, discrepancyRecord)

		} else {
			////
			fmt.Printf("key: %s and associoated usages: own = %f, partner = %f\n", key, ownUsageMap[key], partnerUsageMap[key])
			////
			if ownUsageMap[key] != 0 || partnerUsageMap[key] != 0 {
				var discrepancyRecord = SettlementDiscrepancyRecord{}
				discrepancyRecord.Service = key
				discrepancyRecord.Unit = units

				discrepancyRecord.OwnUsage = ownUsageMap[key]
				discrepancyRecord.PartnerUsage = partnerUsageMap[key]

				discrepancyRecord.DeltaUsageAbs = math.Abs(discrepancyRecord.OwnUsage - discrepancyRecord.PartnerUsage)
				discrepancyRecord.DeltaUsagePercent = calculateRelativeDelta64(discrepancyRecord.OwnUsage, discrepancyRecord.PartnerUsage)
				fmt.Printf("DeltaUsageAbs : %f DeltaUsagePercent %f\n", discrepancyRecord.DeltaUsageAbs, discrepancyRecord.DeltaUsagePercent)

				discrepancyRecord.OwnCalculation = 0
				discrepancyRecord.PartnerCalculation = 0
				fmt.Printf("Own calculation : %f partner calculation %f\n", discrepancyRecord.OwnCalculation, discrepancyRecord.PartnerCalculation)

				discrepancyRecord.DeltaCalculationPercent = 0
				fmt.Printf("DeltaCalculationPercent %f\n", discrepancyRecord.DeltaCalculationPercent)

				*details = append(*details, discrepancyRecord)
			}
		}
	}
}

func createGeneralInformation(ownMap, partnerMap map[string]TelcoService, service, units string, generalInfoArr *[]SettlementDiscrepancyRecord,
	ownUsageMap, partnerUsageMap map[string]float64) (float64, float64) {

	// perform aggregations
	ownCalculationTotalAmount := float64(0)
	for _, telcoService := range ownMap {
		ownCalculationTotalAmount += telcoService.DealValue
	}
	partnerCalculationTotalAmount := float64(0)
	for _, telcoService := range partnerMap {
		partnerCalculationTotalAmount += telcoService.DealValue
	}
	discrepancyRecord := SettlementDiscrepancyRecord{}
	discrepancyRecord.Service = service

	if service == "MOC" || service == "MTC" {
		discrepancyRecord.Unit = "min"
	} else {
		discrepancyRecord.Unit = units
	}

	// usages
	discrepancyRecord.OwnUsage = ownUsageMap[units]
	discrepancyRecord.PartnerUsage = partnerUsageMap[units]
	discrepancyRecord.DeltaUsageAbs = math.Abs(discrepancyRecord.OwnUsage - discrepancyRecord.PartnerUsage)
	discrepancyRecord.DeltaUsagePercent = calculateRelativeDelta64(discrepancyRecord.OwnUsage, discrepancyRecord.PartnerUsage)
	// calculations
	discrepancyRecord.OwnCalculation = ownCalculationTotalAmount
	discrepancyRecord.PartnerCalculation = partnerCalculationTotalAmount
	discrepancyRecord.DeltaCalculationPercent = calculateRelativeDelta64(ownCalculationTotalAmount, partnerCalculationTotalAmount)
	*generalInfoArr = append(*generalInfoArr, discrepancyRecord)

	return ownCalculationTotalAmount, partnerCalculationTotalAmount
}

func calculateRelativeDelta64(A, B float64) float64 {
	// relative delta
	zero := float64(0)
	if A == zero && B == zero {
		return zero
	}

	if A == zero || B == zero {
		return float64(100)
	}

	// [ (A-B) / A] x 100
	C := ((A - B) / A) * 100.0 // C = percent value
	return C
}

func calculateDelta(services map[string]TelcoService) float64 {
	delta := 0.0

	for _, telcoService := range services {
		delta = delta + telcoService.ShortOfCommitment
	}

	return delta
}

func recalculateDealValues(servicesMap *map[string]TelcoService) {

	if len(*servicesMap) == 0 {
		return
	}

	for _, telcoService := range *servicesMap {

		shortOfCommitment := telcoService.ShortOfCommitment
		if shortOfCommitment <= 0 {
			break
		}

		dealValue := telcoService.DealValue
		if dealValue > 0 {
			telcoService.DealValue = dealValue + shortOfCommitment
		}
	}
	return
}

func createMTCServicesMap(input SettlementServices) map[string]TelcoService {
	voiceServicesMap := make(map[string]TelcoService, 0)
	MTC := input.Services.Voice.MTC

	if MTC != nil {
		fmt.Printf("MTC: %+v\n", *MTC)
		voiceServicesMap["MTC"] = *MTC
	}

	return voiceServicesMap
}

func createMOCServicesMap(input SettlementServices) map[string]TelcoService {
	fmt.Println("MOC services values:")

	voiceServicesMap := make(map[string]TelcoService, 0)

	backHome := input.Services.Voice.MOC.BackHome
	local := input.Services.Voice.MOC.Local
	premium := input.Services.Voice.MOC.Premium
	international := input.Services.Voice.MOC.International
	ROW := input.Services.Voice.MOC.ROW
	specialDestinations := input.Services.Voice.MOC.SpecialDestinations
	EU := input.Services.Voice.MOC.EU
	EEA := input.Services.Voice.MOC.EEA
	satellite := input.Services.Voice.MOC.Satellite
	videoTelephony := input.Services.Voice.MOC.VideoTelephony

	if backHome != nil {
		fmt.Printf("backHome: %+v\n", *backHome)
		voiceServicesMap["MOC Back Home"] = *backHome
	}
	if local != nil {
		fmt.Printf("local: %+v\n", *local)
		voiceServicesMap["MOC Local"] = *local
	}
	if premium != nil {
		fmt.Printf("premium: %+v\n", *premium)
		voiceServicesMap["MOC Premium"] = *premium
	}
	if international != nil {
		fmt.Printf("international: %+v\n", *international)
		voiceServicesMap["MOC International"] = *international
	}
	if ROW != nil {
		fmt.Printf("ROW: %+v\n", *ROW)
		voiceServicesMap["MOC Row"] = *ROW
	}

	if EU != nil {
		fmt.Printf("EU: %+v\n", *EU)
		voiceServicesMap["MOC EU"] = *EU
	}

	if EEA != nil {
		fmt.Printf("EEA: %+v\n", *EEA)
		voiceServicesMap["MOC EEA"] = *EEA
	}

	if specialDestinations != nil {
		fmt.Printf("specialDestinations: %+v\n", *specialDestinations)
		voiceServicesMap["MOC Special Destinations"] = *specialDestinations
	}

	if satellite != nil {
		fmt.Printf("satellite: %+v\n", *satellite)
		voiceServicesMap["MOC Satellite"] = *satellite
	}

	if videoTelephony != nil {
		fmt.Printf("videoTelephony: %+v\n", *videoTelephony)
		voiceServicesMap["MOC videoTelephony"] = *videoTelephony
	}

	return voiceServicesMap
}

func createSMSServicesMap(input SettlementServices) map[string]TelcoService {
	smsMO := input.Services.SMS.MO
	smsMT := input.Services.SMS.MT

	smsServicesMap := make(map[string]TelcoService, 0)

	if smsMO != nil {
		smsServicesMap["SMSMO"] = *smsMO
	}
	if smsMT != nil {
		smsServicesMap["SMSMT"] = *smsMT
	}

	return smsServicesMap
}

func createDataServicesMap(input SettlementServices) map[string]TelcoService {
	dataServicesMap := make(map[string]TelcoService, 0)

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
