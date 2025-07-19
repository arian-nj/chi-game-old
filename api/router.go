package api

import (
	"net/http"

	"github.com/arian-nj/chibazi/frontend"
)

func (app *Application) createRouter() *http.ServeMux {
	distFS := frontend.GetDistFS()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/status", app.statusHandler)
	mux.HandleFunc("POST /api/validate-init/", app.validateInitdata)
	mux.HandleFunc("POST /api/validate-dummy/", app.dummyValidate)

	// mux.HandleFunc("GET /api/game/xo/ws/", app.websocketUpgrader)
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServerFS(distFS)))

	return mux
}
