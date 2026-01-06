package models

import (
	"time"
)

type LoggedExchange struct {
	Id       string         `json:"id"`
	Time     time.Time      `json:"time"`
	Request  StoredRequest  `json:"request"`
	Response StoredResponse `json:"response"`
}

type StoredRequest struct {
	Method string              `json:"method"`
	URL    string              `json:"url"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

type StoredResponse struct {
	Status     int                 `json:"status"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
	DurationMs int64               `json:"duration_ms"` // time taken by the backend to send its full response to your replayer proxy
}

type ReplayedLoggedExchange struct {
	Id               string         `json:"id"`
	Time             time.Time      `json:"time"`
	Request          StoredRequest  `json:"request"`
	Response         StoredResponse `json:"response"`
	ReplayedResponse StoredResponse `json:"replayedResponse"`
}

type ResponseDiff struct {
	ExchangeId string `json:"exchange_id"`
	Status     struct {
		Original int  `json:"original"`
		Replayed int  `json:"replayed"`
		Equal    bool `json:"equal"`
	} `json:"status"`
	Body struct {
		//add more fields here later as needed
		Equal bool `json:"equal"`
	} `json:"body"`
	Duration struct {
		OriginalMs int64 `json:"original_ms"`
		ReplayedMs int64 `json:"replayed_ms"`
		DeltaMs    int64 `json:"delta_ms"`
	} `json:"duration"`
	OverallEqual bool `json:"original_equal"`
}
