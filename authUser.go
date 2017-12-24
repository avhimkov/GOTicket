package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/night-codes/summer.v1"
	"gopkg.in/night-codes/types.v1"
	"time"
	//"strings"
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type (
	AuthStruct struct {
		ID          uint64    `form:"id"  json:"id"  bson:"_id"`
		Name        string    `form:"name" json:"name" bson:"name" valid:"required"`
		Description string    `form:"description" json:"description" bson:"description"`
		Created     time.Time `form:"-" json:"created" bson:"created"`
		Updated     time.Time `form:"-" json:"updated" bson:"updated"`
		Deleted     bool      `form:"-" json:"deleted" bson:"deleted"`
	}
	AuthModule struct {
		summer.Module
	}
)

var (
	auth = panel.AddModule(
		&summer.ModuleSettings{
			Name:           "authUser",
			CollectionName: "Auth",
			Title:          "Auth in system",
			MenuOrder:      0,
			MenuTitle:      "Auth in system",
			Rights:         summer.Rights{Groups: []string{"all"}},

			Menu: panel.MainMenu,
		},
		&UsersModule{},
	)
)

//add User
func AddUser(u *summer.Users, c *gin.Context) {
	login, e1 := c.GetPostForm("user-z-login")
	password, e2 := c.GetPostForm("user-z-password")
	password2, e3 := c.GetPostForm("user-z-password-2")
	if e1 && e2 && e3 {
		if err := u.Users.Add(summer.UsersStruct{
			Login:     login,
			Password:  password,
			Password2: password2,
			Name:      strings.Title(login),
			Root:      true,
			//Rights:    summer.Rights{Groups: []string{"root"}, Actions: []string{"all"}},
			Rights:    summer.Rights{Groups: []string{"user"}, Actions: []string{"all"}},
			Settings:  obj{},
		}); err != nil {
			c.String(400, err.Error())
			return
		}
	}
}

// Add new record
func (m *AuthModule) Add(c *gin.Context) {
	var result AuthStruct
	if !summer.PostBind(c, &result) {
		return
	}
	result.ID = panel.AI.Next("authUser")
	result.Created = time.Now()
	result.Updated = time.Now()
	result.Name = sanitize.HTML(result.Name)
	result.Description = sanitize.HTML(result.Description)

	if err := m.Collection.Insert(result); err != nil {
		c.String(400, "DB error")
		return
	}
	c.JSON(200, obj{"data": result})
}

// Edit record
func (m *AuthModule) Edit(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))
	var result UsersStruct
	var newValue UsersStruct
	if !summer.PostBind(c, &newValue) {
		return
	}
	if err := m.Collection.FindId(id).One(&result); err == nil {
		result.Name = sanitize.HTML(newValue.Name)
		result.Description = sanitize.HTML(newValue.Description)
		result.Updated = time.Now()
		if err := m.Collection.UpdateId(newValue.ID, obj{"$set": result}); err != nil {
			c.String(400, "DB error")
			return
		}
	}
	c.JSON(200, obj{"data": result})
}

// Get record from DB
func (m *AuthModule) Get(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))
	result := AuthStruct{}
	if err := m.Collection.FindId(id).One(&result); err != nil {
		c.String(404, "Not found")
	}
	c.JSON(200, obj{"data": result})
}

// GetAll records
func (m *AuthModule) GetAll(c *gin.Context) {
	filter := struct {
		Search  string `form:"search"  json:"search"`
		Page    int    `form:"page"  json:"page"`
		Deleted bool   `form:"deleted"  json:"deleted"`
	}{}
	summer.PostBind(c, &filter)
	results := []AuthStruct{}
	request := obj{"deleted": filter.Deleted}

	// search engine
	if len(filter.Search) > 0 {
		regex := bson.RegEx{Pattern: filter.Search, Options: "i"}
		request["$or"] = arr{
			obj{"name": regex},
			obj{"description": regex},
		}
	}

	// records pagination
	count, _ := m.Collection.Find(request).Count()
	limit := 0
	skip := 0
	if filter.Page > 0 {
		limit = 50
		skip = limit * (filter.Page - 1)
	}

	// request to DB
	if err := m.Collection.Find(request).Sort("-_id").Limit(limit).Skip(skip).All(&results); err != nil {
		c.String(404, "Not found")
		return
	}

	c.JSON(200, obj{"data": results, "page": filter.Page, "count": count, "limit": limit})
}

// Action - remove/restore record
func (m *AuthModule) Action(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))

	if err := m.Collection.UpdateId(id, obj{"$set": obj{"deleted": c.PostForm("action") == "remove"}}); err != nil {
		c.String(404, "Not found")
		return
	}
	c.JSON(200, obj{"data": obj{"id": id}})
}
