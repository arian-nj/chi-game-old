package api

import (
	"log"
	"net/http"

	"github.com/arian-nj/chibazi/backend/internals/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var CORS_PATTERNS = []string{"http://localhost:5173", "https://localhost:5173", "localhost:5173"}

func (app *ApiApplication) createRouter() *chi.Mux {
	mux := chi.NewMux()

	if app.Config.ReleaseMode == config.Develop {
		mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   CORS_PATTERNS,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
		}))
	}

	mux.Get("/api/status", app.statusHandler)
	mux.Post("/api/auth/validate-init/", app.validateInitdata)
	if app.Config.ReleaseMode == config.Develop {
		mux.Post("/api/auth/validate-dummy/", app.dummyValidate)
	}
	// mux.Get("/api/auth/me/", app.getMe)

	mux.Group(func(authHRouter chi.Router) {
		authHRouter.Use(app.Authenticate)

		authHRouter.Get("/api/auth/me", app.getMe)
		authHRouter.Get("/api/session/chat/history", app.getChatHistoryHandler)
	})
	mux.Group(func(authQRouter chi.Router) {
		authQRouter.Use(app.AuthenticateQuery)

		authQRouter.Get("/api/session/", app.gameSessionWS)
		authQRouter.Get("/api/match_making/ticket/", app.makeMatchMakingTicket)

	})
	return mux
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URI:", r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
