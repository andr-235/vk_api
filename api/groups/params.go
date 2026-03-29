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
	GetAddressesFieldCity    = "city"
	GetAddressesFieldCountry = "country"
)

type GetByIDParams struct {
	GroupIDs []string `url:"group_ids,comma,omitempty"`
	GroupID  string   `url:"group_id,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
}

// Validate проверяет валидность параметров метода GetByID
func (p GetByIDParams) Validate() error {
	if p.GroupID == "" && len(p.GroupIDs) == 0 {
		return errors.New("group_id или group_ids обязателен")
	}
	return nil
}

type GetMembersParams struct {
	GroupID string   `url:"group_id,omitempty"`
	Sort    string   `url:"sort,omitempty"`
	Offset  int      `url:"offset,omitempty"`
	Count   int      `url:"count,omitempty"`
	Fields  []string `url:"fields,comma,omitempty"`
	Filter  string   `url:"filter,omitempty"`
}

// Validate проверяет валидность параметров метода GetMembers
func (p GetMembersParams) Validate() error {
	if p.GroupID == "" {
		return errors.New("group_id обязателен")
	}
	if p.Count < 0 {
		return errors.New("count не может быть отрицательным")
	}
	if p.Offset < 0 {
		return errors.New("offset не может быть отрицательным")
	}
	return nil
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

// Validate проверяет валидность параметров метода AddAddress
func (p AddAddressParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	if p.Title == "" {
		return errors.New("title обязателен")
	}
	if p.Latitude != "" {
		if _, err := parseFloat(p.Latitude); err != nil {
			return fmt.Errorf("invalid latitude: %w", err)
		}
	}
	if p.Longitude != "" {
		if _, err := parseFloat(p.Longitude); err != nil {
			return fmt.Errorf("invalid longitude: %w", err)
		}
	}
	return nil
}

type AddCallbackServerParams struct {
	GroupID     int    `url:"group_id,omitempty"`
	URL         string `url:"url,omitempty"`
	Title       string `url:"title,omitempty"`
	SecretKey   string `url:"secret_key,omitempty"`
}

// Validate проверяет валидность параметров метода AddCallbackServer
func (p AddCallbackServerParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	if p.URL == "" {
		return errors.New("url обязателен")
	}
	if p.Title == "" {
		return errors.New("title обязателен")
	}
	return nil
}

type DeleteAddressParams struct {
	GroupID    int `url:"group_id,omitempty"`
	AddressID  int `url:"address_id,omitempty"`
}

// Validate проверяет валидность параметров метода DeleteAddress
func (p DeleteAddressParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	if p.AddressID <= 0 {
		return errors.New("address_id обязателен и должен быть положительным")
	}
	return nil
}

type DeleteCallbackServerParams struct {
	GroupID  int `url:"group_id,omitempty"`
	ServerID int `url:"server_id,omitempty"`
}

// Validate проверяет валидность параметров метода DeleteCallbackServer
func (p DeleteCallbackServerParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	if p.ServerID <= 0 {
		return errors.New("server_id обязателен и должен быть положительным")
	}
	return nil
}

type DisableOnlineParams struct {
	GroupID  int `url:"group_id,omitempty"`
}

// Validate проверяет валидность параметров метода DisableOnline
func (p DisableOnlineParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	return nil
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

// Validate проверяет валидность параметров метода EditAddress
func (p EditAddressParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	if p.AddressID <= 0 {
		return errors.New("address_id обязателен и должен быть положительным")
	}
	if p.Title == "" {
		return errors.New("title обязателен")
	}
	if p.Latitude != "" {
		if _, err := parseFloat(p.Latitude); err != nil {
			return fmt.Errorf("invalid latitude: %w", err)
		}
	}
	if p.Longitude != "" {
		if _, err := parseFloat(p.Longitude); err != nil {
			return fmt.Errorf("invalid longitude: %w", err)
		}
	}
	return nil
}

type EditCallbackServerParams struct {
	GroupID    int    `url:"group_id,omitempty"`
	ServerID   int    `url:"server_id,omitempty"`
	URL        string `url:"url,omitempty"`
	Title      string `url:"title,omitempty"`
	SecretKey  string `url:"secret_key,omitempty"`
}

// Validate проверяет валидность параметров метода EditCallbackServer
func (p EditCallbackServerParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	if p.ServerID <= 0 {
		return errors.New("server_id обязателен и должен быть положительным")
	}
	if p.URL == "" {
		return errors.New("url обязателен")
	}
	return nil
}

type EnableOnlineParams struct {
	GroupID  int `url:"group_id,omitempty"`
}

// Validate проверяет валидность параметров метода EnableOnline
func (p EnableOnlineParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	return nil
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

// Validate проверяет валидность параметров метода GetBanned
func (p GetBannedParams) Validate() error {
	if p.GroupID <= 0 {
		return errors.New("group_id обязателен и должен быть положительным")
	}
	if p.Count < 0 {
		return errors.New("count не может быть отрицательным")
	}
	if p.Offset < 0 {
		return errors.New("offset не может быть отрицательным")
	}
	return nil
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
		lat, err := parseFloat(p.Latitude)
		if err != nil {
			return fmt.Errorf("invalid latitude: %w", err)
		}
		if lat < -90 || lat > 90 {
			return errors.New("latitude должен быть от -90 до 90")
		}
	}
	if p.Longitude != "" {
		lon, err := parseFloat(p.Longitude)
		if err != nil {
			return fmt.Errorf("invalid longitude: %w", err)
		}
		if lon < -180 || lon > 180 {
			return errors.New("longitude должен быть от -180 до 180")
		}
	}
	return nil
}

func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}
	return result, nil
}
