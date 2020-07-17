package scrapper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// TMDBSearch type
type TMDBSearch struct {
	Results []struct {
		PosterPath  string `json:"poster_path"`
		BannerPath  string `json:"backdrop_path"`
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Name        string `json:"name,omitempty"`
		ReleaseDate string `json:"release_date,omitempty"`
	} `json:"results"`
}

// SearchMovie Method
func SearchMovie(w http.ResponseWriter, r *http.Request) {
	var tmdbSearch TMDBSearch

	query := getHTTPRequestQuery(r, "query")
	typeName := getHTTPRequestQuery(r, "type") // movie or tv

	if len(typeName) <= 0 {
		typeName = "movie"
	}

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/search/%s?api_key=%s&query=%s",
		typeName,
		os.Getenv("TMDB_KEY"),
		url.PathEscape(query),
	)

	body, _ := getHTTPRequestGetBody(url)
	err = json.Unmarshal(body, &tmdbSearch)
	respondJSON(w, http.StatusOK, nil, tmdbSearch.Results)
}
