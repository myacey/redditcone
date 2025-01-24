package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/myacey/redditclone/internal/customerror/errhandler"
	"github.com/myacey/redditclone/internal/models"
)

func (h *Handler) UnvotePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usr, err := h.extractUserFromRequestContext(r)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	postID := mux.Vars(r)["id"]

	ctx := r.Context()
	changedPost, err := h.service.UnvotePostWithID(ctx, postID, usr.ID)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledPost, err := changedPost.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal posts", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledPost)
}

func (h *Handler) VotePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usr, err := h.extractUserFromRequestContext(r)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	postID := mux.Vars(r)["id"]

	ctx := r.Context()
	newVote := models.NewVote(usr.ID, 1)
	changedPost, err := h.service.VotePostWithID(ctx, postID, newVote)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledPost, err := changedPost.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal posts", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledPost)
}

func (h *Handler) DownvotePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usr, err := h.extractUserFromRequestContext(r)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	postID := mux.Vars(r)["id"]

	ctx := r.Context()
	newVote := models.NewVote(usr.ID, -1)
	changedPost, err := h.service.VotePostWithID(ctx, postID, newVote)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledPost, err := changedPost.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal posts", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledPost)
}
