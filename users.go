package vk

import "context"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BDate     string `json:"bdate,omitempty"`
}

type UsersGetParams struct {
	UserIDs  []int    `url:"user_ids,comma,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
	NameCase string   `url:"name_case,omitempty"`
}

func (c *Client) UsersGet(ctx context.Context, params UsersGetParams) ([]User, error) {
	var out []User
	if err := c.Call(ctx, "users.get", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
