package controllers

import (
	"encoding/json"
	"net/http"

	"gses2-app/services"
)

type AppController struct {
	ExchangeRateService      services.ExchangeRateService
	EmailSubscriptionService services.EmailSubscriptionService
	EmailSenderService       services.EmailSenderService
}

func NewAppController(
	exchangeRateService services.ExchangeRateService,
	emailSubscriptionService services.EmailSubscriptionService,
	emailSenderService services.EmailSenderService,
) *AppController {
	return &AppController{
		ExchangeRateService:      exchangeRateService,
		EmailSubscriptionService: emailSubscriptionService,
		EmailSenderService:       emailSenderService,
	}
}

func (ac *AppController) GetRate(w http.ResponseWriter, r *http.Request) {
	exchangeRate, err := ac.ExchangeRateService.GetExchangeRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(exchangeRate)
}

func (ac *AppController) SubscribeEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	err := ac.EmailSubscriptionService.Subscribe(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ac *AppController) SendEmails(w http.ResponseWriter, r *http.Request) {
	exchangeRate, err := ac.ExchangeRateService.GetExchangeRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subscribers, err := ac.EmailSubscriptionService.GetSubscriptions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statusCode := ac.EmailSenderService.SendExchangeRate(exchangeRate, subscribers)
	if statusCode != 200 {
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}
