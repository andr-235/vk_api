package users

import (
	"context"

	vk "github.com/andr-235/vk_api"
)

func Get(ctx context.Context, c vk.Caller, params GetParams) ([]User, error) {
	var out []User
	if err := c.Call(ctx, "users.get", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func GetFollowers(ctx context.Context, c vk.Caller, params GetFollowersParams) (*GetFollowersResponse, error) {
	var out GetFollowersResponse
	if err := c.Call(ctx, "users.getFollowers", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetSubscriptions(ctx context.Context, c vk.Caller, params GetSubscriptionsParams) (*GetSubscriptionsResponse, error) {
	params.Extended = false

	var out GetSubscriptionsResponse
	if err := c.Call(ctx, "users.getSubscriptions", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetSubscriptionsExtended(ctx context.Context, c vk.Caller, params GetSubscriptionsParams) (*GetSubscriptionsExtendedResponse, error) {
	params.Extended = true

	var out GetSubscriptionsExtendedResponse
	if err := c.Call(ctx, "users.getSubscriptions", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func Search(ctx context.Context, c vk.Caller, params SearchParams) (*SearchResponse, error) {
	var out SearchResponse
	if err := c.Call(ctx, "users.search", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
