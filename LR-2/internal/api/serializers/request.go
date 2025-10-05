package serializers

import (
	"smartdevices/internal/models"
	"time"
)

type RequestResponse struct {
	ID            uint                  `json:"id"`
	Status        string                `json:"status"`
	Address       string                `json:"address"`
	TotalTraffic  float64               `json:"total_traffic"`
	ClientID      uint                  `json:"client_id"`
	ClientName    string                `json:"client_name"`
	FormedAt      *time.Time            `json:"formed_at,omitempty"`
	CompletedAt   *time.Time            `json:"completed_at,omitempty"`
	ModeratorID   *uint                 `json:"moderator_id,omitempty"`
	ModeratorName string                `json:"moderator_name,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
	Items         []RequestItemResponse `json:"items"`
}

type RequestItemResponse struct {
	ServiceID   uint    `json:"service_id"`
	ServiceName string  `json:"service_name"`
	Quantity    int     `json:"quantity"`
	DataPerHour float64 `json:"data_per_hour"`
	ImageURL    string  `json:"image_url"`
}

type RequestUpdateRequest struct {
	Address string `json:"address"`
}

type RequestFilter struct {
	Status   string    `form:"status"`
	DateFrom time.Time `form:"date_from"`
	DateTo   time.Time `form:"date_to"`
}

// Исправленная функция - проверяем по ModeratorID вместо указателя
func RequestToJSON(request models.Request, items []RequestItemResponse) RequestResponse {
	response := RequestResponse{
		ID:           request.ID,
		Status:       request.Status,
		Address:      request.Address,
		TotalTraffic: request.TotalTraffic,
		ClientID:     request.ClientID,
		ClientName:   request.Client.Username, // Client всегда должен быть загружен
		FormedAt:     request.FormedAt,
		CompletedAt:  request.CompletedAt,
		ModeratorID:  request.ModeratorID,
		CreatedAt:    request.CreatedAt,
		Items:        items,
	}

	// Если есть модератор и он загружен, берем его имя
	if request.ModeratorID != nil && request.Moderator.ID != 0 {
		response.ModeratorName = request.Moderator.Username
	}

	return response
}
