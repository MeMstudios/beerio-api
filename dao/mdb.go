//Database connection
package dao

import (
	"log"

	. "apiserver/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION      = "beers"
	USER_COLLECTION = "users"
)

//Database functions
func (m *DAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *DAO) FindAll() ([]Beer, error) {
	var beers []Beer
	err := db.C(COLLECTION).Find(bson.M{}).All(&beers)
	return beers, err
}

func (m *DAO) FindById(id string) (Beer, error) {
	var beer Beer
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&beer)
	return beer, err
}

func (m *DAO) Insert(beer Beer) error {
	err := db.C(COLLECTION).Insert(&beer)
	return err
}

func (m *DAO) Delete(beer Beer) error {
	err := db.C(COLLECTION).Remove(&beer)
	return err
}

func (m *DAO) Update(beer Beer) error {
	err := db.C(COLLECTION).UpdateId(beer.ID, &beer)
	return err
}

//User Functions
func (m *DAO) InsertUser(user User) error {
	err := db.C(USER_COLLECTION).Insert(&user)
	return err
}

func (m *DAO) FindAllUsers() ([]User, error) {
	var users []User
	err := db.C(USER_COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

func (m *DAO) FindUserById(id string) (User, error) {
	var user User
	err := db.C(USER_COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func (m *DAO) FindUserByEmail(email string) (User, error) {
	var user User
	err := db.C(USER_COLLECTION).Find(bson.M{"email": email}).One(&user)
	return user, err
}

func (m *DAO) AddFavorite(fav Favorite) error {
	//find the user from the favorite id
	if user, err := m.FindUserById(fav.UserID); err != nil {
		return err
	} else {
		//add to the user's favorites and update the user
		user.Favorites = append(user.Favorites, fav)
		err := db.C(USER_COLLECTION).Update(fav.UserID, &user)
		return err
	}
}

func (m *DAO) GetFavorites(userId string) ([]Favorite, error) {
	//find the user by their id passed from a cookie and return their favorites
	if user, err := m.FindUserById(userId); err != nil {
		return nil, err
	} else {
		return user.Favorites, err
	}
}
