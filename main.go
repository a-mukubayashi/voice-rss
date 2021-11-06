package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	apiClient *VoiceRssApiClient
)

type VoiceRssApiClient struct {
	Endpoint string
	ApiKey   string
}

type VoiceRssResponse struct {
	Message string
	Result  []byte
}

type VoiceRssErrorResponse struct {
	StatusCode string
	Message    string
}

func (e *VoiceRssErrorResponse) Error() string {
	return fmt.Sprintf("statusCode: %s, message: %s", e.StatusCode, e.Message)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	VOICE_RSS_API_KEY := os.Getenv("VOICE_RSS_API_KEY")

	if len(VOICE_RSS_API_KEY) == 0 {
		log.Fatal("VOICE_RSS_API_KEY is empty")
	}
	fmt.Println(VOICE_RSS_API_KEY)

	apiClient = &VoiceRssApiClient{
		Endpoint: "http://api.voicerss.org",
		ApiKey:   VOICE_RSS_API_KEY,
	}

	http.HandleFunc("/", getVoiceRss)
	http.ListenAndServe(":8090", nil)
}

func (client VoiceRssApiClient) requestToVoiceRss() (*VoiceRssResponse, error) {
	req, err := http.NewRequest("GET", client.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("key", client.ApiKey)
	q.Add("hl", "en-us")
	q.Add("src", "Hello, world!")

	fmt.Println(req.URL.String())

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// todo: apiがstatuscode200でもerrorになる
	if resp.StatusCode != 200 {
		fmt.Println("error")
		var errorResponse *VoiceRssErrorResponse
		errorResponse = &VoiceRssErrorResponse{
			StatusCode: fmt.Sprint(resp.StatusCode),
			Message:    "Api Error",
		}
		return nil, errorResponse
	}

	var voice *VoiceRssResponse
	body, _ := io.ReadAll(resp.Body)
	voice = &VoiceRssResponse{
		Message: "Success",
		Result:  body,
	}

	return voice, nil
}

func getVoiceRss(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	voice, err := apiClient.requestToVoiceRss()
	if err != nil {
		fmt.Println("get voice rss error")
		http.Error(w, "Internet Server Error", 503)
		return
	}
	json.NewEncoder(w).Encode(voice)
}
