package controllers

import (
	"html/template"
	"net/http"
)

type HomeController struct {
	tmpl *template.Template
}

func NewHomeController(tmpl *template.Template) *HomeController {
	return &HomeController{tmpl: tmpl}
}

func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	// Create a temporary template that includes home blocks
	// We'll execute layout which will use the blocks from home.html
	if err := executeTemplate(w, r, c.tmpl, "home", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

