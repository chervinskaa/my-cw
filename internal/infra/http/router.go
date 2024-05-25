package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BohdanBoriak/boilerplate-go-back/config"
	"github.com/BohdanBoriak/boilerplate-go-back/config/container"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/controllers"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5/middleware"
)

func Router(cont container.Container) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "capacitor://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		// Health
		apiRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())
			healthRouter.Handle("/*", NotFoundJSON())
		})

		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			// Public routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, cont.AuthController, cont.AuthMw)
				})
			})

			// Protected routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Use(cont.AuthMw)

				UserRouter(apiRouter, cont.UserController)
				OrganizationRouter(apiRouter, cont.OrganizationController, cont.OrganizationService)
				RoomRouter(apiRouter, cont.RoomController, cont.RoomService)
				apiRouter.Handle("/*", NotFoundJSON())
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, config.GetConfiguration().FileStorageLocation))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	return router
}

func AuthRouter(r chi.Router, ac controllers.AuthController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			ac.Register(),
		)
		apiRouter.Post(
			"/login",
			ac.Login(),
		)
		apiRouter.With(amw).Post(
			"/logout",
			ac.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController) {
	r.Route("/users", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			uc.FindMe(),
		)
		apiRouter.Put(
			"/",
			uc.Update(),
		)
		apiRouter.Delete(
			"/",
			uc.Delete(),
		)
	})
}

func OrganizationRouter(r chi.Router, oc controllers.OrganizationController, os app.OrganizationService) {
	opom := middlewares.PathObject("orgId", controllers.OrgKey, os)
	r.Route("/organizations", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/",
			oc.Save(),
		)
		apiRouter.Get(
			"/",
			oc.FindForUser(),
		)
		apiRouter.With(opom).Get(
			"/{orgId}",
			oc.Find(),
		)
		apiRouter.With(opom).Put(
			"/{orgId}",
			oc.Update(),
		)
		apiRouter.With(opom).Delete(
			"/{orgId}",
			oc.Delete(),
		)
	})
}

func RoomRouter(r chi.Router, rc controllers.RoomController, rs app.RoomService) {
	rOpom := middlewares.PathObject("roomId", controllers.RoKey, rs)
	r.Route("/rooms", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/",
			rc.Save(),
		)
		apiRouter.Get(
			"/",
			rc.FindForOrganization(),
		)
		apiRouter.With(rOpom).Get(
			"/{roomId}",
			rc.Find(),
		)
		apiRouter.With(rOpom).Put(
			"/{roomId}",
			rc.Update(),
		)
		apiRouter.With(rOpom).Delete(
			"/{roomId}",
			rc.Delete(),
		)
	})
}

func DeviceRouter(r chi.Router, dc controllers.DeviceController, ds app.DeviceService) {
	dOpom := middlewares.PathObject("deviceId", controllers.DevKey, ds)
	r.Route("/devices", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/",
			dc.Save(),
		)
		apiRouter.Get(
			"/",
			dc.FindAll(),
		)
		apiRouter.With(dOpom).Get(
			"/{deviceId}",
			dc.Find(),
		)
		apiRouter.With(dOpom).Put(
			"/{deviceId}",
			dc.Update(),
		)
		apiRouter.With(dOpom).Delete(
			"/{deviceId}",
			dc.Delete(),
		)
		apiRouter.With(dOpom).Post(
			"/{deviceId}/install",
			dc.InstallDevice(),
		)
		apiRouter.With(dOpom).Post(
			"/{deviceId}/uninstall",
			dc.UninstallDevice(),
		)
	})
}

func NotFoundJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode("Ok")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
