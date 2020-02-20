package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/condrowiyono/living-rooms-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var err error

func GetAllSampleLink(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	sampleLink := []model.SampleLink{}
	query := db.Limit(limitInt)
	query = query.Offset(offsetInt)

	if err := query.Find(&sampleLink).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Count all data
	var count int64
	query = query.Offset(0)
	query.Table("sample_links").Count(&count)

	// Write Response
	meta := Meta{limitInt, offsetInt, pageInt, count}
	respondJSON(w, http.StatusOK, meta, sampleLink)
}

func CreateSampleLink(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	sampleLink := model.SampleLink{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sampleLink); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&sampleLink).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, nil, sampleLink)
}

func GetSampleLink(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	sampleLink := getSampleLinkOr404(db, id, w, r)
	if sampleLink == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, sampleLink)
}

func UpdateSampleLink(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	sampleLink := getSampleLinkOr404(db, id, w, r)
	if sampleLink == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sampleLink); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&sampleLink).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, nil, sampleLink)
}

func DeleteSampleLink(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	sampleLink := getSampleLinkOr404(db, id, w, r)
	if sampleLink == nil {
		return
	}
	if err := db.Delete(&sampleLink).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil, nil)
}

// getSampleLinkOr404 gets a instance if exists, or respond the 404 error otherwise
func getSampleLinkOr404(db *gorm.DB, id int64, w http.ResponseWriter, r *http.Request) *model.SampleLink {
	sampleLink := model.SampleLink{}
	if err := db.First(&sampleLink, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &sampleLink
}
