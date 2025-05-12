package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler is the exported function that Vercel will call for each HTTP request.
func Handler(w http.ResponseWriter, r *http.Request) error {
	// Initialize bot with token from environment variable
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Printf("Failed to initialize bot: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	// Parse incoming Telegram update
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Error decoding update: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}

	// Ignore non-message updates or non-command messages
	if update.Message == nil || !update.Message.IsCommand() {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	// Prepare response message
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		// Get current time in IST
		ist, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			msg.Text = "Error loading IST timezone"
		} else {
			currentTime := time.Now().In(ist).Format("3:04:05 PM MST")
			msg.Text = "Hello! Welcome to the bot.\nThe current time in India is: " + currentTime
		}
	default:
		msg.Text = "Unknown command. Please use /start to see the current time in India."
	}

	// Send response
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	// Acknowledge receipt
	w.WriteHeader(http.StatusOK)
	return nil
}
