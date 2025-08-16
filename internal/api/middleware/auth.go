package middleware

import (
	"context"
	"net/http"
	"strings"

	pkgerr "book-management/internal/pkg/errors"
	"book-management/internal/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ctxKey string

const ctxUserID ctxKey = "userID"

func GetUserID(r *http.Request) (string, bool) {
	v := r.Context().Value(ctxUserID)
	if v == nil {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}

func AuthJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			writeErr(w, pkgerr.E(pkgerr.ErrUnauthorized, "missing authorization header", nil))
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeErr(w, pkgerr.E(pkgerr.ErrUnauthorized, "invalid authorization header", nil))
			return
		}
		secret := viper.GetString("auth.jwtsecret")
		if secret == "" {
			writeErr(w, pkgerr.E(pkgerr.ErrInternal, "jwt secret not configured", nil))
			return
		}

		tok, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			writeErr(w, pkgerr.E(pkgerr.ErrUnauthorized, "invalid token", err))
			return
		}
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			writeErr(w, pkgerr.E(pkgerr.ErrUnauthorized, "invalid claims", nil))
			return
		}
		sub, _ := claims["sub"].(string)
		if sub == "" {
			writeErr(w, pkgerr.E(pkgerr.ErrUnauthorized, "missing sub", nil))
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeErr(w http.ResponseWriter, e *pkgerr.AppError) {
	logger.Log.Warn("request failed", zap.String("code", string(e.Code)), zap.String("msg", e.Message), zap.Error(e.Err))
	http.Error(w, `{"error":"`+e.Message+`","code":"`+string(e.Code)+`"}`, pkgerr.HTTPStatus(e.Code))
}
