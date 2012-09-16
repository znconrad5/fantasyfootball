package main

import (
	"fmt"
	"github.com/znconrad5/fantasyfootball"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var dataSourceTestDir = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/parsed")
var dataSourceTestStartWeek = 2
var dataSourceTestEndWeek = 14

var funcMap = template.FuncMap{
	"getWeekHeaders": weekHeaders,
}

var templates = template.Must(template.New("index").Funcs(funcMap).ParseFiles(os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/playerviewer/templ/index.html")))

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/js/", addResponseHeader(fileHandler, "Content-Type", "text/javascript"))
	http.HandleFunc("/css/", addResponseHeader(fileHandler, "Content-Type", "text/css"))
	http.HandleFunc("/css/overcast/images/", fileHandler)
	http.ListenAndServe(":8080", nil)
}

var dataDir = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/playerviewer/data")
var statsName = "stats.txt"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	loader := fantasyfootball.NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	loader.LoadAll()
	err := templates.ExecuteTemplate(w, "index.html", loader)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fantasyfootball.HandleError(err)
}

var validFilePath = regexp.MustCompile("^/(\\w+/)*[\\w+\\.-]+$")

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if validFilePath.MatchString(r.URL.Path) {
		//strip leading slash
		file, err := ioutil.ReadFile(r.URL.Path[1:])
		if err != nil {
			fmt.Printf("%v", err)
			http.NotFound(w, r)
			return
		} else {
			_, err = w.Write(file)
			if err != nil {
				fmt.Printf("%v", err)
			}
		}
	} else {
		http.NotFound(w, r)
		return
	}
}

func addResponseHeader(fn func(w http.ResponseWriter, r *http.Request), key string, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(key, value)
		fn(w, r)
	}
}

func weekHeaders(dataSource fantasyfootball.DataSource) []string {
	weeksSlice := make([]string, dataSource.EndWeek)
	for i := range weeksSlice {
		weeksSlice[i] = fmt.Sprintf("Week %v", i+1)
	}
	return weeksSlice
}
