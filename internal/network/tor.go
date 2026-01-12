package network

import (
	"fmt"
	"net"
	"time"
)

// RotateIP mengirim sinyal NEWNYM ke port Tor Control (9051) dengan autentikasi
func RotateIP(password string) error {
	conn, err := net.Dial("tcp", "127.0.0.1:9051")
	if err != nil {
		return fmt.Errorf("failed to connect to Tor Control Port: %v", err)
	}
	defer conn.Close()

	// Authenticate using the provided password
	// If password is empty (and CookieAuth is 0), it acts as no-auth or empty string auth
	fmt.Fprintf(conn, "AUTHENTICATE \"%s\"\r\n", password)
	
	// Send Signal to rotate IP
	fmt.Fprintf(conn, "SIGNAL NEWNYM\r\n")
	
	// Properly close session
	fmt.Fprintf(conn, "QUIT\r\n")
	
	return nil
}

// StartRotator menjalankan rotasi IP setiap interval tertentu
func StartRotator(intervalMinutes int, password string) {
	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("[Tor] Requesting new IP...")
		err := RotateIP(password)
		if err != nil {
			fmt.Printf("[Tor] Error rotating IP: %v\n", err)
		} else {
			fmt.Println("[Tor] IP Rotation Signal Sent.")
		}
	}
}
