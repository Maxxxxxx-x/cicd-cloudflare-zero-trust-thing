package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DeploymentConfig struct {
	CIDR       string
	VLanName   string
	BaseDomain string
}

type CFConfig struct {
	AuthToken string
	AccountId string
	TunnelId  string
}

type Config struct {
	Deployment DeploymentConfig
	Cloudflare CFConfig
}

func missingEnv(envName string) error {
	return fmt.Errorf("Missing ENV Variable: %s\n", envName)
}

func getDeploymentConfig() (DeploymentConfig, error) {
	CIDR, found := os.LookupEnv("DEPLOY_CIDR")
	if !found {
		return DeploymentConfig{}, missingEnv("DEPLOY_CIDR")
	}

	vlanName, found := os.LookupEnv("DEPLOY_VLAN_NAME")
	if !found {
		return DeploymentConfig{}, missingEnv("DEPLOY_VLAN_NAME")
	}

	baseDomain, found := os.LookupEnv("DEPLOY_BASE_DOMAIN")
	if !found {
		return DeploymentConfig{}, missingEnv("DEPLOY_BASE_DOMAIN")
	}

	return DeploymentConfig{
		CIDR:       CIDR,
		VLanName:   vlanName,
		BaseDomain: baseDomain,
	}, nil
}

func getCFConfig() (CFConfig, error) {
	authToken, found := os.LookupEnv("CF_AUTH_TOKEN")
	if !found {
		return CFConfig{}, missingEnv("CF_AUTH_TOKEN")
	}

	accountId, found := os.LookupEnv("CF_ACCOUNT_ID")
	if !found {
		return CFConfig{}, missingEnv("CF_ACCOUNT_ID")
	}

	tunnelId, found := os.LookupEnv("CF_TUNNEL_ID")
	if !found {
		return CFConfig{}, missingEnv("CF_TUNNEL_ID")
	}

	return CFConfig{
		AuthToken: authToken,
		AccountId: accountId,
		TunnelId:  tunnelId,
	}, nil
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}
	deployConf, err := getDeploymentConfig()
	if err != nil {
		return Config{}, err
	}

	cfConf, err := getCFConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Deployment: deployConf,
		Cloudflare: cfConf,
	}, nil
}
