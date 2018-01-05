package main

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"sync"
	"github.com/night-codes/mgo-wrapper"
	"time"
)


type (
	UsersStruct struct {
		ID     uint64 `form:"id" json:"id" bson:"_id"`
		Login  string `form:"login" json:"login" bson:"login" valid:"required,min(3)"`
		Name   string `form:"name" json:"name" bson:"name" valid:"max(200)"`
		Notice string `form:"notice" json:"notice" bson:"notice" valid:"max(1000)"`

		// Is Root-user? Similar as Rights.Groups = ["root"]
		Root bool `form:"-" json:"-" bson:"root"`

		// Information field, if needs auth by email set Login == Email
		Email string `form:"email" json:"email" bson:"email" valid:"email"`

		// sha512 hash of password (but from form can be received string password value)
		Password string `form:"password" json:"-" bson:"password" valid:"min(5)"`

		// from form can be received string password value)
		Password2 string `form:"password2" json:"-" bson:"password2"`

		// Times of creating or editing (or loading from mongoDB)
		Created int64 `form:"-" json:"created" bson:"created"`
		Updated int64 `form:"-" json:"updated" bson:"updated"`
		Loaded  int64 `form:"-" json:"-" bson:"-"`

		// Fields for users auth limitation
		Disabled bool `form:"-" json:"disabled" bson:"disabled"`
		Deleted  bool `form:"-" json:"deleted" bson:"deleted"`

		// IP control fields (coming soon)
		LastIP   uint32 `form:"-" json:"lastIP" bson:"lastIP"`
		IP       uint32 `form:"-" json:"-" bson:"ip"`
		StringIP string `form:"-" json:"ip" bson:"-"`

		// custom data map
		Settings map[string]interface{} `form:"-" json:"settings" bson:"settings"`

		// user without authentication
		Demo bool `form:"-" json:"demo" bson:"-"`
	}
	Users struct {
		rawList    map[string]*bson.Raw    // key - login
		rawListID  map[uint64]*bson.Raw    // key - id
		list       map[string]*UsersStruct // key - login
		listID     map[uint64]*UsersStruct // key - id
		count      int
		collection *mgo.Collection
		sync.Mutex
		mutUsers sync.Mutex
	}
)

func (u *Users) init() {
	u.Mutex = sync.Mutex{}
	u.collection = mongo.DB("Users").C("UsersCollection"/*Edit*/)
	u.rawList = map[string]*bson.Raw{}
	u.rawListID = map[uint64]*bson.Raw{}
	u.list = map[string]*UsersStruct{}
	u.listID = map[uint64]*UsersStruct{}
	u.count, _ = u.collection.Count()

	go func() {
		for range time.Tick(time.Second * 10) {
			u.count, _ = u.collection.Count()
			u.loadUsers()
			u.clearUsers()
		}
	}()
}

//type user struct {
//	Username string `json:"username"`
//	Password string `json:"-"`
//}

// For this demo, we're storing the user list in memory
// We also have some users predefined.
// In a real application, this list will most likely be fetched
// from a database. Moreover, in production settings, you should
// store passwords securely by salting and hashing them instead
// of using them as we're doing in this demo
//var userList = []user{
//	user{Username: "user1", Password: "pass1"},
//	user{Username: "user2", Password: "pass2"},
//	user{Username: "user3", Password: "pass3"},
//}

// Check if the username and password combination is valid
func (u *UsersStruct)isUserValid(username, password string) bool {
	for _, u := range u.listID {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

// Register a new user with the given username and password
// NOTE: For this demo, we
func registerNewUser(username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	u := user{Username: username, Password: password}

	userList = append(userList, u)

	return &u, nil
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	for _, u := range userList {
		if u.Username == username {
			return false
		}
	}
	return true
}