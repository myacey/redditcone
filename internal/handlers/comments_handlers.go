package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/myacey/redditclone/internal/customerror/errhandler"
	"github.com/myacey/redditclone/internal/models"
)

type AddCommentRequest struct {
	Comment string `json:"comment"`
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postID := mux.Vars(r)["id"]
	if postID == "" {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "invalid postID", "nil postID", nil))
		return
	}

	var addCommentRequest AddCommentRequest
	err := json.NewDecoder(r.Body).Decode(&addCommentRequest)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "invalid json", "invalid addCommentRequest: "+err.Error(), nil))
		return
	}

	// get User to add Author to Comment
	usr, err := h.extractUserFromRequestContext(r)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	newComment := models.NewComment(addCommentRequest.Comment, usr, postID)

	ctx := r.Context()
	updatedPost, err := h.service.AddCommentToPost(ctx, postID, *newComment)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "invalid args", "cant add comment to post: "+err.Error(), nil))
		return
	}

	marshalledPost, err := updatedPost.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "internal", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusCreated, marshalledPost)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDCtx := r.Context().Value(UserIDCtxKeyValue)
	if userIDCtx == nil {
		h.jsonError(w, errhandler.New(http.StatusUnauthorized, "nil userID", "userID not found in context", nil))
		return
	}
	_, ok := userIDCtx.(string)
	if !ok {
		h.jsonError(w, errhandler.New(http.StatusUnauthorized, "invalid user ID type", "userID is not a string", nil))
		return
	}

	postID := mux.Vars(r)["postID"]
	if postID == "" {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "invalid postID", "postID is empty", nil))
		return
	}

	commentID := mux.Vars(r)["commentID"]
	if commentID == "" {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "invalid commentID", "commentID is empty", nil))
		return
	}

	ctx := r.Context()
	updatedPost, err := h.service.RemoveComment(ctx, postID, commentID)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledPost, err := updatedPost.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "internal error", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledPost)
}
