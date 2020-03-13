package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

type PlayerResult struct {
	Type      string `json:"type"`
	PlayerURL string `json:"player_url"`
	Title     string `json:"title"`
	ID        string `json:"ID"`
}

func GetPlayer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	id := string(vars.Get("id"))
	playerType := string(vars.Get("type"))

	idInt, _ := strconv.ParseInt(id, 10, 64)

	player := getPlayerOr404(db, idInt, playerType, w, r)
	if player == nil {
		return
	}
	respondJSON(w, http.StatusOK, nil, player)
}

// getPlayerOr404 gets a instance if exists, or respond the 404 error otherwise
func getPlayerOr404(db *gorm.DB, id int64, playerType string, w http.ResponseWriter, r *http.Request) *PlayerResult {
	result := PlayerResult{}

	queryString := fmt.Sprintf(
		`SELECT players.type, players.player_url, %s.title, %s.id
		FROM players
		LEFT JOIN %s on (players.show_id = %s.id AND players.show_type = '%s')
		WHERE (players.id = %d)`, playerType, playerType, playerType, playerType, playerType, id)

	query := db.Raw(queryString).Scan(&result)
	if err := query.Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &result
}
