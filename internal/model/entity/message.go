package entity

import "time"

type Message struct {
	ID        uint       `gorm:"primaryKey;column:id" json:"id"`
	ConvID    uint       `gorm:"column:conversation_id" json:"conversation_id"`
	UserID    uint       `gorm:"column:user_id" json:"user_id"`
	Type      string     `gorm:"column:type" json:"type"`
	Content   string     `gorm:"column:content" json:"content"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	ReadAt    *time.Time `gorm:"column:read_at" json:"read_at"`
}

func (m *Message) ToResponse() *Message {
	return &Message{
		ID:        m.ID,
		ConvID:    m.ConvID,
		UserID:    m.UserID,
		Type:      m.Type,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		ReadAt:    m.ReadAt,
	}
}

func (m *Message) FromResponse(msg *Message) interface{} {
	return &Message{
		ID:        msg.ID,
		ConvID:    msg.ConvID,
		UserID:    msg.UserID,
		Type:      msg.Type,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		ReadAt:    msg.ReadAt,
	}
}
