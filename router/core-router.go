package router

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/theoerp/bham-server/route/security"
	"github.com/SoftwareUndagi/go-common-libs/common"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

//KeyForReqHeaderEditorToken key untuk editor token. token untuk mengindari double submit
const KeyForReqHeaderEditorToken = "X-Custom-Editor-Token"

//optionAllowedRequestHeaderKey request header yang di ijinkan untuk option request
var optionAllowedRequestHeaderKey = []string{"X-Requested-With", "content-type", "authorization", KeyForReqHeaderEditorToken}

//optionAllowedRequestHeaderKeyFlat variable versi flat
var optionAllowedRequestHeaderKeyFlat = strings.Join(optionAllowedRequestHeaderKey, ",")

//AppendAllowedRequestHeader add allowed key to header
func AppendAllowedRequestHeader(headerKey string) {
	optionAllowedRequestHeaderKey = append(optionAllowedRequestHeaderKey, headerKey)
	optionAllowedRequestHeaderKeyFlat = strings.Join(optionAllowedRequestHeaderKey, ",")
}

//flushLogDataCommand command untuk flush log
var flushLogDataCommand func()

//AssignFlushLogCommand assign flush log command
func AssignFlushLogCommand(command func()) {
	flushLogDataCommand = command
}

//RouteLoggerPredefinedParameterFiller parameter filler custom. ini mungkin akan spesifik pada app.
// executionID = id eksekusi. dalam kasus dengan cloud function ini akan di isi dengan id dari function
type RouteLoggerPredefinedParameterFiller func(executionID string, routePath string, req *http.Request, routeParameter Parameter, username string, userUUID string, logEntry *log.Entry) (modifiedLogEntry *log.Entry)

//LoginInformationProviderFunction login handler definition
type LoginInformationProviderFunction func(DatabaseReference *gorm.DB, req *http.Request) (userData security.SimpleUserData, err common.ErrorWithCodeData)

//CORSAllowedPaths path yang di injinkan cors
var CORSAllowedPaths = make(map[string][]string)

//domain yang di ijinkan cors
var corsAllowedDomains = make(map[string]string)

//disableGzipResponse flag disable gzip atau tidak . di disable dengan memanggil method
var disableGzipResponse = false

//LoggerPredefinedParameterFiller filler log data. for generic logger data filler
type LoggerPredefinedParameterFiller func(executionID string, routePath string, req *http.Request, routeParameter Parameter, username string, userUUID string, logEntry *log.Entry) (modifiedLogEntry *log.Entry)

//defaultLoggerPredefinedParameterFiller default filler dari data logger
var defaultLoggerPredefinedParameterFiller = func(executionID string, routePath string, req *http.Request, routeParameter Parameter, username string, userUUID string, logEntry *log.Entry) (modifiedLogEntry *log.Entry) {
	return logEntry.WithField("executionId", executionID).WithField("routePath", routePath).WithField("username", username).WithField("userUuid", userUUID)
}

//CustomDaoAndLoggerAttributeGeneratorFunction used for var customDaoAndLoggerAttributeGenerator.
// this to put additional parameter on log and on dao . Log need entries to ease logging . for example on multi tenant scenario, this will put tenant id etc on data log
// same sample case for dao. app relied on gorm database. for case multi tenant database, table name will be prefixed with schema name for example
type CustomDaoAndLoggerAttributeGeneratorFunction func(executionID string, routePath string, req *http.Request, routeParameter Parameter, username string, userUUID string, baseLogEntry *log.Entry) (logEntry *log.Entry, daoAttribute map[string]interface{})

//customDaoAndLoggerAttributeGenerator if app need custom attribute for dao and log
var customDaoAndLoggerAttributeGenerator CustomDaoAndLoggerAttributeGeneratorFunction = func(executionID string, routePath string, req *http.Request, routeParameter Parameter, username string, userUUID string, baseLogEntry *log.Entry) (logEntry *log.Entry, daoAttribute map[string]interface{}) {
	logEntry = defaultLoggerPredefinedParameterFiller(executionID, routePath, req, routeParameter, username, userUUID, baseLogEntry)

	return
}

