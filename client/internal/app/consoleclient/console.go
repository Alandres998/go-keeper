package consoleclient

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ClearConsole очищает консоль в зависимости от операционной системы.
func ClearConsole() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "linux", "darwin":
		cmd = exec.Command("clear")
	default:
		fmt.Println("Неизвестная операционная система")
		return
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Ошибка при очистке консоли:", err)
	}
}
