package users

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/potatowhite/books/file-service/pkg/service"
	"log"
	"os"
	"strconv"
)

var (
	// logger with thread id
	logger = log.New(os.Stdout, "handler", log.LstdFlags|log.Lshortfile)
)

type UserEventHandler struct {
	FileSvc   service.FileService
	FolderSvc service.FolderService
}

// Topic() returns the topic that this users is subscribed to
func (h *UserEventHandler) Topic() string {
	return "users"
}

// HandleMessage() handles a message from the topic that this users is subscribed to
func (h *UserEventHandler) HandleMessage(msg *kafka.Message) error {
	// extract eventType from header
	var eventType string
	for _, header := range msg.Headers {
		if string(header.Key) == "eventType" {
			eventType = string(header.Value)
		}
	}

	// extract userID from key
	userID := string(msg.Key)

	// Handle the event based on its type
	switch eventType {
	case "UserCreatedEvent":
		if err := h.createRootFolder(userID); err != nil {
			return fmt.Errorf("failed to create root folder for users %s: %v", userID, err)
		}

		payload := string(msg.Value)
		logger.Printf("payload: %s", payload)
	case "UserDeletedEvent":
		if err := h.deleteAllFolders(userID); err != nil {
			return fmt.Errorf("failed to delete folders for users %s: %v", userID, err)
		}
		if err := h.deleteAllFiles(userID); err != nil {
			return fmt.Errorf("failed to delete files for users %s: %v", userID, err)
		}
	default:
		return fmt.Errorf("unknown users event type: %s", eventType)
	}

	return nil
}

type UserEvent struct {
	UserID    string `json:"userID"`
	EventType string `json:"eventType"`
}

func (h *UserEventHandler) createRootFolder(userID string) error {
	// make userID to uint(not unit64)
	userIDUInt, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		logger.Printf("Failed to convert userID to uint64: %v", err)
		return err
	}

	// create root folder
	folder, err := h.FolderSvc.CreateRootFolder(uint(userIDUInt))
	if err != nil {
		return err
	}
	// logging folder info and userID
	logger.Printf("Created root folder for users %s: %v\n", userID, folder)
	return nil
}

func (h *UserEventHandler) deleteAllFolders(userID string) error {
	logger.Printf("Deleting all folders for users %s\n", userID)
	return nil
}

func (h *UserEventHandler) deleteAllFiles(userID string) error {
	logger.Printf("Deleting all files for users %s\n", userID)
	return nil
}