//AssignCustomDaoAndLoggerAttributeGenerator replace custom dao attribute generator
//for example on multi tenant database, you need to set user data with user schema name
func AssignCustomDaoAndLoggerAttributeGenerator(f CustomDaoAndLoggerAttributeGeneratorFunction) {
	customDaoAndLoggerAttributeGenerator = f
}

//AssignDefaultLoggerPredefinedParameterFiller replace log entry filler
func AssignDefaultLoggerPredefinedParameterFiller(f LoggerPredefinedParameterFiller) {
	defaultLoggerPredefinedParameterFiller = f
}

//HTTPCommonParameter common http parameter
type HTTPCommonParameter struct {
	//Username username dari current
	Username string
	//UserUUID auth user uuid(firebase thing)
	UserUUID string
	//IPAddress ip address current user
	IPAddress string
	//RequestPath path request
	RequestPath string
	//RawRequest raw reqest parameter
	RawRequest *http.Request
	//DatabaseReference reference to GORM
	DatabaseReference *gorm.DB
	//LogEntry log entry untuk kemudahan logging. common item di inject di awal
	LogEntry *log.Entry
	//PathParameters parameter dalam path. misal path = /alpha/{omega}, parameter omega akan di taruh dalam map
	PathParameters map[string]string
}

//Parameter register route parameter
type Parameter struct {
	//CORSEnabledDomains domain yang di ijinkan CORS
	CORSEnabledDomains []string
	//DatabaseReference reference to GORM
	DatabaseReference *gorm.DB
	//MuxRouter mux router untuk register path
	MuxRouter *mux.Router
	//LoginInformationProvider provider logi information
	LoginInformationProvider LoginInformationProviderFunction
}

//Clone clone data except for database reference
func (p *Parameter) Clone(db *gorm.DB) (cloneResult Parameter) {
	rslt := Parameter{CORSEnabledDomains: p.CORSEnabledDomains, MuxRouter: p.MuxRouter, LoginInformationProvider: p.LoginInformationProvider, DatabaseReference: p.DatabaseReference}
	if db != nil {
		rslt.DatabaseReference = db
	}
	return rslt
}

//ResultJSONOKWrapper wrapper for result OK
type ResultJSONOKWrapper struct {
	HaveError bool        `json:"haveError"`
	Data      interface{} `json:"data"`
}

