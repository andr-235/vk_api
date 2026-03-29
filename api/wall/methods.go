package wall

import (
	"context"

	vk "github.com/andr-235/vk_api"
)

func Get(ctx context.Context, c *vk.Client, params WallGetParams) (*WallGetResponse, error) {
	var out WallGetResponse
	if err := c.Call(ctx, "wall.get", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
