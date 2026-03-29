package groups

type GetMembersResponse struct {
	Count int         `json:"count"`
	Items []MemberRef `json:"items"`
}

type Address struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	Address           string `json:"address"`
	AdditionalAddress string `json:"additional_address,omitempty"`
	CountryID         int    `json:"country_id"`
	CityID            int    `json:"city_id"`
	MetroID           int    `json:"metro_id,omitempty"`
	Latitude          string `json:"latitude"`
	Longitude         string `json:"longitude"`
	Phone             string `json:"phone,omitempty"`
	WorkInfoStatus    string `json:"work_info_status,omitempty"`
	Timetable         string `json:"timetable,omitempty"`
	IsMainAddress     bool   `json:"is_main_address"`
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

type AddCallbackServerResponse struct {
	ServerID int `json:"server_id"`
}

type GetResponse struct {
	Count int     `json:"count"`
	Items []Group `json:"items"`
}
