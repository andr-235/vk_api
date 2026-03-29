package groups

import (
	"context"

	vk "github.com/andr-235/vk_api"
)

func GetByID(ctx context.Context, c *vk.Client, params GetByIDParams) ([]Group, error) {
	var out []Group
	if err := c.Call(ctx, "groups.getById", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func GetMembers(ctx context.Context, c *vk.Client, params GetMembersParams) (*GetMembersResponse, error) {
	var out GetMembersResponse
	if err := c.Call(ctx, "groups.getMembers", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func AddAddress(ctx context.Context, c *vk.Client, params AddAddressParams) (*Address, error) {
	var out Address
	if err := c.Call(ctx, "groups.addAddress", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
