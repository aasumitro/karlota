package entity

import "time"

type Participant struct {
	ID        uint       `gorm:"primaryKey;column:id" json:"id"`
	ConvID    uint       `gorm:"column:conversation_id" json:"conversation_id"`
	UserID    uint       `gorm:"column:user_id" json:"user_id"`
	Role      string     `gorm:"column:role" json:"role"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (p *Participant) ToResponse() *Participant {
	return &Participant{
		ID:        p.ID,
		ConvID:    p.ConvID,
		UserID:    p.UserID,
		Role:      p.Role,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p *Participant) FromResponse(part *Participant) interface{} {
	return &Participant{
		ID:        part.ID,
		ConvID:    part.ConvID,
		UserID:    part.UserID,
		Role:      part.Role,
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}
