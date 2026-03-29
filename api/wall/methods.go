package wall

import (
	"context"
	"fmt"

	"github.com/andr-235/vk_api/pkg/client"
)

func Get(ctx context.Context, c client.Caller, params WallGetParams) (*WallGetResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("wall.Get: invalid params: %w", err)
	}
	var out WallGetResponse
	if err := c.Call(ctx, "wall.get", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
