package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/condrowiyono/living-rooms-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	image := []model.Image{}
	query := db.Limit(limitInt)
	query = query.Offset(offsetInt)

	if err := query.Find(&image).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Count all data
	var count int64
	query = query.Offset(0)
	query.Model(&model.Image{}).Count(&count)

	// Write Response
	meta := Meta{limitInt, offsetInt, pageInt, count}
	respondJSON(w, http.StatusOK, meta, image)
}

func CreateImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	image := model.Image{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&image); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&image).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, nil, image)
}

func GetImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	image := getImageOr404(db, id, w, r)
	if image == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, image)
}

func UpdateImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	image := getImageOr404(db, id, w, r)
	if image == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&image); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&image).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, nil, image)
}

func DeleteImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	image := getImageOr404(db, id, w, r)
	if image == nil {
		return
	}
	if err := db.Delete(&image).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil, nil)
}

func UploadImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)

	tempFile, err := ioutil.TempFile("uploads", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	//Save to DB
	image := model.Image{
		Path:    fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), tempFile.Name()),
		Source:  "upload",
		Keyword: "keyword",
	}

	if err := db.Save(&image).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nil, image)
}

// getImageOr404 gets a instance if exists, or respond the 404 error otherwise
func getImageOr404(db *gorm.DB, id int64, w http.ResponseWriter, r *http.Request) *model.Image {
	image := model.Image{}
	if err := db.First(&image, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &image
}
