package api

import (
	"log"
	"net/http"

	"github.com/arian-nj/chibazi/frontend"
	"github.com/go-chi/chi/v5"
)

func (app *Application) createRouter() *chi.Mux {
	distFS := frontend.GetDistFS()
	mux := chi.NewMux()

	mux.Get("/api/status", app.statusHandler)
	mux.Post("/api/validate-init/", app.validateInitdata)
	mux.Post("/api/validate-dummy/", app.dummyValidate)

	mux.Get("/api/game/xo/ws/", app.websocketUpgrader)
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServerFS(distFS)))

	mux.Group(func(authRouter chi.Router) {
		authRouter.Use(app.Authenticate)
		authRouter.Use(app.RequireAuthenticatedUser)

		mux.HandleFunc("GET /api/game/match_making/ticket", app.makeMatchMakingTicket)
	})
	return mux
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URI:", r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
