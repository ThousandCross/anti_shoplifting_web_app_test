package main

import (
	"anti-shoplifting-webapp/internal/config"
	"anti-shoplifting-webapp/internal/handlers"
	"anti-shoplifting-webapp/internal/helpers"
	"anti-shoplifting-webapp/internal/models"
	"anti-shoplifting-webapp/internal/render"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	// what am I going to put in the session
	gob.Register(models.Signin{})
	gob.Register(models.CompanySignupConfirm{})
	gob.Register(models.CompanyBasicInfo{})
	gob.Register(models.CompanyBasicInfoContd{})
	gob.Register(models.CompanyPayment{})

	gob.Register(models.StoreSignupConfirm{})
	gob.Register(models.StoreBasicInfo{})
	gob.Register(models.StoreBasicInfoContd{})
	gob.Register(models.StorePassword{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	// LaxはサードパーティCookieを基本的にブロックします。
	session.Cookie.SameSite = http.SameSiteLaxMode
	// 2022/08/26 fcmの動作を確認するため一時的に変更
	//session.Cookie.SameSite = http.SameSiteNoneMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	fmt.Printf("Staring application on port %s", portNumber)

	// init from data
	handlers.InitFormData()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
