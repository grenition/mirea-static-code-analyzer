package main

import (
	"fmt"
	"os"      // Неиспользуемый импорт
	"strings" // Неиспользуемый импорт
)

// Пример проблемного Go кода для тестирования golangci-lint
// Ожидаемые проблемы: необработанные ошибки, неиспользуемые переменные/импорты

func readFile(filename string) {
	file, _ := os.Open(filename) // Игнорирование ошибки (errcheck)
	defer file.Close()

	data := make([]byte, 100)
	file.Read(data) // Игнорирование ошибки (errcheck)

	fmt.Println(string(data))
}

func divide(a, b int) int {
	return a / b // Возможное деление на ноль
}

func unusedFunction() {
	var unused int = 10 // Неиспользуемая переменная (unused)
	fmt.Println("test")
}

func main() {
	// Неиспользуемая переменная
	unused := 5
	fmt.Println("Hello")
}

