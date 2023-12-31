package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aspandyar/forum/internal/models"
	"github.com/aspandyar/forum/internal/validator"
)

type forumCreateForm struct {
	Title   string
	Content string
	Tags    string
	Expires int
	validator.Validator
}

type userSingupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

type forumLikeForm struct {
	LikeStatus int
	ForumID    int
	UserID     int
	CommentID  int
	validator.Validator
}

type forumCommentForm struct {
	ForumID int
	UserID  int
	Comment string
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	forums, err := app.forums.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)

	data.Forums = forums

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) allForum(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/showAll" {
		app.notFound(w)
		return
	}

	forums, err := app.forums.ShowAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)

	data.Forums = forums

	app.render(w, http.StatusOK, "allForums.tmpl.html", data)
}

func (app *application) forumAllUserPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/forum/allPosts" {
		app.notFound(w)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, _, err := app.sessions.GetSession(cookie.Value)
	if err != nil {
		app.serverError(w, err)
		return
	}

	forums, err := app.forums.ShowAllUserPosts(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)

	data.Forums = forums

	app.render(w, http.StatusOK, "allForums.tmpl.html", data)
}

func (app *application) forumAllUserLikes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/forum/allLikes" {
		app.notFound(w)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, _, err := app.sessions.GetSession(cookie.Value)
	if err != nil {
		app.serverError(w, err)
		return
	}

	forums, err := app.forums.ShowAllUserLikes(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)

	data.Forums = forums

	app.render(w, http.StatusOK, "allForums.tmpl.html", data)
}

func (app *application) forumCategory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/forum/category" {
		app.notFound(w)
		return
	}

	var forum []*models.Forum
	var err error
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		selectedTags := r.Form["tags"]
		customTagsStr := r.PostForm.Get("custom_tags")

		tags := app.processTags(selectedTags, customTagsStr)

		forum, err = app.forums.ShowCategory(tags)
		if err != nil {
			app.serverError(w, err)
			return
		}
	} else {
		forum, err = app.forums.ShowAll()
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	// tagsStr := strings.Join(tags, ", ")

	data := app.newTemplateData(r)

	data.Forums = forum

	app.render(w, http.StatusOK, "category.tmpl.html", data)
}

func (app *application) forumView(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) != 4 || parts[1] != "forum" || parts[2] != "view" {
		http.NotFound(w, r)
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	var userID int

	cookie, err := r.Cookie("session")
	if err != nil {
		userID = 0
	} else {
		userID, _, err = app.sessions.GetSession(cookie.Value)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	forum, err := app.forums.Get(id, userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)

	data.Forum = forum

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) handleForumCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.ForumCreateGet(w, r)
	case http.MethodPost:
		app.ForumCreatePost(w, r)
	default:
		w.Header().Set("Allow", http.MethodPost+", "+http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *application) ForumCreateGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = forumCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) ForumCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadGateway)
		return
	}

	selectedTags := r.Form["tags"]
	customTagsStr := r.PostForm.Get("custom_tags")

	tags := app.processTags(selectedTags, customTagsStr)
	tagsStr := strings.Join(tags, ", ")

	form := forumCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Tags:    tagsStr,
		Expires: expires,
	}

	form.CheckField(validator.IncorrectInput(form.Tags), "tags", "Incorrect tags formation")
	form.CheckField(validator.MaxChars(form.Tags, 50), "tags", "This field cannot be more than 50 characters long")
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		app.serverError(w, err)
		return
	}
	userID, _, err := app.sessions.GetSession(cookie.Value)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.forums.Insert(form.Title, form.Content, form.Tags, expires, userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/forum/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.userSignupGet(w, r)
	case http.MethodPost:
		app.userSignupPost(w, r)
	default:
		w.Header().Set("Allow", http.MethodPost+", "+http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *application) userSignupGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSingupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := userSingupForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email or name address is already in use")
			form.AddFieldError("name", "Email or name address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.userLoginGet(w, r)
	case http.MethodPost:
		app.userLoginPost(w, r)
	default:
		w.Header().Set("Allow", http.MethodPost+", "+http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *application) userLoginGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	userID, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {

			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	session, err := app.sessions.CreateSession(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	models.SetSessionCookie(w, session.Token, session.Expiry)

	http.Redirect(w, r, "/forum/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.sessions.InvalidateSession(cookie.Value)
	if err != nil {
		app.serverError(w, err)
		return
	}

	models.ClearSessionCookie(w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	_, expiry, err := app.sessions.GetSession(cookie.Value)
	if err != nil {
		return false
	}

	return time.Now().Before(expiry)
}

func (app *application) forumIsLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var likeStatus int
	button := r.PostForm.Get("button")
	switch button {
	case "like":
		likeStatus = 1
	case "dislike":
		likeStatus = -1
	default:
		app.clientError(w, http.StatusBadRequest)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) != 4 || parts[1] != "forum" || parts[2] != "like" {
		http.NotFound(w, r)
		return
	}

	idStr := parts[3]
	forumId, err := strconv.Atoi(idStr)
	if err != nil || forumId < 1 {
		http.NotFound(w, r)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, _, err := app.sessions.GetSession(cookie.Value)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := forumLikeForm{
		LikeStatus: likeStatus,
		ForumID:    forumId,
		UserID:     userID,
	}

	id, err := app.forumLike.LikeOrDislike(form.ForumID, form.UserID, form.LikeStatus)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/forum/view/%d", id), http.StatusSeeOther)
}

func (app *application) forumIsLikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var likeStatus int
	button := r.PostForm.Get("button")
	switch button {
	case "like":
		likeStatus = 1
	case "dislike":
		likeStatus = -1
	default:
		app.clientError(w, http.StatusBadRequest)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) != 4 || parts[1] != "forum" || parts[2] != "likeComment" {
		http.NotFound(w, r)
		return
	}

	idStr := parts[3]
	commentId, err := strconv.Atoi(idStr)
	if err != nil || commentId < 1 {
		http.NotFound(w, r)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, _, err := app.sessions.GetSession(cookie.Value)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := forumLikeForm{
		LikeStatus: likeStatus,
		CommentID:  commentId,
		UserID:     userID,
	}

	id, err := app.forumLike.LikeOrDislikeComment(form.CommentID, form.UserID, form.LikeStatus)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/forum/view/%d", id), http.StatusSeeOther)
}

func (app *application) handleForumComment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.forumView(w, r)
	case http.MethodPost:
		app.ForumCommentPost(w, r)
	default:
		w.Header().Set("Allow", http.MethodPost+", "+http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}

func (app *application) ForumCommentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	comment := r.PostForm.Get("comment")

	path := r.URL.Path
	parts := strings.Split(path, "/")

	idStr := parts[3]
	forumId, err := strconv.Atoi(idStr)
	if err != nil || forumId < 1 {
		http.NotFound(w, r)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, _, err := app.sessions.GetSession(cookie.Value)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := forumCommentForm{
		ForumID: forumId,
		UserID:  userID,
		Comment: comment,
	}

	form.CheckField(validator.NotBlank(form.Comment), "comment", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "view.tmpl.html", data)
		return
	}

	id, err := app.forumComment.CommentPost(form.ForumID, form.UserID, form.Comment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/forum/view/%d", id), http.StatusSeeOther)
}
