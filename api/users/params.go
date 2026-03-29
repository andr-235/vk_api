package users

import "errors"

type GetParams struct {
	UserIDs     []string `url:"user_ids,comma,omitempty"`
	Fields      []string `url:"fields,comma,omitempty"`
	NameCase    string   `url:"name_case,omitempty"`
	FromGroupID int      `url:"from_group_id,omitempty"`
}

// Validate проверяет валидность параметров метода Get.
func (p GetParams) Validate() error {
	if len(p.UserIDs) == 0 {
		return errors.New("user_ids обязателен")
	}
	return nil
}

type GetFollowersParams struct {
	UserID   int      `url:"user_id,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Count    int      `url:"count,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
	NameCase string   `url:"name_case,omitempty"`
}

// Validate проверяет валидность параметров метода GetFollowers.
func (p GetFollowersParams) Validate() error {
	if p.UserID <= 0 {
		return errors.New("user_id обязателен и должен быть положительным")
	}
	if p.Count < 0 {
		return errors.New("count не может быть отрицательным")
	}
	if p.Offset < 0 {
		return errors.New("offset не может быть отрицательным")
	}
	return nil
}

type GetSubscriptionsParams struct {
	UserID   int      `url:"user_id,omitempty"`
	Extended bool     `url:"extended,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Count    int      `url:"count,omitempty"`
	Fields   []string `url:"fields,comma,omitempty"`
}

// Validate проверяет валидность параметров метода GetSubscriptions.
func (p GetSubscriptionsParams) Validate() error {
	if p.UserID <= 0 {
		return errors.New("user_id обязателен и должен быть положительным")
	}
	if p.Count < 0 {
		return errors.New("count не может быть отрицательным")
	}
	if p.Offset < 0 {
		return errors.New("offset не может быть отрицательным")
	}
	return nil
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

// Validate проверяет валидность параметров метода Search.
func (p SearchParams) Validate() error {
	if p.Count < 0 {
		return errors.New("count не может быть отрицательным")
	}
	if p.Offset < 0 {
		return errors.New("offset не может быть отрицательным")
	}
	return nil
}
