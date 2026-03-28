package groups

type Group struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name,omitempty"`
	Type       string `json:"type,omitempty"`

	IsClosed int `json:"is_closed,omitempty"`

	IsAdmin      int `json:"is_admin,omitempty"`
	IsMember     int `json:"is_member,omitempty"`
	IsAdvertiser int `json:"is_advertiser,omitempty"`

	MembersCount int `json:"members_count,omitempty"`

	Photo50  string `json:"photo_50,omitempty"`
	Photo100 string `json:"photo_100,omitempty"`
	Photo200 string `json:"photo_200,omitempty"`

	Description string `json:"description,omitempty"`
	Site        string `json:"site,omitempty"`
}
