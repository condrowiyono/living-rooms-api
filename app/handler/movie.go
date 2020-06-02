package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/condrowiyono/living-rooms-api/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllMovie(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	page := string(vars.Get("page"))
	limit := string(vars.Get("limit"))
	genre := string(vars.Get("genre"))
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
	movie := []model.Movie{}
	query := db.Model(model.Movie{})

	queryWhereIn := db.Model(model.Movie{}).Select("DISTINCT(movies.id)")

	if len(genre) != 0 {
		queryWhereIn = queryWhereIn.
			Joins("join movies_genres on movies_genres.movie_id = movies.id").
			Joins("join genres on genres.id = movies_genres.genre_id AND genres.name = ?", genre)
	}

	if len(title) != 0 {
		query = query.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}

	query = query.
		Where("id IN (?)", queryWhereIn.QueryExpr()).
		Limit(limitInt).
		Offset(offsetInt).
		Preload("Genres").
		Preload("Banners").
		Preload("Posters").
		Preload("Languages").
		Preload("Countries").
		Preload("Productions").
		Preload("Actors").
		Preload("Crews").
		Preload("Player").
		Preload("Videos")

	if err := query.Find(&movie).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var count int64
	query = query.Offset(0).
		Count(&count)

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
	countries := []model.Country{}
	productions := []model.Production{}

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
		db.Where(model.Country{Name: v.Name}).FirstOrCreate(&language)
		languages = append(languages, language)
	}

	for _, v := range movie.Countries {
		country := model.Country{}
		db.Where(model.Country{Name: v.Name}).FirstOrCreate(&country)
		countries = append(countries, country)
	}

	for _, v := range movie.Productions {
		production := model.Production{}
		db.Where(model.Production{Name: v.Name}).FirstOrCreate(&production)
		productions = append(productions, production)
	}

	db.Model(&movie).Association("Genres").Replace(genres)
	db.Model(&movie).Association("Actors").Replace(actors)
	db.Model(&movie).Association("Crews").Replace(crews)
	db.Model(&movie).Association("Languages").Replace(languages)
	db.Model(&movie).Association("Countries").Replace(countries)
	db.Model(&movie).Association("Productions").Replace(productions)

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
	countries := []model.Country{}
	productions := []model.Production{}

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
		db.Where(model.Country{Name: v.Name}).FirstOrCreate(&language)
		languages = append(languages, language)
	}

	for _, v := range movie.Countries {
		country := model.Country{}
		db.Where(model.Country{Name: v.Name}).FirstOrCreate(&country)
		countries = append(countries, country)
	}

	for _, v := range movie.Productions {
		production := model.Production{}
		db.Where(model.Production{Name: v.Name}).FirstOrCreate(&production)
		productions = append(productions, production)
	}

	db.Model(&movie).Association("Genres").Replace(genres)
	db.Model(&movie).Association("Actors").Replace(actors)
	db.Model(&movie).Association("Crews").Replace(crews)
	db.Model(&movie).Association("Languages").Replace(languages)
	db.Model(&movie).Association("Country").Replace(countries)
	db.Model(&movie).Association("Productions").Replace(productions)
	db.Model(&movie).Association("Videos").Replace(movie.Videos)
	db.Model(&movie).Association("Player").Replace(movie.Player)

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
		Preload("Countries").
		Preload("Productions").
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
