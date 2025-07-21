package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arian-nj/chibazi/database"
	"github.com/golang-jwt/jwt/v5"
)

func (app *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")

			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				token := headerParts[1]
				new_request := app.ValidateToken(w, r, token)
				if new_request == nil {
					return
				}
				r = new_request
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) AuthenticateQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("auth_token")
		if token != "" {
			new_request := app.ValidateToken(w, r, token)
			if new_request == nil {
				return
			}
			r = new_request
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := ContextGetAuthenticatedUser(r)

		if authenticatedUser == nil {
			app.AuthenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func ContextSetAuthenticatedUser(r *http.Request, user *database.TelegramUser) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetAuthenticatedUser(r *http.Request) *database.TelegramUser {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*database.TelegramUser)
	if !ok {
		return nil
	}

	return user
}

func (app *Application) ValidateToken(w http.ResponseWriter, r *http.Request, tokenString string) *http.Request {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return app.Config.Jwt.SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		app.InvalidAuthenticationToken(w, r)
		return nil
	}

	expireAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}
	if expireAt.Time.Unix() < time.Now().Unix() {

		app.InvalidAuthenticationToken(w, r)
		return nil
	}

	notBefore, err := token.Claims.GetNotBefore()
	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}

	if notBefore.Time.Unix() > time.Now().Unix() {
		app.InvalidAuthenticationToken(w, r)
		return nil
	}

	sub, err := token.Claims.GetSubject()

	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}

	userID, err := strconv.Atoi(sub)
	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}

	user, err := app.Queries.GetTgUser(context.Background(), userID)
	if err != nil {
		app.ServerError(w, r, err)
		return nil
	}

	if user.ID == 0 {
		app.InvalidAuthenticationToken(w, r)
	}
	return ContextSetAuthenticatedUser(r, &user)

}
