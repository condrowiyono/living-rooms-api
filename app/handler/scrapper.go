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
	"time"

	"github.com/PuerkitoBio/goquery"
)

type CommonImageResult struct {
	Thumbnail string `json:"thumbnail"`
	Image     string `json:"image"`
}

type TheMovieDBCredits struct {
	Cast []TheMovieDBPerson `json:"cast"`
	Crew []TheMovieDBPerson `json:"crew"`
}

type TheMovieDBPerson struct {
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
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
	Status      string             `json:"status"`
	Tagline     string             `json:"tagline"`
	Title       string             `json:"title"`
	Video       bool               `json:"video"`
	VoteAverage float64            `json:"vote_average"`
	VoteCount   int                `json:"vote_count"`
	Actors      []TheMovieDBPerson `json:"actors"`
	Crews       []TheMovieDBPerson `json:"crews"`
	Videos      []TheMovieDBVideo  `json:"videos"`
}

type TheMovieDBVideos struct {
	Results []TheMovieDBVideo `json:"results"`
}
type TheMovieDBVideo struct {
	Key  string `json:"key"`
	Site string `json:"site"`
	Type string `json:"type"`
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
	Results []DuckDuckGoImageResultData `json:"results"`
}

type DuckDuckGoImageResultData struct {
	Height    int    `json:"height"`
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Source    string `json:"source"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Image     string `json:"image"`
}

type TMDBMovieImage struct {
	ID        int `json:"id"`
	Backdrops []struct {
		AspectRatio float64     `json:"aspect_ratio"`
		FilePath    string      `json:"file_path"`
		Height      int         `json:"height"`
		Iso6391     interface{} `json:"iso_639_1"`
		VoteAverage float64     `json:"vote_average"`
		VoteCount   int         `json:"vote_count"`
		Width       int         `json:"width"`
	} `json:"backdrops"`
	Posters []struct {
		AspectRatio float64 `json:"aspect_ratio"`
		FilePath    string  `json:"file_path"`
		Height      int     `json:"height"`
		Iso6391     string  `json:"iso_639_1"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		Width       int     `json:"width"`
	} `json:"posters"`
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

	creditURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s/credits?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))
	resp, err = http.Get(creditURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	var theMovieDBCredits TheMovieDBCredits
	err = json.Unmarshal(body, &theMovieDBCredits)

	theMovieDB.Actors = theMovieDBCredits.Cast

	videosURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s/videos?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))
	resp, err = http.Get(videosURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	var theMovieDBVideos TheMovieDBVideos
	err = json.Unmarshal(body, &theMovieDBVideos)

	theMovieDB.Videos = theMovieDBVideos.Results
	respondJSON(w, http.StatusOK, nil, theMovieDB)
}

func GetMovieImage(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	tmdbID := string(vars.Get("tmdb"))
	imageType := string(vars.Get("type"))

	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s/images?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))
	resp, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var tmdbMovieImage TMDBMovieImage
	err = json.Unmarshal(body, &tmdbMovieImage)

	var commonImageResultArray []CommonImageResult

	if imageType == "banners" {
		for _, v := range tmdbMovieImage.Backdrops {
			commonImageResult := CommonImageResult{
				Thumbnail: fmt.Sprintf("https://image.tmdb.org/t/p/w500_and_h282_face%s", v.FilePath),
				Image:     fmt.Sprintf("https://image.tmdb.org/t/p/original%s", v.FilePath),
			}
			commonImageResultArray = append(commonImageResultArray, commonImageResult)
		}
	} else if imageType == "posters" {
		for _, v := range tmdbMovieImage.Posters {
			commonImageResult := CommonImageResult{
				Thumbnail: fmt.Sprintf("https://image.tmdb.org/t/p/w220_and_h330_face%s", v.FilePath),
				Image:     fmt.Sprintf("https://image.tmdb.org/t/p/original%s", v.FilePath),
			}
			commonImageResultArray = append(commonImageResultArray, commonImageResult)
		}
	}

	respondJSON(w, http.StatusOK, nil, commonImageResultArray)
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

	var commonImageResultArray []CommonImageResult

	for _, v := range duckDuckGoImageResult.Results {
		commonImageResult := CommonImageResult{
			Thumbnail: v.Thumbnail,
			Image:     v.Image,
		}
		commonImageResultArray = append(commonImageResultArray, commonImageResult)
	}

	respondJSON(w, http.StatusOK, nil, commonImageResultArray)
}

func GetGoogleImage(w http.ResponseWriter, r *http.Request) {
	// Request the HTML page.
	vars := r.URL.Query()
	query := string(vars.Get("query"))
	queryEscaped := url.QueryEscape(query)
	googleURL := fmt.Sprintf("https://www.google.co.id/search?q=%s&source=lnms&tbm=isch", queryEscaped)
	req, err := http.NewRequest("GET", googleURL, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.116 Safari/537.36 Edg/80.0.361.57")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}

	var resultImage []string

	// Load the HTML document
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		s.Find("script").Each(func(i int, sc *goquery.Selection) {
			js := sc.Text()
			regexpImageURL := regexp.MustCompile(`(http(s?):)([/|.|\w|\s|-])*\.(?:jpg|gif|png)`)
			regexpImageURLFound := regexpImageURL.FindAllString(js, -1)
			if len(regexpImageURLFound) > 0 {
				resultImage = regexpImageURLFound
			}
		})
	})

	var commonImageResultArray []CommonImageResult

	for _, v := range resultImage {
		commonImageResult := CommonImageResult{
			Thumbnail: v,
			Image:     v,
		}
		commonImageResultArray = append(commonImageResultArray, commonImageResult)
	}

	respondJSON(w, http.StatusOK, nil, commonImageResultArray)
}
