package main

import (
	"log"
	"net/http"
)

// Функция для обработки главной страницы (форма ввода данных)
func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FTP Client</title>
</head>
<body>
    <h2>FTP Client Connection</h2>
    <form id="ftpForm">
        <label for="ftpServer">FTP Server:</label><br>
        <input type="text" id="ftpServer" name="ftpServer" required><br><br>
        <label for="ftpLogin">FTP Login:</label><br>
        <input type="text" id="ftpLogin" name="ftpLogin" required><br><br>
        <label for="ftpPassword">FTP Password:</label><br>
        <input type="password" id="ftpPassword" name="ftpPassword" required><br><br>
        <button type="submit">Connect</button>
    </form>

    <script>
        const form = document.getElementById('ftpForm');
        form.addEventListener('submit', function(event) {
            event.preventDefault();
            
            const ftpServer = document.getElementById('ftpServer').value;
            const ftpLogin = document.getElementById('ftpLogin').value;
            const ftpPassword = document.getElementById('ftpPassword').value;

            // Сохраняем данные в localStorage для передачи на следующую страницу
            localStorage.setItem('ftpCredentials', JSON.stringify({
                server: ftpServer,
                login: ftpLogin,
                password: ftpPassword
            }));

            // Переход на страницу работы с FTP
            window.location.href = '/work';
        });
    </script>
</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// Функция для обработки страницы /work (консоль)
func handleWork(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FTP Client Work</title>
</head>
<body>
    <h2>FTP Client Console</h2>
    <form id="consoleForm">
        <label for="console">Enter Command:</label><br>
        <input type="text" id="console" name="console" required><br><br>
        <button type="submit">Send Command</button>
    </form>
    <h3>Result:</h3>
    <pre id="result"></pre>

    <script>
        // Получаем сохраненныеCredentails
        const credentials = JSON.parse(localStorage.getItem('ftpCredentials'));

        // WebSocket подключение к удаленному серверу
        const ws = new WebSocket("ws://localhost:8081/ws");

        ws.onopen = function() {
            console.log("WebSocket connection established.");
            
            // Первым делом отправляем данные для подключения к FTP
            const connectMessage = JSON.stringify({
                server: credentials.server,
                login: credentials.login,
                password: credentials.password,
                command: "connect"
            });
            ws.send(connectMessage);
        };

        const form = document.getElementById('consoleForm');
        form.addEventListener('submit', function(event) {
            event.preventDefault();
            
            const consoleInput = document.getElementById('console').value;
            
            // Формируем JSON команду
            const commandMessage = JSON.stringify({
                server: credentials.server,
                login: credentials.login,
                password: credentials.password,
                command: consoleInput.split(' ')[0],
                args: consoleInput.split(' ').slice(1)
            });

            ws.send(commandMessage);
        });

        ws.onmessage = function(event) {
            document.getElementById('result').textContent = event.data;
        };

        ws.onerror = function(error) {
            console.error('WebSocket Error:', error);
            document.getElementById('result').textContent = 'WebSocket connection error';
        };

        ws.onclose = function() {
            document.getElementById('result').textContent = 'WebSocket connection closed';
        };
    </script>
</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/work", handleWork)

	log.Println("HTTP сервер запущен на порту 8080.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
