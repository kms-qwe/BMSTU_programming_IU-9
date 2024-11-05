package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	address := flag.String("address", "localhost", "SSH server address")
	port := flag.String("port", "22", "SSH server port")
	user := flag.String("user", "testuser", "SSH user")
	password := flag.String("password", "password123", "SSH password")

	flag.Parse()
	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.Password(*password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", *address+":"+*port, config)
	if err != nil {
		fmt.Println("Failed to dial SSH:", err)
		return
	}
	defer client.Close()

	fmt.Println("Enter commands to execute on the server ('exit' to quit):")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		command := scanner.Text()
		if command == "exit" {
			break
		}

		session, err := client.NewSession()
		if err != nil {
			fmt.Println("Failed to create session:", err)
			continue
		}

		output, err := session.CombinedOutput(command)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
		fmt.Println(string(output))

		session.Close()
	}
}
