package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iyorozuya/neohooks/api/services"
	"github.com/iyorozuya/neohooks/api/structs"
)

type WebhookCoreController struct {
	WebhookService services.WebhookService
}

// Routes List of routes supported by webhook controller
func (wc *WebhookCoreController) Routes() []structs.Route {
	var baseURL string = "/api/webhook"
	return []structs.Route{
		{
			Method:  http.MethodGet,
			Path:    baseURL,
			Handler: wc.list,
		},
		{
			Method:  http.MethodPost,
			Path:    baseURL,
			Handler: wc.create,
		},
		{
			Method:  http.MethodGet,
			Path:    fmt.Sprintf("%s/{id}", baseURL),
			Handler: wc.retrieve,
		},
		{
			Method:  http.MethodDelete,
			Path:    fmt.Sprintf("%s/{id}", baseURL),
			Handler: wc.remove,
		},
	}
}

// GET /webhook - List all webhooks
func (wc *WebhookCoreController) list(w http.ResponseWriter, r *http.Request) {
	webhooks, err := wc.WebhookService.List()
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(
			structs.ErrorResponse{
				Errors: []string{"unable to fetch webhooks"},
			})
		return
	}
	json.NewEncoder(w).Encode(structs.ListWebhooksResponse{
		Webhooks: webhooks,
	})
}

// POST /webhook - Create new webhook
func (wc *WebhookCoreController) create(w http.ResponseWriter, r *http.Request) {
	webhook, err := wc.WebhookService.Save()
	if err != nil {
		json.NewEncoder(w).Encode(
			structs.ErrorResponse{
				Errors: []string{"error creating webhook"},
			},
		)
	}
	json.NewEncoder(w).Encode(structs.CreateWebhookResponse{
		ID: webhook,
	})
}

// GET /webhook/{id} - Get webhook by id
func (wc *WebhookCoreController) retrieve(w http.ResponseWriter, r *http.Request) {
	webhookId := chi.URLParam(r, "id")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}
	webhook, err := wc.WebhookService.Retrieve(webhookId, int64(page))
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(
			structs.ErrorResponse{
				Errors: []string{err.Error()},
			},
		)
	}
	json.NewEncoder(w).Encode(
		structs.RetrieveWebhookResponse{
			ID:       webhook.ID,
			Requests: webhook.Requests,
			Page:     webhook.Page,
			Total:    webhook.Total,
			Rows:     webhook.Rows,
		},
	)
}

// DELETE /webhook/{id} - Remove existing webhook
func (wc *WebhookCoreController) remove(w http.ResponseWriter, r *http.Request) {
	webhookId := chi.URLParam(r, "id")
	webhook := wc.WebhookService.Remove(webhookId)
	json.NewEncoder(w).Encode(structs.RemoveWebhookResponse{
		ID: webhook,
	})
}
