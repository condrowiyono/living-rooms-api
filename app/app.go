package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/condrowiyono/ruangtengah-api/app/handler"
	"github.com/condrowiyono/ruangtengah-api/app/handler/scrapper"
	"github.com/condrowiyono/ruangtengah-api/app/model"
	"github.com/condrowiyono/ruangtengah-api/config"
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
	dbURI := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
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

	// Static file
	a.Static()

	//Auth
	a.Post("/login", a.Login)
	a.Post("/register", a.Register)
	a.GetWithAuth("/me", a.GetMyDetail)

	// Partner 3rd party provider, thanks them
	a.GetWithAuth("/partner/tmdb/search", a.GetSearchMovie)
	a.GetWithAuth("/partner/tmdb/movie-detail", a.GetMovieDetail)

	a.GetWithAuth("/partner/tmdb/tv-detail", a.GetTvDetail)
	a.GetWithAuth("/partner/tmdb/tv-season", a.GetTvSeason)
	a.GetWithAuth("/partner/tmdb/tv-episode", a.GetTvEpisode)

	a.GetWithAuth("/partner/tmdb/movie-image", a.GetMovieImage)
	a.GetWithAuth("/partner/tmdb/tv-image", a.GetTvImage)
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
	a.Post("/image/upload-image", a.UploadImage)

	// Player Resources
	a.Get("/player", a.GetPlayer)
	a.Get("/video", a.GetVideo)

	// Concert Resource
	a.Get("/concert", a.GetAllConcert)
	a.PostWithAuth("/concert", a.CreateConcert)
	a.Get("/concert/{id}", a.GetConcert)
	a.PutWithAuth("/concert/{id}", a.UpdateConcert)
	a.DeleteWithAuth("/concert/{id}", a.DeleteConcert)
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

// PARTNER

// GetMovieDetail handler
func (a *App) GetMovieDetail(w http.ResponseWriter, r *http.Request) {
	scrapper.GetMovieDetail(w, r)
}

// GetMovieImage handler
func (a *App) GetMovieImage(w http.ResponseWriter, r *http.Request) {
	scrapper.GetMovieImage(w, r)
}

// GetSearchMovie handler
func (a *App) GetSearchMovie(w http.ResponseWriter, r *http.Request) {
	scrapper.SearchMovie(w, r)
}

// GetDuckDuckGoImage handler
func (a *App) GetDuckDuckGoImage(w http.ResponseWriter, r *http.Request) {
	scrapper.GetDuckDuckGoImage(w, r)
}

// GetGoogleImage handler
func (a *App) GetGoogleImage(w http.ResponseWriter, r *http.Request) {
	scrapper.GetGoogleImage(w, r)
}

// GetTvImage s
func (a *App) GetTvImage(w http.ResponseWriter, r *http.Request) {
	scrapper.GetTvImage(w, r)
}

func (a *App) GetTvDetail(w http.ResponseWriter, r *http.Request) {
	scrapper.GetTvDetail(w, r)
}

func (a *App) GetTvSeason(w http.ResponseWriter, r *http.Request) {
	scrapper.GetTvSeason(w, r)
}

func (a *App) GetTvEpisode(w http.ResponseWriter, r *http.Request) {
	scrapper.GetTvEpisode(w, r)
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

// UploadImage handler
func (a *App) UploadImage(w http.ResponseWriter, r *http.Request) {
	handler.UploadImage(a.DB, w, r)
}

// GetPlayer hanlder
func (a *App) GetPlayer(w http.ResponseWriter, r *http.Request) {
	handler.GetPlayer(a.DB, w, r)
}

// GetVideo hanlder
func (a *App) GetVideo(w http.ResponseWriter, r *http.Request) {
	handler.GetVideo(a.DB, w, r)
}

// Concert

// GetAllConcert handler
func (a *App) GetAllConcert(w http.ResponseWriter, r *http.Request) {
	handler.GetAllConcert(a.DB, w, r)
}

// CreateConcert handler
func (a *App) CreateConcert(w http.ResponseWriter, r *http.Request) {
	handler.CreateConcert(a.DB, w, r)
}

// GetConcert handler
func (a *App) GetConcert(w http.ResponseWriter, r *http.Request) {
	handler.GetConcert(a.DB, w, r)
}

// UpdateConcert handler
func (a *App) UpdateConcert(w http.ResponseWriter, r *http.Request) {
	handler.UpdateConcert(a.DB, w, r)
}

// DeleteConcert handler
func (a *App) DeleteConcert(w http.ResponseWriter, r *http.Request) {
	handler.DeleteConcert(a.DB, w, r)
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
