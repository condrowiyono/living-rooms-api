package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OMDBStruct struct {
	Title      string
	Year       string
	Rated      string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actor      string
	Plot       string
	Language   string
	Country    string
	Awards     string
	Poster     string
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string
	DVD        string
	BoxOffice  string
	Production string
	Website    string
}

type TheMovieDB struct {
	Adult               bool        `json:"adult"`
	BackdropPath        string      `json:"backdrop_path"`
	BelongsToCollection interface{} `json:"belongs_to_collection"`
	Budget              int         `json:"budget"`
	Genres              []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string  `json:"homepage"`
	ID                  int     `json:"id"`
	ImdbID              string  `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Overview            string  `json:"overview"`
	Popularity          float64 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso3166_1 string `json:"iso_3166_1"`
		Name      string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		Iso639_1 string `json:"iso_639_1"`
		Name     string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

func GetMovieDetail(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	tmdbID := string(vars.Get("tmdb"))
	resp, err := http.Get("https://api.themoviedb.org/3/movie/496243?api_key=b9436778204535120f245d0f0457439e")
	if err != nil {
		return
	}
	fmt.Printf("%v", tmdbID)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var theMovieDB TheMovieDB
	err = json.Unmarshal(body, &theMovieDB)
	fmt.Printf("%v", theMovieDB)
	respondJSON(w, http.StatusOK, nil, theMovieDB)
}
