package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	util "github.com/shashankC011/replayr/package"
)

var (
	LOGFILENAME = os.Getenv("LOGFILENAME")
	PORT        = os.Getenv("PORT")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env variables.\nerr: ", err)
	}
	fmt.Println("Server started on PORT: ", PORT)
	http.Handle("/", http.HandlerFunc(HealthzHandler))
	//	util.ReplayReqFromFile(LOGFILENAME)
	err = http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal("server closed.\n err: ", err)
	}
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	switch r.Method {
	case "GET":
		start := time.Now()
		// TODO: ADD THE  HTTP://LOCALHOST:8080 PART TO .ENV so it can be changed for production
		fmt.Println(r.URL.String())
		resp, err := client.Get("http://localhost:8080" + r.URL.String())
		duration := time.Since(start).Nanoseconds()
		if err != nil {
			log.Fatal("error GETTING to PORT 8080.\nerr: ", err)
		}
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Fatal("error reading into buf.\nerr: ", err)
			}
		}
		util.DumpReqAndResToFile(r, resp, duration, LOGFILENAME)
		fmt.Fprintf(w, "%s", string(buf[:n]))
	case "POST":
		resp, err := client.Post(r.URL.String(), "application/json", r.Body)
		if err != nil {
			log.Fatal("error in POST req to PORT 8080.\nerr: ", err)
		}
		fmt.Println(resp)
	}
}
