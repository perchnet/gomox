package util

import (
	"github.com/luthermonson/go-proxmox"
)

func InstantiateClient(pveUrl string, credentials proxmox.Credentials) proxmox.Client {

	client := proxmox.NewClient(
		pveUrl,
		proxmox.WithCredentials(&credentials),
	)
	return *client
}
