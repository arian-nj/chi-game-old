package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/arian-nj/chibazi/backend/database"
	authv1 "github.com/arian-nj/chibazi/backend/gen/auth/v1"
	"github.com/golang-jwt/jwt/v5"
)

const JWTExpiryDuration = 1 * time.Hour

type JWTTokenOutput struct {
	Token string `json:"token"`
}

func createToken(userId int) *jwt.Token {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userId),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTExpiryDuration)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}

func (app *ApiApplication) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return app.Config.Jwt.SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		slog.Error("error parsing token", "err", err)
		return 0, err
	}

	expireAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		slog.Error("error getting expiration time", "err", err)
		return 0, err
	}
	if expireAt.Time.Unix() < time.Now().Unix() {
		slog.Error("token expired", "expireAt", expireAt)
		return 0, err
	}

	notBefore, err := token.Claims.GetNotBefore()
	if err != nil {
		slog.Error("error getting not before", "err", err)
		return 0, err
	}

	if notBefore.Time.Unix() > time.Now().Unix() {
		slog.Error("token not before is in the future", "notBefore", notBefore)
		return 0, err
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		slog.Error("error getting subject", "err", err)
		return 0, err
	}

	userID, err := strconv.Atoi(sub)
	if err != nil {
		slog.Error("error converting subject to int", "err", err)
		return 0, err
	}
	// if user.ID == 0 {
	// 	app.InvalidAuthenticationToken(w, r)
	// }

	return userID, nil
}

// func (app *ApiApplication) refreshToken(w http.ResponseWriter, r *http.Request) {
// 	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
// 	if err != nil {
// 		app.ServerError(w, r, err)
// 		return
// 	}
//
// 	newToken := createToken(tgUser.ID)
// 	tokenString, err := newToken.SignedString(app.Config.Jwt.SecretKey)
// 	if err != nil {
// 		app.ServerError(w, r, err)
// 		return
// 	}
//
// 	err = response.JSON(w, http.StatusOK, JWTTokenOutput{
// 		Token: tokenString,
// 	})
// 	if err != nil {
// 		app.ServerError(w, r, err)
// 		return
// 	}
// }

func (app *ApiApplication) DummyValidate(
	ctx context.Context,
	req *connect.Request[authv1.DummyValidateRequest],
) (*connect.Response[authv1.DummyValidateResponse], error) {
	_, err := app.Queries.GetTgUserByID(ctx, int(req.Msg.GetId()))
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not user found"))
	}

	token := createToken(int(req.Msg.Id))
	tokenString, err := token.SignedString(app.Config.Jwt.SecretKey)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, errors.New("internal"))
	}
	return connect.NewResponse(&authv1.DummyValidateResponse{
		Token: tokenString,
	}), nil
}

func (app *ApiApplication) ValidateTelegramInitData(
	ctx context.Context,
	req *connect.Request[authv1.ValidateTelegramInitDataRequest],
) (*connect.Response[authv1.ValidateTelegramInitDataResponse], error) {
	parsedQuery, _ := url.ParseQuery(req.Msg.InitData)
	user, ok := ValidateWebappRequest(parsedQuery, app.Config.BotToken)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("can't validate init data"))
	}

	botUserRow, err := app.Queries.GetTgUserByTgID(ctx, int(user.ID))
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, errors.New("internal"))
	}

	var tgUserRow database.Person
	if botUserRow.ID == 0 {
		tgUserRow, err = app.CreateBrandNewPerson(int(user.ID))
		if err != nil {
			return nil, connect.NewError(connect.CodeUnknown, errors.New("internal"))
		}
	} else {
		tgUserRow, err = app.Queries.GetTgUserByID(ctx, int(botUserRow.ID))
		if err != nil {
			return nil, connect.NewError(connect.CodeUnknown, errors.New("internal"))
		}
	}

	token := createToken(tgUserRow.ID)
	tokenString, err := token.SignedString(app.Config.Jwt.SecretKey)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, errors.New("internal"))
	}

	return connect.NewResponse(&authv1.ValidateTelegramInitDataResponse{
		Token: tokenString,
	}), nil
}

func (app *ApiApplication) CreateBrandNewPerson(tgId int) (database.Person, error) {

	tgUserRow, err := app.Queries.CreateTgUser(context.Background(), tgId)

	if err != nil {
		return tgUserRow, err
	}
	// _, err = app.Queries.InsertUserStatistic(context.Background(), personRow.ID)
	return tgUserRow, err
}

type WebAppUser struct {
	ID                    int64  `json:"id"`
	IsBot                 bool   `json:"is_bot"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Username              string `json:"username"`
	LanguageCode          string `json:"language_code"`
	IsPremium             bool   `json:"is_premium"`
	AddedToAttachmentMenu bool   `json:"added_to_attachment_menu"`
	AllowsWriteToPM       bool   `json:"allows_write_to_pm"`
	PhotoURL              string `json:"photo_url"`
}

func ValidateWebappRequest(values url.Values, token string) (user *WebAppUser, ok bool) {
	h := values.Get("hash")
	values.Del("hash")

	var vals []string

	var u WebAppUser

	for k, v := range values {
		vv, _ := url.QueryUnescape(v[0])
		vals = append(vals, k+"="+vv)
		if k == "user" {
			errDecodeUser := json.Unmarshal([]byte(vv), &u)
			if errDecodeUser != nil {
				return nil, false
			}
		}
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})

	hmac1 := hmac.New(sha256.New, []byte("WebAppData"))
	hmac1.Write([]byte(token))
	r1 := hmac1.Sum(nil)

	data := []byte(strings.Join(vals, "\n"))

	hmac2 := hmac.New(sha256.New, r1)
	hmac2.Write(data)
	r2 := hmac2.Sum(nil)

	if h != fmt.Sprintf("%x", r2) {
		return nil, false
	}

	return &u, true
}
