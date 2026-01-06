package replay

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shashankC011/replayr/internal/models"
	util "github.com/shashankC011/replayr/package"
)

func ReplayReqFromFile(captureDirAddr, replayFromFile, replaysDirAddr string) {
	// name := util.GenerateReplayFileName()
	// fileName := fileDir + name
	replayFromFileFullAddr := captureDirAddr + "/" + replayFromFile
	channel := make(chan *models.LoggedExchange, 50)
	go ReadReqFromFile(replayFromFileFullAddr, channel)
	go ResendReq(replayFromFile, replaysDirAddr, channel)
}

func ReadReqFromFile(replayFromFileFullAddr string, channel chan *models.LoggedExchange) {
	file, err := os.Open(replayFromFileFullAddr)
	if err != nil {
		if os.IsNotExist(err) {
			//TODO: ADD ERR channel, TO HANDLE THESE BETTER MAYBE
			log.Fatal("The given capture file does not exist.\n", err)
		}
		//TODO: ADD ERR channel, TO HANDLE THESE BETTER MAYBE
		log.Fatal("error opening file: ", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() { // reads the next Line
		fmt.Println("read req no: ", count)
		reqLine := scanner.Text() // return the line as string
		var loggedExchange models.LoggedExchange
		err = json.Unmarshal([]byte(reqLine), &loggedExchange)
		if err != nil {
			log.Fatal("error decoding reqLine read from file: ", replayFromFileFullAddr)
		}
		//send the request from the loggedExchange to the channel(TODO: take the response and use it for comparison later)
		channel <- &loggedExchange
		count++
	}
	close(channel)
	if scanner.Err() != nil {
		log.Fatal("error while reading file.\nerr: ", err)
	}
}

func ResendReq(replayFromFile, replaysDirAddr string, channel chan *models.LoggedExchange) {
	//TODO: MAKE THIS EVEN FASTER BY USING MULTIPLE GO ROUTINES TO ITERATE OVER READ REQUESTS(IN THE CHANNEL) AS RIGHT NOW, REQUESTS ARE ONLY PROCESSED IF THE PREVIOUS ONE IS READ. WHICH DOESNT HAVE TO BE THE CASE.
	count := 0
	for storedLoggedExchange := range channel {
		storedReq := storedLoggedExchange.Request
		fmt.Println("resending req num: ", count)
		body := bytes.NewReader([]byte(storedReq.Body))
		// TODO: ADD THE  HTTP://LOCALHOST:8080 PART TO .ENV so it can be changed for production
		req, err := http.NewRequest(storedReq.Method, "http://localhost:8080"+storedReq.URL, body)
		if err != nil {
			fmt.Println("error created request from stored request\nerr: ", err)
		}
		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// TODO: MAKE AN ERROR CHANNEL to handle these errors better
			fmt.Println("error sending request\nerr: ", err)
		}
		duration := time.Since(start).Nanoseconds()
		//Not calling this in a goroutine for now, AS FOR NOW marhsallRespToStoredResp is a fairly synchronous and fast function... will do later IF DB CALLS ETC are added here.
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			//TODO: handle these errors in a separate channel
			log.Fatal("error reading response body.\nerr: ", err)
		}
		replayedRes, err := models.MarshallRespToStoredResp(resp, bodyBytes, duration)
		if err != nil {
			//TODO: handle these errors in a separate channel
			log.Fatal(err)
		}
		fmt.Println(replayedRes)
		// go compareResAndReplayedRes()
		go storeReplayedRequestResponse(storedLoggedExchange, replayedRes, replaysDirAddr, replayFromFile)
		count++
	}
}

// TODO: TAKE THE GENERATION OF FILE AND FOLDER OUT OF THIS FUNCTION AS THEY SHOULD NOT BE CALLED EVERYTIME
func storeReplayedRequestResponse(originalLoggedExchange *models.LoggedExchange, replayedRes *models.StoredResponse, replaysDirAddr, replayFromFile string) {
	req := &originalLoggedExchange.Request
	res := &originalLoggedExchange.Response
	replayedLoggedExchange := &models.ReplayedLoggedExchange{
		Id:               originalLoggedExchange.Id,
		Request:          *req,
		Response:         *res,
		ReplayedResponse: *replayedRes,
	}
	generatedReplayDirAddr, err := util.GenerateReplaySessionFolder(replaysDirAddr)
	if err != nil {
		log.Fatal(err)
	}
	replayFileName := util.GenerateReplayFileName(replayFromFile)
	fullReplayFilePath := generatedReplayDirAddr + "/" + replayFileName
	file, err := os.OpenFile(fullReplayFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644) // 0644 → 0 = octal, 6(owner)=read+write, 4(group)=read, 4(others)=read
	if err != nil {
		log.Fatal("Error opening given file in APPEND, WRITE ONLY, CREATE mode\ner: " + err.Error())
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(replayedLoggedExchange)
	if err != nil {
		log.Fatal(err)
	}
}

func CompareAndStoreRespAndReplayedResp(originalLoggedExchange *models.LoggedExchange, replayedLoggedExchange *models.ReplayedLoggedExchange) *models.ResponseDiff {
	return &models.ResponseDiff{
		ExchangeId: originalLoggedExchange.Id,
	}

}
