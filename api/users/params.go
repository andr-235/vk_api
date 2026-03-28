package users

type GetParams struct {
	UserIDs     []string `url:"user_ids,comma,omitempty"`
	Fields      []string `url:"fields,comma,omitempty"`
	NameCase    string   `url:"name_case,omitempty"`
	FromGroupID int      `url:"from_group_id,omitempty"`
}

type GetFollowersParams struct {
	UserID   int      `url:"user_id,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Count    int      `url:"count,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
	NameCase string   `url:"name_case,omitempty"`
}

type GetSubscriptionsParams struct {
	UserID   int      `url:"user_id,omitempty"`
	Extended bool     `url:"extended,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Count    int      `url:"count,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
}

type SearchParams struct {
	Q                 string   `url:"q,omitempty"`
	Sort              int      `url:"sort"`
	Offset            int      `url:"offset,omitempty"`
	Count             int      `url:"count,omitempty"`
	Fields            []string `url:"fields,comma,omitempty"`
	City              int      `url:"city,omitempty"`
	CityID            int      `url:"city_id,omitempty"`
	Country           int      `url:"country,omitempty"`
	CountryID         int      `url:"country_id,omitempty"`
	Hometown          string   `url:"hometown,omitempty"`
	UniversityCountry int      `url:"university_country,omitempty"`
	University        int      `url:"university,omitempty"`
	UniversityYear    int      `url:"university_year,omitempty"`
	UniversityFaculty int      `url:"university_faculty,omitempty"`
	UniversityChair   int      `url:"university_chair,omitempty"`
	Sex               int      `url:"sex,omitempty"`
	Status            int      `url:"status,omitempty"`
	AgeFrom           int      `url:"age_from,omitempty"`
	AgeTo             int      `url:"age_to,omitempty"`
	BirthDay          int      `url:"birth_day,omitempty"`
	BirthMonth        int      `url:"birth_month,omitempty"`
	BirthYear         int      `url:"birth_year,omitempty"`
	Online            bool     `url:"online,omitempty"`
	HasPhoto          bool     `url:"has_photo,omitempty"`
	SchoolCountry     int      `url:"school_country,omitempty"`
	SchoolCity        int      `url:"school_city,omitempty"`
	SchoolClass       int      `url:"school_class,omitempty"`
	School            int      `url:"school,omitempty"`
	SchoolYear        int      `url:"school_year,omitempty"`
	Religion          string   `url:"religion,omitempty"`
	Company           string   `url:"company,omitempty"`
	Position          string   `url:"position,omitempty"`
	GroupID           int      `url:"group_id,omitempty"`
	FromList          []string `url:"from_list,comma,omitempty"`
	ScreenRef         string   `url:"screen_ref,omitempty"`
	FromGroupID       int      `url:"from_group_id,omitempty"`
}
