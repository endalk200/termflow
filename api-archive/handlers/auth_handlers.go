package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/endalk200/termflow-api/data"
	jwtutils "github.com/endalk200/termflow-api/utils"
)

var SECRET_KEY = []byte("your-secret-key")

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := data.GetUsers()

	responseData, err := json.Marshal(users)
	if err != nil {
		log.Printf("Error during marshalling users slice")
		http.Error(
			w,
			"Something went wrong while handling the request",
			http.StatusInternalServerError,
		)
	}

	w.Write(responseData)
	return
}

type LoginRequestBodySchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequestBodySchema struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignupResponseBodySchema struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody SignupRequestBodySchema
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error while marshalling user request body %s", err)
		http.Error(
			w,
			"Something went wrong while trying to create user account",
			http.StatusBadRequest,
		)
		return
	}

	user, err := data.CreateUser(data.User{
		Id:        2,
		FirstName: requestBody.FirstName,
		LastName:  requestBody.LastName,
		Email:     requestBody.Email,
		Password:  requestBody.Password,
	})
	if err != nil {
		log.Printf("Error while creating user account with provided credentials")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).
			Encode(map[string]string{"message": fmt.Sprintf("%s", err)})
		return
	}

	responseBody, err := json.Marshal(SignupResponseBodySchema{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	})
	if err != nil {
		log.Printf("Something went wrong while marshalling response body %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
	return
}

type LoginResponseBodySchema struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody LoginRequestBodySchema
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	user, err := data.FindUserByEmail(requestBody.Email)
	if err != nil {
		log.Printf("There is no user with the specified email addresss: %s", requestBody.Email)
		http.Error(w, "Incorrect credentials provided", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)) != nil {
		log.Printf("Incorrect password provided for user with [email]: %s", requestBody.Email)
		http.Error(w, "Incorrect credentials provided", http.StatusUnauthorized)
		return
	}

	_userId := user.Id

	log.Printf("UserId: %v", _userId)
	claims := jwt.RegisteredClaims{
		Issuer:    "twoMatchesCorp",
		Subject:   "1",
		Audience:  jwt.ClaimStrings{"admin"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	accessToken, err := jwtutils.GenerateJWT(claims)
	if err != nil {
		log.Printf("Incorrect password provided for user with [email]: %s", requestBody.Email)
		http.Error(w, "Incorrect credentials provided", http.StatusUnauthorized)
		return
	}

	responseData, err := json.Marshal(LoginResponseBodySchema{
		Token:        accessToken,
		RefreshToken: "Refresh Token",
	})
	if err != nil {
		log.Printf("Error during marshalling users slice")
		http.Error(
			w,
			"Something went wrong while handling the request",
			http.StatusInternalServerError,
		)
	}

	w.Write(responseData)
	return
}

type GetSessionResponseBodySchema struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("middleware.auth.userID").(int)
	if !ok {
		log.Printf("Invalid user id provided as context")
		w.WriteHeader(http.StatusBadRequest)
	}

	user, err := data.FindUserById(userId)
	if err != nil {
		log.Printf("There is no user with the specified [Id]: %d", userId)
		http.Error(w, "Incorrect credentials provided", http.StatusUnauthorized)
		return
	}

	responseData, err := json.Marshal(GetSessionResponseBodySchema{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	})
	if err != nil {
		log.Printf("Error during marshalling users slice")
		http.Error(
			w,
			"Something went wrong while handling the request",
			http.StatusInternalServerError,
		)
	}

	w.Write(responseData)
	return
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("Something went wrong while parsing [id] path value from URL")
		http.Error(w, "Provide [id] path value", http.StatusBadRequest)
	}

	users, err := data.FindUserById(id)
	if err != nil {
		log.Printf("Something wen't wrong while trying to find a user")
		http.Error(w, "", http.StatusNotFound)
		return
	}

	responseData, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("Error during marshalling users slice")
	}

	w.Write(responseData)
	return
}
