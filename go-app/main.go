package main

import (
	"context"
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

func GetSecretWithAppRole() (string, error) {
	// vault client config
	// ref https://pkg.go.dev/github.com/hashicorp/vault/api@v1.14.0#DefaultConfig
	config := vault.DefaultConfig()
	config.Address = "http://vault.vault.svc.cluster.local:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	// Vault's approle auth role_id & secret_id
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

	// get secret from the default mount path for KV v2 "go-app/secret"
	// ref https://pkg.go.dev/github.com/hashicorp/vault/api@v1.14.0#KVv2
	secret, err := client.KVv2("go-app/secret").Get(context.Background(), "creds")
	if err != nil {
		return "", fmt.Errorf("unable to read secret: %w", err)
	}

	// data map can contain more than one key-value pair,
	// in this case we're just grabbing `username`
	value, ok := secret.Data["usename"].(string)
	if !ok {
		return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data["username"], secret.Data["username"])
	}

	return fmt.Sprintf("username is %s", value), nil
}

func main() {

	secret, err := GetSecretWithAppRole()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(secret)
	}
	time.Sleep(1 * time.Hour)
}
