package serializers

import "smartdevices/internal/models"

type UserResponse struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	IsModerator bool   `json:"is_moderator"`
	IsActive    bool   `json:"is_active"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserToJSON(user models.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		IsModerator: user.IsModerator,
		IsActive:    user.IsActive,
	}
}
