package client_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andr-235/vk_api/pkg/client"
	"github.com/andr-235/vk_api/pkg/config"
	"github.com/andr-235/vk_api/pkg/ratelimit"
	"github.com/andr-235/vk_api/pkg/retry"
)

func ExampleNew() {
	cfg := config.DefaultConfig()
	cfg.Token = "your_token"

	c := client.New(cfg,
		client.WithRateLimiter(ratelimit.NewTokenBucketRateLimiter(3.0)),
		client.WithRetryer(retry.NewSimpleRetryer(3, retry.DefaultPolicy())),
	)

	_ = c
}

func ExampleNewBuilder() {
	c, err := client.NewBuilder().
		WithToken("your_token").
		WithVersion("5.199").
		WithRateLimiter(ratelimit.NewTokenBucketRateLimiter(3.0)).
		Build()
	if err != nil {
		panic(err)
	}

	_ = c
}

func ExampleClient_Call() {
	c := client.New(config.Config{Token: "token"})

	var result interface{}
	err := c.Call(context.Background(), "users.get", map[string]any{
		"user_ids": "1",
		"fields":   "bdate",
	}, &result)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Result: %+v\n", result)
}

func ExampleClient_CallWithRawHandler() {
	c := client.New(config.Config{Token: "token"})

	err := c.CallWithRawHandler(context.Background(), "users.get", map[string]any{
		"user_ids": "1",
	}, func(raw json.RawMessage) error {
		// Кастомная обработка JSON
		fmt.Printf("Raw response: %s\n", raw)
		return nil
	})

	if err != nil {
		panic(err)
	}
}
