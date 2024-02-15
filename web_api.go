package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var NOHTTPS = flag.Bool("nohttps", false, "Disable https server")
var NOHTTP = flag.Bool("nohttp", false, "Disable http server")
var portflag = flag.Int("port", 8000, "Port to listen for http connections on")
var sslportflag = flag.Int("sslport", 8001, "Port to listen for https connections on")

func web_startup() {

	mux := http.NewServeMux()

	// Here we set up all the http targets and assign them to functions.
	//mux.HandleFunc("/", api_version)
	mux.HandleFunc("/read/{driver}/{tag...}", api_read)
	mux.HandleFunc("/read_multi/", api_read_multi)
	mux.HandleFunc("/write/{driver}/{tag...}", api_write)
	mux.HandleFunc("/view/{screen}", api_view)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	var handler http.Handler = mux

	if *basicHTTPAuthEnabled {

		a := basicHTTPAuth{
			User: *basicUser,
			Pass: *basicPass,
		}
		handler = a.Handler(handler)
	}

	if !*NOHTTP {
		port := fmt.Sprintf(":%v", *portflag)
		// And finally this starts the server
		go func() {
			http_server := &http.Server{
				Addr:              port,
				ReadHeaderTimeout: time.Minute,
				Handler:           handler,
			}
			err := http_server.ListenAndServe()
			if err != nil {
				log.Printf("Problem with http server. %v", err)
			}
		}()
	} else {
		log.Printf("Not starting HTTP server (NOHTTP = true)")
	}

	if !*NOHTTPS {
		go func() {
			tls_port := fmt.Sprintf(":%v", *sslportflag)
			certFile := "cert.crt"
			keyFile := "key.pem"

			https_server := &http.Server{
				Addr:              tls_port,
				ReadHeaderTimeout: time.Minute,
				Handler:           handler,
			}
			log.Printf("Starting HTTPS server on %s", tls_port)
			err := https_server.ListenAndServeTLS(certFile, keyFile)
			if err != nil {
				log.Printf("Problem with https server. %v", err)
			}
		}()
	} else {
		log.Printf("Not starting HTTPS server (NOHTTPS = true)")
	}

}
