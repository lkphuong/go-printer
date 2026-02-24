package main

import (
	"fmt"
	"go-printer/internal/app"
	"net"
)

func printIP() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Lỗi lấy interfaces:", err)
		return
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP

			if ip.To4() != nil && !ip.IsLoopback() {
				fmt.Printf("Interface: %s - IPv4: %s\n", iface.Name, ip.String())
			}
		}
	}
}

func main() {
	a := app.NewApp()
	printIP()
	a.Run()
}
