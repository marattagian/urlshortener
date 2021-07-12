package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if destination, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, destination, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsUrls, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathsUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(parsedPaths []pathToUrl) map[string]string {
	mapedUrls := make(map[string]string)
	for _, parsedValue := range parsedPaths {
		mapedUrls[parsedValue.Path] = parsedValue.Url
	}
	return mapedUrls
}

func parseYAML(byteYAML []byte) ([]pathToUrl, error) {
	var pathToUrl []pathToUrl
	err := yaml.Unmarshal(byteYAML, &pathToUrl)
	if err != nil {
		return nil, err
	}
	return pathToUrl, nil
}

type pathToUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
