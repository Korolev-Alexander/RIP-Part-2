package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"smartdevices/internal/middleware"
	"smartdevices/internal/models"

	"gorm.io/gorm"
)

type OrderItemAPIHandler struct {
	db             *gorm.DB
	authMiddleware *middleware.AuthMiddleware
}

func NewOrderItemAPIHandler(db *gorm.DB) *OrderItemAPIHandler {
	return &OrderItemAPIHandler{
		db:             db,
		authMiddleware: middleware.NewAuthMiddleware(db),
	}
}

// PUT /api/order-items/{deviceId} - изменение количества
func (h *OrderItemAPIHandler) UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.URL.Path[len("/api/order-items/"):]
	deviceID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	// Находим текущую корзину пользователя
	var order models.SmartOrder
	result := h.db.Where("status = ? AND client_id = ?", "draft", currentUser.ClientID).First(&order)
	if result.Error != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
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

	// Ищем устройство ИМЕННО в этой корзине
	var orderItem models.OrderItem
	result = h.db.Where("order_id = ? AND device_id = ?", order.ID, deviceID).First(&orderItem)
	if result.Error != nil {
		http.Error(w, "Device not found in cart", http.StatusNotFound)
		return
	}

	orderItem.Quantity = request.Quantity
	h.db.Save(&orderItem)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"device_id": orderItem.DeviceID,
		"quantity":  orderItem.Quantity,
		"updated":   true,
	})
}

// DELETE /api/order-items/{deviceId} - удаление из заявки
func (h *OrderItemAPIHandler) DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
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

	// ДОБАВИМ ОТЛАДКУ
	path := r.URL.Path
	log.Printf("🛠️ DeleteOrderItem path: %s", path)

	idStr := r.URL.Path[len("/api/order-items/"):]
	log.Printf("🛠️ DeleteOrderItem idStr: %s", idStr)

	deviceID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("❌ Error converting deviceID: %v", err)
		http.Error(w, "Invalid device ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("🛠️ DeleteOrderItem deviceID: %d", deviceID)

	// Находим текущую корзину пользователя
	var order models.SmartOrder
	result := h.db.Where("status = ? AND client_id = ?", "draft", currentUser.ClientID).First(&order)
	if result.Error != nil {
		log.Printf("❌ Cart not found: %v", result.Error)
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	log.Printf("🛠️ Found cart: ID=%d", order.ID)

	// Удаляем устройство ИЗ ЭТОЙ КОРЗИНЫ
	var orderItem models.OrderItem
	result = h.db.Where("order_id = ? AND device_id = ?", order.ID, deviceID).First(&orderItem)
	if result.Error != nil {
		log.Printf("❌ Device %d not found in cart %d: %v", deviceID, order.ID, result.Error)
		http.Error(w, "Device not found in cart", http.StatusNotFound)
		return
	}

	log.Printf("🛠️ Deleting device %d from cart %d", deviceID, order.ID)
	h.db.Delete(&orderItem)

	w.WriteHeader(http.StatusNoContent)
}
