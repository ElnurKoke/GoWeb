package handler

import (
	"log"
	"net/http"
	"strconv"

	"forum/models"
)

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
	var user1 models.User
	user := r.Context().Value(ctxUserKey).(models.User)
	if user == user1 {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodPost:
		var PostID string
		var url string
		postID1 := r.FormValue("like1")
		postID2 := r.FormValue("like2")
		if postID1 == "" {
			PostID = postID2
			url = "/post/" + PostID
		} else {
			PostID = postID1
			url = "/"
		}
		id, err := strconv.Atoi(PostID)
		if err != nil {
			h.ErrorPage(w, http.StatusNotFound, "Sorry, but no pages were found for your request =(")
			return
		}
		like := models.Like{
			UserID: user.ID,
			PostID: id,
		}
		if err = h.services.Like.SetPostLike(like); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) dislikePost(w http.ResponseWriter, r *http.Request) {
	var user1 models.User
	user := r.Context().Value(ctxUserKey).(models.User)
	if user == user1 {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodPost:
		var PostID string
		var url string
		postID1 := r.FormValue("dislike1")
		postID2 := r.FormValue("dislike2")
		if postID1 == "" {
			PostID = postID2
			url = "/post/" + PostID
		} else {
			PostID = postID1
			url = "/"
		}
		id, err := strconv.Atoi(PostID)
		if err != nil {
			log.Printf("id = %d", id)
			h.ErrorPage(w, http.StatusNotFound, "Sorry, but no pages were found for your request =(")
			return
		}
		dislike := models.DisLike{
			UserID: user.ID,
			PostID: id,
		}
		if err = h.services.Dislike.SetPostDislike(dislike); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
	var user1 models.User
	user := r.Context().Value(ctxUserKey).(models.User)
	if user == user1 {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodPost:
		commentID := r.FormValue("like")
		id, err := strconv.Atoi(commentID)
		if err != nil {
			h.ErrorPage(w, http.StatusNotFound, "Sorry, but no pages were found for your request =(")
			return
		}
		like := models.Like{
			UserID:    user.ID,
			CommentID: id,
		}
		postID, err := h.services.Like.SetCommentLike(like)
		if err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	var user1 models.User
	user := r.Context().Value(ctxUserKey).(models.User)
	if user == user1 {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodPost:
		commentID := r.FormValue("dislike")
		id, err := strconv.Atoi(commentID)
		if err != nil {
			log.Printf("id = %d", id)
			h.ErrorPage(w, http.StatusNotFound, "Sorry, but no pages were found for your request =(")
			return
		}
		dislike := models.DisLike{
			UserID:    user.ID,
			CommentID: id,
		}
		postID, err := h.services.Dislike.SetCommentDislike(dislike)
		if err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
