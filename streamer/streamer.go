package streamer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type StreamSource struct {
	Parameters   map[string]string `json:"parameters"`
	AnalysisDate string            `json:"analysisDate"`
	Streams      []Stream          `json:"streams"`
}

type Stream struct {
	Type       int      `json:"type"`
	Provider   string   `json:"provider"`
	Query      string   `json:"query"`
	DateValues []string `json:"dateValues"`
	InfoValues []string `json:"infoValues"`
}

func ReadStreamConfig(file string) (StreamSource, error) {
	var source StreamSource

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return source, fmt.Errorf("stream config read error")
	}

	json.Unmarshal(raw, &source)

	return source, nil
}
