package urlshortener

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if len(pathsToUrls[path]) > 0 {
			http.Redirect(w, r, pathsToUrls[path], http.StatusSeeOther)
		}

		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(db *bolt.DB,yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var routes []Route

	err := yaml.Unmarshal(yml, &routes)

	if err != nil {

		fmt.Println("failed yaml")
		return nil, err
	}

	for _,route  :=range routes {
		SetRoute(db,route)

	}

	return BoltHandler(db, fallback), nil
}

func JsonHandler(db *bolt.DB  ,js []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var routes []Route
	err := json.Unmarshal(js,&routes)

	if err != nil {

		fmt.Println("failed json")
		return nil, err
	}
	fmt.Println(routes)
	for  i := 1;  i<= len(routes); i++  {

	}

	fmt.Println("range of routes")
	fmt.Println(len(routes))
	for _,route  :=range routes {
		SetRoute(db,route)

	}
	return BoltHandler(db, fallback), nil
	}

type Route struct {
	P string `yaml:"path" json:"path"`
	U string `yaml:"url" json:"url"`
}

func buildMap(routes []Route) map[string]string {

	routeMap := make(map[string]string)
	for _, v := range routes {
		routeMap[v.P] = v.U
	}

return routeMap
}

func SetRoute(db *bolt.DB, route Route) error {

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("ROUTES")).Put([]byte(route.P), []byte(route.U))
		if err != nil {
			return fmt.Errorf("could not set url: %v", err)
		}
		return nil
	})
	fmt.Println("Set Urls : " + route.U)

	return err


}

func BoltHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url []byte
		path := r.URL.Path
		//Check bolt db
		// Get url from db

		 db.View(func(tx *bolt.Tx) error {
			url = tx.Bucket([]byte("DB")).Bucket([]byte("ROUTES")).Get([]byte(path))
			fmt.Printf("Url: %s\n", url)
			return  nil
		})


		if len(url) > 0 {
			http.Redirect(w, r, string(url), http.StatusSeeOther)
		}

		fallback.ServeHTTP(w, r)
	}
}
