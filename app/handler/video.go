package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

type VideoResult struct {
	Type     string `json:"type"`
	VideoURL string `json:"video_url"`
	Title    string `json:"title"`
	ID       string `json:"ID"`
}

func GetVideo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	id := string(vars.Get("id"))
	showType := string(vars.Get("type"))

	idInt, _ := strconv.ParseInt(id, 10, 64)

	video := getVideoOr404(db, idInt, showType, w, r)
	if video == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, video)
}

// getVideoOr404 gets a instance if exists, or respond the 404 error otherwise
func getVideoOr404(db *gorm.DB, id int64, showType string, w http.ResponseWriter, r *http.Request) *VideoResult {
	result := VideoResult{}
	queryString := fmt.Sprintf(
		`SELECT videos.type, videos.video_url, %s.title, %s.id
		FROM videos
		LEFT JOIN %s on (videos.show_id = %s.id AND videos.show_type = '%s')
		WHERE (videos.id = %d)`, showType, showType, showType, showType, showType, id)

	query := db.Raw(queryString).Scan(&result)

	if err := query.Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &result
}
