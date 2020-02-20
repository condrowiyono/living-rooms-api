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
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// App initialize with predefined configuration
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
	a.Get("/healthz", a.HeathzCheck)

	a.Get("/sample-link", a.GetAllSampleLink)
	a.Post("/sample-link", a.CreateSampleLink)
	a.Get("/sample-link/{id}", a.GetSampleLink)
	a.Put("/sample-link/{id}", a.UpdateSampleLink)
	a.Delete("/sample-link/{id}", a.DeleteSampleLink)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Handler HealthzCheck
func (a *App) HeathzCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok"))
}

// GetAllSampleLink get add Sample Link
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

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
