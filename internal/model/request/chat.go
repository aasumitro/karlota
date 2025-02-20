package request

import (
	"context"

	"github.com/golodash/galidator/v2"
)

type (
	NewGroupRequest struct {
		CreatorID          uint                    `json:"-" form:"-"`
		CreatorDisplayName string                  `json:"-" form:"-"`
		Name               string                  `json:"name" form:"name"`
		Members            []NewGroupMemberRequest `json:"members" form:"members"`
		Message            string                  `json:"message" form:"message"`
	}

	NewGroupMemberRequest struct {
		ID          int    `json:"id" from:"id"`
		DisplayName string `json:"display_name" form:"display_name"`
	}

	WebsocketPayload struct {
		Action         string `json:"action"`
		UserID         uint   `json:"user_id"`
		UserDN         string `json:"user_display_name"`
		ConversationID uint   `json:"conversation_id"`
		RecipientID    []uint `json:"recipient_id"`
		TargetID       uint   `json:"target_id"`
		TypingStatus   bool   `json:"typing_status"`
		TextMessage    string `json:"text_message"`
		Call           *Call  `json:"call"`
	}

	Call struct {
		Audio  bool   `json:"audio"`
		Video  bool   `json:"video"`
		PeerID string `json:"peer_id"`
		Action string `json:"action"`
	}
)

func (r *NewGroupMemberRequest) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"ID":          g.R("id").Required(),
		"DisplayName": g.R("display_name").Required(),
	}).Validate(context.TODO(), r)
}

func (r *WebsocketPayload) ValidateActionRequest() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"Action": g.R("action").Required(),
		"UserID": g.R("user_id").Required(),
	}).Validate(context.TODO(), r)
}

func (r *WebsocketPayload) ValidateConversationIDRequest() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"ConversationID": g.R("conversation_id").Required(),
	}).Validate(context.TODO(), r)
}

func (r *WebsocketPayload) ValidateTypingRequest() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"RecipientID":    g.R("recipient_id").Required(),
		"ConversationID": g.R("conversation_id").Required(),
	}).Validate(context.TODO(), r)
}

func (r *WebsocketPayload) ValidateNewMessageRequest() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"RecipientID":    g.R("recipient_id").Required(),
		"ConversationID": g.R("conversation_id").Required(),
		"TextMessage":    g.R("text_message").Required(),
	}).Validate(context.TODO(), r)
}

func (r *WebsocketPayload) ValidateOnlineStatusRequest() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"TargetID":       g.R("target_id").Required(),
		"ConversationID": g.R("conversation_id").Required(),
	}).Validate(context.TODO(), r)
}

func (r *WebsocketPayload) ValidateNewCallRequest() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"RecipientID": g.R("recipient_id").Required(),
	}).Validate(context.TODO(), r)
}

func (r *Call) ValidateCallRequest() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"PeerID": g.R("peer_id").Required(),
	}).Validate(context.TODO(), r)
}

func (r *Call) ValidateAnswerCallRequest() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"PeerID": g.R("peer_id").Required(),
		"Action": g.R("action").Required().Choices("accepted", "rejected"),
	}).Validate(context.TODO(), r)
}
