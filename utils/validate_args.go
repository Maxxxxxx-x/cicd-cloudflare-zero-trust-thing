package utils

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Maxxxxxx-x/cicd-cloudflare-zero-trust-thing/config"
	"github.com/Maxxxxxx-x/cicd-cloudflare-zero-trust-thing/models"
)

func validateIPAddr(config config.Config, IpAddr string) error {
	ip := net.ParseIP(IpAddr)
	if ip.To4() == nil && ip.To16() == nil {
		return fmt.Errorf("IP is not a valid IP address. Received: %s\n", IpAddr)
	}
	if !ip.IsPrivate() {
		return fmt.Errorf("%s is not a private IP  addres\n", IpAddr)
	}

	_, subnet, err := net.ParseCIDR(config.Deployment.CIDR)
	if err != nil {
		return fmt.Errorf("Invalid DEPLOY_CIDR. Received: %s\n", config.Deployment.CIDR)
	}
	if !subnet.Contains(ip) {
		return fmt.Errorf("%s is not in subnet %s\n", IpAddr, subnet.String())
	}

	splitted := strings.Split(IpAddr, ".")
	lastOctet, err := strconv.Atoi(splitted[len(splitted)-1])
	if err != nil {
		return err
	}

	if lastOctet <= 1 || lastOctet >= 255 {
		return fmt.Errorf("%s is not allowed! Contact Admin for allowed IP addresses.\n", IpAddr)
	}

	return nil
}

func validatePortNumber(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("Failed to convert %s to number. Error: %s\n", port, err.Error())
	}

	if portNum < 1 || portNum > 65535 {
		return errors.New("Invalid port number. Port number should be in range 1 - 65535")
	}

	return nil
}

func validateBaseDomain(domain string) error {
	splitted := strings.Split(domain, ".")
	if len(splitted) < 2 {
		return errors.New("DEPLOY_BASE_DOMAIN does not contain a valid TLD")
	}

	return nil
}

func ValidateInputs(config config.Config, inputs *models.Input) error {
	if err := validateBaseDomain(config.Deployment.BaseDomain); err != nil {
		return err
	}

	if err := validateIPAddr(config, inputs.DockerIP); err != nil {
		return err
	}

	if err := validatePortNumber(inputs.DockerPort); err != nil {
		return err
	}

	return nil
}
