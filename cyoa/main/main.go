package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)
type Story struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}
func main() {

	file, err := ioutil.ReadFile("./story.json")
	fmt.Println(err)
	jsonMap := make(map[string]Story)

	err = json.Unmarshal(file,&jsonMap)
	if err != nil {

		fmt.Println("failed to read json")
	}
	tmpl := template.Must(template.ParseFiles("index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		pathh := strings.Trim(string(path),"/")
		fmt.Println(pathh)
		tmpl.Execute(w, jsonMap[pathh])
	})

	http.ListenAndServe(":9090", nil)
}




