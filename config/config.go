package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AuthConfig struct {
    AuthEmail string
    AuthKey string
    AuthToken string
}

type AccountConfig struct {
    AccountId string
    TunnelId string
}

type Config struct {
    Account AccountConfig
    Auth AuthConfig
    BaseUrl string
}

func missingEnv(name string) string {
    return fmt.Sprintf("Missing ENV variable: %s", name)
}


func getAccountConfig() AccountConfig {
    accountId, found := os.LookupEnv("")
}


func getAuthConfig() AuthConfig {
    email, found := os.LookupEnv("AUTH_EMAIL")
    if !found {
        log.Fatal(missingEnv("AUTH_EMAIL"))
    }

    key, found := os.LookupEnv("AUTH_KEY")
    if !found {
        log.Fatal(missingEnv("AUTH_KEY"))
    }

    token, found := os.LookupEnv("AUTH_TOKEN")
    if !found {
        log.Fatal(missingEnv("AUTH_TOKEN"))
    }

    return AuthConfig{
        AuthEmail: email,
        AuthKey: key,
        AuthToken: token,
    }
}


func GetConfig() Config {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Failed to load env %s", err.Error())
    }

}
