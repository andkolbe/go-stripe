package main

import (
	"encoding/json"
	"myapp/internal/cards"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v72"
)

// the payload we are receiving from the front end
type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	ExpiryMonth   string `json:"exp_month"`
	ExpiryYear    string `json:"exp_year"`
	LastFour      string `json:"last_four"`
	Plan          string `json:"plan"`
	ProductID     string `json:"product_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

// the response we are sending back the front end
// ,omitempty means if the field is blank, don't show it
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
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

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
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
		j := jsonResponse{
			OK:      false,
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

func (app *application) CreateCustomerAndSubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	// there is a json payload passed to this handler from the client
	var data stripePayload
	// get the values out of the request body and into the data variable
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Println(data.Email, data.LastFour, data.PaymentMethod, data.Plan)

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: data.Currency,
	}

	okay := true
	var subscription *stripe.Subscription

	// get the customer
	// we pull payment method and email out of the json we received from the payload
	stripeCustomer, msg, err := card.CreateCustomer(data.PaymentMethod, data.Email)
	if err != nil {
		app.errorLog.Println(err)
		okay = false
	}

	if okay {
		// subscribe a customer to a plan
		subscription, err = card.SubscribeToPlan(stripeCustomer, data.Plan, data.Email, data.LastFour, "")
		if err != nil {
			app.errorLog.Println(err)
			okay = false
		}

		app.infoLog.Println("subscription id is", subscription.ID)
	}

	if okay {
		
	}
	// msg := ""

	resp := jsonResponse{
		OK:      okay,
		Message: msg,
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
