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
	a.DB.LogMode(true)
	// Routing for handling the projects
	a.Get("/healthz", a.HealthzCheck)

	//Auth
	a.Post("/login", a.Login)
	a.Post("/register", a.Register)
	a.GetWithAuth("/me", a.GetMyDetail)

	// Partner 3rd party provider, thanks them
	a.GetWithAuth("/partner/tmdb/movie-detail", a.GetMovieDetail)
	a.GetWithAuth("/partner/tmdb/search-movie", a.GetSearchMovie)
	a.GetWithAuth("/partner/duckduckgo/image-search", a.GetDuckDuckGoImage)
	a.GetWithAuth("/partner/google/image-search", a.GetGoogleImage)

	// Person resources
	a.Get("/person", a.GetAllPerson)
	a.PostWithAuth("/person", a.CreatePerson)
	a.Get("/person/{id}", a.GetPerson)
	a.PutWithAuth("/person/{id}", a.UpdatePerson)
	a.DeleteWithAuth("/person/{id}", a.DeletePerson)

	// Movie Resource
	a.Get("/movie", a.GetAllMovie)
	a.PostWithAuth("/movie", a.CreateMovie)
	a.Get("/movie/{id}", a.GetMovie)
	a.PutWithAuth("/movie/{id}", a.UpdateMovie)
	a.DeleteWithAuth("/movie/{id}", a.DeleteMovie)

	// Genre Resources
	a.Get("/genre", a.GetAllGenre)
	a.PostWithAuth("/genre", a.CreateGenre)
	a.Get("/genre/{id}", a.GetGenre)
	a.PutWithAuth("/genre/{id}", a.UpdateGenre)
	a.DeleteWithAuth("/genre/{id}", a.DeleteGenre)

	// Image Resources
	a.Get("/image", a.GetAllImage)
	a.PostWithAuth("/image", a.CreateImage)
	a.Get("/image/{id}", a.GetImage)
	a.PutWithAuth("/image/{id}", a.UpdateImage)
	a.DeleteWithAuth("/image/{id}", a.DeleteImage)
}

// HealthzCheck handler
func (a *App) HealthzCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok"))
}

// AUTH

// Login handle user login
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	handler.Login(a.DB, w, r)
}

// Register handle user register
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	handler.Register(a.DB, w, r)
}

// GetMyDetail handle logged in user detail
func (a *App) GetMyDetail(w http.ResponseWriter, r *http.Request) {
	handler.GetMyDetail(a.DB, w, r)
}

// GetMovieDetail handler
func (a *App) GetMovieDetail(w http.ResponseWriter, r *http.Request) {
	handler.GetMovieDetail(w, r)
}

// PARTNER

// GetSearchMovie handler
func (a *App) GetSearchMovie(w http.ResponseWriter, r *http.Request) {
	handler.SearchMovie(w, r)
}

// GetDuckDuckGoImage handler
func (a *App) GetDuckDuckGoImage(w http.ResponseWriter, r *http.Request) {
	handler.GetDuckDuckGoImage(w, r)
}

// GetGoogleImage handler
func (a *App) GetGoogleImage(w http.ResponseWriter, r *http.Request) {
	handler.GetGoogleImage(w, r)
}

// PERSON

// GetAllPerson handle getAll people
func (a *App) GetAllPerson(w http.ResponseWriter, r *http.Request) {
	handler.GetAllPerson(a.DB, w, r)
}

// CreatePerson handler
func (a *App) CreatePerson(w http.ResponseWriter, r *http.Request) {
	handler.CreatePerson(a.DB, w, r)
}

// GetPerson Handler
func (a *App) GetPerson(w http.ResponseWriter, r *http.Request) {
	handler.GetPerson(a.DB, w, r)
}

// UpdatePerson Handler
func (a *App) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	handler.UpdatePerson(a.DB, w, r)
}

// DeletePerson Handler
func (a *App) DeletePerson(w http.ResponseWriter, r *http.Request) {
	handler.DeletePerson(a.DB, w, r)
}

// MOVIE

// GetAllMovie handler
func (a *App) GetAllMovie(w http.ResponseWriter, r *http.Request) {
	handler.GetAllMovie(a.DB, w, r)
}

// CreateMovie handler
func (a *App) CreateMovie(w http.ResponseWriter, r *http.Request) {
	handler.CreateMovie(a.DB, w, r)
}

// GetMovie handler
func (a *App) GetMovie(w http.ResponseWriter, r *http.Request) {
	handler.GetMovie(a.DB, w, r)
}

// UpdateMovie handler
func (a *App) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	handler.UpdateMovie(a.DB, w, r)
}

// DeleteMovie handler
func (a *App) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	handler.DeleteMovie(a.DB, w, r)
}

// GENRE

// GetAllGenre handler
func (a *App) GetAllGenre(w http.ResponseWriter, r *http.Request) {
	handler.GetAllGenre(a.DB, w, r)
}

// CreateGenre handler
func (a *App) CreateGenre(w http.ResponseWriter, r *http.Request) {
	handler.CreateGenre(a.DB, w, r)
}

// GetGenre handler
func (a *App) GetGenre(w http.ResponseWriter, r *http.Request) {
	handler.GetGenre(a.DB, w, r)
}

// UpdateGenre handler
func (a *App) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	handler.UpdateGenre(a.DB, w, r)
}

// DeleteGenre handler
func (a *App) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	handler.DeleteGenre(a.DB, w, r)
}

// IMAGE

// GetAllImage handler
func (a *App) GetAllImage(w http.ResponseWriter, r *http.Request) {
	handler.GetAllImage(a.DB, w, r)
}

// CreateImage handler
func (a *App) CreateImage(w http.ResponseWriter, r *http.Request) {
	handler.CreateImage(a.DB, w, r)
}

// GetImage handler
func (a *App) GetImage(w http.ResponseWriter, r *http.Request) {
	handler.GetImage(a.DB, w, r)
}

// UpdateImage handler
func (a *App) UpdateImage(w http.ResponseWriter, r *http.Request) {
	handler.UpdateImage(a.DB, w, r)
}

// DeleteImage handler
func (a *App) DeleteImage(w http.ResponseWriter, r *http.Request) {
	handler.DeleteImage(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		// Debug: true,
	})
	handler := c.Handler(a.Router)
	log.Fatal(http.ListenAndServe(host, handler))
}
