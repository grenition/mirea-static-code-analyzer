package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"reflect"
	"github.com/google/uuid"
	"webapp/internal/app/middleware"
)

func executeTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageName string, data interface{}) error {
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
	
	// Add user info to data if available
	var dataMap map[string]interface{}
	if data != nil {
		if m, ok := data.(map[string]interface{}); ok {
			dataMap = m
		} else {
			// For non-map data, use reflection to copy fields directly
			dataMap = make(map[string]interface{})
			dataMap["Data"] = data
			
			// Use reflection to copy struct fields directly to dataMap
			// This allows templates to access fields directly (e.g., .Issues instead of .Data.Issues)
			rv := reflect.ValueOf(data)
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			if rv.Kind() == reflect.Struct {
				rt := rv.Type()
				for i := 0; i < rv.NumField(); i++ {
					field := rt.Field(i)
					fieldValue := rv.Field(i)
					// Only export public fields (capitalized)
					if field.PkgPath == "" && fieldValue.CanInterface() {
						dataMap[field.Name] = fieldValue.Interface()
					}
				}
			}
		}
	} else {
		dataMap = make(map[string]interface{})
	}
	
	// Add user info
	if username, ok := middleware.GetUsername(r); ok {
		dataMap["Username"] = username
		dataMap["IsAuthenticated"] = true
	} else {
		dataMap["IsAuthenticated"] = false
	}

	// Execute the layout which will use blocks from the page template
	return pageTemplate.ExecuteTemplate(w, "layout.html", dataMap)
}

