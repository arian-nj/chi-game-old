package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/arian-nj/chibazi/database"
	"github.com/arian-nj/chibazi/pkg/request"
	"github.com/arian-nj/chibazi/pkg/response"
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

type getMeOut struct {
	ID int `json:"id"`
}

func (app *ApiApplication) getMe(w http.ResponseWriter, r *http.Request) {
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
	out := &getMeOut{
		ID: tgUser.ID,
	}

	err = response.JSON(w, http.StatusOK, out)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

func (app *ApiApplication) refreshToken(w http.ResponseWriter, r *http.Request) {
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	newToken := createToken(tgUser.ID)
	tokenString, err := newToken.SignedString(app.Config.Jwt.SecretKey)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, JWTTokenOutput{
		Token: tokenString,
	})
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

type dummyValidate struct {
	UserID int `json:"user_id"`
}

func (app *ApiApplication) dummyValidate(w http.ResponseWriter, r *http.Request) {
	input := &dummyValidate{}
	request.DecodeJSON(w, r, input)
	slog.Info("dummy auth", "userid", input.UserID)
	_, err := app.Queries.GetTgUserByID(r.Context(), input.UserID)
	if err != nil {
		app.NotFound(w, r)
		return
	}

	token := createToken(input.UserID)
	tokenString, err := token.SignedString(app.Config.Jwt.SecretKey)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, JWTTokenOutput{
		Token: tokenString,
	})
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}

func (app *ApiApplication) validateInitdata(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InitData string `json:"init_data"`
	}
	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.ServerError(w, r, err)
		slog.Error("can't decode init data input ", "err", err)
		return
	}

	parsedQuery, _ := url.ParseQuery(input.InitData)
	user, ok := ValidateWebappRequest(parsedQuery, app.Config.BotToken)
	if !ok {
		app.invalidAuthenticationCreds(w, r)
		return
	}

	botUserRow, err := app.Queries.GetTgUserBtTgID(r.Context(), int(user.ID))
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	var tgUserRow database.TelegramUser
	if botUserRow.ID == 0 {
		tgUserRow, err = app.CreateBrandNewPerson(int(user.ID))
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
	} else {
		tgUserRow, err = app.Queries.GetTgUserByID(r.Context(), int(botUserRow.ID))
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
	}

	token := createToken(tgUserRow.ID)
	tokenString, err := token.SignedString(app.Config.Jwt.SecretKey)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, JWTTokenOutput{
		Token: tokenString,
	})
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

}

func (app *ApiApplication) CreateBrandNewPerson(tgId int) (database.TelegramUser, error) {

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
