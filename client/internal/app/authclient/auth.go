package authclient

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	configclient "github.com/Alandres998/go-keeper/client/internal/config"
	"github.com/Alandres998/go-keeper/proto/auth"
	"github.com/manifoldco/promptui"
	"google.golang.org/grpc"
)

func StartSession(conn *grpc.ClientConn) {
	client := auth.NewAuthServiceClient(conn)
	for {
		prompt := promptui.Select{
			Label: "Выберите опцию",
			Items: []string{"Войти", "Зарегстироваться"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			log.Fatalf("Не правильный текст: %v", err)
		}

		switch result {
		case "Войти":
			// Логика авторизации
			loginUser(client)
		case "Зарегстироваться":
			// Логика регистрации
			registerUser(client)
		}
		if configclient.Options.UserToken != "" {
			fmt.Println("Вы успешно авторизовались")
			break
		}
		fmt.Println("неуспешная попытка авторизации, повторяем попытку...")
		time.Sleep(2 * time.Second)
	}
}

func loginUser(client auth.AuthServiceClient) {
	var username, password string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Логин: ")
	if scanner.Scan() {
		username = scanner.Text()
	}

	fmt.Print("Пароль: ")
	if scanner.Scan() {
		password = scanner.Text()
	}

	ctx, cancel := context.WithTimeout(context.Background(), configclient.Options.TimeOut)
	defer cancel()

	// Вызываем метод Login
	res, err := client.Login(ctx, &auth.LoginRequest{
		Login:    username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Ошибка авторизации: %v", err)
	}

	if res.Success {
		configclient.Options.UserToken = res.Token
	} else {
		fmt.Printf("Ошибка авторизации: %s\n", res.Message)
	}
}

func registerUser(client auth.AuthServiceClient) {
	var username, password, email string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Логин: ")
	if scanner.Scan() {
		username = scanner.Text()
	}

	fmt.Print("Пароль: ")
	if scanner.Scan() {
		password = scanner.Text()
	}

	fmt.Print("Email: ")
	if scanner.Scan() {
		email = scanner.Text()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.Register(ctx, &auth.RegisterRequest{
		Login:    username,
		Password: password,
		Email:    email,
	})
	if err != nil {
		log.Fatalf("Ошибка регистрации: %v", err)
	}

	if res.Success {
		configclient.Options.UserToken = res.Token
	} else {
		fmt.Printf("Ошибка регистрации: %s\n", res.Message)
	}
}
