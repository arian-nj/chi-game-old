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

const JWTExpiryDuration = 10 * time.Minute

type JWTTokenOutput struct {
	Token string `json:"token"`
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
