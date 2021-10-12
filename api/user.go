package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	gorm.Model
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

var SECRET_KEY = []byte("gosecretkey")

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

//set error message in Error struct
func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}

//compare plain password with hash password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(email string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email

	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		fmt.Errorf("%s", err.Error())
		return "", err
	}

	return tokenString, nil
}

//check whether user is authorized or not
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header["Authorization"]
		fmt.Println("auth", auth)
		if len(auth) == 0 {
			var err Error
			fmt.Println("notoken")
			err = SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}
		tokens := strings.Split(auth[0], " ")
		var mySigningKey = []byte(SECRET_KEY)
		fmt.Println("tokens", tokens)
		tokenString := tokens[len(tokens)-1]
		fmt.Println("tokenstring", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("parse", token)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing token")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err Error
			err = SetError(err, "Your Token has been expired.")
			json.NewEncoder(w).Encode(err)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			handler.ServeHTTP(w, r)
			return

		}

		var reserr Error
		reserr = SetError(reserr, "Not Authorized.")
		json.NewEncoder(w).Encode(reserr)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {

	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	var user User

	//get request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, " enter data to create user")
	}

	json.Unmarshal(reqBody, &user)
	user.Password = getHash([]byte(user.Password))

	db.Create(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var authuser Authentication
	var dbUser User
	json.NewDecoder(r.Body).Decode(&authuser)
	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()
	err := db.Where("email=?", authuser.Email).First(&dbUser).Error
	fmt.Print(dbUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	check := CheckPasswordHash(authuser.Password, dbUser.Password)

	if !check {
		var err Error
		err = SetError(err, "Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := GenerateJWT(authuser.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	var token Token
	token.Email = authuser.Email

	token.TokenString = validToken
	db.Create(&token)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(token)

}

func Logout(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("TOKEN")
	var dbtoken Token
	//db.Where("email = ?",)
	var db *gorm.DB = openDataBase()
	defer db.Close()
	err := db.Where("token = ?", token).Delete(&dbtoken).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"message":"loggedout "}`))

}

func GetUserByMail(email string) User {
	var dbUser User
	var db *gorm.DB = openDataBase()
	defer db.Close()
	db.Where("email = ?", email).First(&dbUser)
	return dbUser

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get user")
	email := mux.Vars(r)["email"]
	fmt.Println(email)
	user := GetUserByMail(email)
	fmt.Println("get user", user)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)

}
