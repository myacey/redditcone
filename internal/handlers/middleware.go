package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/myacey/redditclone/internal/customerror/errhandler"
	"go.uber.org/zap"
)

var UserIDCtxKeyValue = "userID"

func (h *Handler) JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AuthMiddleware(next http.Handler, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Fields(r.Header.Get("authorization"))
		logger.Infow("check auth token",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
		)

		// token should be "Bearer {token}", so check len
		if len(authHeader) != 2 {
			h.jsonError(w, errhandler.New(http.StatusUnauthorized, "invalid token length", "authHeader length is not 2", nil))
			return
		}

		logger.Infow("token",
			"authHeader", authHeader,
		)

		userID, err := h.tokenMaker.ExtractUserID(authHeader[1])
		if err != nil {
			h.jsonError(w, errhandler.New(http.StatusUnauthorized, "token verification failed", err.Error(), err))
			return
		}
		h.logger.Infow("extracted token info",
			"userID", userID,
		)

		if err = h.service.CheckUserSession(context.Background(), userID, authHeader[1]); err != nil {
			h.jsonError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDCtxKeyValue, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) LoggingMiddleware(next http.Handler, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infow("request received",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)
	})
}
