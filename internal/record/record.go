package record

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	models "github.com/shashankC011/replayr/internal/models"
)

func (handler *Handler) HealthzHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	switch r.Method {
	case "GET":

		start := time.Now()
		resp, err := client.Get(handler.cfg.SERVER_ADDR + r.URL.String())
		if err != nil {
			log.Fatal("error GETTING to ", handler.cfg.SERVER_ADDR, "\nerr: ", err)
		}
		duration := time.Since(start).Nanoseconds()

		//read resp as bytes
		bodyBytes, err := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		if err != nil {
			log.Fatal("error reading response body.\nerr: ", err)
		}

		//dump req and res to file
		err = DumpReqAndResToFile(r, resp, bodyBytes, duration, handler.file)
		if err != nil {
			log.Fatal("error dumping req and res to file.\nerr: ", err)
		}

		//send response
		fmt.Fprintf(w, "%v", string(bodyBytes))

	case "POST":
		resp, err := client.Post(r.URL.String(), "application/json", r.Body)
		if err != nil {

		}
		fmt.Println(resp)
	}
}

func DumpReqAndResToFile(r *http.Request, resp *http.Response, bodyBytes []byte, duration int64, file *os.File) error {
	// fileName := util.GenerateFileName()
	// filePath := dirName + "/" + fileName
	// file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644) // 0644 → 0 = octal, 6(owner)=read+write, 4(group)=read, 4(others)=read
	// if err != nil {
	// 	return errors.New("Error opening given file in APPEND, WRITE ONLY, CREATE mode\ner: " + err.Error())
	// }
	// defer file.Close()
	storedReq, err := models.MarshalReqToStoredReq(r)
	if err != nil {
		return err
	}
	storedResp, err := models.MarshallRespToStoredResp(resp, bodyBytes, duration)
	if err != nil {
		return err
	}
	loggedExchange := &models.LoggedExchange{
		Id:       uuid.New().String(),
		Time:     time.Now(),
		Request:  *storedReq,
		Response: *storedResp,
	}
	err = json.NewEncoder(file).Encode(loggedExchange)
	if err != nil {
		return errors.New("Error encoding loggedExchange struct\nerr: " + err.Error())
	}
	fmt.Println(loggedExchange)
	return nil
}
