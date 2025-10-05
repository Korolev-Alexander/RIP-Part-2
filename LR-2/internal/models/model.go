package models

import (
	"time"
)

// Пользователь (аналог системной таблицы Django)
type User struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Username    string     `gorm:"uniqueIndex;size:150;not null" json:"username"`
	Password    string     `gorm:"size:128;not null" json:"-"` // совместимость с Django, скрываем в JSON
	IsModerator bool       `gorm:"default:false" json:"is_moderator"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	LastLogin   *time.Time `json:"last_login,omitempty"`
	DateJoined  time.Time  `gorm:"autoCreateTime" json:"date_joined"`
}

// Услуга (твои устройства)
type Service struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:200;not null" json:"name"`
	Model          string    `gorm:"size:100" json:"model"`
	AvgDataRate    float64   `json:"avg_data_rate"`
	DataPerHour    float64   `json:"data_per_hour"`
	ImageURL       string    `gorm:"size:500;null" json:"image_url"` // URL из MinIO
	Description    string    `json:"description"`
	DescriptionAll string    `gorm:"type:text" json:"description_all"`
	Protocol       string    `gorm:"size:50" json:"protocol"`
	IsActive       bool      `gorm:"default:true" json:"is_active"` // статус удален/действует
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Заявка (расширенная)
type Request struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Status    string    `gorm:"type:varchar(20);default:'draft';check:status IN ('draft','deleted','formed','completed','rejected')" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	ClientID  uint      `gorm:"not null" json:"client_id"`
	Client    User      `gorm:"foreignKey:ClientID;constraint:OnDelete:RESTRICT" json:"client"`

	// Дополнительные поля по ТЗ
	FormedAt    *time.Time `json:"formed_at,omitempty"`    // дата формирования
	CompletedAt *time.Time `json:"completed_at,omitempty"` // дата завершения
	ModeratorID *uint      `json:"moderator_id,omitempty"` // модератор
	Moderator   User       `gorm:"foreignKey:ModeratorID;constraint:OnDelete:RESTRICT" json:"moderator,omitempty"`

	// Поля по предметной области
	Address      string  `gorm:"size:500" json:"address"`
	TotalTraffic float64 `json:"total_traffic"` // рассчитывается при завершении

	// Индекс для ограничения "одна черновая заявка на пользователя"
	// Добавим через миграцию отдельно
}

// М-М Заявки-Услуги с составным ключом
type RequestService struct {
	RequestID uint      `gorm:"primaryKey" json:"request_id"` // часть составного ключа
	ServiceID uint      `gorm:"primaryKey" json:"service_id"` // часть составного ключа
	Quantity  int       `gorm:"default:1;not null" json:"quantity"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Связи
	Request Request `gorm:"foreignKey:RequestID;constraint:OnDelete:RESTRICT" json:"request"`
	Service Service `gorm:"foreignKey:ServiceID;constraint:OnDelete:RESTRICT" json:"service"`
}

// Уникальный индекс для ограничения одной черновой заявки
// Добавим в миграции
