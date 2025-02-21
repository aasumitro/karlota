package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/model/entity"
	"karlota.aasumitro.id/internal/model/request"
	"karlota.aasumitro.id/internal/utils/cache"
	"karlota.aasumitro.id/internal/utils/http/middleware"
	"karlota.aasumitro.id/internal/utils/http/wrapper"
)

var wsu = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		for _, origin := range config.AllowOrigins {
			if r.Header.Get("Origin") == origin {
				return true
			}
		}
		return false
	},
}

type conversationHandler struct {
	mu  sync.RWMutex
	srv IConversationService
	rmq *amqp.Connection
	wsc map[uint]*websocket.Conn
}

func (h *conversationHandler) add(ctx *gin.Context) {
	rawID := ctx.MustGet("id").(float64)
	rawDN := ctx.MustGet("display_name").(string)
	var body request.NewGroupRequest
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	for _, m := range body.Members {
		if err := m.Validate(); err != nil {
			wrapper.NewHTTPRespondWrapper(
				ctx, http.StatusUnprocessableEntity, err)
			return
		}
	}
	body.CreatorID, body.CreatorDisplayName = uint(rawID), rawDN
	id, err := h.srv.New(ctx, body)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Return the response
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusCreated, id)
}

func (h *conversationHandler) interact(ctx *gin.Context) {
	// Upgrade the connection to WebSocket
	ws, err := wsu.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	defer func(ws *websocket.Conn) { _ = ws.Close() }(ws)
	// Get chat ID and user ID from the context
	unsafeUID, ok := ctx.MustGet("id").(float64)
	if !ok {
		wrapper.NewHTTPRespondWrapper(ctx,
			http.StatusBadRequest, "invalid user id")
		return
	}
	uid := uint(unsafeUID)
	unsafeUDN, ok := ctx.MustGet("display_name").(string)
	if !ok {
		wrapper.NewHTTPRespondWrapper(ctx,
			http.StatusBadRequest, "invalid user display_name")
		return
	}
	udn := unsafeUDN
	// Add the WebSocket & WebRTC client
	h.onOpen(uid, ws)
	defer h.onClose(uid)
	// Read messages in a loop
	h.onAction(ctx, ws, uid, udn)
}

func (h *conversationHandler) onOpen(
	userID uint, conn *websocket.Conn,
) {
	// set Websocket conn
	h.mu.Lock()
	defer h.mu.Unlock()
	h.wsc[userID] = conn
	// set user online
	onlineStatusKey := fmt.Sprintf("%s%d", common.OnlineStatusKeyState, userID)
	lastOnlineKey := fmt.Sprintf("%s%d", common.LastOnlineKeyState, userID)
	cache.Instance().Set(onlineStatusKey, true, cache.NoExpiration)
	cache.Instance().Set(lastOnlineKey, time.Now().Unix(), cache.NoExpiration)
	// check notify from queue table and sent to user when online
	go h.dequeueNotify(userID)
}

func (h *conversationHandler) onAction(
	ctx context.Context, ws *websocket.Conn, uid uint, udn string,
) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		// When the client sends a message with the registered case
		// call the specified function and do the action
		var payload request.WebsocketPayload
		if err := json.Unmarshal(message, &payload); err != nil {
			log.Printf("Failed to unmarshal WSAction payload: %v", err)
			continue
		}
		payload.UserID = uid
		payload.UserDN = udn
		if err := payload.ValidateActionRequest(); err != nil {
			log.Printf("Failed to validate WSAction payload: %v", err)
			continue
		}
		// proceed action
		switch payload.Action {
		// messaging
		case common.WSEventActionChats:
			h.chats(ctx, &payload)
		case common.WSEventActionChatMessages:
			h.messages(ctx, &payload)
		case common.WSEventActionUserOnlineStatus:
			h.onlineStatus(&payload)
		case common.WSEventActionUserTypingState:
			h.typingState(&payload)
		case common.WSEventActionNewTextMessage:
			h.newTextMessage(ctx, &payload)
		case common.WSEventActionDeleteGroup:
			h.deleteGroup(ctx, &payload)
		case common.WSEventActionLeaveGroup:
			h.leaveGroup(ctx, &payload)
		case common.WSEventCalling:
			h.call(&payload)
		case common.WSEventAnswerCall:
			h.answer(&payload)
		}
	}
}

