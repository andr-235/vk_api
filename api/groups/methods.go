package groups

import (
	"context"
	"fmt"

	"github.com/andr-235/vk_api/pkg/client"
)

func GetByID(ctx context.Context, c client.Caller, params GetByIDParams) ([]Group, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("groups.GetByID: invalid params: %w", err)
	}
	var out []Group
	if err := c.Call(ctx, "groups.getById", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func GetMembers(ctx context.Context, c client.Caller, params GetMembersParams) (*GetMembersResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("groups.GetMembers: invalid params: %w", err)
	}
	var out GetMembersResponse
	if err := c.Call(ctx, "groups.getMembers", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func AddAddress(ctx context.Context, c client.Caller, params AddAddressParams) (*Address, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("groups.AddAddress: invalid params: %w", err)
	}
	var out Address
	if err := c.Call(ctx, "groups.addAddress", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func AddCallbackServer(ctx context.Context, c client.Caller, params AddCallbackServerParams) (*AddCallbackServerResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("groups.AddCallbackServer: invalid params: %w", err)
	}
	var out AddCallbackServerResponse
	if err := c.Call(ctx, "groups.addCallbackServer", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func DeleteAddress(ctx context.Context, c client.Caller, params DeleteAddressParams) (bool, error) {
	if err := params.Validate(); err != nil {
		return false, fmt.Errorf("groups.DeleteAddress: invalid params: %w", err)
	}
	var out int
	if err := c.Call(ctx, "groups.deleteAddress", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func DeleteCallbackServer(ctx context.Context, c client.Caller, params DeleteCallbackServerParams) (bool, error) {
	if err := params.Validate(); err != nil {
		return false, fmt.Errorf("groups.DeleteCallbackServer: invalid params: %w", err)
	}
	var out int
	if err := c.Call(ctx, "groups.deleteCallbackServer", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func DisableOnline(ctx context.Context, c client.Caller, params DisableOnlineParams) (bool, error) {
	if err := params.Validate(); err != nil {
		return false, fmt.Errorf("groups.DisableOnline: invalid params: %w", err)
	}
	var out int
	if err := c.Call(ctx, "groups.disableOnline", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func EditAddress(ctx context.Context, c client.Caller, params EditAddressParams) (*Address, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("groups.EditAddress: invalid params: %w", err)
	}
	var out Address
	if err := c.Call(ctx, "groups.editAddress", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func EditCallbackServer(ctx context.Context, c client.Caller, params EditCallbackServerParams) (bool, error) {
	if err := params.Validate(); err != nil {
		return false, fmt.Errorf("groups.EditCallbackServer: invalid params: %w", err)
	}
	var out int
	if err := c.Call(ctx, "groups.editCallbackServer", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func EnableOnline(ctx context.Context, c client.Caller, params EnableOnlineParams) (bool, error) {
	if err := params.Validate(); err != nil {
		return false, fmt.Errorf("groups.EnableOnline: invalid params: %w", err)
	}
	var out int
	if err := c.Call(ctx, "groups.enableOnline", params, &out); err != nil {
		return false, err
	}
	return out == 1, nil
}

func Get(ctx context.Context, c client.Caller, params GetParams) (*GetResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("groups.Get: invalid params: %w", err)
	}
	var out GetResponse
	if err := c.Call(ctx, "groups.get", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetAddresses(ctx context.Context, c client.Caller, params GetAddressesParams) (*GetAddressesResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("groups.GetAddresses: invalid params: %w", err)
	}
	var out GetAddressesResponse
	if err := c.Call(ctx, "groups.getAddresses", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetBanned(ctx context.Context, c client.Caller, params GetBannedParams) (*GetBannedResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("groups.GetBanned: invalid params: %w", err)
	}
	var out GetBannedResponse
	if err := c.Call(ctx, "groups.getBanned", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
