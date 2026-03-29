package groups

import "github.com/andr-235/vk_api/api/users"

// Group представляет сообщество VK.
// https://dev.vk.com/ru/method/groups.getById
type Group struct {
	// Базовые поля
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name,omitempty"`
	Type       string `json:"type,omitempty"`

	// Статусы
	IsClosed     int `json:"is_closed,omitempty"`
	IsAdmin      int `json:"is_admin,omitempty"`
	IsMember     int `json:"is_member,omitempty"`
	IsAdvertiser int `json:"is_advertiser,omitempty"`

	// Счётчики
	MembersCount int `json:"members_count,omitempty"`

	// Фотографии
	Photo50  string `json:"photo_50,omitempty"`
	Photo100 string `json:"photo_100,omitempty"`
	Photo200 string `json:"photo_200,omitempty"`

	// Описание и сайт
	Description string `json:"description,omitempty"`
	Site        string `json:"site,omitempty"`

	// Активность
	Activity string `json:"activity,omitempty"`

	// Информация о бане (требуется право groups)
	BanInfo *BanInfo `json:"ban_info,omitempty"`

	// Права доступа
	CanPost        int `json:"can_post,omitempty"`
	CanSeeAllPosts int `json:"can_see_all_posts,omitempty"`

	// Город и страна
	City    *City    `json:"city,omitempty"`
	Country *Country `json:"country,omitempty"`

	// Контакты
	Contacts *Contacts `json:"contacts,omitempty"`

	// Счётчики (расширенные)
	Counters *Counters `json:"counters,omitempty"`

	// Обложка
	Cover *Cover `json:"cover,omitempty"`

	// Даты
	FinishDate int `json:"finish_date,omitempty"`
	StartDate  int `json:"start_date,omitempty"`

	// Закреплённый пост
	FixedPost int `json:"fixed_post,omitempty"`

	// Ссылки
	Links *Links `json:"links,omitempty"`

	// Магазин
	Market *Market `json:"market,omitempty"`

	// Место
	Place *Place `json:"place,omitempty"`

	// Статус сообщества
	Status string `json:"status,omitempty"`

	// Проверенное сообщество
	Verified int `json:"verified,omitempty"`

	// Wiki-страница
	WikiPage string `json:"wiki_page,omitempty"`
}

// Cover представляет обложку сообщества
type Cover struct {
	Enabled bool   `json:"enabled,omitempty"`
	URL     string `json:"url,omitempty"`
	URL2x   string `json:"url_2x,omitempty"`
	URL3x   string `json:"url_3x,omitempty"`
}

// Contacts представляет контакты сообщества
type Contacts struct {
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Website string `json:"website,omitempty"`
}

// Counters представляет счётчики сообщества
type Counters struct {
	AlbumsPhotos int `json:"albums_photos,omitempty"`
	AlbumsVideos int `json:"albums_videos,omitempty"`
	AlbumsAudios int `json:"albums_audios,omitempty"`
	Photos       int `json:"photos,omitempty"`
	Videos       int `json:"videos,omitempty"`
	Audios       int `json:"audios,omitempty"`
	Docs         int `json:"docs,omitempty"`
	Topics       int `json:"topics,omitempty"`
}

// Links представляет ссылки сообщества
type Links struct {
	Items []LinkItem `json:"items,omitempty"`
}

// LinkItem представляет отдельную ссылку
type LinkItem struct {
	ID          int    `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Photo50     string `json:"photo_50,omitempty"`
	Photo100    string `json:"photo_100,omitempty"`
}

// Market представляет настройки магазина
type Market struct {
	Enabled    int    `json:"enabled,omitempty"`
	ContactID  int    `json:"contact_id,omitempty"`
	Currency    int    `json:"currency,omitempty"`
	CurrencyText string `json:"currency_text,omitempty"`
	Comments    int    `json:"comments,omitempty"`
}

// Place представляет место сообщества
type Place struct {
	ID        int     `json:"id,omitempty"`
	Title     string  `json:"title,omitempty"`
	CountryID int     `json:"country_id,omitempty"`
	CityID    int     `json:"city_id,omitempty"`
	Address   string  `json:"address,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	CreatedAt int     `json:"created_at,omitempty"`
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
