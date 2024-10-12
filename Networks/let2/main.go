package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	// Используем команду tcpdump для захвата всех пакетов (-l для немедленного вывода, -A для ASCII)
	cmd := exec.Command("tcpdump", "-l", "-A")

	// Создаем pipe для чтения данных
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Ошибка создания pipe:", err)
		return
	}

	// Запускаем команду tcpdump
	if err := cmd.Start(); err != nil {
		fmt.Println("Ошибка запуска tcpdump:", err)
		return
	}

	// Чтение данных построчно
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		// Если строка содержит "key" и фамилию "Krasnobaev"
		if strings.Contains(line, "key") && strings.Contains(line, "Krasnobaev") {
			if containsHash(line) {
				hash := extractHash(line)
				if hash != "" {
					fmt.Println("Найден хэш пользователя:", hash)

					// Формируем GET-запрос для получения пароля по хэшу
					password, err := getPasswordByHash(hash)
					if err != nil {
						fmt.Println("Ошибка получения пароля:", err)
						continue
					}
					fmt.Println("Полученный пароль:", password)

					// Отправляем GET-запрос с найденным паролем
					err = sendPassword(password)
					if err != nil {
						fmt.Println("Ошибка отправки пароля:", err)
					} else {
						fmt.Println("Пароль успешно отправлен")
					}
				}
			}
		}
	}

	// Обрабатываем ошибки сканера
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения данных tcpdump:", err)
	}
}

// Функция для отправки запроса и получения пароля по хэшу
func getPasswordByHash(hash string) (string, error) {
	// Формируем URL для получения пароля
	url := fmt.Sprintf("http://pstgu.yss.su/iu9/networks/let1_2024/getkey.php?hash=%s", hash)

	// Отправляем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Предположим, что пароль находится в виде "pass: s8dfxdvx"
	response := string(body)
	if strings.HasPrefix(response, "pass: ") {
		return strings.TrimPrefix(response, "pass: "), nil
	}

	return "", fmt.Errorf("пароль не найден в ответе")
}

// Функция для отправки пароля
func sendPassword(password string) error {
	// Формируем URL для отправки пароля
	fio := "Михаил%20Краснобаев"
	url := fmt.Sprintf("http://pstgu.yss.su/iu9/networks/let1_2024/send_from_go.php?subject=let1&fio=%s&pass=%s", fio, password)

	// Отправляем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем, что запрос был успешным
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось отправить пароль, код ответа: %d", resp.StatusCode)
	}

	return nil
}

// Пример функции для поиска хэша в строке
func containsHash(line string) bool {
	// Логика поиска хэша (упрощенная)
	return true
}

// Пример функции для извлечения хэша из строки
func extractHash(line string) string {
	// Простой пример поиска хэша по длине (MD5 - 32 символа, SHA1 - 40 символов)
	words := bufio.NewScanner(strings.NewReader(line))
	words.Split(bufio.ScanWords)

	for words.Scan() {
		word := words.Text()
		if len(word) == 32 || len(word) == 40 { // MD5 или SHA1
			return word
		}
	}
	return ""
}
