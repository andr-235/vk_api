package wall

type WallGetParams struct {
	OwnerID int `url:"owner_id,omitempty"`
	Offset  int `url:"offset,omitempty"`
	Count   int `url:"count,omitempty"`
}
