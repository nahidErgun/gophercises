package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nahidErgun/gophercises/gophercises/urlshortener"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshortener.MapHandler(pathsToUrls, mux)

	extension := flag.String("file", "yaml", "File extension (yaml, json)")
	flag.Parse()

	handler, err := getHandler(*extension, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)

}


func getRoutesFromFile(extension string) []byte {
	filePath := "./routes/routing." + extension
	output, err := ioutil.ReadFile(filePath)

	if !errors.Is(err, nil) {
		fmt.Println("Error: " + fmt.Sprintf("failed at reading file at %s", filePath))
		os.Exit(1)
	}

	return output
}

func getHandler(fileType string, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	routes := getRoutesFromFile(fileType)
	switch fileType {
	case "yaml":
		return urlshortener.YAMLHandler(routes, fallback)
	case "json":
		return urlshortener.JsonHandler(routes, fallback)
	default:
		return nil, errors.New("handler type not found")
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
