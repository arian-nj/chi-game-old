package api

import (
	"fmt"
	"log"
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/arian-nj/chibazi/backend/gen/account/v1/accountv1connect"
	"github.com/arian-nj/chibazi/backend/gen/auth/v1/authv1connect"
	"github.com/arian-nj/chibazi/backend/gen/session/v1/sessionv1connect"
	"github.com/rs/cors"
)

var CORS_PATTERNS = []string{"http://localhost:5173", "https://localhost:5173", "localhost:5173"}

func (app *ApiApplication) createRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// if app.Config.ReleaseMode == config.Develop {
	// 	mux.Use(func(next http.Handler) http.Handler {
	// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			log.Printf("Request Received: Method=%s, Path=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)
	// 			next.ServeHTTP(w, r)
	// 		})
	// 	})
	// }

	sessionPath, sessionHandler := sessionv1connect.NewSessionServiceHandler(app)
	mux.Handle(sessionPath, withCORS(sessionHandler))

	authPath, authHandler := authv1connect.NewAuthServiceHandler(app)
	mux.Handle(authPath, withCORS(authHandler))
	fmt.Println("Auth path:", authPath)

	accountPath, accountHandler := accountv1connect.NewAccountServiceHandler(app)
	mux.Handle(accountPath, withCORS(accountHandler))

	mux.Handle("/api/session/", withCORS(app.AuthenticateQuery(http.HandlerFunc(app.gameSessionWS))))
	mux.Handle("/api/match_making/ticket/", withCORS(app.AuthenticateQuery(http.HandlerFunc(app.makeMatchMakingTicketWS))))
	return mux
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URI:", r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
func withCORS(next http.Handler) http.Handler {

	allowdHeaders := []string{"Accept", "Authorization", "Content-Type"}
	allowdHeaders = append(allowdHeaders, connectcors.AllowedHeaders()...)

	c := cors.New(cors.Options{
		AllowedOrigins: CORS_PATTERNS,
		// AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", , connectcors.AllowedMethods(),
		// AllowedHeaders:   []string{},
		// ExposedHeaders:   []string{"Link"},
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   allowdHeaders,
		ExposedHeaders:   connectcors.ExposedHeaders(),
		MaxAge:           7200,
		AllowCredentials: true,
	})

	return c.Handler(next)

}
