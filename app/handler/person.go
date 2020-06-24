package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/condrowiyono/ruangtengah-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllPerson(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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
	query.Model(&model.Person{}).Count(&count)

	person := []model.Person{}
	query = db.Limit(limitInt).
		Offset(offsetInt).
		Joins("LEFT JOIN persons_pictures on persons_pictures.person_id=people.id").
		Joins("LEFT JOIN images on persons_pictures.image_id=images.id").
		Group("people.id").
		Preload("Pictures").
		Find(&person)

	if err := query.Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write Response
	meta := Meta{limitInt, offsetInt, pageInt, count}
	respondJSON(w, http.StatusOK, meta, person)
}

func CreatePerson(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	person := model.Person{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&person); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Create(&person).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, nil, person)
}

func GetPerson(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	person := getPersonOr404(db, id, w, r)
	if person == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, person)
}

func UpdatePerson(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	person := getPersonOr404(db, id, w, r)
	if person == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&person); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&person).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, nil, person)
}

func DeletePerson(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	person := getPersonOr404(db, id, w, r)
	if person == nil {
		return
	}
	if err := db.Delete(&person).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil, nil)
}

// getPersonOr404 gets a instance if exists, or respond the 404 error otherwise
func getPersonOr404(db *gorm.DB, id int64, w http.ResponseWriter, r *http.Request) *model.Person {
	person := model.Person{}
	if err := db.Preload("Pictures").First(&person, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &person
}
