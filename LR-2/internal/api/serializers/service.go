package serializers

import (
	"smartdevices/internal/models"
	"time"
)

type ServiceResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Model          string    `json:"model"`
	AvgDataRate    float64   `json:"avg_data_rate"`
	DataPerHour    float64   `json:"data_per_hour"`
	ImageURL       string    `json:"image_url"`
	Description    string    `json:"description"`
	DescriptionAll string    `json:"description_all"`
	Protocol       string    `json:"protocol"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type ServiceCreateRequest struct {
	Name           string  `json:"name" binding:"required"`
	Model          string  `json:"model"`
	AvgDataRate    float64 `json:"avg_data_rate"`
	DataPerHour    float64 `json:"data_per_hour"`
	Description    string  `json:"description"`
	DescriptionAll string  `json:"description_all"`
	Protocol       string  `json:"protocol"`
}

func ServiceToJSON(service models.Service) ServiceResponse {
	return ServiceResponse{
		ID:             service.ID,
		Name:           service.Name,
		Model:          service.Model,
		AvgDataRate:    service.AvgDataRate,
		DataPerHour:    service.DataPerHour,
		ImageURL:       service.ImageURL,
		Description:    service.Description,
		DescriptionAll: service.DescriptionAll,
		Protocol:       service.Protocol,
		IsActive:       service.IsActive,
		CreatedAt:      service.CreatedAt,
	}
}
