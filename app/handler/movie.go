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
	query.Table("people").Count(&count)

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

	if err := db.Create(&movie).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
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

	if err := db.Save(&movie).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
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
		First(&movie, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &movie
}
