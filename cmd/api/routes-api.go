package main

import (
	"net/http"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// restricts who can and cannot access our backend
	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		// ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	  }))

	  mux.Post("/api/payment-intent", app.GetPaymentIntent)

	  mux.Get("/api/widget/{id}", app.GetWidgetByID)

	  mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)

	  mux.Post("/api/authenticate", app.CreateAuthToken)

	return mux
}