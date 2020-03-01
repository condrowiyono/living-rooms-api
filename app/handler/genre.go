package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/condrowiyono/living-rooms-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllGenre(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	genre := []model.Genre{}
	query := db.Limit(limitInt)
	query = query.Offset(offsetInt)

	if err := query.Find(&genre).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Count all data
	var count int64
	query = query.Offset(0)
	query.Table("genres").Count(&count)

	// Write Response
	meta := Meta{limitInt, offsetInt, pageInt, count}
	respondJSON(w, http.StatusOK, meta, genre)
}

func CreateGenre(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	genre := model.Genre{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&genre); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&genre).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, nil, genre)
}

func GetGenre(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	genre := getGenreOr404(db, id, w, r)
	if genre == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, genre)
}

func UpdateGenre(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	genre := getGenreOr404(db, id, w, r)
	if genre == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&genre); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&genre).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, nil, genre)
}

func DeleteGenre(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	genre := getGenreOr404(db, id, w, r)
	if genre == nil {
		return
	}
	if err := db.Delete(&genre).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil, nil)
}

// getGenreOr404 gets a instance if exists, or respond the 404 error otherwise
func getGenreOr404(db *gorm.DB, id int64, w http.ResponseWriter, r *http.Request) *model.Genre {
	genre := model.Genre{}
	if err := db.First(&genre, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &genre
}
