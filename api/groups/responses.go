package groups

type GetMembersResponse struct {
	Count int         `json:"count"`
	Items []MemberRef `json:"items"`
}

type Address struct {
	ID                int         `json:"id"`
	Title             string      `json:"title"`
	Address           string      `json:"address"`
	AdditionalAddress string      `json:"additional_address,omitempty"`
	CountryID         int         `json:"country_id"`
	CityID            int         `json:"city_id"`
	MetroStationID    int         `json:"metro_station_id,omitempty"`
	Latitude          float64     `json:"latitude"`
	Longitude         float64     `json:"longitude"`
	Distance          int         `json:"distance,omitempty"`
	Phone             string      `json:"phone,omitempty"`
	TimeOffset        int         `json:"time_offset,omitempty"`
	WorkInfoStatus    string      `json:"work_info_status,omitempty"`
	Timetable         *Timetable  `json:"timetable,omitempty"`
	IsMainAddress     bool        `json:"is_main_address"`
	City              *City       `json:"city,omitempty"`
	Country           *Country    `json:"country,omitempty"`
}

type Timetable struct {
	Mon *DaySchedule `json:"mon,omitempty"`
	Tue *DaySchedule `json:"tue,omitempty"`
	Wed *DaySchedule `json:"wed,omitempty"`
	Thu *DaySchedule `json:"thu,omitempty"`
	Fri *DaySchedule `json:"fri,omitempty"`
	Sat *DaySchedule `json:"sat,omitempty"`
	Sun *DaySchedule `json:"sun,omitempty"`
}

type DaySchedule struct {
	OpenTime       int `json:"open_time"`
	CloseTime      int `json:"close_time"`
	BreakOpenTime  int `json:"break_open_time"`
	BreakCloseTime int `json:"break_close_time"`
}

type City struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Country struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
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

type GetAddressesResponse struct {
	Count int       `json:"count"`
	Items []Address `json:"items"`
}

type GetBannedResponse struct {
	Count int          `json:"count"`
	Items []BannedItem `json:"items"`
}
