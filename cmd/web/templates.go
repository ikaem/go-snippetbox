// cmd/web/templates.go

package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/ikaem/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// we create a function
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// this is a functions map
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// first we initialize a new map
	cache := map[string]*template.Template{}

	// then we use filepath flob to get slice of all filepaths
	// so glob will return a list of filepaths that match a patter

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	// then we loop over all pages, to geenerate template for each page
	for _, page := range pages {
		// first, we need the file name
		// we use base funciton to get last segment of the path
		name := filepath.Base(page)

		// then we want to parse the page tmeplate file into the tempalte set
		// ts, err := template.ParseFiles(page)
		// this is now to register the new functions map
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// then we also use parseGlob method to add any layout files to the template set
		// parse glob will get a list of paths? is it
		// so filepath join will create a single string
		// and then parse glob will go to that locaton, and will find all items that match the string
		// and if there is * inside file name, it will take all items that match the pattern?
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		// and we do the same for partials
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// and then we just add this particular tempalte to the cache map
		cache[name] = ts
	}

	return cache, nil

}
