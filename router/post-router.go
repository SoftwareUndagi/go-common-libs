package router

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/SoftwareUndagi/go-common-libs/common"
	"github.com/gorilla/mux"
)

//PostHandlerFunction handler http post request
type PostHandlerFunction func(parameter HTTPPostParameter) (result interface{}, error common.ErrorWithCodeData)

//PostRouterDefinition definisi post router
type PostRouterDefinition struct {
	//RoutePath path route
	RoutePath string
	//Secured flag secured request atau tidak
	Secured bool
	//DisableGzip disable gzip response. per request. ini bisa di override pada level app
	DisableGzip bool
	//Handler request handler
	Handler PostHandlerFunction

	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
	//CheckForDoubleSubmit check for request that double submit. scan by editor token data
	CheckForDoubleSubmit *NeedDoubleSubmitProtectionDefinition
}

//HTTPPostParameter parameter for post method
type HTTPPostParameter struct {
	HTTPCommonParameter
}

//RegisterPostJSONHandlerParameter parameter post
type RegisterPostJSONHandlerParameter struct {
	//route data
	RouteParameter Parameter
	//RoutePath path of route
	RoutePath string
	//Secured secured flag. if no user then request is rejected
	Secured bool
	//DisableGzip disable gzip response. per request. ini bisa di override pada level app
	DisableGzip bool
	//Handler handler post task
	Handler PostHandlerFunction

	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
	//CheckForDoubleSubmit check for request that double submit. scan by editor token data
	CheckForDoubleSubmit *NeedDoubleSubmitProtectionDefinition
}

//RegisterPostJSONHandler register post route handler
//routePath = path of route to serve
func RegisterPostJSONHandler(param RegisterPostJSONHandlerParameter) *mux.Route {
	routeParameter := param.RouteParameter
	routePath := param.RoutePath
	secured := param.Secured
	handler := param.Handler
	// register for cors

	if !param.DoNotAllowCORS {
		if CORSAllowedPaths[routePath] == nil {
			var containerBaru []string
			CORSAllowedPaths[routePath] = containerBaru
		}
		CORSAllowedPaths[routePath] = append(CORSAllowedPaths[routePath], "POST")
		appendOptionRouter(param.RouteParameter.MuxRouter, routePath)
	}

	theRoute := routeParameter.MuxRouter.HandleFunc(routePath, func(w http.ResponseWriter, req *http.Request) {
		var routePathActual = req.URL.Path
		var userUUID, username string
		executionID := generateSimpleRequestExecutionID(req)
		logEntry := logrus.WithField("executionId", executionID)
		if !param.DoNotAllowCORS {
			if !checkForCors(routePath, false, w, req) {
				return
			}
		}

		if secured {
			if routeParameter.LoginInformationProvider == nil {
				sendErrorResponseWithStatusCode(w, resultJSONErrorWrapper{ErrorCode: "APP_CONFIG_ERROR", ErrorMessage: "Application configuration not ok. login information provider is missing"}, http.StatusInternalServerError)
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
		v := HTTPPostParameter{HTTPCommonParameter: commParam}
		if flushLogDataCommand != nil {
			defer flushLogDataCommand()
		}
		// v.
		getRslt, errGet := handler(v)
		if errGet == nil {
			sendOkData(getRslt, param.DisableGzip, w, req)
			return
		}
		errorCode := errGet.GetErrorCode()
		errorMessg := errGet.Error()
		sendErrorResponse(w, resultJSONErrorWrapper{ErrorCode: errorCode, ErrorMessage: errorMessg})
	})
	theRoute.Methods("POST")
	return theRoute
}
