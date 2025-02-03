package entity

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/utils/cache"
)

type UserChat struct {
	ID          uint           `gorm:"column:conversation_id" json:"id"`
	Type        string         `gorm:"column:conversation_type" json:"type"`
	Name        string         `gorm:"column:conversation_name" json:"name"`
	UsersRaw    string         `gorm:"column:users" json:"-"`
	MessagesRaw string         `gorm:"column:messages" json:"-"`
	Users       []*ChatUser    `gorm:"-"  json:"users"`
	Messages    []*ChatMessage `gorm:"-"  json:"messages"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

type ChatUser struct {
	ID         uint   `json:"id"`
	DP         string `json:"display_name"`
	Email      string `json:"email"`
	Role       string `json:"role,omitempty"`
	IsOnline   bool   `json:"is_online"`
	LastOnline int64  `json:"last_online"`
}

type ChatMessage struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *UserChat) ToResponse() *UserChat {
	uc := &UserChat{
		ID:        c.ID,
		Type:      c.Type,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if err := json.Unmarshal([]byte(c.UsersRaw), &uc.Users); err != nil {
		log.Println("failed to decode users conv:", uc.ID)
	}
	for _, u := range uc.Users {
		if u.Role == "none" {
			u.Role = ""
		}
		onlineStatusKey := fmt.Sprintf("%s%d", common.OnlineStatusKeyState, u.ID)
		lastOnlineKey := fmt.Sprintf("%s%d", common.LastOnlineKeyState, u.ID)
		if data, ok := cache.Instance().Get(onlineStatusKey); ok && data != nil {
			if status, ok := data.(bool); ok {
				u.IsOnline = status
			}
		}
		if data, ok := cache.Instance().Get(lastOnlineKey); ok && data != nil {
			if lastOnlineValue, ok := data.(int64); ok {
				u.LastOnline = lastOnlineValue
			}
		}
	}
	if err := json.Unmarshal([]byte(c.MessagesRaw), &uc.Messages); err != nil {
		log.Println("failed to decode messages conv:", uc.ID)
	}
	return uc
}

func (c *UserChat) FromResponse(uc *UserChat) interface{} {
	return &UserChat{
		ID:   uc.ID,
		Type: uc.Type,
		Name: uc.Name,
	}
}
