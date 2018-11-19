package router

import (
	"net/http"

	"github.com/chinx/cobweb"
	"github.com/chinx/wepub/handler"
)

func InitRouter() (http.Handler, error) {
	mux := cobweb.New()
	mux.Post("/v1/pages", handler.CreatePage)

	mux.Group("/editor", func() {
		mux.Get("/v1/actions", handler.GetAction)
		mux.Post("/v1/actions", handler.PostAction)
		mux.Get("/", handler.EditHandler)
		mux.Get("/*filename", handler.StaticHandler)
	})

	mux.Group("/v1/admin", func() {
		mux.Post("/login", handler.AdminLogin)
		mux.Post("/signup", handler.RegistryHandler)
		mux.Post("/logout", handler.UserLogout)
	})

	return mux.Build()
}
