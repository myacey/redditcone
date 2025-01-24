package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/myacey/redditclone/internal/customerror/errhandler"
	"github.com/myacey/redditclone/internal/models"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {

		h.jsonError(w, errhandler.New(http.StatusBadRequest, "bad json", "failed to decode request body: "+err.Error(), nil))
		return
	}

	usr := models.NewUser(registerRequest.Username, registerRequest.Password)

	ctx := r.Context()
	session, err := h.service.CreateNewUser(ctx, usr)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledSession, err := session.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal session", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusCreated, marshalledSession)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "bad json", "failed to decode request body: "+err.Error(), nil))
		return
	}

	ctx := r.Context()
	session, err := h.service.LoginUser(ctx, loginRequest.Username)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledSession, err := session.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal session", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledSession)
}
