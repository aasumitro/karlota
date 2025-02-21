package entity

import (
	"time"
)

type Queue struct {
	ID        uint       `gorm:"primaryKey;column:id" json:"id"`
	Target    string     `gorm:"column:target" json:"target"`
	TargetID  uint       `gorm:"column:target_id" json:"target_id"`
	Trigger   string     `gorm:"column:trigger" json:"trigger"`
	Payload   string     `gorm:"column:payload" json:"payload"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"-"`
}

func (p *Queue) ToResponse() *Queue {
	return &Queue{
		ID:        p.ID,
		Target:    p.Target,
		TargetID:  p.TargetID,
		Trigger:   p.Trigger,
		Payload:   p.Payload,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p *Queue) FromResponse(part *Queue) interface{} {
	return &Queue{
		ID:        part.ID,
		Target:    part.Target,
		TargetID:  part.TargetID,
		Trigger:   part.Trigger,
		Payload:   part.Payload,
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}
