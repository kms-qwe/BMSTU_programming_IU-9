package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

func main() {
	c, err := ftp.Dial("students.yss.su:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal("Ошибка подключения к серверу:", err)
	}
	defer c.Quit()

	err = c.Login("ftpiu8", "3Ru7yOTA")
	if err != nil {
		log.Fatal("Ошибка входа на сервер:", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("FTP клиент готов к работе. Введите команду (help для списка команд).")

	for {
		fmt.Print("$ ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "upload":
			if len(args) != 2 {
				fmt.Println("Неверные аргументы: upload <local_file> <remote_file>")
				continue
			} else {
				uploadFile(c, args[0], args[1])
			}
		case "download":
			if len(args) != 2 {
				fmt.Println("Неверные аргументы: download <remote_file> <local_file>")
			} else {
				downloadFile(c, args[0], args[1])
			}
		case "mkdir":
			if len(args) != 1 {
				fmt.Println("Неверные аргументы: mkdir <directory_name>")
			} else {
				createDirectory(c, args[0])
			}
		case "ls":
			if len(args) != 1 {
				fmt.Println("Неверные аргументы: ls <directory_name>")
			} else {
				listDirectoryContents(c, args[0])
			}
		case "cd":
			if len(args) != 1 {
				fmt.Println("Неверные аргументы: cd <directory_name>")
			} else {
				changeDirectory(c, args[0])
			}
		case "delete":
			if len(args) != 1 {
				fmt.Println("Использование: delete <file_name>")
			} else {
				deleteFile(c, args[0])
			}
		case "rmdir":
			if len(args) != 1 {
				fmt.Println("Неверные аргументы: rmdir <directory_name>")
			} else {
				removeDirectory(c, args[0])
			}
		case "rmdir_rec":
			if len(args) != 1 {
				fmt.Println("Неверные аргументы: rmdir_rec <directory_name>")
			} else {
				err := removeDirectoryRecursive(c, args[0])
				if err != nil {
					fmt.Println("Ошибка при удалении директории рекурсивно:", err)
				}
			}
		case "help":
			printHelp()
		case "exit":
			fmt.Println("Выход из FTP клиента.")
			return
		default:
			fmt.Println("Неизвестная команда. Введите 'help' для списка команд.")
		}
	}
}

func printHelp() {
	fmt.Println("Список команд:")
	fmt.Println(" upload <local_file> <remote_file> - Загрузить файл на сервер")
	fmt.Println(" download <remote_file> <local_file> - Скачать файл с сервера")
	fmt.Println(" mkdir <directory_name> - Создать директорию")
	fmt.Println(" ls <directory_name> - Показать содержимое директории")
	fmt.Println(" cd <directory_name> - Перейти в директорию")
	fmt.Println(" delete <file_name> - Удалить файл")
	fmt.Println(" rmdir <directory_name> - Удалить пустую директорию")
	fmt.Println(" rmdir_rec <directory_name> - Рекурсивное удаление директории")
	fmt.Println(" exit - Выйти из программы")
}

func uploadFile(c *ftp.ServerConn, localFileName, remoteFileName string) {
	file, err := os.Open(localFileName)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	err = c.Stor(remoteFileName, file)
	if err != nil {
		fmt.Println("Ошибка загрузки файла:", err)
		return
	}
	fmt.Println("Файл успешно загружен:", remoteFileName)
}

func downloadFile(c *ftp.ServerConn, remoteFileName, localFileName string) {
	r, err := c.Retr(remoteFileName)
	if err != nil {
		fmt.Println("Ошибка скачивания файла:", err)
		return
	}
	defer r.Close()

	localFile, err := os.Create(localFileName)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, r)
	if err != nil {
		fmt.Println("Ошибка копирования содержимого:", err)
		return
	}
	fmt.Println("Файл успешно скачан:", localFileName)
}

func createDirectory(c *ftp.ServerConn, dirName string) {
	err := c.MakeDir(dirName)
	if err != nil {
		fmt.Println("Ошибка создания директории:", err)
		return
	}
	fmt.Println("Директория создана:", dirName)
}

func listDirectoryContents(c *ftp.ServerConn, dirName string) {
	entries, err := c.List(dirName)
	if err != nil {
		fmt.Println("Ошибка получения содержимого директории:", err)
		return
	}
	fmt.Println("Содержимое директории", dirName)
	for _, entry := range entries {
		fmt.Println(entry.Type, ":", entry.Name)
	}
}

func changeDirectory(c *ftp.ServerConn, dirName string) {
	err := c.ChangeDir(dirName)
	if err != nil {
		fmt.Println("Ошибка перехода в директорию:", err)
		return
	}
	fmt.Println("Перешли в директорию:", dirName)
}

func deleteFile(c *ftp.ServerConn, fileName string) {
	err := c.Delete(fileName)
	if err != nil {
		fmt.Println("Ошибка удаления файла:", err)
		return
	}
	fmt.Println("Файл удален:", fileName)
}

func removeDirectory(c *ftp.ServerConn, dirName string) {
	err := c.RemoveDir(dirName)
	if err != nil {
		fmt.Println("Ошибка удаления директории:", err)
		return
	}
	fmt.Println("Директория удалена:", dirName)
}

func removeDirectoryRecursive(c *ftp.ServerConn, dirName string) error {
	entries, err := c.List(dirName)
	if err != nil {
		return fmt.Errorf("ошибка получения содержимого директории: %w", err)
	}

	for _, entry := range entries {
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		path := dirName + "/" + entry.Name
		if entry.Type == ftp.EntryTypeFolder {
			if err := removeDirectoryRecursive(c, path); err != nil {
				return err
			}
		} else {
			if err := c.Delete(path); err != nil {
				return fmt.Errorf("ошибка удаления файла: %w", err)
			}
		}
	}

	err = c.RemoveDir(dirName)
	if err != nil {
		return fmt.Errorf("ошибка удаления директории: %w", err)
	}
	fmt.Println("Директория рекурсивно удалена:", dirName)
	return nil
}
