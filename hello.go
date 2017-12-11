
package main

import (
	"gopkg.in/night-codes/summer.v1"
)


var(

	aboutMenu = panel.MainMenu.Add("About Summer")
	hello = panel.AddModule(
		&summer.ModuleSettings{
			Name:       "hello",
			Title:      "Welcome to Summer panel",
			MenuTitle:  "Hello",
			Menu:       aboutMenu,
			GroupTitle: "Welcome",
		},
		&summer.Module{},
	)
	_ = panel.AddModule(
		&summer.ModuleSettings{
			Name:    "howto",
			Title:   "How To Use",
			GroupTo: hello,
			Menu:    aboutMenu,
		},
		&summer.Module{},
	)

)

