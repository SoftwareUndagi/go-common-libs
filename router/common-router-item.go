package router

//HandlerManagerFunction producer route handler. untuk register route
type HandlerManagerFunction func() (getHandlers []GetRouterDefinition, postHandlers []PostRouterDefinition, putHandlers []PutRouterDefinition, delHandlers []DeleteRouterDefinition)

//RegisterRouteHandlers register semua handlers pada maanager
func RegisterRouteHandlers(routeParameter Parameter, routeManager HandlerManagerFunction) {
	getHandlers, postHandlers, putHandlers, delHandlers := routeManager()
	for _, handler := range getHandlers {
		p := RegisterGetJSONHandlerParam{
			DoNotAllowCORS: handler.DoNotAllowCORS,
			Handler:        handler.Handler,
			RouteParameter: routeParameter,
			RoutePath:      handler.RoutePath,
			Secured:        handler.Secured}
		RegisterGetJSONHandler(p)
	}
	for _, postHandler := range postHandlers {
		p := RegisterPostJSONHandlerParameter{
			DoNotAllowCORS: postHandler.DoNotAllowCORS,
			Handler:        postHandler.Handler,
			RouteParameter: routeParameter,
			RoutePath:      postHandler.RoutePath,
			Secured:        postHandler.Secured}
		RegisterPostJSONHandler(p)
	}
	for _, delHandler := range delHandlers {
		p := RegisterDeleteJSONHandlerParam{
			DoNotAllowCORS: delHandler.DoNotAllowCORS,
			Handler:        delHandler.Handler,
			RouteParameter: routeParameter,
			RoutePath:      delHandler.RoutePath,
			Secured:        delHandler.Secured}
		RegisterDeleteJSONHandler(p)
	}
	for _, putHandler := range putHandlers {
		p := RegisterPutJSONHandlerParam{
			DoNotAllowCORS: putHandler.DoNotAllowCORS,
			Handler:        putHandler.Handler,
			RouteParameter: routeParameter,
			RoutePath:      putHandler.RoutePath,
			Secured:        putHandler.Secured}
		RegisterPutJSONHandler(p)
	}
}
