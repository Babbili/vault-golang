package main

import (
	"context"
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

func LoginWithAppRole() (string, error) {
	// vault client config
	// ref https://pkg.go.dev/github.com/hashicorp/vault/api@v1.14.0#DefaultConfig
	config := vault.DefaultConfig()

	client, err := vault.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	roleID := os.Getenv("ROLE_ID")
	if roleID == "" {
		return "", fmt.Errorf("no role ID was provided in ROLE_ID env var")
	}
	secretID := &auth.SecretID{FromFile: "env/secret-id"}

	appRoleAuth, err := auth.NewAppRoleAuth(
		roleID,
		secretID,
	)
	if err != nil {
		return "", fmt.Errorf("unable to initialize AppRole auth method: %w", err)
	}

	// login to Vault with approle auth method
	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return "", fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return "", fmt.Errorf("no auth info was returned after login")
	}

	return "login successfully", nil
}

func main() {

	login, err := LoginWithAppRole()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("logged in to Vault with approle auth %s", login)
	}
	time.Sleep(1 * time.Hour)
}
