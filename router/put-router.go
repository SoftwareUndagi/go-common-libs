package router

import (
	"net/http"

	"github.com/SoftwareUndagi/go-common-libs/common"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//PutHandlerFunction handler http put function
type PutHandlerFunction func(parameter HTTPPutParameter) (result interface{}, error common.ErrorWithCodeData)

//HTTPPutParameter parameter for put request
type HTTPPutParameter struct {
	HTTPCommonParameter
}

//PutRouterDefinition definisi http put method
type PutRouterDefinition struct {
	//RoutePath path route
	RoutePath string
	//Secured flag secured request atau tidak
	Secured bool
	//DisableGzip disable gzip response. per request. ini bisa di override pada level app
	DisableGzip bool
	//Handler request handler
	Handler PutHandlerFunction

	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
	//CheckForDoubleSubmit check for request that double submit. scan by editor token data
	CheckForDoubleSubmit *NeedDoubleSubmitProtectionDefinition
}

//RegisterPutJSONHandlerParam parameter put
type RegisterPutJSONHandlerParam struct {
	//route data
	RouteParameter Parameter
	//RoutePath path of route
	RoutePath string
	//Secured secured flag. if no user then request is rejected
	Secured bool
	//DisableGzip disable gzip response. per request. ini bisa di override pada level app
	DisableGzip bool
	//Handler handler put task
	Handler func(parameter HTTPPutParameter) (interface{}, common.ErrorWithCodeData)
	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
	//CheckForDoubleSubmit check for request that double submit. scan by editor token data
	CheckForDoubleSubmit *NeedDoubleSubmitProtectionDefinition
}

//RegisterPutJSONHandler register put http handler
func RegisterPutJSONHandler(param RegisterPutJSONHandlerParam) *mux.Route {
	routeParameter := param.RouteParameter
	routePath := param.RoutePath
	secured := param.Secured
	handler := param.Handler
	if !param.DoNotAllowCORS {
		if CORSAllowedPaths[routePath] == nil {
			var containerBaru []string
			CORSAllowedPaths[routePath] = containerBaru
		}
		CORSAllowedPaths[routePath] = append(CORSAllowedPaths[routePath], "PUT")
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
			user, errLogin := routeParameter.LoginInformationProvider(routeParameter.DatabaseReference, logEntry, req) //.GetUserLoginInformation(routeParameter.DatabaseReference, req)
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
				sendErrorResponseWithStatusCode(w, resultJSONErrorWrapper{ErrorCode: errCode, ErrorMessage: errChk.Error()}, http.StatusPreconditionFailed)
				return
			}
		}

		localDB := routeParameter.DatabaseReference.New()
		//defer localDB.Close()
		_, commParam := generateCommonHTTPParam(executionID, logEntry, routePath, req, param.RouteParameter.Clone(localDB), username, userUUID)
		v := HTTPPutParameter{HTTPCommonParameter: commParam}
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
	theRoute.Methods("PUT")
	return theRoute
}
