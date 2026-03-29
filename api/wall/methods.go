package wall

import (
	"context"

	"github.com/andr-235/vk_api/pkg/client"
)

func Get(ctx context.Context, c client.Caller, params WallGetParams) (*WallGetResponse, error) {
	var out WallGetResponse
	if err := c.Call(ctx, "wall.get", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
