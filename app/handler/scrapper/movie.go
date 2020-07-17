package scrapper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TMDBMovie struct {
	Adult               bool                       `json:"adult"`
	BackdropPath        string                     `json:"backdrop_path"`
	Genres              []MovieGenres              `json:"genres"`
	ID                  int                        `json:"id"`
	ImdbID              string                     `json:"imdb_id"`
	Overview            string                     `json:"overview"`
	PosterPath          string                     `json:"poster_path"`
	ProductionCompanies []MovieProductionCompanies `json:"production_companies"`
	ProductionCountries []MovieProductionCountries `json:"production_countries"`
	ReleaseDate         string                     `json:"release_date"`
	Runtime             int                        `json:"runtime"`
	Title               string                     `json:"title"`
	Credits             Credits                    `json:"credits"`
	Videos              Videos                     `json:"videos"`
}
type MovieGenres struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type MovieProductionCompanies struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}
type MovieProductionCountries struct {
	Iso31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}
type MovieCast struct {
	CastID      int    `json:"cast_id"`
	Character   string `json:"character"`
	CreditID    string `json:"credit_id"`
	Gender      int    `json:"gender"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Order       int    `json:"order"`
	ProfilePath string `json:"profile_path"`
}
type MovieCrew struct {
	CreditID    string `json:"credit_id"`
	Department  string `json:"department"`
	Gender      int    `json:"gender"`
	ID          int    `json:"id"`
	Job         string `json:"job"`
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
}
type Credits struct {
	Cast []MovieCast `json:"cast"`
	Crew []MovieCrew `json:"crew"`
}
type Results struct {
	ID       string `json:"id"`
	Iso6391  string `json:"iso_639_1"`
	Iso31661 string `json:"iso_3166_1"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Site     string `json:"site"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
}
type Videos struct {
	Results []Results `json:"results"`
}

func GetMovieDetail(w http.ResponseWriter, r *http.Request) {
	var tmdbMovie TMDBMovie

	tmdbID := getHTTPRequestQuery(r, "tmdb")

	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s&append_to_response=credits,videos", tmdbID, os.Getenv("TMDB_KEY"))

	body, _ := getHTTPRequestGetBody(apiURL)
	err = json.Unmarshal(body, &tmdbMovie)

	// limit cast to only first 10
	if len(tmdbMovie.Credits.Cast) > 10 {
		tmdbMovie.Credits.Cast = tmdbMovie.Credits.Cast[0:10]
	}

	// filter only director
	var movieCrew []MovieCrew = nil

	for _, crew := range tmdbMovie.Credits.Crew {
		if crew.Job == "Director" || crew.Job == "Writer" {
			movieCrew = append(movieCrew, crew)
		}
	}

	tmdbMovie.Credits.Crew = movieCrew

	respondJSON(w, http.StatusOK, nil, tmdbMovie)
}
