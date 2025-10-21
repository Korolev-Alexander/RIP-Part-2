package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"smartdevices/internal/api/serializers"
	"smartdevices/internal/middleware"
	"smartdevices/internal/models"

	"gorm.io/gorm"
)

type SmartOrderAPIHandler struct {
	db             *gorm.DB
	authMiddleware *middleware.AuthMiddleware
}

func NewSmartOrderAPIHandler(db *gorm.DB) *SmartOrderAPIHandler {
	return &SmartOrderAPIHandler{
		db:             db,
		authMiddleware: middleware.NewAuthMiddleware(db),
	}
}

// GET /api/smart-orders/cart - иконка корзины
func (h *SmartOrderAPIHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Получаем текущего пользователя
	currentUser := h.authMiddleware.GetCurrentUser(r)
	if currentUser == nil {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}

	var order models.SmartOrder
	result := h.db.Where("status = ? AND client_id = ?", "draft", currentUser.ClientID).First(&order)

	var response struct {
		OrderID uint `json:"order_id"`
		Count   int  `json:"count"`
	}

	if result.Error != nil {
		response.OrderID = 0
		response.Count = 0
	} else {
		response.OrderID = order.ID
		var totalQuantity struct {
			Total int
		}
		h.db.Model(&models.OrderItem{}).
			Select("SUM(quantity) as total").
			Where("order_id = ?", order.ID).
			Scan(&totalQuantity)

		response.Count = totalQuantity.Total
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/smart-orders - список заявок (кроме удаленных и черновика)
func (h *SmartOrderAPIHandler) GetSmartOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Получаем текущего пользователя
	currentUser := h.authMiddleware.GetCurrentUser(r)
	if currentUser == nil {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")

	var orders []models.SmartOrder
	query := h.db.Preload("Client").Preload("Moderator")

	// Если не модератор - показываем только свои заявки
	if !currentUser.IsModerator {
		query = query.Where("client_id = ?", currentUser.ClientID)
	} else {
		// Модераторы не видят черновики и удаленные
		query = query.Where("status != ? AND status != ?", "deleted", "draft")
	}

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
			query = query.Where("formed_at <= ?", dateTo.AddDate(0, 0, 1))
		}
	}

	result := query.Find(&orders)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	var response []serializers.SmartOrderResponse
	for _, order := range orders {
		var items []models.OrderItem
		h.db.Preload("Device").Where("order_id = ?", order.ID).Find(&items)

		var itemResponses []serializers.SmartOrderItemResponse
		for _, item := range items {
			itemResponses = append(itemResponses, serializers.SmartOrderItemResponse{
				DeviceID:     item.DeviceID,
				DeviceName:   item.Device.Name,
				Quantity:     item.Quantity,
				DataPerHour:  item.Device.DataPerHour,
				NamespaceURL: item.Device.NamespaceURL,
			})
		}

		response = append(response, serializers.SmartOrderToJSON(order, itemResponses))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/smart-orders/{id} - одна заявка
func (h *SmartOrderAPIHandler) GetSmartOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Получаем текущего пользователя
	currentUser := h.authMiddleware.GetCurrentUser(r)
	if currentUser == nil {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.SmartOrder
	result := h.db.Preload("Client").Preload("Moderator").First(&order, id)
	if result.Error != nil || order.Status == "deleted" {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Проверяем права доступа
	if !currentUser.IsModerator && order.ClientID != currentUser.ClientID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var items []models.OrderItem
	h.db.Preload("Device").Where("order_id = ?", order.ID).Find(&items)

	var itemResponses []serializers.SmartOrderItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, serializers.SmartOrderItemResponse{
			DeviceID:     item.DeviceID,
			DeviceName:   item.Device.Name,
			Quantity:     item.Quantity,
			DataPerHour:  item.Device.DataPerHour,
			NamespaceURL: item.Device.NamespaceURL,
		})
	}

	response := serializers.SmartOrderToJSON(order, itemResponses)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PUT /api/smart-orders/{id} - изменение полей заявки
func (h *SmartOrderAPIHandler) UpdateSmartOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Получаем текущего пользователя
	currentUser := h.authMiddleware.GetCurrentUser(r)
	if currentUser == nil {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.SmartOrder
	result := h.db.First(&order, id)
	if result.Error != nil || order.Status == "deleted" {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Проверяем права доступа
	if !currentUser.IsModerator && order.ClientID != currentUser.ClientID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req serializers.SmartOrderUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновляем только разрешенные поля
	if req.Address != "" {
		order.Address = req.Address
	}

	h.db.Save(&order)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializers.SmartOrderToJSON(order, nil))
}

// PUT /api/smart-orders/{id}/form - формирование заявки
func (h *SmartOrderAPIHandler) FormSmartOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Получаем текущего пользователя
	currentUser := h.authMiddleware.GetCurrentUser(r)
	if currentUser == nil {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-orders/")
	idStr = strings.TrimSuffix(idStr, "/form")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.SmartOrder
	result := h.db.First(&order, id)
	if result.Error != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Проверяем права доступа
	if !currentUser.IsModerator && order.ClientID != currentUser.ClientID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Проверка обязательных полей
	if order.Address == "" {
		http.Error(w, "Address is required to form order", http.StatusBadRequest)
		return
	}

	// Установка статуса и даты формирования
	now := time.Now()
	order.Status = "formed"
	order.FormedAt = &now

	h.db.Save(&order)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializers.SmartOrderToJSON(order, nil))
}

// PUT /api/smart-orders/{id}/complete - завершение заявки
func (h *SmartOrderAPIHandler) CompleteSmartOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверяем права модератора
	currentUser := h.authMiddleware.GetCurrentUser(r)
	if currentUser == nil || !currentUser.IsModerator {
		http.Error(w, `{"error": "Moderator access required"}`, http.StatusForbidden)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-orders/")
	idStr = strings.TrimSuffix(idStr, "/complete")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.SmartOrder
	result := h.db.Preload("Client").First(&order, id)
	if result.Error != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Проверяем что заявка сформирована
	if order.Status != "formed" {
		http.Error(w, "Only formed orders can be completed", http.StatusBadRequest)
		return
	}

	// Расчет общего трафика по формуле из лабы 2
	var items []models.OrderItem
	h.db.Preload("Device").Where("order_id = ?", order.ID).Find(&items)

	totalTraffic := 0.0
	for _, item := range items {
		baseTraffic := item.Device.DataPerHour * float64(item.Quantity)

		// Формула расчета с коэффициентами для разных типов устройств
		var coefficient float64
		switch {
		case strings.Contains(item.Device.Name, "Хаб"):
			coefficient = 1.3 // Хабы требуют больше трафика
		case strings.Contains(item.Device.Name, "Датчик"):
			coefficient = 0.7 // Датчики экономят трафик
		case strings.Contains(item.Device.Name, "Лампочка"):
			coefficient = 1.1 // Лампочки немного больше
		case strings.Contains(item.Device.Name, "Розетка"):
			coefficient = 0.9 // Розетки мало трафика
		case strings.Contains(item.Device.Name, "Выключатель"):
			coefficient = 0.8 // Выключатели мало трафика
		default:
			coefficient = 1.0
		}

		traffic := baseTraffic * coefficient
		totalTraffic += traffic
	}

	// Установка статуса, модератора и даты завершения
	now := time.Now()
	order.Status = "completed"
	order.CompletedAt = &now
	order.ModeratorID = &currentUser.ClientID
	order.TotalTraffic = totalTraffic

	h.db.Save(&order)

	// Загружаем items для ответа
	var itemResponses []serializers.SmartOrderItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, serializers.SmartOrderItemResponse{
			DeviceID:     item.DeviceID,
			DeviceName:   item.Device.Name,
			Quantity:     item.Quantity,
			DataPerHour:  item.Device.DataPerHour,
			NamespaceURL: item.Device.NamespaceURL,
		})
	}

	response := serializers.SmartOrderToJSON(order, itemResponses)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DELETE /api/smart-orders/{id} - удаление заявки
func (h *SmartOrderAPIHandler) DeleteSmartOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Получаем текущего пользователя
	currentUser := h.authMiddleware.GetCurrentUser(r)
	if currentUser == nil {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.SmartOrder
	result := h.db.First(&order, id)
	if result.Error != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Проверяем права доступа
	if !currentUser.IsModerator && order.ClientID != currentUser.ClientID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Мягкое удаление - меняем статус
	order.Status = "deleted"
	h.db.Save(&order)

	w.WriteHeader(http.StatusNoContent)
}

// Вспомогательная функция
func uintPtr(i uint) *uint {
	return &i
}
