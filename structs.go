package main

import (
	"github.com/BurntSushi/toml"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net/http"
)

var (
	tomlFile = "config.toml"
	config   tomlConfig
	db       gorm.DB
	app      WebApplication
)

type tomlConfig struct {
	SqlConnection   string
	SetMaxIdleConns int
	SetMaxOpenConns int
	DbLogMode       bool
}

type WebApplication struct {
	M *martini.ClassicMartini
}

// Set up your routing, a.M allows you access Martini

func (a *WebApplication) SetupControllers(args ...interface{}) {
	a.M.Get("/", a.checkUser, controllerHomepage)

	a.M.Group("/signup", func(r martini.Router) {
		r.Get("", controllerSignupGet)
		r.Post("", binding.Bind(ProfileNew{}), controllerSignupPost)
	}, a.checkUser, a.publicOnly)

	a.M.Group("/login", func(r martini.Router) {
		r.Get("", controllerLoginGet)
		r.Post("", binding.Bind(Profile{}), controllerLoginPost)
	}, a.checkUser, a.publicOnly)

	a.M.Get("/logout", controllerLogoutGet)

	a.M.Get("/profile", a.checkUser, a.requireLogin, controllerProfileGet)

	if args[0].(bool) {
		a.M.Run()
	}
}

// Set up your environment, a.M allows you to access Martini to plug in
// any middleware you may want to use

func (a *WebApplication) SetupEnvironment() (err error) {
	if _, err = toml.DecodeFile(tomlFile, &config); err != nil {
		return
	}

	if err = a.setupDB(); err != nil {
		return
	}

	app.M.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	store := sessions.NewCookieStore([]byte("er6LlBrxn96XbWIsLFbDZsLV4911e8CA"))
	app.M.Use(sessions.Sessions("1U9Qs77i3SAHm969G3f489Zev83y88rSgp", store))

	a.M.Use(csrf.Generate(&csrf.Options{
		Secret:     "token123",
		SessionKey: "profileID",
		// Custom error response.
		ErrorFunc: func(w http.ResponseWriter) {
			http.Error(w, "CSRF token validation failed", http.StatusBadRequest)
		},
	}))

	return
}

// Set up your database, config is read for idle conns, open conns and if gorm
// should output query information to the console while debugging.

func (a *WebApplication) setupDB() (err error) {
	db, err = gorm.Open("mysql", config.SqlConnection)

	if err != nil {
		return
	}

	db.DB().SetMaxIdleConns(config.SetMaxIdleConns)
	db.DB().SetMaxOpenConns(config.SetMaxOpenConns)
	db.LogMode(config.DbLogMode)
	return
}

func (a *WebApplication) publicOnly(p Profile, r render.Render) {
	if p.Id != 0 {
		r.Redirect("/")
	}
}

func (a *WebApplication) requireLogin(p Profile, r render.Render) {
	if p.Id == 0 {
		r.Redirect("/login")
	}
}

func (a *WebApplication) checkUser(s sessions.Session, c martini.Context) {
	var profile Profile
	v := s.Get("profileID")

	if v != nil {
		db.Find(&profile, v.(int64))
	}

	c.Map(profile)
}

type ProfileNew struct {
	Id              int64  `form:"id"`
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"password_confirm" sql:"-" binding:"required"`
}

type Profile struct {
	Id       int64  `form:"id"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (p ProfileNew) TableName() (table string) {
	table = "profiles"
	return
}

// This is a struct of data that will be needed on every page
type SiteVD struct {
	CurrentProfile Profile
	Token          string
}
