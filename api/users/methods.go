package users

import (
	"context"
	"fmt"

	"github.com/andr-235/vk_api/pkg/client"
)

func Get(ctx context.Context, c client.Caller, params GetParams) ([]User, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("users.Get: invalid params: %w", err)
	}
	var out []User
	if err := c.Call(ctx, "users.get", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func GetFollowers(ctx context.Context, c client.Caller, params GetFollowersParams) (*GetFollowersResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("users.GetFollowers: invalid params: %w", err)
	}
	var out GetFollowersResponse
	if err := c.Call(ctx, "users.getFollowers", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetSubscriptions(ctx context.Context, c client.Caller, params GetSubscriptionsParams) (*GetSubscriptionsResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("users.GetSubscriptions: invalid params: %w", err)
	}
	params.Extended = false

	var out GetSubscriptionsResponse
	if err := c.Call(ctx, "users.getSubscriptions", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetSubscriptionsExtended(ctx context.Context, c client.Caller, params GetSubscriptionsParams) (*GetSubscriptionsExtendedResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("users.GetSubscriptionsExtended: invalid params: %w", err)
	}
	params.Extended = true

	var out GetSubscriptionsExtendedResponse
	if err := c.Call(ctx, "users.getSubscriptions", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func Search(ctx context.Context, c client.Caller, params SearchParams) (*SearchResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("users.Search: invalid params: %w", err)
	}
	var out SearchResponse
	if err := c.Call(ctx, "users.search", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
