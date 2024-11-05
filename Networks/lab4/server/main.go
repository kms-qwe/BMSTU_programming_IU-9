package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gliderlabs/ssh"
)

func main() {
	server := ssh.Server{
		Addr: "localhost:8080",
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			return ctx.User() == "testuser" && password == "password123"
		},
	}

	server.Handle(func(s ssh.Session) {
		cmd := s.Command()

		if len(cmd) == 0 {
			fmt.Fprintln(s, "No command provided.")
			return
		}

		switch cmd[0] {
		case "mkdir":

			if len(cmd) > 1 {
				err := os.Mkdir(cmd[1], 0755)
				if err != nil {
					fmt.Fprintln(s, "Error creating directory:", err)
				} else {
					fmt.Fprintln(s, "Directory created:", cmd[1])
				}
			}
		case "rmdir":
			if len(cmd) > 1 {
				err := os.Remove(cmd[1])
				if err != nil {
					fmt.Fprintln(s, "Error deleting directory:", err)
				} else {
					fmt.Fprintln(s, "Directory deleted:", cmd[1])
				}
			}
		case "ls":
			dir := "."
			if len(cmd) > 1 {
				dir = cmd[1]
			}
			files, err := os.ReadDir(dir)
			if err != nil {
				fmt.Fprintln(s, "Error reading directory:", err)
			} else {
				for _, file := range files {
					fmt.Fprintln(s, file.Name())
				}
			}
		case "ping":
			if len(cmd) > 1 {
				out, err := exec.Command("ping", "-c", "3", cmd[1]).CombinedOutput()
				if err != nil {
					fmt.Fprintln(s, "Ping error:", err)
				} else {
					fmt.Fprintln(s, string(out))
				}
			}
		default:
			fmt.Fprintln(s, "Unknown command:", cmd)
		}
	})

	fmt.Println("Starting SSH server on locahost...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
