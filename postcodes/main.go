package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"postcodes/config"
	"postcodes/routes"
	"strconv"

	"crypto/tls"
	"golang.org/x/crypto/acme"
	"net/http"
)

func main() {
	loadEnv()
	config.DB = config.InitDB()
	router := routes.LoadRoutes()

	if s, err := strconv.ParseBool(os.Getenv("USE_SSL")); s == false {
		log.Fatal(err, router.Start(os.Getenv("HTTP_HOST")+":"+os.Getenv("HTTP_PORT")))
	} else {
		s := http.Server{
			Addr:    os.Getenv("HTTP_HOST") + ":" + os.Getenv("HTTP_PORT"),
			Handler: router, // set Echo as handler
			TLSConfig: &tls.Config{
				//Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
				MinVersion:       tls.VersionTLS12,
				CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
				NextProtos: []string{
					"h2", "http/1.1", // enable HTTP/2
					acme.ALPNProto, // enable tls-alpn ACME challenges
				},
			},
		}
		if err := s.ListenAndServeTLS("/etc/letsencrypt/live/postcodes.local/fullchain.pem", "/etc/letsencrypt/live/postcodes.local/privkey.pem"); err != http.ErrServerClosed {
			router.Logger.Fatal(err)
		}
	}
}

func loadEnv() {
	godotenv.Load()
}
