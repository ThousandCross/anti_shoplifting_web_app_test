package render

import (
	"anti-shoplifting-webapp/internal/config"
	"anti-shoplifting-webapp/internal/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "jwt") {
		td.IsAuthenticated = 1
	}
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	controllers_pre := []string{"signin", "signup"}
	controllers_main := []string{"blacklists", "dashboard", "incidents", "settings"}

	tc := map[string]*template.Template{}
	tc1, _ := CreateTemplateCacheEach(controllers_pre, "base-pre.layout.tmpl")
	tc2, _ := CreateTemplateCacheEach(controllers_main, "base.layout.tmpl")
	for k, v := range tc1 {
		tc[k] = v
	}
	for k, v := range tc2 {
		tc[k] = v
	}

	return tc, nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCacheEach(controllers []string, layoutFile string) (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages := []string{}
	for _, controller := range controllers {
		//fmt.Println("controller:", controller)
		pages_tmp, err := filepath.Glob("./templates/" + controller + "/*.page.tmpl")
		if err != nil {
			return myCache, err
		}
		pages = append(pages, pages_tmp...)
	}

	for _, page := range pages {
		//fmt.Println("page:", page)
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Fatal("Could not make template!!")
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/" + layoutFile)
		if err != nil {
			log.Fatal("Could not make template2 !!")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/" + layoutFile)
			if err != nil {
				log.Fatal("Could not make template3 !!")
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
