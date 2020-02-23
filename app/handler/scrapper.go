package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/condrowiyono/living-rooms-api/app/model"
)

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

type TMDBSearch struct {
	Results []struct {
		Popularity       float64 `json:"popularity"`
		VoteCount        int     `json:"vote_count"`
		Video            bool    `json:"video"`
		PosterPath       string  `json:"poster_path"`
		ID               int     `json:"id"`
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		GenreIds         []int   `json:"genre_ids"`
		Title            string  `json:"title"`
		VoteAverage      float64 `json:"vote_average"`
		Overview         string  `json:"overview"`
		ReleaseDate      string  `json:"release_date,omitempty"`
	} `json:"results"`
}

type DuckDuckGoImageResult struct {
	Results []struct {
		Height    int    `json:"height"`
		URL       string `json:"url"`
		Width     int    `json:"width"`
		Source    string `json:"source"`
		Title     string `json:"title"`
		Thumbnail string `json:"thumbnail"`
		Image     string `json:"image"`
	} `json:"results"`
}

func GetMovieDetail(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	tmdbID := string(vars.Get("tmdb"))
	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))
	resp, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var theMovieDB TheMovieDB
	err = json.Unmarshal(body, &theMovieDB)
	releasedDate, _ := time.Parse(theMovieDB.ReleaseDate, "2020-02-02")
	releasedYear := releasedDate.Year()

	show := model.Show{
		Title:     theMovieDB.Title,
		Actors:    "",
		Awards:    "",
		BoxOffice: "",
		Country:   "",
		DVD:       "",
		Director:  "",
		Genre:     "",
		Language:  theMovieDB.OriginalLanguage,
		Plot:      theMovieDB.Overview,
		Poster:    fmt.Sprintf("https://image.tmdb.org/t/p/w370_and_h556_bestv2%s", theMovieDB.PosterPath),
		// Production: map[interface]string{}theMovieDB.ProductionCompanies,
		Rated:      "",
		Released:   theMovieDB.ReleaseDate,
		Response:   "",
		Runtime:    strconv.Itoa(theMovieDB.Runtime),
		Type:       "movie",
		Website:    theMovieDB.Homepage,
		Writer:     "",
		Year:       strconv.Itoa(releasedYear),
		ImdbID:     theMovieDB.ImdbID,
		ImdbRating: "",
		ImdbVotes:  "",
		Banner:     fmt.Sprintf("https://image.tmdb.org/t/p/original%s", theMovieDB.BackdropPath),
		Trailer:    "",
		PlayerURL:  "",
	}
	respondJSON(w, http.StatusOK, nil, show)
}

func SearchMovie(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	query := string(vars.Get("query"))
	queryEscaped := url.PathEscape(query)

	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", os.Getenv("TMDB_KEY"), queryEscaped)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var tmdbSearch TMDBSearch
	err = json.Unmarshal(body, &tmdbSearch)
	respondJSON(w, http.StatusOK, nil, tmdbSearch.Results)
}

func GetDuckDuckGoImage(w http.ResponseWriter, r *http.Request) {
	// Request the HTML page.
	vars := r.URL.Query()
	query := string(vars.Get("query"))
	queryEscaped := url.QueryEscape(query)
	duckduckGoURL := fmt.Sprintf("https://duckduckgo.com/?q=%s&iar=images&iax=images&ia=images", queryEscaped)
	res, err := http.Get(duckduckGoURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var vqd string
	var q string
	// Load the HTML document
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	// Find VQD and Query
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		s.Find("script").Each(func(i int, sc *goquery.Selection) {
			js := sc.Text()
			regexVQD := regexp.MustCompile(`vqd=(.*?)&`)
			vqdFound := regexVQD.FindString(js)
			if len(vqdFound) > 0 {
				vqd = string(vqdFound[4:len(vqdFound)])
				vqd = string(vqd[0 : len(vqd)-1])
			}
			regeqQ := regexp.MustCompile(`q=(.*?)&`)
			qFound := regeqQ.FindString(js)
			if len(qFound) > 0 {
				q = string(qFound[2:len(qFound)])
				q = string(q[0 : len(q)-1])
			}
		})
	})

	// Lets find DuckduckGO Image
	duckduckGoJSONImage := fmt.Sprintf("https://duckduckgo.com/i.js?l=us-en&o=json&q=%s&vqd=%s&f=,,,&p=1&v7exp=a", q, vqd)
	resp, err := http.Get(duckduckGoJSONImage)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var duckDuckGoImageResult DuckDuckGoImageResult
	err = json.Unmarshal(body, &duckDuckGoImageResult)
	respondJSON(w, http.StatusOK, nil, duckDuckGoImageResult.Results)
}
