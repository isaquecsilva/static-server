package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/isaquecsilva/static-server/controller/upload"
	"github.com/isaquecsilva/static-server/model/rules"
	"github.com/isaquecsilva/static-server/routes"
	"github.com/isaquecsilva/static-server/utils"
)

var (
	dir             = flag.String("dir", ".", "Directory to serve.")
	port            = flag.Int("port", 8080, "The port server will be listening on.")
	addr            = flag.String("addr", "localhost", "The address the server will be listening.")
	uploadEnabled   = flag.Bool("enableupload", false, "Sets if the server can accept file uploads.")
	uploadRulesFile = flag.String("rulesfile", "", "The xml file to be used as a ruler for file uploading(whether enableupload is true).")
	newRulesFile    = flag.Bool("newrulesfile", false, "Creates a rules.xml file.")
)

func main() {
	flag.Parse()

	if *newRulesFile {
		if err := rules.CreateDefaultRulesFileTemplate(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("rules.xml file created.")
		os.Exit(0)
	}

	if *uploadEnabled {
		rules, err := rules.LoadUploadRulesFromFile(*uploadRulesFile)
		if err != nil {
			log.Fatal(err)
		}

		if size, err := utils.ParseBytesUnit(rules.MaxFileSize); err != nil {
			log.Fatal(err)
		} else {
			upload.MaxUploadSize = size
		}

		uploadController, err := upload.NewUploadController(rules)

		if err != nil {
			log.Fatal(err)
		}

		routes.InitRoutes(uploadController)
	}

	routes.InitDefaultHandler(dir)

	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", *addr, *port),
		Handler:     nil,
		ReadTimeout: 0,
		ErrorLog:    log.New(os.Stderr, "static", 0),
	}

	go listenSignal(server)

	log.Printf("Server running at %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func listenSignal(server *http.Server) {
	defer server.Close()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	defer close(signalChannel)

	sign := <-signalChannel
	fmt.Fprintf(os.Stderr, "SIGNAL: %s\n", sign.String())
}
