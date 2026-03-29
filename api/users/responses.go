package users

import "github.com/andr-235/vk_api/pkg/client"

type GetFollowersResponse = client.ListResponse[Profile]

type SubscriptionIDs struct {
	Count int   `json:"count"`
	Items []int `json:"items"`
}

type GetSubscriptionsResponse struct {
	Users  SubscriptionIDs `json:"users"`
	Groups SubscriptionIDs `json:"groups"`
}

type GetSubscriptionsExtendedResponse = client.ListResponse[SubscriptionItem]

type SubscriptionItem struct {
	Type string `json:"type,omitempty"`
	ID   int    `json:"id"`

	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	CanAccessClosed bool   `json:"can_access_closed,omitempty"`
	IsClosed        bool   `json:"is_closed,omitempty"`

	Name       string `json:"name,omitempty"`
	ScreenName string `json:"screen_name,omitempty"`

	Photo50  string `json:"photo_50,omitempty"`
	Photo100 string `json:"photo_100,omitempty"`

	IsAdmin      int `json:"is_admin,omitempty"`
	IsMember     int `json:"is_member,omitempty"`
	IsAdvertiser int `json:"is_advertiser,omitempty"`
	MembersCount int `json:"members_count,omitempty"`
}

type SearchResponse = client.ListResponse[User]
