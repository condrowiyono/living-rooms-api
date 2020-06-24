package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/condrowiyono/ruangtengah-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllPlaylist(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	playlist := []model.Playlist{}
	query := db.Limit(limitInt)
	query = query.Offset(offsetInt)

	if err := query.Find(&playlist).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Count all data
	var count int64
	query = query.Offset(0)
	query.Model(&model.Playlist{}).Count(&count)

	// Write Response
	meta := Meta{limitInt, offsetInt, pageInt, count}
	respondJSON(w, http.StatusOK, meta, playlist)
}

func CreatePlaylist(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	playlist := model.Playlist{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&playlist); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&playlist).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, nil, playlist)
}

func GetPlaylist(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	playlist := getPlaylistOr404(db, id, w, r)
	db.Model(&playlist).Association("Players").Find(&playlist.Players)
	if playlist == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, playlist)
}

func UpdatePlaylist(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	playlist := getPlaylistOr404(db, id, w, r)
	if playlist == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&playlist); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&playlist).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, nil, playlist)
}

func DeletePlaylist(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	playlist := getPlaylistOr404(db, id, w, r)
	if playlist == nil {
		return
	}
	if err := db.Delete(&playlist).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil, nil)
}

// getPlaylistOr404 gets a instance if exists, or respond the 404 error otherwise
func getPlaylistOr404(db *gorm.DB, id int64, w http.ResponseWriter, r *http.Request) *model.Playlist {
	playlist := model.Playlist{}
	if err := db.First(&playlist, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &playlist
}
