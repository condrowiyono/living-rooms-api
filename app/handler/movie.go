package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/condrowiyono/living-rooms-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllMovie(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	page := string(vars.Get("page"))
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	limit := string(vars.Get("limit"))
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 25
	}

	offsetInt := (pageInt - 1) * limitInt
	var count int64
	query := db.Offset(0)
	query.Model(&model.Movie{}).Count(&count)

	movie := []model.Movie{}
	query = db.Limit(limitInt).
		Offset(offsetInt).
		Preload("Banners").
		Preload("Genres").
		Preload("Posters").
		Preload("Languages").
		Preload("Country").
		Preload("Actors").
		Preload("Crews").
		Preload("Player").
		Find(&movie)

	if err := query.Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write Response
	meta := Meta{limitInt, offsetInt, pageInt, count}
	respondJSON(w, http.StatusOK, meta, movie)
}

func CreateMovie(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	movie := model.Movie{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&movie); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Set("gorm:association_autoupdate", false).Create(&movie).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Save another association
	genres := []model.Genre{}
	actors := []model.Person{}
	crews := []model.Person{}
	languages := []model.Country{}

	for _, v := range movie.Genres {
		genre := model.Genre{}
		db.Where(model.Genre{Name: v.Name}).FirstOrCreate(&genre)
		genres = append(genres, genre)
	}

	for _, v := range movie.Actors {
		actor := model.Person{}
		db.Where(model.Person{Name: v.Name}).FirstOrCreate(&actor)
		actors = append(actors, actor)
	}

	for _, v := range movie.Crews {
		crew := model.Person{}
		db.Where(model.Person{Name: v.Name}).FirstOrCreate(&crew)
		crews = append(crews, crew)
	}

	for _, v := range movie.Languages {
		language := model.Country{}
		db.Where(model.Country{Iso639_1: v.Iso639_1}).FirstOrCreate(&language)
		languages = append(languages, language)
	}

	originalCountry := model.Country{}
	db.Where(model.Country{Iso639_1: movie.Country.Iso639_1}).FirstOrCreate(&originalCountry)

	db.Model(&movie).Association("Genres").Replace(genres)
	db.Model(&movie).Association("Actors").Replace(actors)
	db.Model(&movie).Association("Crews").Replace(crews)
	db.Model(&movie).Association("Languages").Replace(languages)
	db.Model(&movie).Association("Country").Replace(originalCountry)

	respondJSON(w, http.StatusCreated, nil, movie)
}

func GetMovie(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	movie := getMovieOr404(db, id, w, r)
	if movie == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, movie)
}

func UpdateMovie(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	movie := getMovieOr404(db, id, w, r)
	if movie == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&movie); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Set("gorm:association_autoupdate", false).Save(&movie).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Save another association
	genres := []model.Genre{}
	actors := []model.Person{}
	crews := []model.Person{}
	languages := []model.Country{}

	for _, v := range movie.Genres {
		genre := model.Genre{}
		db.Where(model.Genre{Name: v.Name}).FirstOrCreate(&genre)
		genres = append(genres, genre)
	}

	for _, v := range movie.Actors {
		actor := model.Person{}
		db.Where(model.Person{Name: v.Name}).FirstOrCreate(&actor)
		actors = append(actors, actor)
	}

	for _, v := range movie.Crews {
		crew := model.Person{}
		db.Where(model.Person{Name: v.Name}).FirstOrCreate(&crew)
		crews = append(crews, crew)
	}

	for _, v := range movie.Languages {
		language := model.Country{}
		db.Where(model.Country{Iso639_1: v.Iso639_1}).FirstOrCreate(&language)
		languages = append(languages, language)
	}

	originalCountry := model.Country{}
	db.Where(model.Country{Iso639_1: movie.Country.Iso639_1}).FirstOrCreate(&originalCountry)

	db.Model(&movie).Association("Genres").Replace(genres)
	db.Model(&movie).Association("Actors").Replace(actors)
	db.Model(&movie).Association("Crews").Replace(crews)
	db.Model(&movie).Association("Languages").Replace(languages)
	db.Model(&movie).Association("Country").Replace(originalCountry)

	respondJSON(w, http.StatusOK, nil, movie)
}

func DeleteMovie(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	movie := getMovieOr404(db, id, w, r)
	if movie == nil {
		return
	}
	if err := db.Delete(&movie).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil, nil)
}

// getMovieOr404 gets a instance if exists, or respond the 404 error otherwise
func getMovieOr404(db *gorm.DB, id int64, w http.ResponseWriter, r *http.Request) *model.Movie {
	movie := model.Movie{}
	if err := db.
		Preload("Banners").
		Preload("Genres").
		Preload("Posters").
		Preload("Languages").
		Preload("Country").
		Preload("Actors").
		Preload("Crews").
		Preload("Player").
		Preload("Videos").
		First(&movie, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &movie
}
