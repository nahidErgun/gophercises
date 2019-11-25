package urlshortener

import (
	"fmt"
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var routes []route

	err := yaml.Unmarshal(yml, &routes)

	if err != nil {

		fmt.Println("failed yaml")
		return nil, err
	}

	return MapHandler(buildMap(routes), fallback), nil
}

type route struct {
	P string `yaml:"path"`
	U string `yaml:"url"`
}

func buildMap(routes []route) map[string]string {

	routeMap := make(map[string]string)
	for _, v := range routes {
		routeMap[v.P] = v.U
	}

	return routeMap
}