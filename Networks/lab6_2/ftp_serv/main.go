package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jlaffaye/ftp"
)

var ftpClient *ftp.ServerConn
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Структура для входящих JSON-команд
type CommandMessage struct {
	Server   string   `json:"server"`
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Command  string   `json:"command"`
	Args     []string `json:"args"`
}

func main() {
	fmt.Println("WebSocket сервер запущен на порту 8081.")
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка установки WebSocket-соединения:", err)
		return
	}
	defer conn.Close()

	// Ожидаем FTP-соединение
	for {
		// Читаем JSON-команду от клиента
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			break
		}

		var cmdMsg CommandMessage
		err = json.Unmarshal(msg, &cmdMsg)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Ошибка парсинга JSON: "+err.Error()))
			continue
		}

		// Если получено FTP-соединение
		if cmdMsg.Command == "connect" {
			err := connectToFTPServer(cmdMsg.Server, cmdMsg.Login, cmdMsg.Password)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Ошибка подключения или входа в FTP: "+err.Error()))
				return
			}
			conn.WriteMessage(websocket.TextMessage, []byte("Connection established"))
			break
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte("Ожидается команда 'connect' для установления FTP-соединения"))
		}
	}

	// После успешного подключения начинаем выполнять команды FTP
	for {
		// Читаем JSON-команду от клиента
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			break
		}

		var cmdMsg CommandMessage
		err = json.Unmarshal(msg, &cmdMsg)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Ошибка парсинга JSON: "+err.Error()))
			continue
		}

		// Выполняем команду
		response := executeCommand(cmdMsg.Command, cmdMsg.Args)
		err = conn.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			log.Println("Ошибка отправки сообщения:", err)
			break
		}
	}
}

func connectToFTPServer(server, login, password string) error {
	var err error
	// Добавляем порт, если не указан
	if !strings.Contains(server, ":") {
		server += ":21"
	}

	ftpClient, err = ftp.Dial(server, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("не удалось подключиться к серверу: %v", err)
	}

	err = ftpClient.Login(login, password)
	if err != nil {
		return fmt.Errorf("не удалось войти в систему: %v", err)
	}
	return nil
}

func executeCommand(command string, args []string) string {
	// Проверка подключения перед выполнением команды
	if ftpClient == nil {
		return "Нет активного FTP-подключения"
	}

	switch command {
	case "mkdir":
		if len(args) != 1 {
			return "Неверные аргументы: mkdir <directory_name>"
		}
		return createDirectory(ftpClient, args[0])

	case "ls":
		dirPath := "."
		if len(args) > 0 {
			dirPath = args[0]
		}
		return listDirectory(ftpClient, dirPath)

	case "cd":
		if len(args) != 1 {
			return "Неверные аргументы: cd <directory_name>"
		}
		return changeDirectory(ftpClient, args[0])

	case "delete":
		if len(args) != 1 {
			return "Использование: delete <file_name>"
		}
		return deleteFile(ftpClient, args[0])

	case "rmdir":
		if len(args) != 1 {
			return "Неверные аргументы: rmdir <directory_name>"
		}

		return removeDirectory(ftpClient, args[0])

	case "rmdir_rec":
		if len(args) != 1 {
			return "Неверные аргументы: rmdir_rec <directory_name>"
		}
		err := removeDirectoryRecursive(ftpClient, args[0])
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("Директория рекурсивно удалена: %s", args[0])

	case "help":
		return `Список команд:
 mkdir <directory_name> - Создать директорию
 ls [directory_name] - Показать содержимое директории
 cd <directory_name> - Перейти в директорию
 delete <file_name> - Удалить файл
 rmdir <directory_name> - Удалить пустую директорию
 rmdir_rec <directory_name> - Рекурсивное удаление директории
 help - Показать список команд`

	default:
		return "Неизвестная команда. Введите 'help' для списка команд."
	}
}

func createDirectory(c *ftp.ServerConn, dirName string) string {
	err := c.MakeDir(dirName)
	if err != nil {
		return fmt.Sprintf("Ошибка создания директории: %v", err)
	}
	return "Директория создана: " + dirName
}

func listDirectory(c *ftp.ServerConn, dirPath string) string {
	entries, err := c.List(dirPath)
	if err != nil {
		return fmt.Sprintf("Ошибка получения содержимого директории: %v", err)
	}

	var result strings.Builder
	for _, entry := range entries {
		// Пропускаем служебные директории
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		// Определяем тип элемента
		entryType := "File"
		if entry.Type == ftp.EntryTypeFolder {
			entryType = "Directory"
		}

		// Добавляем более подробную информацию
		result.WriteString(fmt.Sprintf(
			"Type: %s, Name: %s, Size: %d bytes, Permissions: %s, Modified: %s\n",
			entryType,
			entry.Name,
			entry.Size,
			formatPermissions(entry),
			entry.Time.Format("2006-01-02 15:04:05"),
		))
	}
	return result.String()
}

// Вспомогательная функция для форматирования прав доступа
func formatPermissions(entry *ftp.Entry) string {
	// Базовая логика определения прав доступа
	permStr := "----------"

	// Определение типа файла
	if entry.Type == ftp.EntryTypeFolder {
		permStr = "d---------"
	} else if entry.Type == ftp.EntryTypeLink {
		permStr = "l---------"
	}

	// Здесь вы можете добавить более сложную логику интерпретации прав доступа,
	// если FTP-библиотека предоставляет такую информацию
	return permStr
}

func changeDirectory(c *ftp.ServerConn, dirName string) string {
	err := c.ChangeDir(dirName)
	if err != nil {
		return fmt.Sprintf("Ошибка перехода в директорию: %v", err)
	}
	return "Перешли в директорию: " + dirName
}

func deleteFile(c *ftp.ServerConn, fileName string) string {
	err := c.Delete(fileName)
	if err != nil {
		return fmt.Sprintf("Ошибка удаления файла: %v", err)
	}
	return "Файл удален: " + fileName
}

func removeDirectory(c *ftp.ServerConn, dirName string) string {
	err := c.RemoveDir(dirName)
	if err != nil {
		return fmt.Sprintf("Ошибка удаления директории: %v", err)
	}
	return fmt.Sprintf("Директория удалена: %s", dirName)
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
	return nil
}
