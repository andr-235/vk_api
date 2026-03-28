package vk

import "context"

type User struct {
	// Базовые поля
	ID              int    `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Deactivated     string `json:"deactivated,omitempty"`
	IsClosed        bool   `json:"is_closed,omitempty"`
	CanAccessClosed bool   `json:"can_access_closed,omitempty"`

	// A-I
	About                  string           `json:"about,omitempty"`
	Activities             string           `json:"activities,omitempty"`
	BDate                  string           `json:"bdate,omitempty"`
	Blacklisted            int              `json:"blacklisted,omitempty"`
	BlacklistedByMe        int              `json:"blacklisted_by_me,omitempty"`
	Books                  string           `json:"books,omitempty"`
	CanPost                int              `json:"can_post,omitempty"`
	CanSeeAllPosts         int              `json:"can_see_all_posts,omitempty"`
	CanSeeAudio            int              `json:"can_see_audio,omitempty"`
	CanSendFriendRequest   int              `json:"can_send_friend_request,omitempty"`
	CanWritePrivateMessage int              `json:"can_write_private_message,omitempty"`
	Career                 []UserCareer     `json:"career,omitempty"`
	City                   *UserCity        `json:"city,omitempty"`
	CommonCount            int              `json:"common_count,omitempty"`
	Connections            *UserConnections `json:"connections,omitempty"`
	Contacts               *UserContacts    `json:"contacts,omitempty"`
	Counters               *UserCounters    `json:"counters,omitempty"`
	Country                *UserCountry     `json:"country,omitempty"`
	CropPhoto              *UserCropPhoto   `json:"crop_photo,omitempty"`
	Domain                 string           `json:"domain,omitempty"`
	Education              *UserEducation   `json:"education,omitempty"`
	Exports                *UserExports     `json:"exports,omitempty"`
	FirstNameNom           string           `json:"first_name_nom,omitempty"`
	FirstNameGen           string           `json:"first_name_gen,omitempty"`
	FirstNameDat           string           `json:"first_name_dat,omitempty"`
	FirstNameAcc           string           `json:"first_name_acc,omitempty"`
	FirstNameIns           string           `json:"first_name_ins,omitempty"`
	FirstNameAbl           string           `json:"first_name_abl,omitempty"`
	FollowersCount         int              `json:"followers_count,omitempty"`
	FriendStatus           int              `json:"friend_status,omitempty"`
	Games                  string           `json:"games,omitempty"`
	HasMobile              int              `json:"has_mobile,omitempty"`
	HasPhoto               int              `json:"has_photo,omitempty"`
	HomeTown               string           `json:"home_town,omitempty"`
	Interests              string           `json:"interests,omitempty"`
	IsFavorite             int              `json:"is_favorite,omitempty"`
	IsFriend               int              `json:"is_friend,omitempty"`
	IsHiddenFromFeed       int              `json:"is_hidden_from_feed,omitempty"`
	IsNoIndex              int              `json:"is_no_index,omitempty"`
	IsVerified             bool             `json:"is_verified,omitempty"`

	// L-R
	LastNameNom     string               `json:"last_name_nom,omitempty"`
	LastNameGen     string               `json:"last_name_gen,omitempty"`
	LastNameDat     string               `json:"last_name_dat,omitempty"`
	LastNameAcc     string               `json:"last_name_acc,omitempty"`
	LastNameIns     string               `json:"last_name_ins,omitempty"`
	LastNameAbl     string               `json:"last_name_abl,omitempty"`
	LastSeen        *UserLastSeen        `json:"last_seen,omitempty"`
	Lists           string               `json:"lists,omitempty"`
	MaidenName      string               `json:"maiden_name,omitempty"`
	Military        []UserMilitary       `json:"military,omitempty"`
	Movies          string               `json:"movies,omitempty"`
	Music           string               `json:"music,omitempty"`
	Nickname        string               `json:"nickname,omitempty"`
	Occupation      *UserOccupation      `json:"occupation,omitempty"`
	Online          int                  `json:"online,omitempty"`
	OnlineMobile    int                  `json:"online_mobile,omitempty"`
	OnlineApp       int                  `json:"online_app,omitempty"`
	Personal        *UserPersonal        `json:"personal,omitempty"`
	Photo50         string               `json:"photo_50,omitempty"`
	Photo100        string               `json:"photo_100,omitempty"`
	Photo200Orig    string               `json:"photo_200_orig,omitempty"`
	Photo200        string               `json:"photo_200,omitempty"`
	Photo400Orig    string               `json:"photo_400_orig,omitempty"`
	PhotoID         string               `json:"photo_id,omitempty"`
	PhotoMax        string               `json:"photo_max,omitempty"`
	PhotoMaxOrig    string               `json:"photo_max_orig,omitempty"`
	Quotes          string               `json:"quotes,omitempty"`
	Relatives       []UserRelative       `json:"relatives,omitempty"`
	Relation        int                  `json:"relation,omitempty"`
	RelationPartner *UserRelationPartner `json:"relation_partner,omitempty"`

	// S-W
	Schools      []UserSchool     `json:"schools,omitempty"`
	ScreenName   string           `json:"screen_name,omitempty"`
	Sex          int              `json:"sex,omitempty"`
	Site         string           `json:"site,omitempty"`
	Status       string           `json:"status,omitempty"`
	StatusAudio  map[string]any   `json:"status_audio,omitempty"`
	Timezone     int              `json:"timezone,omitempty"`
	Trending     int              `json:"trending,omitempty"`
	TV           string           `json:"tv,omitempty"`
	Universities []UserUniversity `json:"universities,omitempty"`
	Verified     int              `json:"verified,omitempty"`
	WallDefault  string           `json:"wall_default,omitempty"`
}

type UserCareer struct {
	GroupID   int    `json:"group_id,omitempty"`
	Company   string `json:"company,omitempty"`
	CountryID int    `json:"country_id,omitempty"`
	CityID    int    `json:"city_id,omitempty"`
	CityName  string `json:"city_name,omitempty"`
	From      int    `json:"from,omitempty"`
	Until     int    `json:"until,omitempty"`
	Position  string `json:"position,omitempty"`
}

type UserCity struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type UserCountry struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type UserConnections struct {
	Skype       string `json:"skype,omitempty"`
	Facebook    string `json:"facebook,omitempty"`
	Twitter     string `json:"twitter,omitempty"`
	Livejournal string `json:"livejournal,omitempty"`
	Instagram   string `json:"instagram,omitempty"`
}

type UserContacts struct {
	MobilePhone string `json:"mobile_phone,omitempty"`
	HomePhone   string `json:"home_phone,omitempty"`
}

type UserCounters struct {
	Albums        int `json:"albums,omitempty"`
	Videos        int `json:"videos,omitempty"`
	Audios        int `json:"audios,omitempty"`
	Photos        int `json:"photos,omitempty"`
	Notes         int `json:"notes,omitempty"`
	Friends       int `json:"friends,omitempty"`
	Gifts         int `json:"gifts,omitempty"`
	Groups        int `json:"groups,omitempty"`
	OnlineFriends int `json:"online_friends,omitempty"`
	MutualFriends int `json:"mutual_friends,omitempty"`
	UserVideos    int `json:"user_videos,omitempty"`
	UserPhotos    int `json:"user_photos,omitempty"`
	Followers     int `json:"followers,omitempty"`
	Pages         int `json:"pages,omitempty"`
	Subscriptions int `json:"subscriptions,omitempty"`
}

type UserCropCoords struct {
	X  float64 `json:"x,omitempty"`
	Y  float64 `json:"y,omitempty"`
	X2 float64 `json:"x2,omitempty"`
	Y2 float64 `json:"y2,omitempty"`
}

type UserCropPhoto struct {
	Photo map[string]any `json:"photo,omitempty"`
	Crop  UserCropCoords `json:"crop"`
	Rect  UserCropCoords `json:"rect"`
}

type UserEducation struct {
	University     int    `json:"university,omitempty"`
	UniversityName string `json:"university_name,omitempty"`
	Faculty        int    `json:"faculty,omitempty"`
	FacultyName    string `json:"faculty_name,omitempty"`
	Graduation     int    `json:"graduation,omitempty"`
}

type UserExports struct {
	Livejournal int `json:"livejournal,omitempty"`
}

type UserLastSeen struct {
	Time     int64 `json:"time,omitempty"`
	Platform int   `json:"platform,omitempty"`
}

type UserMilitary struct {
	Unit      string `json:"unit,omitempty"`
	UnitID    int    `json:"unit_id,omitempty"`
	CountryID int    `json:"country_id,omitempty"`
	From      int    `json:"from,omitempty"`
	Until     int    `json:"until,omitempty"`
}

type UserOccupation struct {
	Type string `json:"type,omitempty"`
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserPersonal struct {
	Political  int      `json:"political,omitempty"`
	Langs      []string `json:"langs,omitempty"`
	Religion   string   `json:"religion,omitempty"`
	InspiredBy string   `json:"inspired_by,omitempty"`
	PeopleMain int      `json:"people_main,omitempty"`
	LifeMain   int      `json:"life_main,omitempty"`
	Smoking    int      `json:"smoking,omitempty"`
	Alcohol    int      `json:"alcohol,omitempty"`
}

type UserRelative struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type UserRelationPartner struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type UserSchool struct {
	ID            string `json:"id,omitempty"`
	Country       int    `json:"country,omitempty"`
	City          int    `json:"city,omitempty"`
	Name          string `json:"name,omitempty"`
	YearFrom      int    `json:"year_from,omitempty"`
	YearTo        int    `json:"year_to,omitempty"`
	YearGraduated int    `json:"year_graduated,omitempty"`
	Class         string `json:"class,omitempty"`
	Speciality    string `json:"speciality,omitempty"`
	Type          int    `json:"type,omitempty"`
	TypeStr       string `json:"type_str,omitempty"`
}

type UserUniversity struct {
	ID              int    `json:"id,omitempty"`
	Country         int    `json:"country,omitempty"`
	City            int    `json:"city,omitempty"`
	Name            string `json:"name,omitempty"`
	Faculty         int    `json:"faculty,omitempty"`
	FacultyName     string `json:"faculty_name,omitempty"`
	Chair           int    `json:"chair,omitempty"`
	ChairName       string `json:"chair_name,omitempty"`
	Graduation      int    `json:"graduation,omitempty"`
	EducationForm   string `json:"education_form,omitempty"`
	EducationStatus string `json:"education_status,omitempty"`
}

type UsersGetParams struct {
	UserIDs     []string `url:"user_ids,comma,omitempty"`
	Fields      []string `url:"fields,comma,omitempty"`
	NameCase    string   `url:"name_case,omitempty"`
	FromGroupID int      `url:"from_group_id,omitempty"`
}

func (c *Client) UsersGet(ctx context.Context, params UsersGetParams) ([]User, error) {
	var out []User
	if err := c.Call(ctx, "users.get", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
