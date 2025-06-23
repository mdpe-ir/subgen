package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"subgen/internal/admin"
	"subgen/internal/db"
	"subgen/internal/userlink"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func printHelp() {
	fmt.Print(`Usage: subgen <command> [options]

Commands:
  create admin         Create a new admin (interactive)
  list admin           List all admins
  delete admin <user>  Delete admin by username
  gen id               Generate and save a new UUID
  help                 Show this help message
`)
}

func Execute() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		printHelp()
		os.Exit(0)
	}

	db.InitDB()
	db.Migrate()

	switch os.Args[1] {
	case "create":
		if len(os.Args) >= 3 && os.Args[2] == "admin" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)
			fmt.Print("Enter password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)
			if username == "" || password == "" {
				fmt.Println("Username and password cannot be empty.")
				return
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println("Failed to hash password:", err)
				return
			}
			adminRec := admin.Admin{Username: username, Password: string(hash)}
			if err := db.DB.Create(&adminRec).Error; err != nil {
				fmt.Println("Failed to create admin:", err)
			} else {
				fmt.Println("Admin created successfully.")
			}
		} else {
			fmt.Println("Unknown create target.")
		}
	case "list":
		if len(os.Args) >= 3 && os.Args[2] == "admin" {
			var admins []admin.Admin
			if err := db.DB.Find(&admins).Error; err != nil {
				fmt.Println("Failed to list admins:", err)
				return
			}
			fmt.Println("Admins:")
			for _, a := range admins {
				fmt.Printf("- %s\n", a.Username)
			}
		} else {
			fmt.Println("Unknown list target.")
		}
	case "delete":
		if len(os.Args) >= 4 && os.Args[2] == "admin" {
			username := os.Args[3]
			if username == "" {
				fmt.Println("Username required.")
				return
			}
			if err := db.DB.Where("username = ?", username).Delete(&admin.Admin{}).Error; err != nil {
				fmt.Println("Failed to delete admin:", err)
			} else {
				fmt.Println("Admin deleted successfully.")
			}
		} else {
			fmt.Println("Unknown delete target.")
		}
	case "gen":
		if len(os.Args) >= 3 && os.Args[2] == "id" {
			newUUID := uuid.NewString()
			rec := userlink.UUID{UUID: newUUID}
			if err := db.DB.Create(&rec).Error; err != nil {
				fmt.Println("Failed to save UUID:", err)
				return
			}
			f, err := os.Create("uuid.txt")
			if err == nil {
				f.WriteString(newUUID)
				f.Close()
			}
			fmt.Println("Generated and saved UUID:", newUUID)
		} else {
			fmt.Println("Unknown gen target.")
		}
	default:
		printHelp()
	}
}
