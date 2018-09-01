package controllers

import (
	"html/template"
)

//Register function registers all the controllers
//and their handlers
func Register(templates *template.Template) {
	gainers(templates)
}