func (h *conversationHandler) onClose(userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conn, ok := h.wsc[userID]; ok {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close WSConn for user %d: %v", userID, err)
		}
		delete(h.wsc, userID)
		// set user online
		onlineStatusKey := fmt.Sprintf("%s%d", common.OnlineStatusKeyState, userID)
		cache.Instance().Set(onlineStatusKey, false, cache.NoExpiration)
	}
}

// chats handler that is trigger triggered when a user connects to load all the conversation messages
func (h *conversationHandler) chats(
	ctx context.Context, payload *request.WebsocketPayload,
) {
	chats, err := h.srv.Chats(ctx, payload.UserID)
	if err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err.Error())
		return
	}
	h.reply(common.WSEventCallbackChats, payload.UserID, chats)
}

// messages handler that is when user try to get the conversation detail
func (h *conversationHandler) messages(
	ctx context.Context, payload *request.WebsocketPayload,
) {
	if err := payload.ValidateConversationIDRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	messages, err := h.srv.Messages(ctx, payload.ConversationID)
	if err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err.Error())
		return
	}
	h.reply(common.WSEventCallbackChatMessages, payload.UserID, messages)
}

func (h *conversationHandler) onlineStatus(
	payload *request.WebsocketPayload,
) {
	if err := payload.ValidateOnlineStatusRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	if payload.TargetID == 0 {
		return
	}
	user := &entity.ChatUser{ID: payload.TargetID}
	onlineStatusKey := fmt.Sprintf("%s%d", common.OnlineStatusKeyState, user.ID)
	lastOnlineKey := fmt.Sprintf("%s%d", common.LastOnlineKeyState, user.ID)
	if data, ok := cache.Instance().Get(onlineStatusKey); ok && data != nil {
		if status, ok := data.(bool); ok {
			user.IsOnline = status
		}
	}
	if data, ok := cache.Instance().Get(lastOnlineKey); ok && data != nil {
		if lastOnlineValue, ok := data.(int64); ok {
			user.LastOnline = lastOnlineValue
		}
	}
	h.reply(common.WSEventCallbackOnlineStatus, payload.UserID, map[string]interface{}{
		"user": user, "conversation_id": payload.ConversationID})
}

// typing handler that is triggered when a user trying to type something
func (h *conversationHandler) typingState(
	payload *request.WebsocketPayload,
) {
	if err := payload.ValidateTypingRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	h.broadcast(payload.RecipientID, func(recipientID uint) {
		h.reply(common.WSEventCallbackTypingState, recipientID, map[string]interface{}{
			"conversation_id": payload.ConversationID, "status": payload.TypingStatus})
	})
}

func (h *conversationHandler) newTextMessage(
	ctx context.Context, payload *request.WebsocketPayload,
) {
	if err := payload.ValidateNewMessageRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	message, err := h.srv.NewTextMessage(ctx, payload)
	if err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err.Error())
		return
	}
	h.broadcast(payload.RecipientID, func(recipientID uint) {
		h.reply(common.WSEventCallbackTypingState, recipientID, map[string]interface{}{
			"conversation_id": payload.ConversationID, "status": false})
		h.reply(common.WSEventCallbackNewMessage, recipientID, message)
	})
}

func (h *conversationHandler) deleteGroup(
	ctx context.Context, payload *request.WebsocketPayload,
) {
	if err := payload.ValidateConversationIDRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	if err := h.srv.RequestDeleteGroup(ctx, payload); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err.Error())
	}
}

