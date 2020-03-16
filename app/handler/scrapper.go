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
	ProfilePath string `json:"-"`
}

type TheMovieDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TheMovieDBProduction struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type TheMovieDBProductionCountry struct {
	Iso3166_1 string `json:"iso_3166_1"`
	Name      string `json:"name"`
}

type TheMovieDBLanguageCountry struct {
	Iso639_1 string `json:"iso_639_1"`
	Name     string `json:"name"`
}

type TheMovieDBVideos struct {
	Results []TheMovieDBVideo `json:"results"`
}

type TheMovieDBVideo struct {
	Key  string `json:"key"`
	Site string `json:"site"`
	Type string `json:"type"`
}

type OMDBAdditionalInfo struct {
	Director   string `json:"Director"`
	Awards     string `json:"Awards"`
	ImdbRating string `json:"imdbRating"`
	Rated      string `json:"Rated"`
}

type TheMovieDB struct {
	BackdropPath        string                        `json:"backdrop_path"`
	Genres              []TheMovieDBGenre             `json:"genres"`
	ID                  int                           `json:"id"`
	ImdbID              string                        `json:"imdb_id"`
	OriginalLanguage    string                        `json:"original_language"`
	Overview            string                        `json:"overview"`
	PosterPath          string                        `json:"poster_path"`
	ProductionCompanies []TheMovieDBProduction        `json:"production_companies"`
	ProductionCountries []TheMovieDBProductionCountry `json:"production_countries"`
	ReleaseDate         string                        `json:"release_date"`
	Runtime             int                           `json:"runtime"`
	SpokenLanguages     []TheMovieDBLanguageCountry   `json:"spoken_languages"`
	Title               string                        `json:"title"`
	Video               bool                          `json:"video"`
	VoteAverage         float64                       `json:"vote_average"`
	Director            string                        `json:"director"`
	Awards              string                        `json:"awards"`
	ImdbRating          float64                       `json:"imdb_rating"`
	Rated               string                        `json:"rated"`
	Actors              []TheMovieDBPerson            `json:"actors"`
	Crews               []TheMovieDBPerson            `json:"crews"`
	Videos              []TheMovieDBVideo             `json:"videos"`
}

type TMDBSearch struct {
	Results []struct {
		PosterPath  string  `json:"poster_path"`
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		VoteAverage float64 `json:"vote_average"`
		Overview    string  `json:"overview"`
		ReleaseDate string  `json:"release_date,omitempty"`
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
		FilePath string `json:"file_path"`
		Height   int    `json:"height"`
		Width    int    `json:"width"`
	} `json:"backdrops"`
	Posters []struct {
		FilePath string `json:"file_path"`
		Height   int    `json:"height"`
		Width    int    `json:"width"`
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

	theMovieDB.Actors = theMovieDBCredits.Cast[0:10]

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

	omdbURL := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", theMovieDB.ImdbID, os.Getenv("OMDB_KEY"))
	resp, err = http.Get(omdbURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	var omdb OMDBAdditionalInfo
	err = json.Unmarshal(body, &omdb)

	theMovieDB.Director = omdb.Director
	theMovieDB.Awards = omdb.Awards
	theMovieDB.Rated = omdb.Rated
	if s, err := strconv.ParseFloat(omdb.ImdbRating, 64); err == nil {
		theMovieDB.ImdbRating = s
	}

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
