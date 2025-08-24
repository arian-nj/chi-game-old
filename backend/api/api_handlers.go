package api

import (
	"log/slog"
	"net/http"

	"github.com/arian-nj/chibazi/backend/pkg/response"
)

func (app *ApiApplication) statusHandler(w http.ResponseWriter, r *http.Request) {
	err := response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	if err != nil {
		slog.Error(err.Error())
	}
}
