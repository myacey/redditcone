package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/myacey/redditclone/internal/customerror/errhandler"
	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/service"
	"github.com/myacey/redditclone/internal/token"
	"go.uber.org/zap"
)

type Handler struct {
	service    service.ServiceInterface
	logger     *zap.SugaredLogger
	tokenMaker token.TokenMaker
}

func NewHandler(s service.ServiceInterface, l *zap.SugaredLogger, tm token.TokenMaker) *Handler {
	return &Handler{
		service:    s,
		logger:     l,
		tokenMaker: tm,
	}
}

func (h *Handler) jsonError(w http.ResponseWriter, err error) {
	statusCode := errhandler.GetStatusCode(err)

	// check if it's a StatusCodedError to get the UserAnswer
	var scErr *errhandler.StatusCodedError
	var message string
	if errors.As(err, &scErr) {
		message = scErr.UserAnswer
		if scErr.Err != nil {
			h.logger.Errorw("internal error",
				"err", scErr.Err,
			)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	} else {
		message = err.Error()
	}

	// log about internal errors

	resp, err := json.Marshal(map[string]interface{}{
		"message": message,
	})
	if err != nil {
		h.logger.Errorw("cant marshal error to json",
			"err", err.Error(),
		)
		http.Error(w, "internal", http.StatusInternalServerError)
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(resp)
	if err != nil {
		h.logger.Errorw("cant write to response",
			"err", err.Error(),
		)
		http.Error(w, "internal", http.StatusInternalServerError)
	}
}

func (h *Handler) WriteToResponse(w http.ResponseWriter, statusCode int, marshalledText []byte) {
	w.WriteHeader(statusCode)
	_, err := w.Write(marshalledText)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "internal", "cant write to response", nil))
	}
}

func (h *Handler) extractUserFromRequestContext(r *http.Request) (*models.User, error) {
	userIDCtx := r.Context().Value(UserIDCtxKeyValue)
	if userIDCtx == nil {
		return nil, errhandler.New(http.StatusUnauthorized,
			"invalid userID",
			"user didnt passed userID field",
			nil)
	}
	userID, ok := userIDCtx.(string)
	if !ok {
		return nil, errhandler.New(http.StatusUnauthorized, "invalid userID type", "user provided invalid userID type", nil)
	}
	ctx := r.Context()
	usr, err := h.service.GetUserFromDBByID(ctx, userID)
	if err != nil {
		// выше должен вызываться h.jsonError(w, http.StatusUnauthorized, err.Error())
		return nil, errhandler.New(http.StatusUnauthorized, "invalid token", "didnt find user in DB, id: "+userID, nil)
	}

	return usr, nil
}
