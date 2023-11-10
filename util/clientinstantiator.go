package util

import (
	"github.com/luthermonson/go-proxmox"
)

func InstantiateClient(pveurl string, credentials proxmox.Credentials) proxmox.Client {

	client := proxmox.NewClient(
		pveurl,
		proxmox.WithCredentials(&credentials),
	)
	return *client
}
