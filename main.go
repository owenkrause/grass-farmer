package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

func prompt() bool {
	choices := "y/n"

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", "Generate accounts?", choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return prompt()
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func main() {

	fmt.Println()
	fmt.Println("  ▄████  ██▀███   ▄▄▄        ██████   ██████ ")
	fmt.Println(" ██▒ ▀█▒▓██ ▒ ██▒▒████▄    ▒██    ▒ ▒██    ▒ ")
	fmt.Println("▒██░▄▄▄░▓██ ░▄█ ▒▒██  ▀█▄  ░ ▓██▄   ░ ▓██▄   ")
	fmt.Println("░▓█  ██▓▒██▀▀█▄  ░██▄▄▄▄██   ▒   ██▒  ▒   ██▒")
	fmt.Println("░▒▓███▀▒░██▓ ▒██▒ ▓█   ▓██▒▒██████▒▒▒██████▒▒")
	fmt.Println(" ░▒   ▒ ░ ▒▓ ░▒▓░ ▒▒   ▓▒█░▒ ▒▓▒ ▒ ░▒ ▒▓▒ ▒ ░")
	fmt.Println("  ░   ░   ░▒ ░ ▒░  ▒   ▒▒ ░░ ░▒  ░ ░░ ░▒  ░ ░")
	fmt.Println("░ ░   ░   ░░   ░   ░   ▒   ░  ░  ░  ░  ░  ░  ")
	//fmt.Println("      ░    ░           ░  ░      ░        ░  ")
	fmt.Println()

	targetURL := "wss://proxy.wynd.network:4650/"
	accounts := getAccounts("config/accounts.txt")
	proxies := getProxies()

	if prompt() {
		if len(accounts) > len(proxies) {
			fmt.Println("Please supply a proxy for each account")
			os.Exit(1)
		}
		fmt.Printf("[0/%v] Generating accounts\n", len(accounts))
		for i := 0; i < len(accounts); i++ {
			err := createAccount(accounts[i], proxies[i])
			if err != nil {
				fmt.Println(err)
				i--
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Print("\033[F\033[K")
			fmt.Printf("[%v/%v] Generating accounts\n", i+1, len(accounts))
			time.Sleep(time.Duration(rand.Intn(60-10+1)+10) * time.Second)
		}
	}
	fmt.Println("Hit enter or ctrl+c to exit")

	var wg sync.WaitGroup
	quit := make(chan struct{})

	accounts = getAccounts("config/dist.txt")
	fmt.Printf("\n%v Accounts loaded\n", len(accounts))

	for i := 0; i < len(accounts); i++ {
		temp := strings.SplitN(accounts[i], " ", 3)
		wg.Add(1)
		go createTask(temp[0], temp[1], temp[2], targetURL, &wg, quit)
	}

	fmt.Scanln()
	close(quit)
	wg.Wait()
}
