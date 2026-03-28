package main

import (
	"context"
	"fmt"
	"log"
	"os"

	vk "github.com/andr-235/vk_api"
)

func main() {
	token := os.Getenv("VK_TOKEN")
	if token == "" {
		log.Fatal("VK_TOKEN env variable is required")
	}

	client := vk.New(
		vk.WithToken(token),
		vk.WithVersion("5.199"),
	)

	users, err := client.UsersGet(context.Background(), vk.UsersGetParams{
		UserIDs: []int{1},
		Fields:  []string{"bdate"},
	})
	if err != nil {
		log.Fatalf("UsersGet error: %v", err)
	}

	fmt.Println("Users:")
	for _, u := range users {
		fmt.Printf("ID=%d %s %s bdate=%s\n",
			u.ID,
			u.FirstName,
			u.LastName,
			u.BDate,
		)
	}
}
