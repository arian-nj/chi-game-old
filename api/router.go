package api

import (
	"log"
	"net/http"

	"github.com/arian-nj/chibazi/frontend"
	"github.com/arian-nj/chibazi/internals/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var COSS_PATTERNS = []string{"http://localhost:*", "https://localhost:*", "localhost:*"}

func (app *ApiApplication) createRouter() *chi.Mux {
	distFS := frontend.GetDistFS()
	mux := chi.NewMux()

	if app.Config.ReleaseMode == config.Develop {
		mux.Use(cors.Handler(cors.Options{
			// AllowedOrigins: []string{"https://*", "http://*"},
			AllowedOrigins: COSS_PATTERNS,
		}))
	}

	mux.Get("/api/status", app.statusHandler)
	mux.Post("/api/auth/validate-init/", app.validateInitdata)
	if app.Config.ReleaseMode == config.Develop {
		mux.Post("/api/auth/validate-dummy/", app.dummyValidate)
	}
	// mux.Get("/api/auth/me/", app.getMe)

	mux.Handle("/web/*", http.StripPrefix("/web/", http.FileServerFS(distFS)))

	mux.Group(func(authRouter chi.Router) {
		authRouter.Use(app.AuthenticateQuery)

		authRouter.Get("/api/game_session/", app.gameSessionWS)
		authRouter.Get("/api/match_making/ticket/", app.makeMatchMakingTicket)
	})
	return mux
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URI:", r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
