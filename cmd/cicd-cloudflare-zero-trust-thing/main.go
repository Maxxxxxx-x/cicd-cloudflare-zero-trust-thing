package main

import (
	"fmt"
	"log"
	"os"

	c "github.com/Maxxxxxx-x/cicd-cloudflare-zero-trust-thing/config"
	"github.com/Maxxxxxx-x/cicd-cloudflare-zero-trust-thing/models"
	"github.com/Maxxxxxx-x/cicd-cloudflare-zero-trust-thing/utils"
)

/*
/bin/cicd-cf-xxx branchName dockerIp dockerPort [hostname](default: branchName)
branchName = gitlab CICD
dockerIP = IP addr
dockerPort = port
hostname (default = branch name)
*/

func main() {
	inputs := models.Input{}
	rawArgs := os.Args[1:]

	lenArgs := len(rawArgs)
	if lenArgs < 3 {
		log.Println("Missing argument!\nUsage: ./docker-cicd branch-name docker-ip docker-port [hostname]")
		return
	}
	if lenArgs > 4 {
		log.Println("Too many arguments!\nUsage: ./docker-cicd branch-name docker-ip docker-port [hostname]")
		return
	}

	inputs.BranchName = rawArgs[0]
	inputs.DockerIP = rawArgs[1]
	inputs.DockerPort = rawArgs[2]

	config, err := c.LoadConfig()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	if err := utils.ValidateInputs(config, &inputs); err != nil {
		log.Fatal(err.Error())
		return
	}

	if lenArgs == 4 {
		inputs.Hostname = fmt.Sprintf("%s.%s", rawArgs[3], config.Deployment.BaseDomain)
	} else {
		inputs.Hostname = fmt.Sprintf("%s.%s", inputs.BranchName, config.Deployment.BaseDomain)
	}

	isInUse, err := utils.IsAddrInUse(config.Deployment.VLanName, inputs.DockerIP)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if isInUse {
		log.Fatalf("IP: %s is in use!\n", inputs.DockerIP)
	}
	log.Printf("IP: %s Not in use!\n", inputs.DockerIP)

	fmt.Printf("Branch Name: %s | Docker IP: %s | Docker Port: %s | Hostname: %s\n", inputs.BranchName, inputs.DockerIP, inputs.DockerPort, inputs.Hostname)
	err = utils.SetZeroTrust(config, inputs)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
