package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"github.com/google/uuid"
)

func executeTemplate(w http.ResponseWriter, tmpl *template.Template, pageName string, data interface{}) error {
	// Create a new template set with only layout + current page
	// This ensures the correct blocks are used
	funcMap := template.FuncMap{
		"uuid": func() string {
			return uuid.New().String()
		},
	}
	pageTemplate := template.New("").Funcs(funcMap)
	
	// Parse layout
	layoutPath := "./internal/app/views/templates/layout.html"
	pageTemplate, err := pageTemplate.ParseFiles(layoutPath)
	if err != nil {
		return err
	}
	
	// Parse the specific page template
	if pageName != "" {
		pagePath := filepath.Join("./internal/app/views/templates", pageName+".html")
		pageTemplate, err = pageTemplate.ParseFiles(pagePath)
		if err != nil {
			return err
		}
	}
	
	// Execute the layout which will use blocks from the page template
	return pageTemplate.ExecuteTemplate(w, "layout.html", data)
}

