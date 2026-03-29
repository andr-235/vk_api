package groups

type GetByIDParams struct {
	GroupIDs []string `url:"group_ids,comma,omitempty"`
	GroupID  string   `url:"group_id,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
}

type GetMembersParams struct {
	GroupID string   `url:"group_id,omitempty"`
	Sort    string   `url:"sort,omitempty"`
	Offset  int      `url:"offset,omitempty"`
	Count   int      `url:"count,omitempty"`
	Fields  []string `url:"fields,comma,omitempty"`
	Filter  string   `url:"filter,omitempty"`
}

type AddAddressParams struct {
	GroupID           int    `url:"group_id,omitempty"`
	Title             string `url:"title,omitempty"`
	Address           string `url:"address,omitempty"`
	AdditionalAddress string `url:"additional_address,omitempty"`
	CountryID         int    `url:"country_id,omitempty"`
	CityID            int    `url:"city_id,omitempty"`
	MetroID           int    `url:"metro_id,omitempty"`
	Latitude          string `url:"latitude,omitempty"`
	Longitude         string `url:"longitude,omitempty"`
	Phone             string `url:"phone,omitempty"`
	WorkInfoStatus    string `url:"work_info_status,omitempty"`
	Timetable         string `url:"timetable,omitempty"`
	IsMainAddress     bool   `url:"is_main_address,omitempty"`
}

type AddCallbackServerParams struct {
	GroupID     int    `url:"group_id,omitempty"`
	URL         string `url:"url,omitempty"`
	Title       string `url:"title,omitempty"`
	SecretKey   string `url:"secret_key,omitempty"`
}

type DeleteAddressParams struct {
	GroupID    int `url:"group_id,omitempty"`
	AddressID  int `url:"address_id,omitempty"`
}
