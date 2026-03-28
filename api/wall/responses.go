package wall

type WallGetResponse struct {
	Count int        `json:"count"`
	Items []WallPost `json:"items"`
}
