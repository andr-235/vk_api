package groups

import "github.com/andr-235/vk_api/api/users"

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

// BannedProfile — тип-алиас для users.Profile, используется в методе GetBanned
type BannedProfile = users.Profile

// MemberRef — тип-алиас для users.Profile, используется в методе GetMembers
type MemberRef = users.Profile

// BannedGroup представляет информацию о забаненном сообществе
type BannedGroup struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ScreenName   string `json:"screen_name,omitempty"`
	Type         string `json:"type,omitempty"`
	IsClosed     int    `json:"is_closed,omitempty"`
	MembersCount int    `json:"members_count,omitempty"`
	Photo50      string `json:"photo_50,omitempty"`
	Photo100     string `json:"photo_100,omitempty"`
	Photo200     string `json:"photo_200,omitempty"`
	Description  string `json:"description,omitempty"`
}

// BanInfo представляет информацию о блокировке
type BanInfo struct {
	AdminID int    `json:"admin_id"`
	Date    int    `json:"date"`
	Reason  int    `json:"reason"`
	Comment string `json:"comment,omitempty"`
	EndDate int    `json:"end_date"`
}

// BannedItem представляет элемент чёрного списка сообщества
type BannedItem struct {
	Type    string         `json:"type"`
	Profile *BannedProfile `json:"profile,omitempty"`
	Group   *BannedGroup   `json:"group,omitempty"`
	BanInfo *BanInfo       `json:"ban_info,omitempty"`
}
