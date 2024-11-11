package utils

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Maxxxxxx-x/cicd-cloudflare-zero-trust-thing/models"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func getNetworkIdByName(apiClient *client.Client, name string) (string, error) {
	networkFilter := network.ListOptions{
		Filters: filters.NewArgs(filters.Arg("name", name)),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	networkList, err := apiClient.NetworkList(ctx, networkFilter)
	if err != nil {
		return "", err
	}

	if len(networkList) != 1 {
		return "", fmt.Errorf("Failed to find network by name: %s\n", name)
	}

	return networkList[0].ID, nil
}

func getContainersFromNetwork(apiClient *client.Client, networkId string) (models.AddressMap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	network, err := apiClient.NetworkInspect(ctx, networkId, network.InspectOptions{})
	if err != nil {
		return nil, err
	}

	addrMap := make(models.AddressMap)
	for containerId, container := range network.Containers {
		if container.IPv4Address == "" && container.IPv6Address == "" {
			continue
		}

		if ip, _, _ := net.ParseCIDR(container.IPv4Address); ip != nil {
			addrMap[ip.String()] = models.ContainerData{
				Id:   containerId,
				Name: container.Name,
			}
		}

		if ip, _, _ := net.ParseCIDR(container.IPv6Address); ip != nil {
			addrMap[ip.String()] = models.ContainerData{

				Id:   containerId,
				Name: container.Name,
			}
		}
	}

	return addrMap, nil
}

func IsAddrInUse(networkName string, ipAddr string) (bool, error) {
	apiClient, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		return false, err
	}
	defer apiClient.Close()

	networkId, err := getNetworkIdByName(apiClient, networkName)
	if err != nil {
		return false, err
	}

	addrMap, err := getContainersFromNetwork(apiClient, networkId)
	if err != nil {
		return false, err
	}

	if len(addrMap) == 0 {
		return false, fmt.Errorf("Address map is empty")
	}

	_, exists := addrMap[ipAddr]
	return exists, nil
}
