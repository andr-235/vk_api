package vk

import "context"

type Group struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Type       string `json:"type"`
}

type GroupsGetByIDParams struct {
	GroupIDs []string `url:"group_ids,comma,omitempty"`
	GroupID  string   `url:"group_id,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
}

func (c *Client) GroupsGetByID(ctx context.Context, params GroupsGetByIDParams) ([]Group, error) {
	var out []Group
	if err := c.Call(ctx, "groups.getById", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
