package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/andr-235/vk_api/pkg/client"
	"github.com/andr-235/vk_api/pkg/config"
	"github.com/andr-235/vk_api/api/groups"
	"github.com/andr-235/vk_api/api/users"
)

func main() {
	token := os.Getenv("VK_TOKEN")
	if token == "" {
		log.Fatal("VK_TOKEN env variable is required")
	}

	client := client.New(
		config.Config{Token: token},
		client.WithVersion("5.199"),
	)

	resp, err := users.Get(context.Background(), client, users.GetParams{
		UserIDs: []string{"1"},
		Fields:  []string{"bdate"},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range resp {
		fmt.Printf("ID=%d %s %s bdate=%s\n", u.ID, u.FirstName, u.LastName, u.BDate)
	}

	items, err := groups.GetByID(context.Background(), client, groups.GetByIDParams{
		GroupIDs: []string{"vk"},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, g := range items {
		fmt.Printf("ID=%d name=%s screen_name=%s\n", g.ID, g.Name, g.ScreenName)
	}
}
