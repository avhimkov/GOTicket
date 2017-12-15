package main

import (
	"github.com/kennygrant/sanitize"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/night-codes/summer.v1"
	"gopkg.in/night-codes/types.v1"
	"time"
)

type (
	ExampleStruct struct {
		ID          uint64    `form:"id"  json:"id"  bson:"_id"`
		Name        string    `form:"name" json:"name" bson:"name" valid:"required"`
		Description string    `form:"description" json:"description" bson:"description"`
		Created     time.Time `form:"-" json:"created" bson:"created"`
		Updated     time.Time `form:"-" json:"updated" bson:"updated"`
		Deleted     bool      `form:"-" json:"deleted" bson:"deleted"`
	}
	ExampleModule struct {
		summer.Module
	}

)


var (

	example = panel.AddModule(
		&summer.ModuleSettings{
			Name:           "example",
			Title:          "My Example",
			MenuOrder:      0,
			MenuTitle:      "My Example",
			Rights:         summer.Rights{Groups: []string{"all"}},

			Menu:           panel.MainMenu,
		},
		&ExampleModule{},
	)

)


// Add new record
func (m *ExampleModule) Add(c *gin.Context) {
	var result TasksStruct
	if !summer.PostBind(c, &result) {
		return
	}
	result.ID = panel.AI.Next("Example")
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
func (m *ExampleModule) Edit(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))
	var result ExampleStruct
	var newValue ExampleStruct
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
func (m *ExampleModule) Get(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))
	result := ExampleStruct{}
	if err := m.Collection.FindId(id).One(&result); err != nil {
		c.String(404, "Not found")
	}
	c.JSON(200, obj{"data": result})
}

// GetAll records
func (m *ExampleModule) GetAll(c *gin.Context) {
	filter := struct {
		Sort    string `form:"sort"  json:"sort"`
		Search  string `form:"search"  json:"search"`
		Page    int    `form:"page"  json:"page"`
		Deleted bool   `form:"deleted"  json:"deleted"`
	}{}
	summer.PostBind(c, &filter)
	results := []ExampleStruct{}
	request := obj{"deleted": filter.Deleted }

	sort := "-_id"

	// sort engine
	if len(filter.Sort) > 0 {
		sort = filter.Sort
	}

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
	if err := m.Collection.Find(request).Sort(sort).Limit(limit).Skip(skip).All(&results); err != nil {
		c.String(404, "Not found")
		return
	}

	c.JSON(200, obj{"data": results, "page": filter.Page, "count": count, "limit": limit })
}

// Action - remove/restore record
func (m *ExampleModule) Action(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))

	if err := m.Collection.UpdateId(id, obj{"$set": obj{"deleted": c.PostForm("action") == "remove" }}); err != nil {
		c.String(404, "Not found")
		return
	}
	c.JSON(200, obj{"data": obj{"id": id}})
}