//ResultJSONErrorWrapper struct for error
type resultJSONErrorWrapper struct {
	HaveError    bool   `json:"haveError"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

//SetDisableGzipResponse disable/enable gzip response
func SetDisableGzipResponse(disable bool) {
	disableGzipResponse = disable
}

//RegisterAllowedCORDomain register domain to allowed cors
func RegisterAllowedCORDomain(domains ...string) {

	if len(domains) > 0 {
		for _, domain := range domains {
			corsAllowedDomains[domain] = domain
		}
	}
}

//RegisterAllowedCORDomains register domain to allowed cors
func RegisterAllowedCORDomains(domains []string) {

	if len(domains) > 0 {
		for _, domain := range domains {
			corsAllowedDomains[domain] = domain
		}
	}
}

//accessDeniedHandler send access denied response. user blm login etc
func accessDeniedHandler(path string, w http.ResponseWriter, req *http.Request) {
	respContent, _ := json.Marshal(resultJSONErrorWrapper{
		HaveError:    true,
		ErrorCode:    "ACCESS_DENIED",
		ErrorMessage: "Not allowed to access path " + path})
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respContent)
}

//sendErrorResponse send http error response
func sendErrorResponse(w http.ResponseWriter, errorData resultJSONErrorWrapper) {
	errorData.HaveError = true
	respContent, _ := json.Marshal(errorData)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respContent)
}

//sendOkData send http ok response
func sendOkData(data interface{}, disableGzipOnParam bool, w http.ResponseWriter, req *http.Request) {
	httpHeader := w.Header()
	httpHeader.Set("Content-Type", "application/json")

	if disableGzipOnParam || disableGzipResponse || !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {

		respContent, _ := json.Marshal(ResultJSONOKWrapper{HaveError: false, Data: data})
		w.Write(respContent)
		w.WriteHeader(http.StatusOK)
		return
	}
	httpHeader.Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	json.NewEncoder(gz).Encode(ResultJSONOKWrapper{HaveError: false, Data: data})
	w.WriteHeader(http.StatusOK)
	gz.Close()
}

func sendNotAllowedRequest(pathToCheck string, w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	respContent, _ := json.Marshal(resultJSONErrorWrapper{HaveError: true, ErrorCode: "CORS_NOT_ALLOWED", ErrorMessage: "Cors request is not allowed for this path"})
	w.Write(respContent)
}

//checkForCors check cors. return false berarti di tolak request
func checkForCors(pathToCheck string, isOption bool, w http.ResponseWriter, req *http.Request) bool {
	if len(req.Header["origin"]) > 0 || len(req.Header["Origin"]) > 0 {
		swapDomain := req.Header["origin"]
		if len(req.Header["Origin"]) > 0 {
			swapDomain = req.Header["Origin"]
		}

		var theDomains []string
		for _, dm := range swapDomain {
			_, isPresent := corsAllowedDomains[dm]
			if isPresent {
				theDomains = append(theDomains, dm)
			}
		}
		if len(theDomains) == 0 {
			sendNotAllowedRequest(pathToCheck, w, req)
			return false
		}
		for _, dmToInvk := range swapDomain {
			w.Header().Add("Access-Control-Allow-Origin", dmToInvk)
		}
		mthd := CORSAllowedPaths[pathToCheck]
		if mthd == nil || len(mthd) == 0 {
			sendNotAllowedRequest(pathToCheck, w, req)
			return false
		}
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		if isOption {
			//w.Header().Add("Access-Control-Max-Age" , "1728000")
			w.Header().Add("Access-Control-Allow-Methods", strings.Join(mthd, ","))
			w.Header().Add("Access-Control-Allow-Headers", optionAllowedRequestHeaderKeyFlat)
		}
	}
	return true
}

//appendOptionRouter add option http handler
func appendOptionRouter(muxRouter *mux.Router, routePath string) {
	theRoute := muxRouter.HandleFunc(routePath, func(w http.ResponseWriter, req *http.Request) {
		log.WithField("path", routePath).Info("Handling options untuk path: " + routePath)
		checkForCors(routePath, true, w, req)
		//w.WriteHeader(http.StatusOK)
		//w.Write([]byte{'O' , 'K'})
	})
	theRoute.Methods("OPTIONS")
}

//generateCommonHTTPParam generate common parameter request
//localDB = db for current connection
func generateCommonHTTPParam(routePath string, req *http.Request, routeParameter Parameter, username string, userUUID string) (requestID string, commonParam HTTPCommonParameter) {
	executionID := fmt.Sprintf("%d", time.Now().UnixNano())
	logEntry := log.WithField("executionId", executionID)
	routeParameter.DatabaseReference.InstantSet("executionId", executionID)
	routeParameter.DatabaseReference.InstantSet("username", username)
	routeParameter.DatabaseReference.InstantSet("userUUID", userUUID)
	requestID = executionID
	if customDaoAndLoggerAttributeGenerator != nil {
		x1, x2 := customDaoAndLoggerAttributeGenerator(executionID, routePath, req, routeParameter, username, userUUID, logEntry)
		logEntry = x1
		if x2 != nil {
			for k, v := range x2 {
				routeParameter.DatabaseReference.InstantSet(k, v)
			}
		}
	}
	logEntry = defaultLoggerPredefinedParameterFiller(executionID, routePath, req, routeParameter, username, userUUID, logEntry)
	commonParam = HTTPCommonParameter{IPAddress: req.RemoteAddr,
		RawRequest:        req,
		RequestPath:       req.URL.Path,
		DatabaseReference: routeParameter.DatabaseReference,
		UserUUID:          userUUID,
		LogEntry:          logEntry,
		Username:          username,
		PathParameters:    nil}
	if strings.Contains(routePath, "{") { // ada route parameter
		commonParam.PathParameters = mux.Vars(req)
	}
	return
}
