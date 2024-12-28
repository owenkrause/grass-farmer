package main

import (
	"bufio"
	"log"
	"os"
)

func getProxies() []string {

	file, err := os.Open("config/proxies.txt")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()
	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
	}
	return proxies
}
