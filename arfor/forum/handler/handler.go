package handler

import (
	"foruum/models"
	"foruum/service"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.MiddleWare(h.home))
	mux.HandleFunc("/signup", h.MiddleWare(h.signUp))
	mux.HandleFunc("/signin", h.MiddleWare(h.signIn))
	mux.HandleFunc("/createpost", h.MiddleWare(h.createpost))
	mux.HandleFunc("/logout", h.MiddleWare(h.logout))
	mux.HandleFunc("/post/", h.MiddleWare(h.post))
	// mux.HandleFunc("/likePost", h.MiddleWare(h.likePost))
	mux.HandleFunc("/post/change/", h.MiddleWare(h.ChangePost))
	mux.HandleFunc("/post/delete/", h.MiddleWare(h.DeletePost))
	mux.HandleFunc("/myPosts", h.MiddleWare(h.myPosts))
	mux.HandleFunc("/myCommentPosts", h.MiddleWare(h.myCommentPosts))
	mux.HandleFunc("/myLikedPosts", h.MiddleWare(h.myLikedPosts))
	mux.Handle("/ui/css/", http.StripPrefix("/ui/css/", http.FileServer(http.Dir("./ui/css/"))))
	return mux
}

type Display struct {
	Username string
	Posts    []models.Post
	Category []string
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxUserKey).(models.User)
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/" {
			ErrorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var posts []models.Post
		posts, err := h.service.Post.ShowAllPosts()
		if err != nil {
			log.Print(err)
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/homepage.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, display)
	case http.MethodPost:
		var category []string
		if r.FormValue("category"+string('5')) != "" {
			var posts []models.Post
			posts, err := h.service.Post.ShowAllPosts()
			if err != nil {
				log.Print(err)
			}
			display := Display{
				Username: user.Username,
				Posts:    posts,
			}
			temp, err := template.ParseFiles("ui/homepage.html")
			if err != nil {
				log.Fatal(err)
			}
			temp.Execute(w, display)
		}
		for i := '0'; i <= '4'; i++ {
			if r.FormValue("category"+string(i)) != "" {
				category = append(category, r.FormValue("category"+string(i)))
			}
		}
		if len(category) != 0 {
			// fmt.Println(category)
			posts, err := h.service.Post.GetPostsByCategoty(category)
			if err != nil {
				log.Print(err)
			}
			display := Display{
				Username: user.Username,
				Posts:    posts,
				Category: category,
			}
			temp, err := template.ParseFiles("ui/homepage.html")
			if err != nil {
				log.Fatal(err)
			}
			temp.Execute(w, display)
		} else {
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
		if r.FormValue("postLike") != "" {
			postId, err := strconv.Atoi(r.FormValue("postLike"))
			if err != nil {
				ErrorHandler(w, "Error", 404)
				return
			}
			like := models.Like{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Like.SetPostLike(like); err != nil {
				log.Print(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, err := strconv.Atoi(r.FormValue("postDislike"))
			if err != nil {
				ErrorHandler(w, "Error", 404)
				return
			}
			dislike := models.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Print(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	}
}

func (h *Handler) myPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxUserKey).(models.User)
	var empty models.User
	if user == empty {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/myPosts" {
			ErrorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var posts []models.Post
		posts, err := h.service.Post.ShowMyPosts(user.Id)
		if err != nil {
			log.Fatal(err)
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/myposts.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, display)
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, err := strconv.Atoi(r.FormValue("postLike"))
			if err != nil {
				ErrorHandler(w, "Error", 404)
				return
			}
			like := models.Like{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Like.SetPostLike(like); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, err := strconv.Atoi(r.FormValue("postDislike"))
			if err != nil {
				ErrorHandler(w, "Error", 404)
				return
			}
			dislike := models.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	}
}

func (h *Handler) myCommentPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxUserKey).(models.User)
	var empty models.User
	if user == empty {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/myCommentPosts" {
			ErrorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var posts []models.Post
		posts, err := h.service.Post.ShowMyCommentPosts(user.Id)
		if err != nil {
			log.Fatal(err)
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/mycommentposts.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, display)
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, err := strconv.Atoi(r.FormValue("postLike"))
			if err != nil {
				ErrorHandler(w, "Error", 404)
				return
			}
			like := models.Like{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Like.SetPostLike(like); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, err := strconv.Atoi(r.FormValue("postDislike"))
			if err != nil {
				ErrorHandler(w, "Error", 404)
				return
			}
			dislike := models.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	}
}

func (h *Handler) myLikedPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxUserKey).(models.User)
	var empty models.User
	if user == empty {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path != "/myLikedPosts" {
			ErrorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var posts []models.Post
		posts, err := h.service.Post.ShowMyLikedPosts(user.Id)
		if err != nil {
			log.Fatal(err)
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/mylikedposts.html")
		if err != nil {
			log.Fatal(err)
		}
		temp.Execute(w, display)
	case http.MethodPost:
		if r.FormValue("postLike") != "" {
			postId, err := strconv.Atoi(r.FormValue("postLike"))
			if err != nil {
				ErrorHandler(w, "Error", 404)
				return
			}
			like := models.Like{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Like.SetPostLike(like); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, err := strconv.Atoi(r.FormValue("postDislike"))
			if err != nil {
				ErrorHandler(w, "Error", 404)
				return
			}
			dislike := models.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	}
}
