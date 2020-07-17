package scrapper

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

func GetMovieImage(w http.ResponseWriter, r *http.Request) {
	tmdbID := getHTTPRequestQuery(r, "tmdb")
	imageType := getHTTPRequestQuery(r, "type")

	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s/images?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))

	body, _ := getHTTPRequestGetBody(apiURL)

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

func GetTvImage(w http.ResponseWriter, r *http.Request) {
	tmdbID := getHTTPRequestQuery(r, "tmdb")
	imageType := getHTTPRequestQuery(r, "type")

	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s/images?api_key=%s", tmdbID, os.Getenv("TMDB_KEY"))

	body, _ := getHTTPRequestGetBody(apiURL)

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

func GetDuckDuckGoImage(w http.ResponseWriter, r *http.Request) {
	// Request the HTML page.
	query := getHTTPRequestQuery(r, "query")

	duckduckGoURL := fmt.Sprintf("https://duckduckgo.com/?q=%s&iar=images&iax=images&ia=images", url.QueryEscape(query))
	res, _ := http.Get(duckduckGoURL)
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
	query := getHTTPRequestQuery(r, "query")
	googleURL := fmt.Sprintf("https://www.google.co.id/search?q=%s&source=lnms&tbm=isch", url.QueryEscape(query))

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
