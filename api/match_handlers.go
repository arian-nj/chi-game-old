package api

import "net/http"

func (app *Application) makeMatchMakingTicket(w http.ResponseWriter, r *http.Request) {
	tgUser := ContextGetAuthenticatedUser(r)
	_ = tgUser
}
