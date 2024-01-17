// Package main provides ...
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/linnv/simdog"
)

func main() {
	simdog.ReadBuildInfo()
	flagDir := flag.String("dir", "./", "server directory")
	flagPort := flag.String("port", "8801", "listen port")
	flag.Parse()

	if !flag.Parsed() {
		os.Stderr.Write([]byte("ERROR: logging before flag.Parse"))
		return
	}
	dirHandler := http.FileServer(http.Dir(*flagDir))
	http.Handle("/", dirHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	port := ":" + *flagPort
	server := http.Server{
		Addr:    port,
		Handler: http.DefaultServeMux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
			} else {
				panic(err)
			}
		}
	}()

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	log.Print("http server listen on port ", port)
	log.Print("use c-c to exit: \n")
	<-sigChan
}
