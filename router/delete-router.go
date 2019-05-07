package router

import (
	"net/http"

	"github.com/SoftwareUndagi/go-common-libs/common"
	"github.com/gorilla/mux"
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
	//ProtectFromDoubleSubmit should app check for double dubmit. double submit check with header as specified on key : KeyForReqHeaderEditorToken
	ProtectFromDoubleSubmit bool
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
		localDB := routeParameter.DatabaseReference.New()
		//defer localDB.Close()
		_, commParam := generateCommonHTTPParam(routePath, req, param.RouteParameter.Clone(localDB), username, userUUID)
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
