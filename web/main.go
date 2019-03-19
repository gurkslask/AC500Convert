package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gurkslask/AC500Convert"
)

//var htmlpath string = os.Getenv("GOPATH") + "src/github.com/gurkslask/AC500Convert/web/"
//#var htmlpath string = "/home/alexander/go/src/github.com/gurkslask/AC500Convert/web/"
var htmlpath string = "./web/"

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(htmlpath + "index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	in := r.FormValue("body")

	stext := strings.Split(in, "\n")
	vars, err := AC500Convert.ExtractData(stext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s := AC500Convert.OutputToText(vars)

	i := &infoslice{Text: s}
	t, err := template.ParseFiles(htmlpath + "view.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, i)
	if err != nil {
		log.Fatal(err)
	}
}

func genHandler(w http.ResponseWriter, r *http.Request) {
	in := r.FormValue("gen")

	stext := strings.Split(in, "\n")
	vars, err := AC500Convert.GenerateAccess(stext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	i := &infoslice{Text: vars}
	t, err := template.ParseFiles("./web/genview.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, i)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/gen", genHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type info struct {
	Text string
}
type infoslice struct {
	Text []string
}
