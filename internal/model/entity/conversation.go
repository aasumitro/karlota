package entity

import (
	"time"

	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/model/request"
)

type Conversation struct {
	ID        uint       `gorm:"primaryKey;column:id" json:"-"`
	Type      string     `gorm:"column:type" json:"type"`
	Name      string     `gorm:"column:name" json:"name,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (c *Conversation) ToResponse() *Conversation {
	return &Conversation{
		ID:        c.ID,
		Type:      c.Type,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (c *Conversation) FromResponse(conv *Conversation) interface{} {
	return &Conversation{
		ID:        conv.ID,
		Type:      conv.Type,
		Name:      conv.Name,
		CreatedAt: conv.CreatedAt,
		UpdatedAt: conv.UpdatedAt,
	}
}

func (c *Conversation) MakeNewGroupConversation(body request.NewGroupRequest) error {
	members := []request.NewGroupMemberRequest{{ID: int(body.CreatorID), DisplayName: body.CreatorDisplayName}}
	members = append(members, body.Members...)
	if body.Name == "" {
		name := "Conversation Between "
		numMembers := len(members)
		for i, member := range members {
			if i == numMembers-1 && numMembers > 1 {
				name += " & " + member.DisplayName
			} else {
				if i > 0 {
					name += ", "
				}
				name += member.DisplayName
			}
		}
		c.Name = name
	} else {
		c.Name = body.Name
	}
	c.Type = func() string {
		if len(members) == 2 {
			return common.ConversationTypePrivate
		}
		return common.ConversationTypeGroup
	}()
	c.CreatedAt = time.Now()
	return nil
}
