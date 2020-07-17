package scrapper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TMDBTv struct {
	BackdropPath        string                `json:"backdrop_path"`
	CreatedBy           []CreatedBy           `json:"created_by"`
	EpisodeRunTime      []int                 `json:"episode_run_time"`
	FirstAirDate        string                `json:"first_air_date"`
	Genres              []Genres              `json:"genres"`
	Name                string                `json:"name"`
	Networks            []Networks            `json:"networks"`
	NumberOfEpisodes    int                   `json:"number_of_episodes"`
	NumberOfSeasons     int                   `json:"number_of_seasons"`
	OriginCountry       []string              `json:"origin_country"`
	Overview            string                `json:"overview"`
	PosterPath          string                `json:"poster_path"`
	ProductionCompanies []ProductionCompanies `json:"production_companies"`
	Seasons             []Seasons             `json:"seasons"`
	Casts               []string              `json:"casts"`
	Videos              []Video               `json:"videos"`
}

type CreatedBy struct {
	ID          int    `json:"id"`
	CreditID    string `json:"credit_id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	ProfilePath string `json:"profile_path"`
}

type Genres struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Networks struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type ProductionCompanies struct {
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type Seasons struct {
	AirDate      string    `json:"air_date"`
	EpisodeCount int       `json:"episode_count"`
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Overview     string    `json:"overview"`
	Episode      []Episode `json:"episodes"`
	PosterPath   string    `json:"poster_path"`
	SeasonNumber int       `json:"season_number"`
}

type TMDBTvCast struct {
	Cast []Cast `json:"cast"`
}

type Cast struct {
	Name string `json:"name"`
}

type TMDBTvVideo struct {
	Results []Video `json:"results"`
}

type Video struct {
	Key  string `json:"key"`
	Site string `json:"site"`
	Type string `json:"type"`
}

type Episode struct {
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ID             int     `json:"id"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}

// GetTvDetail get tv detail from tmdb
func GetTvDetail(w http.ResponseWriter, r *http.Request) {
	var tmdbTv TMDBTv
	var tMDBTvCast TMDBTvCast
	var tMDBTvVideo TMDBTvVideo

	tmdbID := getHTTPRequestQuery(r, "tmdb")

	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))
	creditURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s/credits?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))
	videoURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s/videos?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))

	body, _ := getHTTPRequestGetBody(apiURL)
	err = json.Unmarshal(body, &tmdbTv)

	body, _ = getHTTPRequestGetBody(creditURL)
	err = json.Unmarshal(body, &tMDBTvCast)

	body, _ = getHTTPRequestGetBody(videoURL)
	err = json.Unmarshal(body, &tMDBTvVideo)

	for i := range tMDBTvCast.Cast {
		tmdbTv.Casts = append(tmdbTv.Casts, tMDBTvCast.Cast[i].Name)
	}

	tmdbTv.Videos = tMDBTvVideo.Results
	respondJSON(w, http.StatusOK, nil, tmdbTv)
}

// GetTvEpisode get tv detail from tmdb
func GetTvSeason(w http.ResponseWriter, r *http.Request) {
	var seasons Seasons

	tmdbID := getHTTPRequestQuery(r, "tmdb")
	seasonsNumber := getHTTPRequestQuery(r, "season")

	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s/season/%s?api_key=%s", tmdbID, seasonsNumber, os.Getenv("TMDB_KEY"))

	body, _ := getHTTPRequestGetBody(apiURL)
	err = json.Unmarshal(body, &seasons)

	respondJSON(w, http.StatusOK, nil, seasons)
}

// GetTvEpisode get tv detail from tmdb
func GetTvEpisode(w http.ResponseWriter, r *http.Request) {
	var episode Episode

	tmdbID := getHTTPRequestQuery(r, "tmdb")
	seasonsNumber := getHTTPRequestQuery(r, "season")
	episodeNumber := getHTTPRequestQuery(r, "episode")

	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s/season/%s/episode/%s?api_key=%s", tmdbID, seasonsNumber, episodeNumber, os.Getenv("TMDB_KEY"))

	body, _ := getHTTPRequestGetBody(apiURL)
	err = json.Unmarshal(body, &episode)

	respondJSON(w, http.StatusOK, nil, episode)
}
