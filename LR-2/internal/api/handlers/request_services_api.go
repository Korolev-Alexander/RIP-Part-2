package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"smartdevices/internal/models"

	"gorm.io/gorm"
)

type RequestServiceAPIHandler struct {
	db *gorm.DB
}

func NewRequestServiceAPIHandler(db *gorm.DB) *RequestServiceAPIHandler {
	return &RequestServiceAPIHandler{db: db}
}

// DELETE /api/request-services/{id} - удаление из заявки
func (h *RequestServiceAPIHandler) DeleteRequestService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := r.URL.Path[len("/api/request-services/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid request service ID", http.StatusBadRequest)
		return
	}

	// Находим запись по ID услуги (service_id)
	var requestService models.RequestService
	result := h.db.Where("service_id = ?", id).First(&requestService)
	if result.Error != nil {
		http.Error(w, "Request service not found", http.StatusNotFound)
		return
	}

	h.db.Delete(&requestService)

	w.WriteHeader(http.StatusNoContent)
}

// PUT /api/request-services/{id} - изменение количества
func (h *RequestServiceAPIHandler) UpdateRequestService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := r.URL.Path[len("/api/request-services/"):]
	serviceID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Quantity int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Quantity <= 0 {
		http.Error(w, "Quantity must be positive", http.StatusBadRequest)
		return
	}

	// Находим запись по ID услуги
	var requestService models.RequestService
	result := h.db.Where("service_id = ?", serviceID).First(&requestService)
	if result.Error != nil {
		http.Error(w, "Request service not found", http.StatusNotFound)
		return
	}

	// Обновляем количество
	requestService.Quantity = request.Quantity
	h.db.Save(&requestService)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service_id": requestService.ServiceID,
		"quantity":   requestService.Quantity,
		"updated":    true,
	})
}
