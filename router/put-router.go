package router

import (
	"net/http"

	"github.com/SoftwareUndagi/go-common-libs/common"
	"github.com/gorilla/mux"
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
	//DoNotAllowCORS default false, set to true to disable CORS for this path ( post)
	DoNotAllowCORS bool
	//Handler handler put task
	Handler func(parameter HTTPPutParameter) (interface{}, common.ErrorWithCodeData)
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
