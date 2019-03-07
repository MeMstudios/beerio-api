package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	. "beerio-api/config"
	. "beerio-api/dao"
	. "beerio-api/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var config = Config{}
var dao = DAO{}

// browser imposed preflight OPTIONS methods
func BeerPreFlightResponse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	respondWithJson(w, http.StatusOK, "Do anything with beers!")
}
func PostOnlyPreFlightResponse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	respondWithJson(w, http.StatusOK, "posts allowed")
}
func GetOnlyPreFlightResponse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	respondWithJson(w, http.StatusOK, "gets allowed")
}

// GET list of beers
func AllBeersEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	beers, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, beers)
}

// GET a beer by its ID
func FindBeerEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	params := mux.Vars(r)
	beer, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid beer ID")
		return
	}
	respondWithJson(w, http.StatusOK, beer)
}

// POST a new beer
func CreateBeerEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	defer r.Body.Close()
	var beer Beer
	if err := json.NewDecoder(r.Body).Decode(&beer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	beer.ID = bson.NewObjectId()
	if err := dao.Insert(beer); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, beer)
}

// PUT update an existing beer
func UpdateBeerEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	var beer Beer
	if err := json.NewDecoder(r.Body).Decode(&beer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(beer); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing beer
func DeleteBeerEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	var beer Beer
	if err := json.NewDecoder(r.Body).Decode(&beer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(beer); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// POST A new user
func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.ID = bson.NewObjectId()
	//encrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error encrypting password")
		log.Fatal(err)
		return
	}
	user.Password = string(hash)
	if err := dao.InsertUser(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, user.Name)
}

//POST A user password check
func LoginEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	var check UserCheck
	if err := json.NewDecoder(r.Body).Decode(&check); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	var email = check.Email
	var password = check.Password

	//find user by email
	if user, err := dao.FindUserByEmail(email); err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't find user email.")
	} else {

		//decrypt the password
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Wrong password!")
			return
		}

		//set a logged in cookie
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "logged_in", Value: user.ID.String(), Expires: expiration}
		http.SetCookie(w, &cookie)
		respondWithJson(w, http.StatusOK, user.ID)

	}

}

//POST a favorite beer
func FavoriteEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	var fav Favorite
	if err := json.NewDecoder(r.Body).Decode(&fav); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	//check if you're logged in with helper function
	if checkLoggedIn(fav.UserID, r) {
		//add favorite beer to the database
		if err := dao.AddFavorite(fav); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	}

}

//GET all favorites
func GetFavoritesEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", FrontendURL)
	params := mux.Vars(r)
	if checkLoggedIn(params["id"], r) {
		favorites, err := dao.GetFavorites(params["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error Getting Favorites!")
			return
		}
		respondWithJson(w, http.StatusOK, favorites)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//Check for logged in status with the cookie
func checkLoggedIn(ID string, r *http.Request) bool {
	var loggedIn = false
	cookie, err := r.Cookie("logged_in")
	if err != nil {
		if cookie.Value == ID {
			return true
		}
	}
	return loggedIn
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}
func main() {

	//fire up the router
	r := mux.NewRouter()

	//Beer endpoints
	r.HandleFunc("/beers", AllBeersEndPoint).Methods("GET")
	r.HandleFunc("/beers", CreateBeerEndPoint).Methods("POST")
	r.HandleFunc("/beers", UpdateBeerEndPoint).Methods("PUT")
	r.HandleFunc("/beers", DeleteBeerEndPoint).Methods("DELETE")
	r.HandleFunc("/beers/{id}", FindBeerEndpoint).Methods("GET")

	//User endpoints
	r.HandleFunc("/register", CreateUserEndPoint).Methods("POST")
	r.HandleFunc("/login", LoginEndPoint).Methods("POST")
	r.HandleFunc("/addFavorite", FavoriteEndPoint).Methods("POST")
	r.HandleFunc("/favorites", GetFavoritesEndpoint).Methods("GET")

	//PreFlight Responses
	r.HandleFunc("/beers", BeerPreFlightResponse).Methods("OPTIONS")
	r.HandleFunc("/login", PostOnlyPreFlightResponse).Methods("OPTIONS")
	r.HandleFunc("/register", PostOnlyPreFlightResponse).Methods("OPTIONS")
	r.HandleFunc("/addFavorite", PostOnlyPreFlightResponse).Methods("OPTIONS")
	r.HandleFunc("/favorites", GetOnlyPreFlightResponse).Methods("OPTIONS")
	//Release the hounds
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
