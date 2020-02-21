package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/condrowiyono/living-rooms-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllShow(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	show := []model.Show{}
	query := db.Limit(limitInt)
	query = query.Offset(offsetInt)

	if err := query.Find(&show).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Count all data
	var count int64
	query = query.Offset(0)
	query.Table("shows").Count(&count)

	// Write Response
	meta := Meta{limitInt, offsetInt, pageInt, count}
	respondJSON(w, http.StatusOK, meta, show)
}

func CreateShow(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	show := model.Show{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&show); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&show).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, nil, show)
}

func GetShow(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	show := getShowOr404(db, id, w, r)
	if show == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, show)
}

func UpdateShow(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	show := getShowOr404(db, id, w, r)
	if show == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&show); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&show).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, nil, show)
}

func DeleteShow(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	show := getShowOr404(db, id, w, r)
	if show == nil {
		return
	}
	if err := db.Delete(&show).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil, nil)
}

// getShowOr404 gets a instance if exists, or respond the 404 error otherwise
func getShowOr404(db *gorm.DB, id int64, w http.ResponseWriter, r *http.Request) *model.Show {
	show := model.Show{}
	if err := db.First(&show, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &show
}
