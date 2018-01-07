package main

import (
	"errors"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"sync"

	"github.com/night-codes/govalidator"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

type Users struct {
	rawList    map[string]*bson.Raw    // key - login
	rawListID  map[uint64]*bson.Raw    // key - id
	list       map[string]*user // key - login
	listID     map[uint64]*user // key - id
	count      int
	collection *mgo.Collection
	sync.Mutex
	mutUsers   sync.Mutex
}

//For this demo, we're storing the user list in memory
//We also have some users predefined.
//In a real application, this list will most likely be fetched
//from a database. Moreover, in production settings, you should
//store passwords securely by salting and hashing them instead
//of using them as we're doing in this demo
var userList = []user{
	user{Username: "user1", Password: "pass1"},
	user{Username: "user2", Password: "pass2"},
	user{Username: "user3", Password: "pass3"},
}

//func (u *Users) Validate(user *UsersStruct) error {
//	if _, err := govalidator.ValidateStruct(user); err != nil {
//		ers := []string{}
//		for k, v := range govalidator.ErrorsByField(err) {
//			ers = append(ers, k+": "+v)
//		}
//		return errors.New(strings.Join(ers, " \n"))
//	}
//	if user.Password != user.Password2 {
//		return errors.New("Password mismatch!")
//	}
//	user.Password2 = ""
//	return nil
//}

// Check if the username and password combination is valid
func (user *Users) isUserValid(username, password string) bool {

	//user.collection = mongo.DB("").C("sessions")

	for _, u := range userList {
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
