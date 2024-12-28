package main

import (
	"bufio"
	"log"
	"os"
)

func getAccounts(filePath string) []string {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()
	var accounts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		accounts = append(accounts, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
	}
	return accounts
}
