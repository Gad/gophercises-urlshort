package urlshort

import (
	"net/http"
	"gopkg.in/yaml.v3"
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
		redirect, exists := pathsToUrls[path]
		if exists {
			http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}



type pathToUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYaml(yml []byte) ([]pathToUrl, error) {
	// Parse yaml data into a slice of pathToUrl 
	var t []pathToUrl
	err := yaml.Unmarshal(yml, &t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func buildMap(parsedYaml []pathToUrl) map[string]string {
	// build a map of string:string from the slice of pathToUrl
	pathMap := make(map[string]string)
	for _,pT := range parsedYaml{
		pathMap[pT.Path]=pT.Url
	}
	return pathMap
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// Parse yaml data into a slice of pathToUrl 
	t, err := parseYaml(yml)
	if err!= nil {
		return nil, err
	}
	// build a map of string:string from the slice of pathToUrl 
	pathMap := buildMap(t)
	// pass it to MapHandler with the fallback handler
	return MapHandler(pathMap, fallback), nil
}
