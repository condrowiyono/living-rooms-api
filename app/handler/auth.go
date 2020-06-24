package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"os"

	"github.com/condrowiyono/ruangtengah-api/app/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var err error

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var request model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	user := model.User{}

	if err = db.Where("username = ?", request.Username).First(&user).Error; err != nil {
		respondError(w, http.StatusUnauthorized, "wrong username")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		respondError(w, http.StatusUnauthorized, "wrong password")
		return
	}
	ttl := 720 * time.Hour

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().UTC().Add(ttl).Unix(),
	})
	token, err := sign.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nil, map[string]interface{}{"token": token, "user": user})
}

func Register(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user = model.User{
		Email:    user.Email,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
		Password: string(hashedPassword),
	}

	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, nil, user)
}

func GetMyDetail(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	username := r.Header.Get("username")
	if err = db.Where("username = ?", username).First(&user).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, nil, user)
}
