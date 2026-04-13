package rest

import (
	"github.com/gin-gonic/gin"
)

type UserHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	GetUser(c *gin.Context)
	SearchUsers(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UpdateStatus(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type ChatHandlerInterface interface {
	CreatePrivateChat(c *gin.Context)
	CreateGroupChat(c *gin.Context)
	GetUserChats(c *gin.Context)
	GetChatByID(c *gin.Context)
	UpdateChat(c *gin.Context)
	DeleteChat(c *gin.Context)
	GetAllChats(c *gin.Context)
}

type ParticipantHandlerInterface interface {
	AddParticipants(c *gin.Context)
	RemoveParticipant(c *gin.Context)
	LeaveChat(c *gin.Context)
	GetChatParticipants(c *gin.Context)
	UpdateRole(c *gin.Context)
	UpdateLastRead(c *gin.Context)
	GetUnreadCount(c *gin.Context)
}

type MessageHandlerInterface interface {
	SendMessage(c *gin.Context)
	UploadFile(c *gin.Context)
	GetMessages(c *gin.Context)
	EditMessage(c *gin.Context)
	DeleteMessage(c *gin.Context)
}