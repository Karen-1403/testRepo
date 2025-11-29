package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/zGate-Team/zGate-Platform/internal/store"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// Initialize store
	keyHex := os.Getenv("ZGATE_STORE_KEY")
	if keyHex == "" {
		fmt.Println("Error: ZGATE_STORE_KEY environment variable is required")
		os.Exit(1)
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		fmt.Printf("Error decoding key: %v\n", err)
		os.Exit(1)
	}

	dbPath := os.Getenv("ZGATE_STORE_PATH")
	if dbPath == "" {
		dbPath = "data/zgate.db"
	}

	dataStore, err := store.NewStore(dbPath, key)
	if err != nil {
		fmt.Printf("Error initializing store: %v\n", err)
		os.Exit(1)
	}
	defer dataStore.Close()

	// Check command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		if len(os.Args) < 5 {
			fmt.Println("Usage: go run test_admin.go create <username> <password> <name> [email]")
			os.Exit(1)
		}
		username := os.Args[2]
		password := os.Args[3]
		name := os.Args[4]
		email := ""
		if len(os.Args) > 5 {
			email = os.Args[5]
		}

		err := dataStore.CreateAdmin(username, password, name, email)
		if err != nil {
			fmt.Printf("Error creating admin: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Admin '%s' created successfully!\n", username)

	case "list":
		admins, err := dataStore.ListAdmins()
		if err != nil {
			fmt.Printf("Error listing admins: %v\n", err)
			os.Exit(1)
		}

		if len(admins) == 0 {
			fmt.Println("No admins found")
			return
		}

		fmt.Printf("\nFound %d admin(s):\n\n", len(admins))
		fmt.Printf("%-20s %-30s %-30s %-20s\n", "USERNAME", "NAME", "EMAIL", "CREATED AT")
		fmt.Println("----------------------------------------------------------------------------------------------------")

		for _, admin := range admins {
			email := admin.Email
			if email == "" {
				email = "-"
			}
			fmt.Printf("%-20s %-30s %-30s %-20s\n",
				admin.Username,
				admin.Name,
				email,
				admin.CreatedAt.Format("2006-01-02 15:04:05"),
			)
		}
		fmt.Println()

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run test_admin.go delete <username>")
			os.Exit(1)
		}
		username := os.Args[2]

		err := dataStore.DeleteAdmin(username)
		if err != nil {
			fmt.Printf("Error deleting admin: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Admin '%s' deleted successfully!\n", username)

	case "verify":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run test_admin.go verify <username> <password>")
			os.Exit(1)
		}
		username := os.Args[2]
		password := os.Args[3]

		valid, err := dataStore.VerifyAdminPassword(username, password)
		if err != nil {
			fmt.Printf("Error verifying password: %v\n", err)
			os.Exit(1)
		}

		if valid {
			fmt.Printf("✓ Password is correct for admin '%s'\n", username)
		} else {
			fmt.Printf("✗ Password is incorrect for admin '%s'\n", username)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Admin CRUD Test Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  go run test_admin.go <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  create <username> <password> <name> [email]  - Create a new admin")
	fmt.Println("  list                                          - List all admins")
	fmt.Println("  delete <username>                             - Delete an admin")
	fmt.Println("  verify <username> <password>                  - Verify admin credentials")
	fmt.Println("\nExamples:")
	fmt.Println("  go run test_admin.go create admin mypass123 \"Administrator\" admin@example.com")
	fmt.Println("  go run test_admin.go list")
	fmt.Println("  go run test_admin.go verify admin mypass123")
	fmt.Println("  go run test_admin.go delete admin")
}
