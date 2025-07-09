package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"jelly/pkg/api"
)

func main() {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))

	// pprof web server. See: https://golang.org/pkg/net/http/pprof/
	if env == "local" || env == "dev" {
		go func() {
			pprof := "0.0.0.0:6060"
			slog.Info("Server is listening (pprof)", "address", pprof)
			log.Fatal(http.ListenAndServe(pprof, nil))
		}()
	}

	err := api.Run()
	if err != nil {
		os.Exit(1)
	}
}
