package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/myacey/redditclone/internal/customerror/errhandler"
	"github.com/myacey/redditclone/internal/models"
)

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	posts, err := h.service.GetAllPosts(ctx)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	data, err := json.Marshal(posts)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal posts", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, data)
}

type AddPostRequest struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Type     string `json:"type"`

	Text string `json:"text,omitempty"`
	URL  string `json:"url,omitempty"`
}

func (h *Handler) AddPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usr, err := h.extractUserFromRequestContext(r)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	var addPostRequest AddPostRequest
	err = json.NewDecoder(r.Body).Decode(&addPostRequest)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "bad json", "failed to decode request body: "+err.Error(), nil))
		return
	}

	post := models.NewPost(
		usr,
		addPostRequest.Category,
		addPostRequest.Title,
		addPostRequest.Type,
		addPostRequest.Text,
		addPostRequest.URL,
	)

	ctx := r.Context()
	err = h.service.AddPost(ctx, post)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "failed to add post", err.Error(), nil))
		return
	}

	marshalledPost, err := post.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal post", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusCreated, marshalledPost)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postID := mux.Vars(r)["id"]

	ctx := r.Context()
	post, err := h.service.GetPostByID(ctx, postID, true)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledPost, err := post.GetMarshal()
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal post", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledPost)
}

func (h *Handler) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	category := mux.Vars(r)["category"]

	ctx := r.Context()
	posts, err := h.service.GetPostsByCategory(ctx, category)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledPosts, err := json.Marshal(posts)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal posts", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledPosts)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postID := mux.Vars(r)["id"]
	ctx := r.Context()
	err := h.service.DeletePostWithID(ctx, postID)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	ans := map[string]string{"message": "success"}
	marshalledAns, err := json.Marshal(ans)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal response", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledAns)
}

func (h *Handler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	username := mux.Vars(r)["username"]
	if username == "" {
		h.jsonError(w, errhandler.New(http.StatusBadRequest, "invalid username", "username is empty", nil))
		return
	}

	ctx := r.Context()
	posts, err := h.service.GetPostsByAuthor(ctx, username)
	if err != nil {
		h.jsonError(w, err)
		return
	}

	marshalledPosts, err := json.Marshal(posts)
	if err != nil {
		h.jsonError(w, errhandler.New(http.StatusInternalServerError, "failed to marshal posts", err.Error(), err))
		return
	}

	h.WriteToResponse(w, http.StatusOK, marshalledPosts)
}
