package router

import (
	"net/http"

	"github.com/SoftwareUndagi/go-common-libs/common"
	"github.com/gorilla/mux"
)

//HTTPGetParameter parameter for http( GET)
type HTTPGetParameter struct {
	HTTPCommonParameter
}

//GetRouterDefinition definisi route
type GetRouterDefinition struct {
	//RoutePath path route
	RoutePath string
	//Secured flag secured request atau tidak
	Secured bool
	//DisableGzip disable gzip response. per request. ini bisa di override pada level app
	DisableGzip bool
	//Handler request handler
	Handler GetHandlerFunction

	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
}

//GetHandlerFunction route handler function
type GetHandlerFunction func(parameter HTTPGetParameter) (result interface{}, error common.ErrorWithCodeData)

//RegisterGetJSONHandlerParam Param get
type RegisterGetJSONHandlerParam struct {
	//RouteParameter param route
	RouteParameter Parameter
	//RoutePath path route
	RoutePath string
	//Secured flag secured request atau tidak
	Secured bool
	//DisableGzip disable gzip response. per request. ini bisa di override pada level app
	DisableGzip bool
	//Handler request handler
	Handler GetHandlerFunction

	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
}

//RegisterGetJSONHandler register get handler
//
func RegisterGetJSONHandler(param RegisterGetJSONHandlerParam) *mux.Route {
	routeParameter := param.RouteParameter
	routePath := param.RoutePath
	secured := param.Secured
	handler := param.Handler
	if !param.DoNotAllowCORS {
		if CORSAllowedPaths[routePath] == nil {
			var containerBaru []string
			CORSAllowedPaths[routePath] = containerBaru
		}
		CORSAllowedPaths[routePath] = append(CORSAllowedPaths[routePath], "GET")
		appendOptionRouter(param.RouteParameter.MuxRouter, routePath)
	}
	theRoute := routeParameter.MuxRouter.HandleFunc(routePath, func(w http.ResponseWriter, req *http.Request) {
		//param.RouteParameter.DatabaseReference.Set("requestId")
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
		v := HTTPGetParameter{
			HTTPCommonParameter: commParam}
		// v.
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
	return theRoute
}
