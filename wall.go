package vk

import "context"

type WallPost struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	FromID  int    `json:"from_id"`
	Date    int64  `json:"date"`
	Text    string `json:"text"`
}

type WallGetResponse struct {
	Count int        `json:"count"`
	Items []WallPost `json:"items"`
}

type WallGetParams struct {
	OwnerID int `url:"owner_id,omitempty"`
	Offset  int `url:"offset,omitempty"`
	Count   int `url:"count,omitempty"`
}

func (c *Client) WallGet(ctx context.Context, params WallGetParams) (*WallGetResponse, error) {
	var out WallGetResponse
	if err := c.Call(ctx, "wall.get", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
