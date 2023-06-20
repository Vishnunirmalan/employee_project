package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Request structure for incoming JSON data
type Request struct {
	Ev  string               `json:"ev"`
	Et  string               `json:"et"`
	ID  string               `json:"id"`
	UID string               `json:"uid"`
	MID string               `json:"mid"`
	T   string               `json:"t"`
	P   string               `json:"p"`
	L   string               `json:"l"`
	SC  string               `json:"sc"`
	ATR map[string]Attribute `json:"-"`
	UAT map[string]Trait     `json:"-"`
}

type Attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Trait struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type TransformedRequest struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppID           string               `json:"app_id"`
	UserID          string               `json:"user_id"`
	MessageID       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageURL         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]Attribute `json:"attributes"`
	Traits          map[string]Trait     `json:"traits"`
}

func main() {
	//channel creation
	channel := make(chan Request)

	go worker(channel)

	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {

		var request Request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Send the request to the worker
		channel <- request

		w.WriteHeader(http.StatusOK)
	})
}

func worker(channel chan Request) {
	for request := range channel {

		transformedRequest := TransformedRequest{
			Event:           request.Ev,
			EventType:       request.Et,
			AppID:           request.ID,
			UserID:          request.UID,
			MessageID:       request.MID,
			PageTitle:       request.T,
			PageURL:         request.P,
			BrowserLanguage: request.L,
			ScreenSize:      request.SC,
			Attributes:      request.ATR,
			Traits:          request.UAT,
		}

		payload, err := json.Marshal(transformedRequest)
		if err != nil {
			fmt.Println("Error marshaling transformed request:", err)
			continue
		}
		//send request to wb
		err = sendToWebhook(payload)
		if err != nil {
			fmt.Println("Error sending transformed request to webhook:", err)
			continue
		}

		fmt.Println("Request processed successfully")
	}
}

func sendToWebhook(payload []byte) error {
	// Send the payload to the webhook (replace the URL with your actual webhook URL)
	resp, err := http.Post("https://webhook.site/1193b825-f63e-45f2-9c1d-6c7e7b3b5a4e", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook request failed with status code %d", resp.StatusCode)
	}
	log.Fatal(http.ListenAndServe(":3000", nil))
	return nil

}
