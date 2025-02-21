package common

const (
	SwaggerDefaultModelsExpandDepth = 4
	EmptyPath                       = ""
)

const (
	OnlineStatusKeyState = "online_status_"
	LastOnlineKeyState   = "last_online_"
)

const (
	EmailEventTopic = "email.notification.event"
)

const (
	ConversationTypePrivate = "private"
	ConversationTypeGroup   = "group"
)

const (
	ParticipantRoleNone   = "none"
	ParticipantRoleAdmin  = "admin"
	ParticipantRoleMember = "member"
)

const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
	MessageTypeFile  = "file"
	MessageTypeVoice = "voice"
	MessageTypeVoIP  = "voice_call"
	MessageTypeVC    = "video_call"
)

const (
	WSEventActionChats            = "chats"
	WSEventActionChatMessages     = "messages"
	WSEventActionUserOnlineStatus = "online_status"
	WSEventActionUserTypingState  = "typing_state"
	WSEventActionNewTextMessage   = "new_text_message"
	WSEventActionDeleteGroup      = "delete_group"
	WSEventActionLeaveGroup       = "leave_group"
	WSEventCalling                = "calling"
	WSEventAnswerCall             = "answering"
)

const (
	WSEventCallbackErr          = "error"
	WSEventCallbackChats        = "chats"
	WSEventCallbackChatMessages = "messages"
	WSEventCallbackOnlineStatus = "online_status"
	WSEventCallbackTypingState  = "typing_state"
	WSEventCallbackNewMessage   = "new_message"
	WSEventCallbackRefreshChat  = "refresh_chat"
	WSEventCallbackIncomingCall = "incoming_call"
	WSEventCallbackAnswerCall   = "answer_call"
	WSEventCallbackNotification = "notification"
)
