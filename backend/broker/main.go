package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/common-nighthawk/go-figure"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Servers []struct {
		Name   string `yaml:"name"`
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Prefix string `yaml:"prefix"`
	} `yaml:"servers"`
}

func newProxy(target *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host
			// Path and query untouched
		},
	}
}

func normalizePrefix(p string) string {
	p = "/" + strings.Trim(p, "/")
	return p + "/"
}

func main() {

	figure.NewColorFigure("Broker Service", "", "blue", true).Print()
	fmt.Println("\n")

	configContent, err := os.ReadFile("config.yml")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to read config file")
	}

	var config Config
	if err := yaml.Unmarshal(configContent, &config); err != nil {
		logger.Fatal().Err(err).Msg("failed to unmarshal config file")
	}

	type route struct {
		prefix string
		proxy  *httputil.ReverseProxy
	}

	var routes []route

	for _, server := range config.Servers {
		target, err := url.Parse(
			fmt.Sprintf("http://%s:%d", server.Host, server.Port),
		)
		if err != nil {
			logger.Fatal().Err(err).Msg("invalid target url")
		}

		prefix := normalizePrefix(server.Prefix)

		routes = append(routes, route{
			prefix: prefix,
			proxy:  newProxy(target),
		})

		logger.Info().
			Msgf("proxy %s → %s%s", server.Name, target, prefix)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, rt := range routes {
			if strings.HasPrefix(r.URL.Path, rt.prefix) {
				rt.proxy.ServeHTTP(w, r)
				return
			}
		}
		http.NotFound(w, r)
	})

	logger.Info().
		Msgf("broker started on %s:%d", config.Host, config.Port)

	logger.Fatal().
		Err(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler)).
		Msg("failed to start broker")
}
