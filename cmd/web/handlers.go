package main

import (
	"errors"
	"fmt"
	"github.com/djcurill/snippetbox/internal/models"
	"github.com/djcurill/snippetbox/internal/validators"
	"net/http"
	"strconv"
)

// use go-playground for automatic form parsing.
// tags indicate which form fields map to struct
// attributes. a "-" mean to ignore this value
type snippetCreateForm struct {
	Title                string `form:"title"`
	Content              string `form:"content"`
	Expires              int    `form:"expires"`
	validators.Validator `form:"-"`
}

type userSignupForm struct {
	Name                 string `form:"name"`
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

type userLoginForm struct {
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

// Create a home pathfor the webpages
func (app *app) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	err = app.render(w, r, "home.gohtml", http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Snippet view
func (app *app) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			http.NotFound(w, r)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// populate templateData
	data := app.newTemplateData(r)
	data.Snippet = snippet

	err = app.render(w, r, "snippetView.gohtml", http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Create snippet
func (app *app) createSnippet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	err := app.render(w, r, "create.gohtml", http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// POST endpoint for creating the snipper
func (app *app) createSnippetPost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm
	err := app.decodeForm(r, &form)

	if err != nil {
		app.logger.Error("error parsing form", "err", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.NotBlank(form.Title), "title", "this field cannot be empty")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "this field cannot be longer than 100 chars")
	form.CheckField(form.NotBlank(form.Content), "content", "this field cannot be empty")
	form.CheckField(validators.PermittedValue(form.Expires, 1, 7, 365), "expires", "this field must equal 1,7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		err = app.render(w, r, "create.gohtml", http.StatusUnprocessableEntity, data)
		if err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// if there is no existing session, a new one will be created
	// this is handled by the middleware
	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully!")

	// Redirect to the view of the snippet after it has been created
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *app) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	err := app.render(w, r, "login.gohtml", http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *app) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodeForm(r, &form)
	if err != nil {
		// app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validate email
	form.CheckField(form.NotBlank(form.Email), "email", "This field cannot be empty")
	form.CheckField(form.MaxChars(form.Email, 255), "email", "This field has a 255 character limit")
	form.CheckField(form.Matches(form.Email, validators.EmailRX), "email", "Invalid email")

	// validate password
	form.CheckField(form.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(form.MaxBytes(form.Password, 72), "password", "This field must not be more than 72 bytes long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		err := app.render(w, r, "login.gohtml", http.StatusUnprocessableEntity, data)
		if err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	// authenticate user
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		data := app.newTemplateData(r)
		if errors.Is(err, models.ErrInavlidCredentials) {
			app.logger.Info("login attempt failed", "userId", id)
			form.NonFieldErrors = append(form.NonFieldErrors, "Invalid email address or password")
			data.Form = form
			err = app.render(w, r, "login.gohtml", http.StatusUnprocessableEntity, data)
			if err != nil {
				app.internalServerError(w, r, err)
			}
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}
	app.logger.Info("login successful", "userId", id)

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// add the id to the current user session so they are currently
	// logged in
	app.sessionManager.Put(r.Context(), "authenticatedUserId", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *app) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	form := userSignupForm{}

	data.Form = form

	err := app.render(w, r, "signup.gohtml", http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *app) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodeForm(r, &form)
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
	}

	// validate fields
	form.CheckField(form.NotBlank(form.Email), "email", "This field cannot be empty")
	form.CheckField(form.MaxChars(form.Email, 255), "email", "This field has a 255 character limit")
	form.CheckField(form.Matches(form.Email, validators.EmailRX), "email", "Invalid email")

	form.CheckField(form.NotBlank(form.Name), "name", "This field cannot be empty")
	form.CheckField(form.MaxChars(form.Name, 255), "name", "This field has a 255 character limit")

	form.CheckField(form.MinChars(form.Password, 8), "password", "Field must be at least 8 characters long")
	form.CheckField(form.MaxBytes(form.Password, 72), "password", "This field cannot be more than 72 bytes. Make it shorter.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		err = app.render(w, r, "signup.gohtml", http.StatusUnprocessableEntity, data)
		if err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDpuplicateEmail) {
			form.AddFieldError("email", "Email address already in use")
			data := app.newTemplateData(r)
			data.Form = form
			err = app.render(w, r, "signup.gohtml", http.StatusUnprocessableEntity, data)
			if err != nil {
				app.internalServerError(w, r, err)
			}
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful! You may now login.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *app) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// it is good practice to renew a session before logging out
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
	}

	// remove the authenticated user id context from the users session
	app.sessionManager.Remove(r.Context(), "authenticatedUserId")

	// flash the user they logged out
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// return to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
