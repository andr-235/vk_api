package groups

import (
	"errors"
	"fmt"
)

const (
	GetCountMax = 1000
)

// Значения для параметра fields в методе GetAddresses
const (
	AddressFieldCity    = "city"
	AddressFieldCountry = "country"
)

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

type DeleteCallbackServerParams struct {
	GroupID  int `url:"group_id,omitempty"`
	ServerID int `url:"server_id,omitempty"`
}

type DisableOnlineParams struct {
	GroupID  int `url:"group_id,omitempty"`
}

type EditAddressParams struct {
	GroupID           int    `url:"group_id,omitempty"`
	AddressID         int    `url:"address_id,omitempty"`
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

type EditCallbackServerParams struct {
	GroupID    int    `url:"group_id,omitempty"`
	ServerID   int    `url:"server_id,omitempty"`
	URL        string `url:"url,omitempty"`
	Title      string `url:"title,omitempty"`
	SecretKey  string `url:"secret_key,omitempty"`
}

type EnableOnlineParams struct {
	GroupID  int `url:"group_id,omitempty"`
}

type GetParams struct {
	UserID   int      `url:"user_id,omitempty"`
	Extended bool     `url:"extended,omitempty"`
	Filter   []string `url:"filter,comma,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Count    int      `url:"count,omitempty"`
}

// Validate проверяает валидность параметров метода Get
func (p GetParams) Validate() error {
	if p.Count > GetCountMax {
		return errors.New("count не может превышать 1000")
	}
	if p.Count < 0 {
		return errors.New("count не может быть отрицательным")
	}
	if p.Offset < 0 {
		return errors.New("offset не может быть отрицательным")
	}
	return nil
}

// Значения для параметра filter в методе Get
const (
	FilterAdmin      = "admin"
	FilterEditor     = "editor"
	FilterModer      = "moder"
	FilterAdvertiser = "advertiser"
	FilterGroups     = "groups"
	FilterPublics    = "publics"
	FilterEvents     = "events"
	FilterHasAddress = "hasAddress"
)

// Значения для параметра fields в методе Get
const (
	FieldActivity         = "activity"
	FieldCanCreateTopic   = "can_create_topic"
	FieldCanPost          = "can_post"
	FieldCanSeeAllPosts   = "can_see_all_posts"
	FieldCity             = "city"
	FieldContacts         = "contacts"
	FieldCounters         = "counters"
	FieldCountry          = "country"
	FieldDescription      = "description"
	FieldFinishDate       = "finish_date"
	FieldFixedPost        = "fixed_post"
	FieldLinks            = "links"
	FieldMembersCount     = "members_count"
	FieldPlace            = "place"
	FieldSite             = "site"
	FieldStartDate        = "start_date"
	FieldStatus           = "status"
	FieldVerified         = "verified"
	FieldWikiPage         = "wiki_page"
)

type GetAddressesParams struct {
	GroupID     int    `url:"group_id,omitempty"`
	AddressIDs  []int  `url:"address_ids,comma,omitempty"`
	Latitude    string `url:"latitude,omitempty"`
	Longitude   string `url:"longitude,omitempty"`
	Offset      int    `url:"offset,omitempty"`
	Count       int    `url:"count,omitempty"`
	Fields      []string `url:"fields,comma,omitempty"`
}

type GetBannedParams struct {
	GroupID int      `url:"group_id,omitempty"`
	Offset  int      `url:"offset,omitempty"`
	Count   int      `url:"count,omitempty"`
	Fields  []string `url:"fields,comma,omitempty"`
	OwnerID int      `url:"owner_id,omitempty"`
}

// Validate проверяет валидность параметров метода GetAddresses
func (p GetAddressesParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	if p.Count < 0 {
		return errors.New("count не может быть отрицательным")
	}
	if p.Offset < 0 {
		return errors.New("offset не может быть отрицательным")
	}
	if p.Latitude != "" {
		lat := parseFloat(p.Latitude)
		if lat < -90 || lat > 90 {
			return errors.New("latitude должен быть от -90 до 90")
		}
	}
	if p.Longitude != "" {
		lon := parseFloat(p.Longitude)
		if lon < -180 || lon > 180 {
			return errors.New("longitude должен быть от -180 до 180")
		}
	}
	return nil
}

func parseFloat(s string) float64 {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	if err != nil {
		return 0
	}
	return result
}
