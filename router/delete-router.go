package router

import (
	"net/http"

	"github.com/SoftwareUndagi/go-common-libs/common"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//HTTPDeleteParameter parameter for delete
type HTTPDeleteParameter struct {
	HTTPCommonParameter
}

//DeleteRouterDefinition handler http delete request
type DeleteRouterDefinition struct {
	//RoutePath path route
	RoutePath string
	//Secured flag secured request atau tidak
	Secured bool
	//DisableGzip disable gzip response. per request. ini bisa di override pada level app
	DisableGzip bool
	//Handler request handler
	Handler DelHandlerFunction

	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
	//CheckForDoubleSubmit check for request that double submit. scan by editor token data
	CheckForDoubleSubmit *NeedDoubleSubmitProtectionDefinition
}

//DelHandlerFunction handler http delete request
type DelHandlerFunction func(parameter HTTPDeleteParameter) (result interface{}, error common.ErrorWithCodeData)

//RegisterDeleteJSONHandlerParam parameter JSON delete
type RegisterDeleteJSONHandlerParam struct {
	//route data
	RouteParameter Parameter
	//RoutePath path of route
	RoutePath string
	//DisableGzip disable gzip response. per request. ini bisa di override pada level app
	DisableGzip bool
	//Secured secured flag. if no user then request is rejected
	Secured bool
	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
	//Handler handler post task
	Handler DelHandlerFunction
	///CheckForDoubleSubmit check for request that double submit. scan by editor token data
	CheckForDoubleSubmit *NeedDoubleSubmitProtectionDefinition
}

//RegisterDeleteJSONHandler register method delete
func RegisterDeleteJSONHandler(param RegisterDeleteJSONHandlerParam) *mux.Route {
	routeParameter := param.RouteParameter
	routePath := param.RoutePath
	secured := param.Secured
	handler := param.Handler
	if !param.DoNotAllowCORS {
		if CORSAllowedPaths[routePath] == nil {
			var containerBaru []string
			CORSAllowedPaths[routePath] = containerBaru
		}
		CORSAllowedPaths[routePath] = append(CORSAllowedPaths[routePath], "DELETE")
		appendOptionRouter(param.RouteParameter.MuxRouter, routePath)
	}

	theRoute := routeParameter.MuxRouter.HandleFunc(routePath, func(w http.ResponseWriter, req *http.Request) {
		executionID := generateSimpleRequestExecutionID(req)
		logEntry := logrus.WithField("executionId", executionID)

		if param.DoNotAllowCORS {
			if checkForCors(routePath, false, w, req) {
				return
			}
		}
		var routePathActual = req.URL.Path
		var userUUID, username string
		if secured {
			if routeParameter.LoginInformationProvider == nil {
				sendErrorResponse(w, resultJSONErrorWrapper{ErrorCode: "APP_CONFIG_ERROR", ErrorMessage: "Application configuration not ok. login information provider is missing"})
				return
			}
			user, errLogin := routeParameter.LoginInformationProvider(routeParameter.DatabaseReference, req) //.GetUserLoginInformation(routeParameter.DatabaseReference, req)
			if errLogin != nil {
				if secured {
					accessDeniedHandler(routePathActual, w, req)
					return
				}
			}
			userUUID = user.UserUUID
			username = user.Username
		}
		if param.CheckForDoubleSubmit != nil {
			var checkerToken EditorTokenValidationCheckerFunction
			if param.CheckForDoubleSubmit.CustomChecker != nil {
				checkerToken = *param.CheckForDoubleSubmit.CustomChecker
			} else {
				checkerToken = DefaultEditorTokenValidationChecker
			}
			editorToken := DefaultGetterEditorToken(logEntry, req)
			ok, errCode, errChk := checkerToken(editorToken, username, param.CheckForDoubleSubmit.ObjectName, param.CheckForDoubleSubmit.BusinessObjectName, param.RouteParameter.DatabaseReference, logEntry, req)
			if !ok {
				sendErrorResponse(w, resultJSONErrorWrapper{ErrorCode: errCode, ErrorMessage: errChk.Error()})
				return
			}
		}

		localDB := routeParameter.DatabaseReference.New()
		//defer localDB.Close()
		_, commParam := generateCommonHTTPParam(executionID, logEntry, routePath, req, param.RouteParameter.Clone(localDB), username, userUUID)
		v := HTTPDeleteParameter{HTTPCommonParameter: commParam}
		if flushLogDataCommand != nil {
			defer flushLogDataCommand()
		}
		getRslt, errGet := handler(v)
		if errGet == nil {
			sendOkData(getRslt, param.DisableGzip, w, req)
			return
		}
		errorCode := errGet.GetErrorCode()
		errorMessg := errGet.Error()
		sendErrorResponse(w, resultJSONErrorWrapper{ErrorCode: errorCode, ErrorMessage: errorMessg})
	})
	theRoute.Methods("DELETE")
	return theRoute
}
