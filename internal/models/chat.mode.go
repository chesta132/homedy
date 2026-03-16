package models

import "time"

type ChatRoom struct {
	BaseRecyclable
	Name     string    `json:"name" gorm:"type:varchar(100)" example:"General" validate:"required,min=1,max=100"`
	IsGroup  bool      `json:"is_group" gorm:"default:false" example:"true"`
	Members  []User    `json:"members" gorm:"many2many:chat_room_members;"`
	Messages []Message `json:"messages,omitempty" gorm:"foreignKey:RoomID"`
}

type Message struct {
	BaseRecyclable
	RoomID    string      `json:"room_id" gorm:"not null" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid4"`
	SenderID  string      `json:"sender_id" gorm:"not null" example:"550e8400-e29b-41d4-a716-446655440001" validate:"required,uuid4"`
	Sender    User        `json:"sender" gorm:"foreignKey:SenderID"`
	Content   string      `json:"content" gorm:"type:text;not null" example:"Hello!" validate:"required,min=1"`
	Type      MessageType `json:"type" gorm:"default:'text'" example:"text" validate:"required,message_type"`
	ReadBy    []User      `json:"read_by,omitempty" gorm:"many2many:message_reads;"`
	ReplyToID *string     `json:"reply_to_id" gorm:"default:null" example:"550e8400-e29b-41d4-a716-446655440002" validate:"omitempty,uuid4"`
	ReplyTo   *Message    `json:"reply_to,omitempty" gorm:"foreignKey:ReplyToID"`
}

type MessageRead struct {
	MessageID string    `json:"message_id" gorm:"primaryKey" example:"550e8400-e29b-41d4-a716-446655440002" validate:"required,uuid4"`
	UserID    string    `json:"user_id" gorm:"primaryKey" example:"550e8400-e29b-41d4-a716-446655440001" validate:"required,uuid4"`
	ReadAt    time.Time `json:"read_at" example:"2024-01-15T10:30:00Z"`
}

type ChatRoomMember struct {
	RoomID   string     `json:"room_id" gorm:"primaryKey" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid4"`
	UserID   string     `json:"user_id" gorm:"primaryKey" example:"550e8400-e29b-41d4-a716-446655440001" validate:"required,uuid4"`
	JoinedAt time.Time  `json:"joined_at" example:"2024-01-15T08:00:00Z"`
	Role     MemberRole `json:"role" gorm:"default:'member'" example:"member" validate:"required,member_role"`
}

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

var MessageTypes = []MessageType{MessageTypeText, MessageTypeImage, MessageTypeFile}

type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

var MemberRoles = []MemberRole{MemberRoleAdmin, MemberRoleMember}
