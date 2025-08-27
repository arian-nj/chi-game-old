package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/arian-nj/chibazi/backend/database"
)

func (app *ApiApplication) AuthenticateHeader(ctx context.Context, header http.Header) *database.Person {
	header.Add("Vary", "Authorization")

	authorizationHeader := header.Get("Authorization")
	if authorizationHeader == "" {
		return nil
	}

	headerParts := strings.Split(authorizationHeader, " ")

	if (len(headerParts) == 2 && headerParts[0] == "Bearer") == false {
		return nil
	}
	token := headerParts[1]
	userID, err := app.ValidateToken(token)
	if err != nil {
		slog.Error("authorize header ", "error", err)
		return nil
	}

	person, err := app.Queries.GetTgUserByID(ctx, userID)

	if err != nil {
		slog.Error("no user found auth header middleware", "err", err)
		return nil
	}
	return &person
}

//	func (app *ApiApplication) Authenticate(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.Header().Add("Vary", "Authorization")
//
//			authorizationHeader := r.Header.Get("Authorization")
//			if authorizationHeader == "" {
//				app.invalidAuthenticationCreds(w, r)
//				return
//			}
//
//			headerParts := strings.Split(authorizationHeader, " ")
//
//			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
//				token := headerParts[1]
//				newRequest := app.ValidateToken(w, r, token)
//				if newRequest == nil {
//					return
//				}
//				r = newRequest
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
func (app *ApiApplication) AuthenticateQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("auth_token")
		if token == "" {
			app.invalidAuthenticationCreds(w, r)
			return
		}
		userID, err := app.ValidateToken(token)
		if err != nil {
			app.invalidAuthenticationCreds(w, r)
			return
		}

		reqUser := &ReqContextUser{UserID: userID}
		newRequest := ContextSetAuthenticatedUser(r, reqUser)

		next.ServeHTTP(w, newRequest)
	})
}

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

type ReqContextUser struct {
	UserID int
}

func ContextSetAuthenticatedUser(r *http.Request, user *ReqContextUser) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetAuthenticatedUser(queries *database.Queries, r *http.Request) (*database.Person, error) {
	val := r.Context().Value(authenticatedUserContextKey)
	reqConUser, ok := val.(*ReqContextUser)
	if !ok || reqConUser == nil {
		return nil, errors.New("authenticated user missing or invalid type in context")
	}

	var user database.Person
	var err error
	user, err = queries.GetTgUserByID(r.Context(), reqConUser.UserID)

	if err != nil {
		slog.Error("no user found in context get authticated user", "err", err)
		return nil, err
	}

	return &user, nil
}
