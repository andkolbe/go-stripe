package main

import (
	"encoding/json"
	"myapp/internal/cards"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// the payload we are receiving from the front end
type stripePayload struct {
	Currency string `json:"currency"`
	Amount string `json:"amount"`
}

// the response we are sending back the front end
// ,omitempty means if the field is blank, don't show it
type jsonResponse struct {
	OK bool `json:"ok"`
	Message string 	`json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID int `json:"id,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {

	var payload stripePayload
	// get the body of the request and decode it into the payload variable
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// cast amount from string to int
	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	card := cards.Card {
		Secret: app.config.stripe.secret,
		Key: app.config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true

	// connects to stripe for us
	// we get a payment intent back
	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		// if the charge succeeded and the payment intent comes back with no error, send json back
		// create the json we want to send back
		out, err := json.MarshalIndent(pi, "", "   ")
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		// set the header and send it out
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		// if the charge failed
		j := jsonResponse {
			OK: false,
			Message: msg,
			Content: "",

		}
		// convert the response to json and send it back
		out, err := json.MarshalIndent(j, "", "   ")
		if err != nil {
		app.errorLog.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

func (app *application) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, err := json.MarshalIndent(widget, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}