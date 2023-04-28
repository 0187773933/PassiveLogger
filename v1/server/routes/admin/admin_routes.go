package adminroutes

import (
	fiber "github.com/gofiber/fiber/v2"
	types "github.com/0187773933/PassiveLogger/v1/types"
)

var GlobalConfig *types.ConfigFile

var ui_html_pages = map[ string ]string {
	"/": "./v1/server/html/index.html" ,
}

func RegisterRoutes( fiber_app *fiber.App , config *types.ConfigFile ) {
	GlobalConfig = config
	admin_route_group := fiber_app.Group( "/admin" )

	// HTML UI Pages
	admin_route_group.Get( "/login" , ServeLoginPage )
	for url , _ := range ui_html_pages {
		admin_route_group.Get( url , ServeAuthenticatedPage )
	}

	// API Routes
	admin_route_group.Get( "/logout" , Logout )
	admin_route_group.Post( "/login" , HandleLogin )

}