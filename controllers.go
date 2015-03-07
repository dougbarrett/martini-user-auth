package main

import (
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func controllerHomepage(r render.Render, p Profile) {
	var viewData struct {
		SiteVD
	}

	viewData.CurrentProfile = p

	r.HTML(http.StatusOK, "homepage", viewData)
}

func controllerSignupGet(r render.Render, x csrf.CSRF) {
	var viewData struct {
		SiteVD
		Profile Profile
		Message string
	}

	viewData.Token = x.GetToken()

	r.HTML(http.StatusOK, "profile/edit", viewData)
}

func controllerSignupPost(r render.Render, p ProfileNew) (errorMessage string) {
	var profile Profile
	if p.Password != p.PasswordConfirm {
		errorMessage = "Passwords do not match"
		return
	}

	db.Where("email = ?", p.Email).
		Find(&profile)

	if profile.Id != 0 {
		errorMessage = "Email already exists in system"
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(p.Password), 10)
	p.Password = string(hashedPassword)

	db.Save(&p)
	r.Redirect("/")
	return
}

func controllerLoginGet(r render.Render, x csrf.CSRF) {
	var viewData struct {
		SiteVD
		Profile Profile
	}

	viewData.Token = x.GetToken()

	r.HTML(http.StatusOK, "profile/login", viewData)
}

func controllerLoginPost(r render.Render, p Profile, s sessions.Session) (errorMessage string) {
	var profile Profile

	db.Where("email = ?", p.Email).
		Find(&profile)

	errorMessage = "Email or password is incorrect"

	if profile.Id == 0 {
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(p.Password))

	if err != nil {
		return
	}

	s.Set("profileID", profile.Id)

	r.Redirect("/")
	return
}

func controllerLogoutGet(s sessions.Session, r render.Render) {
	op := sessions.Options{
		MaxAge: -1,
	}
	s.Options(op)

	s.Clear()
	r.Redirect("/")
}

func controllerProfileGet(r render.Render, p Profile) {
	var viewData struct {
		SiteVD
	}

	viewData.CurrentProfile = p

	r.HTML(http.StatusOK, "profile/view", viewData)
}
