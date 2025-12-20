package controllers

import (
	"html/template"
	"net/http"
)

type AboutController struct {
	tmpl *template.Template
}

func NewAboutController(tmpl *template.Template) *AboutController {
	return &AboutController{tmpl: tmpl}
}

func (c *AboutController) Index(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, r, c.tmpl, "about", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

