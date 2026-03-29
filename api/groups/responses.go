package groups

import vk "github.com/andr-235/vk_api"

type GetMembersResponse = vk.ListResponse[MemberRef]

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

type AddCallbackServerResponse struct {
	ServerID int `json:"server_id"`
}

type GetResponse = vk.ListResponse[Group]

type GetAddressesResponse = vk.ListResponse[Address]

type GetBannedResponse struct {
	Count int          `json:"count"`
	Items []BannedItem `json:"items"`
}
