package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"smartdevices/internal/api/serializers"
	"smartdevices/internal/models"

	"gorm.io/gorm"
)

type RequestAPIHandler struct {
	db *gorm.DB
}

func NewRequestAPIHandler(db *gorm.DB) *RequestAPIHandler {
	return &RequestAPIHandler{db: db}
}

// GET /api/requests/cart - иконка корзины
func (h *RequestAPIHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Фиксированный пользователь для демо (как в ТЗ)
	clientID := uint(1)

	var request models.Request
	result := h.db.Where("status = ? AND client_id = ?", "draft", clientID).First(&request)

	var response struct {
		RequestID uint `json:"request_id"`
		Count     int  `json:"count"`
	}

	if result.Error != nil {
		// Если корзины нет - возвращаем нули
		response.RequestID = 0
		response.Count = 0
	} else {
		response.RequestID = request.ID
		// Считаем общее количество товаров (сумма quantity)
		var totalQuantity struct {
			Total int
		}
		h.db.Model(&models.RequestService{}).
			Select("SUM(quantity) as total").
			Where("request_id = ?", request.ID).
			Scan(&totalQuantity)

		response.Count = totalQuantity.Total
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/requests - список заявок (кроме удаленных и черновика)
func (h *RequestAPIHandler) GetRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	status := r.URL.Query().Get("status")
	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")

	var requests []models.Request
	query := h.db.Preload("Client").Preload("Moderator").
		Where("status != ? AND status != ?", "deleted", "draft")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			query = query.Where("formed_at >= ?", dateFrom)
		}
	}

	if dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			query = query.Where("formed_at <= ?", dateTo.AddDate(0, 0, 1)) // включая весь день
		}
	}

	result := query.Find(&requests)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	var response []serializers.RequestResponse
	for _, request := range requests {
		// Загружаем items для каждой заявки
		var items []models.RequestService
		h.db.Preload("Service").Where("request_id = ?", request.ID).Find(&items)

		var itemResponses []serializers.RequestItemResponse
		for _, item := range items {
			itemResponses = append(itemResponses, serializers.RequestItemResponse{
				ServiceID:   item.ServiceID,
				ServiceName: item.Service.Name,
				Quantity:    item.Quantity,
				DataPerHour: item.Service.DataPerHour,
				ImageURL:    item.Service.ImageURL,
			})
		}

		response = append(response, serializers.RequestToJSON(request, itemResponses))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/requests/{id} - одна заявка
func (h *RequestAPIHandler) GetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/requests/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	var request models.Request
	result := h.db.Preload("Client").Preload("Moderator").First(&request, id)
	if result.Error != nil || request.Status == "deleted" {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}

	// Загружаем items
	var items []models.RequestService
	h.db.Preload("Service").Where("request_id = ?", request.ID).Find(&items)

	var itemResponses []serializers.RequestItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, serializers.RequestItemResponse{
			ServiceID:   item.ServiceID,
			ServiceName: item.Service.Name,
			Quantity:    item.Quantity,
			DataPerHour: item.Service.DataPerHour,
			ImageURL:    item.Service.ImageURL,
		})
	}

	response := serializers.RequestToJSON(request, itemResponses)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
