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

func AddCallbackServer(ctx context.Context, c *vk.Client, params AddCallbackServerParams) (*AddCallbackServerResponse, error) {
	var out AddCallbackServerResponse
	if err := c.Call(ctx, "groups.addCallbackServer", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func DeleteAddress(ctx context.Context, c *vk.Client, params DeleteAddressParams) (bool, error) {
	var out int
	if err := c.Call(ctx, "groups.deleteAddress", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func DeleteCallbackServer(ctx context.Context, c *vk.Client, params DeleteCallbackServerParams) (bool, error) {
	var out int
	if err := c.Call(ctx, "groups.deleteCallbackServer", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func DisableOnline(ctx context.Context, c *vk.Client, params DisableOnlineParams) (bool, error) {
	var out int
	if err := c.Call(ctx, "groups.disableOnline", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func EditAddress(ctx context.Context, c *vk.Client, params EditAddressParams) (*Address, error) {
	var out Address
	if err := c.Call(ctx, "groups.editAddress", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func EditCallbackServer(ctx context.Context, c *vk.Client, params EditCallbackServerParams) (bool, error) {
	var out int
	if err := c.Call(ctx, "groups.editCallbackServer", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func EnableOnline(ctx context.Context, c *vk.Client, params EnableOnlineParams) (bool, error) {
	var out int
	if err := c.Call(ctx, "groups.enableOnline", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func Get(ctx context.Context, c *vk.Client, params GetParams) (*GetResponse, error) {
	var out GetResponse
	if err := c.Call(ctx, "groups.get", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
