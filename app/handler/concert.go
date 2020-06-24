package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/condrowiyono/ruangtengah-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllConcert(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	page := string(vars.Get("page"))
	limit := string(vars.Get("limit"))
	title := string(vars.Get("title"))

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 25
	}

	offsetInt := (pageInt - 1) * limitInt
	concert := []model.Concert{}
	query := db.Model(model.Concert{})

	queryWhereIn := db.Model(model.Concert{}).Select("DISTINCT(concerts.id)")

	if len(title) != 0 {
		query = query.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}

	query = query.
		Where("id IN (?)", queryWhereIn.QueryExpr()).
		Limit(limitInt).
		Offset(offsetInt).
		Preload("Banners").
		Preload("Posters").
		Preload("Artist").
		Preload("Player").
		Preload("Videos")

	if err := query.Find(&concert).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var count int64
	query = query.Offset(0).
		Count(&count)

	// Write Response
	meta := Meta{limitInt, offsetInt, pageInt, count}
	respondJSON(w, http.StatusOK, meta, concert)
}

func CreateConcert(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	concert := model.Concert{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&concert); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Set("gorm:association_autoupdate", false).Create(&concert).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, nil, concert)
}

func GetConcert(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	concert := getConcertOr404(db, id, w, r)
	if concert == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, concert)
}

func UpdateConcert(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	concert := getConcertOr404(db, id, w, r)
	if concert == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&concert); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Set("gorm:association_autoupdate", false).Save(&concert).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nil, concert)
}

func DeleteConcert(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	concert := getConcertOr404(db, id, w, r)
	if concert == nil {
		return
	}
	if err := db.Delete(&concert).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil, nil)
}

// getConcertOr404 gets a instance if exists, or respond the 404 error otherwise
func getConcertOr404(db *gorm.DB, id int64, w http.ResponseWriter, r *http.Request) *model.Concert {
	concert := model.Concert{}
	if err := db.
		Preload("Banners").
		Preload("Posters").
		Preload("Player").
		Preload("Videos").
		First(&concert, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &concert
}
