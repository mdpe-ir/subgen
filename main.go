package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Read content from configs.txt
	configContent, err := os.ReadFile("configs.txt")
	if err != nil {
		log.Fatalf("Error reading configs.txt: %v", err)
	}
	base64Config := base64.StdEncoding.EncodeToString(configContent)

	// Read content from configs_desktop.txt
	configDesktopContent, err := os.ReadFile("configs_desktop.txt")
	if err != nil {
		log.Fatalf("Error reading configs_desktop.txt: %v", err)
	}
	base64ConfigDesktop := base64.StdEncoding.EncodeToString(configDesktopContent)

	// Read UUID from uuid.txt
	uuidBytes, err := os.ReadFile("uuid.txt")
	if err != nil {
		log.Fatalf("Error reading uuid.txt: %v", err)
	}
	uuid := strings.TrimSpace(string(uuidBytes))

	// Load whitelist.txt
	whitelist := loadWhitelist("whitelist.txt")

	// Set up HTTP handler for UUID endpoint
	http.HandleFunc("/"+uuid, func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		uaLower := strings.ToLower(ua)
		log.Printf("Config requested from %s (User-Agent: %s)", r.RemoteAddr, ua)
		w.Header().Set("Content-Type", "text/plain")

		// Check if UA is whitelisted
		isWhitelisted := false
		for entry := range whitelist {
			if strings.Contains(ua, entry) {
				isWhitelisted = true
				break
			}
		}

		// Check if UA is a desktop
		isDesktop := (strings.Contains(uaLower, "nekoray") ||
			(strings.Contains(uaLower, "v2rayn") && !strings.Contains(uaLower, "v2rayng")) ||
			strings.Contains(uaLower, "windows") ||
			(strings.Contains(uaLower, "linux") && !strings.Contains(uaLower, "android")) ||
			strings.Contains(uaLower, "mac os") || strings.Contains(uaLower, "macos") || strings.Contains(uaLower, "darwin")) &&
			!strings.Contains(uaLower, "android") &&
			!strings.Contains(uaLower, "ios") &&
			!strings.Contains(uaLower, "hiddify")

		if isDesktop || isWhitelisted {
			fmt.Fprint(w, base64ConfigDesktop)
		} else {
			fmt.Fprint(w, base64Config)
		}
	})

	// Default route returns 404 for security
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Start the server
	port := "8095"
	log.Printf("Server starting on :%s", port)
	log.Printf("Access your config at http://localhost:%s/%s", port, uuid)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Load whitelist.txt into a map
func loadWhitelist(filename string) map[string]bool {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("No whitelist.txt found or unable to read it: %v", err)
		return map[string]bool{}
	}
	defer file.Close()

	whitelist := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			whitelist[line] = true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading whitelist.txt: %v", err)
	}
	return whitelist
}
