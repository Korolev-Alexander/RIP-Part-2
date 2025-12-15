package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"smartdevices/internal/api/serializers"
	"smartdevices/internal/middleware"
	"smartdevices/internal/models"
	"smartdevices/internal/storage"

	"gorm.io/gorm"
)

type SmartDeviceAPIHandler struct {
	db             *gorm.DB
	authMiddleware *middleware.AuthMiddleware
}

func NewSmartDeviceAPIHandler(db *gorm.DB) *SmartDeviceAPIHandler {
	return &SmartDeviceAPIHandler{
		db:             db,
		authMiddleware: middleware.NewAuthMiddleware(db),
	}
}

// GET /api/smart-devices - список с фильтрацией
func (h *SmartDeviceAPIHandler) GetSmartDevices(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	protocol := r.URL.Query().Get("protocol")

	var devices []models.SmartDevice
	query := h.db.Where("is_active = ?", true)

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	if protocol != "" {
		query = query.Where("protocol = ?", protocol)
	}

	result := query.Find(&devices)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	var response []serializers.SmartDeviceResponse
	for _, device := range devices {
		response = append(response, serializers.SmartDeviceToJSON(device))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/smart-devices/{id} - одна запись
func (h *SmartDeviceAPIHandler) GetSmartDevice(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-devices/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var device models.SmartDevice
	result := h.db.First(&device, id)
	if result.Error != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializers.SmartDeviceToJSON(device))
}

// POST /api/smart-devices - добавление устройства
func (h *SmartDeviceAPIHandler) CreateSmartDevice(w http.ResponseWriter, r *http.Request) {
	var req serializers.SmartDeviceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	device := models.SmartDevice{
		Name:           req.Name,
		Model:          req.Model,
		AvgDataRate:    req.AvgDataRate,
		DataPerHour:    req.DataPerHour,
		NamespaceURL:   req.NamespaceURL,
		Description:    req.Description,
		DescriptionAll: req.DescriptionAll,
		Protocol:       req.Protocol,
		IsActive:       true,
	}

	result := h.db.Create(&device)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(serializers.SmartDeviceToJSON(device))
}

// PUT /api/smart-devices/{id} - изменение устройства
func (h *SmartDeviceAPIHandler) UpdateSmartDevice(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-devices/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var device models.SmartDevice
	result := h.db.First(&device, id)
	if result.Error != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	var req serializers.SmartDeviceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	device.Name = req.Name
	device.Model = req.Model
	device.AvgDataRate = req.AvgDataRate
	device.DataPerHour = req.DataPerHour
	device.NamespaceURL = req.NamespaceURL
	device.Description = req.Description
	device.DescriptionAll = req.DescriptionAll
	device.Protocol = req.Protocol

	h.db.Save(&device)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializers.SmartDeviceToJSON(device))
}

// DELETE /api/smart-devices/{id} - удаление устройства (БЕЗ удаления изображения из MinIO)
func (h *SmartDeviceAPIHandler) DeleteSmartDevice(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-devices/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var device models.SmartDevice
	result := h.db.First(&device, id)
	if result.Error != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	// ТОЛЬКО деактивация устройства, без удаления изображения из MinIO
	device.IsActive = false
	h.db.Save(&device)

	fmt.Printf("✅ Device deactivated: %s (ID: %d)\n", device.Name, device.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// POST /api/smart-devices/{id}/image - добавление изображения
func (h *SmartDeviceAPIHandler) UploadDeviceImage(w http.ResponseWriter, r *http.Request) {
	// Парсим multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем файл в память
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-devices/")
	idStr = strings.TrimSuffix(idStr, "/image")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var device models.SmartDevice
	result := h.db.First(&device, id)
	if result.Error != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	// Генерируем имя файла на латинице
	fileExt := ".png"
	if strings.Contains(handler.Filename, ".") {
		fileExt = filepath.Ext(handler.Filename)
	}
	newFileName := fmt.Sprintf("device_%d_%d%s", device.ID, time.Now().Unix(), fileExt)

	// Загружаем файл в MinIO
	minioClient := storage.NewMinIOClient()
	err = minioClient.UploadFile(newFileName, fileData)
	if err != nil {
		fmt.Printf("❌ MinIO upload failed: %v\n", err)
		http.Error(w, "Failed to upload image to storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Обновляем URL в БД
	namespaceURL := minioClient.GetImageURL(newFileName)
	device.NamespaceURL = namespaceURL
	h.db.Save(&device)

	fmt.Printf("✅ Image uploaded: %s (%d bytes)\n", newFileName, len(fileData))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Image uploaded successfully",
		"image_url": namespaceURL,
		"file_name": newFileName,
		"file_size": len(fileData),
	})
}

// DELETE /api/smart-devices/{id}/image - удаление изображения устройства
func (h *SmartDeviceAPIHandler) DeleteDeviceImage(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/smart-devices/")
	idStr = strings.TrimSuffix(idStr, "/image")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var device models.SmartDevice
	result := h.db.First(&device, id)
	if result.Error != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	// Удаляем изображение из MinIO если есть
	if device.NamespaceURL != "" && strings.Contains(device.NamespaceURL, "localhost:9000") {
		filename := filepath.Base(device.NamespaceURL)
		minioClient := storage.NewMinIOClient()
		err := minioClient.DeleteFile(filename)
		if err != nil {
			fmt.Printf("⚠️ Failed to delete image from MinIO: %v\n", err)
			http.Error(w, "Failed to delete image from storage", http.StatusInternalServerError)
			return
		} else {
			fmt.Printf("✅ Image deleted from MinIO: %s\n", filename)
		}

		// Очищаем URL в БД
		device.NamespaceURL = ""
		h.db.Save(&device)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Image deleted successfully",
	})
}
