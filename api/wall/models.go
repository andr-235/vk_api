package wall

type WallPost struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	FromID  int    `json:"from_id"`
	Date    int64  `json:"date"`
	Text    string `json:"text"`
}
