package api

import "net/http"

func (app *ApiApplication) makeMatchMakingTicket(w http.ResponseWriter, r *http.Request) {
	tgUser := ContextGetAuthenticatedUser(r)
	_ = tgUser
}
