package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"smartdevices/internal/api/serializers"
	"smartdevices/internal/models"

	"gorm.io/gorm"
)

type ServiceAPIHandler struct {
	db *gorm.DB
}

func NewServiceAPIHandler(db *gorm.DB) *ServiceAPIHandler {
	return &ServiceAPIHandler{db: db}
}

// GET /api/services - список с фильтрацией
func (h *ServiceAPIHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	// Добавляем CORS headers для работы с Postman
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	search := r.URL.Query().Get("search")
	protocol := r.URL.Query().Get("protocol")

	var services []models.Service
	query := h.db.Where("is_active = ?", true)

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	if protocol != "" {
		query = query.Where("protocol = ?", protocol)
	}

	result := query.Find(&services)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Сериализуем в JSON
	var response []serializers.ServiceResponse
	for _, service := range services {
		response = append(response, serializers.ServiceToJSON(service))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/services/{id} - одна запись
func (h *ServiceAPIHandler) GetService(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/services/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	var service models.Service
	result := h.db.First(&service, id)
	if result.Error != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializers.ServiceToJSON(service))
}

// POST /api/services - добавление услуги
func (h *ServiceAPIHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req serializers.ServiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	service := models.Service{
		Name:           req.Name,
		Model:          req.Model,
		AvgDataRate:    req.AvgDataRate,
		DataPerHour:    req.DataPerHour,
		Description:    req.Description,
		DescriptionAll: req.DescriptionAll,
		Protocol:       req.Protocol,
		IsActive:       true,
	}

	result := h.db.Create(&service)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(serializers.ServiceToJSON(service))
}

// PUT /api/services/{id} - изменение услуги
func (h *ServiceAPIHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/services/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	var service models.Service
	result := h.db.First(&service, id)
	if result.Error != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	var req serializers.ServiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновляем поля
	service.Name = req.Name
	service.Model = req.Model
	service.AvgDataRate = req.AvgDataRate
	service.DataPerHour = req.DataPerHour
	service.Description = req.Description
	service.DescriptionAll = req.DescriptionAll
	service.Protocol = req.Protocol

	h.db.Save(&service)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializers.ServiceToJSON(service))
}

// DELETE /api/services/{id} - удаление услуги
func (h *ServiceAPIHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/services/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	var service models.Service
	result := h.db.First(&service, id)
	if result.Error != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Мягкое удаление - устанавливаем is_active = false
	service.IsActive = false
	h.db.Save(&service)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
