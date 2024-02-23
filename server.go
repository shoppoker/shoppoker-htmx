package main

import (
	"fmt"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/w1png/go-htmx-ecommerce-template/config"
	"github.com/w1png/go-htmx-ecommerce-template/handlers"
	admin_handlers "github.com/w1png/go-htmx-ecommerce-template/handlers/admin"
	user_handlers "github.com/w1png/go-htmx-ecommerce-template/handlers/user"
	"github.com/w1png/go-htmx-ecommerce-template/middleware"
)

type HTTPServer struct {
	echo *echo.Echo
}

func NewHTTPServer() *HTTPServer {
	server := &HTTPServer{
		echo: echo.New(),
	}

	user_page_group := server.echo
	user_api_group := server.echo.Group("/api")

	admin_page_group := server.echo.Group("/admin")
	admin_api_group := admin_page_group.Group("/api")
	admin_page_group.Use(middleware.UseAdmin)

	server.echo.Use(echoMiddleware.Logger())
	server.echo.Use(echoMiddleware.Recover())
	server.echo.Use(middleware.UseAuth)
	server.echo.Use(middleware.UseUrl)
	server.echo.Use(middleware.UseCategories)
	server.echo.Use(middleware.UseCart)

	server.echo.Use(echoprometheus.NewMiddleware("shoppoker"))
	server.echo.GET("/metrics", echoprometheus.NewHandler(), middleware.UseAdmin)
	server.echo.Static("/static", "static")

	gather_funcs := []func(*echo.Echo, *echo.Group, *echo.Group, *echo.Group){
		user_handlers.GatherIndexHandlers,
		admin_handlers.GatherIndexRoutes,
		user_handlers.GatherLoginRoutes,
		admin_handlers.GatherUsersRoutes,
		user_handlers.GatherProductsRoutes,
		admin_handlers.GatherProductsRoutes,
		user_handlers.GatherCategoriesRoutes,
		admin_handlers.GatherCategoriesRoutes,
		user_handlers.GatherOrdersRoutes,
		// admin_handlers.GatherOrdersRoutes,
		user_handlers.GatherCartRoutes,
		handlers.GatherFilesHandler,
		admin_handlers.GatherSettingsRoutes,
		admin_handlers.GatherExamplesRoutes,
		admin_handlers.GatherCustomChipsRoutes,
	}

	for _, f := range gather_funcs {
		f(user_page_group, user_api_group, admin_page_group, admin_api_group)
	}

	return server
}

func (s *HTTPServer) Run() error {
	return s.echo.Start(fmt.Sprintf(":%s", config.ConfigInstance.Port))
}