func (h *conversationHandler) leaveGroup(
	ctx context.Context, payload *request.WebsocketPayload,
) {
	if err := payload.ValidateConversationIDRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	message, err := h.srv.RequestLeaveGroup(ctx, payload)
	if err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	if message == nil {
		return
	}
	h.broadcast(payload.RecipientID, func(recipientID uint) {
		h.reply(common.WSEventCallbackNewMessage, recipientID, message)
		h.reply(common.WSEventCallbackRefreshChat, recipientID, "")
	})
}

func (h *conversationHandler) call(payload *request.WebsocketPayload) {
	if err := payload.ValidateNewCallRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	if err := payload.Call.ValidateCallRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
	}
	h.broadcast(payload.RecipientID, func(recipientID uint) {
		h.reply(common.WSEventCallbackIncomingCall, recipientID, map[string]interface{}{
			"payload": payload.Call, "recipient": payload.UserID})
	})
}

func (h *conversationHandler) answer(payload *request.WebsocketPayload) {
	if err := payload.ValidateNewCallRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
		return
	}
	if err := payload.Call.ValidateAnswerCallRequest(); err != nil {
		h.reply(common.WSEventCallbackErr, payload.UserID, err)
	}
	h.broadcast(payload.RecipientID, func(recipientID uint) {
		h.reply(common.WSEventCallbackAnswerCall, recipientID,
			map[string]interface{}{"payload": payload.Call})
	})
}

func (h *conversationHandler) broadcast(
	recipients []uint, fn func(recipientID uint),
) {
	if len(recipients) == 0 {
		return
	}
	// Define a channel to hold the tasks (recipientIDs)
	taskChan := make(chan uint, len(recipients))
	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ { // Start 10 workers (adjust based on your needs)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for recipientID := range taskChan {
				fn(recipientID)
			}
		}()
	}
	// Send tasks (recipientIDs) to workers
	for _, id := range recipients {
		taskChan <- id
	}
	// Close the channel after all tasks are sent
	close(taskChan)
	// Wait for all workers to complete
	wg.Wait()
}

// reply helper function, used to send a message to a specified client by given id
func (h *conversationHandler) reply(
	replyType string, userID uint, data any,
) {
	// Retrieve and lock the connection at once
	h.mu.Lock()
	conn, ok := h.wsc[userID]
	h.mu.Unlock()
	if !ok {
		h.enqueueNotify(replyType, userID, data)
		log.Printf("No connection found to reply to user %d", userID)
		return
	}
	// Marshal the message data into JSON
	message, err := json.Marshal(map[string]any{
		"type": replyType, "data": data})
	if err != nil {
		log.Printf("Failed to encode message: %v", err)
		return
	}
	// Send the message to the client over the WebSocket connection
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Printf("Failed to write message to user %d: %v", userID, err)
		return
	}
}

func (h *conversationHandler) enqueueNotify(
	replyType string, userID uint, data any,
) {
	if replyType != common.WSEventCallbackNewMessage {
		return
	}
	msg, ok := data.(*entity.Message)
	if !ok {
		return
	}
	h.srv.AddMessageQueue(context.Background(), userID, msg)
}

func (h *conversationHandler) dequeueNotify(userID uint) {
	ctx := context.Background()
	notify, err := h.srv.LoadQueue(ctx, userID)
	if err != nil {
		return
	}
	for _, msg := range notify {
		h.reply(common.WSEventCallbackNotification, msg.TargetID, msg.Payload)
		now := time.Now()
		msg.UpdatedAt = &now
		h.srv.DeleteQueue(ctx, msg)
	}
}

func NewConversationHandler(
	router gin.IRoutes,
	service IConversationService,
	rabbitMQ *amqp.Connection,
) {
	wsm := make(map[uint]*websocket.Conn)
	handler := &conversationHandler{srv: service, rmq: rabbitMQ, wsc: wsm}
	router.POST(common.EmptyPath, middleware.Auth(), handler.add)
	router.GET(common.EmptyPath, middleware.WSAuth(), handler.interact)
}
