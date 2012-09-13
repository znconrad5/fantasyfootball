package main

import (
	"fmt"
	"github.com/znconrad5/fantasyfootball"
	"html/template"
	"net/http"
	"os"
)

var dataSourceTestDir = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/parsed")
var dataSourceTestStartWeek = 2
var dataSourceTestEndWeek = 14

var templates = template.Must(template.ParseFiles(os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/playerviewer/templ/index.html")))

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	loader := fantasyfootball.NewNormalizedDataSource(fantasyfootball.NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek))
	err := templates.ExecuteTemplate(w, "index.html", loader)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fantasyfootball.HandleError(err)
}
