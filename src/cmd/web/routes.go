package main

import (
	"anti-shoplifting-webapp/internal/config"
	"anti-shoplifting-webapp/internal/handlers"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Post("/register_fcm_token", handlers.Repo.PostFcmToken)

	// signin
	mux.Get("/", handlers.Repo.Signin)            // Get  /signin
	mux.Get("/signin", handlers.Repo.Signin)      // Get  /signin
	mux.Post("/signin", handlers.Repo.PostSignin) // Post /signin

	// signup
	mux.Get("/signup", handlers.Repo.Signup)                                                   // GET  /signup
	mux.Get("/signup/company", handlers.Repo.SignupCompany)                                    // GET  /signup/company
	mux.Post("/signup/company/basicinfo", handlers.Repo.PostSignupCompanyBasicInfo)            // POST /signup/company/basicinfo
	mux.Get("/signup/company/basicinfo/contd", handlers.Repo.SignupCompanyBasicInfoContd)      // GET /signup/company/basicinfo/contd
	mux.Post("/signup/company/basicinfo/contd", handlers.Repo.PostSignupCompanyBasicInfoContd) // POST /signup/company/basicinfo/contd
	mux.Get("/signup/company/payment", handlers.Repo.SignupCompanyPayment)                     // GET /signup/company/payment
	mux.Post("/signup/company/payment", handlers.Repo.PostSignupCompanyPayment)                // POST /signup/company/payment
	mux.Get("/signup/company/confirm", handlers.Repo.SignupCompanyConfirm)                     // GET /signup/company/confirm
	mux.Post("/signup/company/confirm", handlers.Repo.PostSignupCompanyConfirm)                // POST /signup/company/confirm
	mux.Get("/signup/company/complete", handlers.Repo.SignupCompanyComplete)                   // GET /signup/company/complete
	mux.Get("/signup/company/verify-email", handlers.Repo.SignupCompanyVeryfyEmail)            // GET /signup/company/veryfy-email

	mux.Get("/signup/store", handlers.Repo.SignupStore)                                    // GET  /signup/store
	mux.Post("/signup/store/basicinfo", handlers.Repo.PostSignupStoreBasicInfo)            // POST /signup/store/basicinfo
	mux.Get("/signup/store/basicinfo/contd", handlers.Repo.SignupStoreBasicInfoContd)      // GET /signup/store/basicinfo/contd
	mux.Post("/signup/store/basicinfo/contd", handlers.Repo.PostSignupStoreBasicInfoContd) // POST /signup/store/basicinfo/contd
	mux.Get("/signup/store/password", handlers.Repo.SignupStorePassword)                   // GET /signup/store/password
	mux.Post("/signup/store/password", handlers.Repo.PostSignupStorePassword)              // POST /signup/store/password
	mux.Get("/signup/store/confirm", handlers.Repo.SignupStoreConfirm)                     // GET /signup/store/confirm
	mux.Post("/signup/store/confirm", handlers.Repo.PostSignupStoreConfirm)                // POST /signup/store/confirm
	mux.Get("/signup/store/complete", handlers.Repo.SignupStoreComplete)                   // GET /signup/store/complete
	mux.Get("/signup/store/verify-email", handlers.Repo.SignupStoreVeryfyEmail)            // GET /signup/store/veryfy-email
	mux.Get("/signup/store/approve-new-store", handlers.Repo.SignupStoreApproveNewStore)   // GET /signup/store/approve-new-store
	mux.Get("/signup/store/reset-password", handlers.Repo.SignupStoreResetPassword)        // GET /signup/store/reset-password
	mux.Post("/signup/store/reset-password", handlers.Repo.PostSignupStoreResetPassword)   // POST /signup/store/reset-password
	// signout
	mux.Route("/signout", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.Signout) // Get  /signout
	})

	mux.Route("/blacklists", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.Blacklists)
		mux.Post("/register", handlers.Repo.BlacklistRegister) // POST /incidents/global_id
	})

	mux.Route("/dashboard", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.Dashboard)
	})

	mux.Route("/incidents", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.Incidents)
		mux.Post("/global_id", handlers.Repo.IncidentsByGlobalId) // POST /incidents/global_id
		mux.Post("/sales_items", handlers.Repo.IncidentsSalesItem) // POST /incidents/sales_items
	})

	mux.Route("/settings", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.Settings)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	// 2022.08.27 add for firebase
	// mux.Handle("/templates/*", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates/"))))
	mux.Handle("/firebase-messaging-sw.js", http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	mux.Handle("/manifest.json", http.StripPrefix("/", http.FileServer(http.Dir("./"))))

	return mux
}
