package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arian-nj/chibazi/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

func (app *ApiApplication) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.invalidAuthenticationCreds(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) == 2 && headerParts[0] == "Bearer" {
			token := headerParts[1]
			newRequest := app.ValidateToken(w, r, token)
			if newRequest == nil {
				return
			}
			r = newRequest
		}

		next.ServeHTTP(w, r)
	})
}

func (app *ApiApplication) AuthenticateQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("auth_token")
		if token == "" {
			app.invalidAuthenticationCreds(w, r)
			return
		}

		newRequest := app.ValidateToken(w, r, token)

		if newRequest == nil {
			return
		}

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

func ContextGetAuthenticatedUser(queries *database.Queries, r *http.Request) (*database.TelegramUser, error) {
	val := r.Context().Value(authenticatedUserContextKey)
	reqConUser, ok := val.(*ReqContextUser)
	if !ok || reqConUser == nil {
		return nil, errors.New("authenticated user missing or invalid type in context")
	}

	var user database.TelegramUser
	var err error
	user, err = queries.GetTgUserByID(r.Context(), reqConUser.UserID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user, err = queries.CreateTgUser(r.Context(), reqConUser.UserID)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	return &user, nil
}

func (app *ApiApplication) ValidateToken(w http.ResponseWriter, r *http.Request, tokenString string) *http.Request {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return app.Config.Jwt.SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		slog.Error("error parsing token", "err", err)
		app.InvalidAuthenticationToken(w, r)
		return nil
	}

	expireAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		slog.Error("error getting expiration time", "err", err)
		app.ServerError(w, r, err)
		return nil
	}
	if expireAt.Time.Unix() < time.Now().Unix() {
		slog.Error("token expired", "expireAt", expireAt)
		app.InvalidAuthenticationToken(w, r)
		return nil
	}

	notBefore, err := token.Claims.GetNotBefore()
	if err != nil {
		slog.Error("error getting not before", "err", err)
		app.ServerError(w, r, err)
		return nil
	}

	if notBefore.Time.Unix() > time.Now().Unix() {
		slog.Error("token not before is in the future", "notBefore", notBefore)
		app.InvalidAuthenticationToken(w, r)
		return nil
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		slog.Error("error getting subject", "err", err)
		app.ServerError(w, r, err)
		return nil
	}

	userID, err := strconv.Atoi(sub)
	if err != nil {
		slog.Error("error converting subject to int", "err", err)
		app.invalidAuthenticationCreds(w, r)
		return nil
	}
	// if user.ID == 0 {
	// 	app.InvalidAuthenticationToken(w, r)
	// }
	return ContextSetAuthenticatedUser(r, &ReqContextUser{UserID: userID})
}
