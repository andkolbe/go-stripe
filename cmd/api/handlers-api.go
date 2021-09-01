package main

import (
	"encoding/json"
	"net/http"
)

// the payload we are receiving from the front end
type stripePayload struct {
	Currency string `json:"currency"`
	Amount string `json:"amount"`
}

// the response we are sending back the front end
type jsonResponse struct {
	OK bool `json:"ok"`
	Message string 	`json:"message"`
	Content string `json:"content"`
	ID int `json:"id"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	// send some manually built json back to the end user
	j := jsonResponse {
		OK: true,
	}

	// convert the response to json and send it back
	out, err := json.MarshalIndent(j, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}