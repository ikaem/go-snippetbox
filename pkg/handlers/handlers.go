// pkg/handlers/handlers.go

package handlers

import (
	"net/http"
	"text/template"

	"github.com/ikaem/snippetbox/pkg/config"
)

// note that we have to specify return value
func Home(app *config.Application) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// if incorrect path, just send not found status code
		// if r.URL.Path != "/" {
		// 	http.NotFound(w, r)
		// 	return
		// }

		files := []string{
			"./ui/html/home.page.html",
			"./ui/html/base.layout.html",
			"./ui/html/footer.partial.html",
		}

		// so with  this we read the template file into a template set
		// template set is the ts variable
		// ts, err := template.ParseFiles("./ui/html/home.page.html")
		ts, err := template.ParseFiles(files...)

		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// now, we will use Execute to write the template as the response body
		// the last variable, for which we passed nil
		// is any dynamic data we want to pass in
		// we have no dynamic data for now

		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		// w.Write([]byte("Hello from snippetbox"))
	}
}
