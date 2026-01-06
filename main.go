package main

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/shashankC011/replayr/internal/config"
	"github.com/shashankC011/replayr/internal/replay"
	util "github.com/shashankC011/replayr/package"
)

// TODO: add cli with the following features: start, stop recording, replay one or multiple capture files, merge capture files
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env variables.\nerr: ", err)
	}

	cfg := config.Load()

	//for now i am embending the file descriptor in the handler however i am unsure if it is a correct approach. TODO: confirm this
	// file, err := start(cfg.CAPTURE_DIR_ADDR)
	//TODO: maybe move this to the start function itself, basically dry out the main func
	// h := record.New(cfg, file)

	// fmt.Println("Server started on PORT: ", cfg.PORT)

	// http.HandleFunc("/", h.HealthzHandler)t
	// err = http.ListenAndServe(":"+cfg.PORT, nil)
	// if err != nil {
	// 	//TODO: stop should not be implemented here,works for testing but have to change it when cli is made
	// 	err := stop(file)
	// 	log.Fatal("server closed.\n err: ", err)
	// }

	//TODO:hard coded file name for now(should be user inupt from cli)
	replay.ReplayReqFromFile(cfg.CAPTURE_DIR_ADDR, "capture_2025-12-17_22-36-06.jsonl", cfg.REPLAYS_DIR_ADDR)
	for true {
	}
}

func stop(file *os.File) error {
	err := file.Close()
	if err != nil {
		return err
	}
	return nil
}

func start(dirName string) (*os.File, error) {
	fileName := util.GenerateCaptureFileName()
	filePath := dirName + "/" + fileName
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644) // 0644 → 0 = octal, 6(owner)=read+write, 4(group)=read, 4(others)=read
	if err != nil {
		return nil, errors.New("Error opening given file in APPEND, WRITE ONLY, CREATE mode\ner: " + err.Error())
	}
	return file, nil
}
