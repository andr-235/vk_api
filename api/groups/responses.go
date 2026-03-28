package groups

type GetMembersResponse struct {
	Count int         `json:"count"`
	Items []MemberRef `json:"items"`
}

type MemberRef struct {
	ID              int    `json:"id"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	CanAccessClosed bool   `json:"can_access_closed,omitempty"`
	IsClosed        bool   `json:"is_closed,omitempty"`
	Deactivated     string `json:"deactivated,omitempty"`

	Photo50  string `json:"photo_50,omitempty"`
	Photo100 string `json:"photo_100,omitempty"`
	Photo200 string `json:"photo_200,omitempty"`

	Online int `json:"online,omitempty"`
	Sex    int `json:"sex,omitempty"`
}
