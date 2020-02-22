package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/condrowiyono/living-rooms-api/app/handler"
	"github.com/condrowiyono/living-rooms-api/app/model"
	"github.com/condrowiyono/living-rooms-api/config"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// Set all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/healthz", a.HealthzCheck)

	a.Get("/sample-link", a.GetAllSampleLink)
	a.Post("/sample-link", a.CreateSampleLink)
	a.Get("/sample-link/{id}", a.GetSampleLink)
	a.Put("/sample-link/{id}", a.UpdateSampleLink)
	a.Delete("/sample-link/{id}", a.DeleteSampleLink)

	a.Get("/show", a.GetAllShow)
	a.Post("/show", a.CreateShow)
	a.Get("/show/{id}", a.GetShow)
	a.Put("/show/{id}", a.UpdateShow)
	a.Delete("/show/{id}", a.DeleteShow)

	a.Get("/playlist", a.GetAllPlaylist)
	a.Post("/playlist", a.CreatePlaylist)
	a.Get("/playlist/{id}", a.GetPlaylist)
	a.Put("/playlist/{id}", a.UpdatePlaylist)
	a.Delete("/playlist/{id}", a.DeletePlaylist)

	a.GetWithAuth("/movie-detail", a.GetMovieDetail)

	a.Post("/login", a.Login)
	a.Post("/register", a.Register)
	a.GetWithAuth("/me", a.GetMyDetail)
}

// HealthzCheck handler
func (a *App) HealthzCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok"))
}

// GetAllSampleLink handler
func (a *App) GetAllSampleLink(w http.ResponseWriter, r *http.Request) {
	handler.GetAllSampleLink(a.DB, w, r)
}

func (a *App) CreateSampleLink(w http.ResponseWriter, r *http.Request) {
	handler.CreateSampleLink(a.DB, w, r)
}

func (a *App) GetSampleLink(w http.ResponseWriter, r *http.Request) {
	handler.GetSampleLink(a.DB, w, r)
}

func (a *App) UpdateSampleLink(w http.ResponseWriter, r *http.Request) {
	handler.UpdateSampleLink(a.DB, w, r)
}

func (a *App) DeleteSampleLink(w http.ResponseWriter, r *http.Request) {
	handler.DeleteSampleLink(a.DB, w, r)
}

// GetAllShow handler
func (a *App) GetAllShow(w http.ResponseWriter, r *http.Request) {
	handler.GetAllShow(a.DB, w, r)
}

func (a *App) CreateShow(w http.ResponseWriter, r *http.Request) {
	handler.CreateShow(a.DB, w, r)
}

func (a *App) GetShow(w http.ResponseWriter, r *http.Request) {
	handler.GetShow(a.DB, w, r)
}

func (a *App) UpdateShow(w http.ResponseWriter, r *http.Request) {
	handler.UpdateShow(a.DB, w, r)
}

func (a *App) DeleteShow(w http.ResponseWriter, r *http.Request) {
	handler.DeleteShow(a.DB, w, r)
}

func (a *App) GetAllPlaylist(w http.ResponseWriter, r *http.Request) {
	handler.GetAllPlaylist(a.DB, w, r)
}

func (a *App) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	handler.CreatePlaylist(a.DB, w, r)
}

func (a *App) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	handler.GetPlaylist(a.DB, w, r)
}

func (a *App) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	handler.UpdatePlaylist(a.DB, w, r)
}

func (a *App) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	handler.DeletePlaylist(a.DB, w, r)
}

func (a *App) GetMovieDetail(w http.ResponseWriter, r *http.Request) {
	handler.GetMovieDetail(w, r)
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	handler.Login(a.DB, w, r)
}

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	handler.Register(a.DB, w, r)
}

func (a *App) GetMyDetail(w http.ResponseWriter, r *http.Request) {
	handler.GetMyDetail(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	handler := c.Handler(a.Router)
	log.Fatal(http.ListenAndServe(host, handler))
}
