package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/nahidErgun/gophercises/gophercises/urlshortener"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback

	//mapHandler := urlshortener.MapHandler(pathsToUrls, mux)

	extension := flag.String("file", "json", "File extension (yaml, json)")
	flag.Parse()

	boltHandler := urlshortener.BoltHandler(db, mux)
	handler, err := getHandler(db, *extension, boltHandler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :9090")
	err = http.ListenAndServe(":9090", handler)

	if err != nil {
		log.Fatal(err)
	}

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

func getHandler(db *bolt.DB, fileType string, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	routes := getRoutesFromFile(fileType)
	switch fileType {
	case "yaml":
		return urlshortener.YAMLHandler(db, routes, fallback)
	case "json":
		return urlshortener.JsonHandler(db, routes, fallback)
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

func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		fmt.Printf("could not open db, %v", err)
		return nil, fmt.Errorf("could not open db, %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("ROUTES"))
		if err != nil {
			return fmt.Errorf("could not create days bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}

	p1 := "/urlshort-godoc"
	p2 := "/yaml-godoc"
	u1 := "https://godoc.org/github.com/gophercises/urlshort"
	u2 := "https://godoc.org/gopkg.in/yaml.v2"

	urlshortener.SetRoute(db, urlshortener.Route{P: p1, U: u1})
	urlshortener.SetRoute(db, urlshortener.Route{P: p2, U: u2})

	fmt.Println("DB Setup Done")
	return db, nil
}
