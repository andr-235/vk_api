package groups

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/andr-235/vk_api/pkg/client"
	"github.com/andr-235/vk_api/pkg/config"
)

func TestGetByID(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": [
				{
					"id": 1,
					"name": "VK Test Group",
					"screen_name": "vk_test",
					"type": "group",
					"is_closed": 0,
					"members_count": 1000,
					"photo_100": "https://example.com/group_100.jpg",
					"description": "test group"
				}
			]
		}`))
	}))
	defer srv.Close()

	client := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	items, err := GetByID(context.Background(), client, GetByIDParams{
		GroupIDs: []string{"vk_test"},
		Fields:   []string{"members_count", "photo_100", "description"},
	})
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if gotPath != "/groups.getById" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("group_ids") != "vk_test" {
		t.Fatalf("unexpected group_ids: %q", gotForm.Get("group_ids"))
	}
	if gotForm.Get("fields") != "members_count,photo_100,description" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}

	if len(items) != 1 {
		t.Fatalf("unexpected groups len: %d", len(items))
	}
	if items[0].Name != "VK Test Group" {
		t.Fatalf("unexpected name: %q", items[0].Name)
	}
	if items[0].MembersCount != 1000 {
		t.Fatalf("unexpected members_count: %d", items[0].MembersCount)
	}
	if items[0].Photo100 != "https://example.com/group_100.jpg" {
		t.Fatalf("unexpected photo_100: %q", items[0].Photo100)
	}
	if items[0].Description != "test group" {
		t.Fatalf("unexpected description: %q", items[0].Description)
	}
}

func TestGetMembers(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"count": 2,
				"items": [
					{
						"id": 10,
						"first_name": "Ivan",
						"last_name": "Ivanov",
						"photo_100": "https://example.com/u1.jpg",
						"online": 1,
						"can_access_closed": true,
						"is_closed": false
					},
					{
						"id": 20,
						"first_name": "Petr",
						"last_name": "Petrov",
						"sex": 2,
						"can_access_closed": true,
						"is_closed": false
					}
				]
			}
		}`))
	}))
	defer srv.Close()

	client := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	resp, err := GetMembers(context.Background(), client, GetMembersParams{
		GroupID: "vk_test",
		Sort:    "id_asc",
		Offset:  10,
		Count:   2,
		Fields:  []string{"photo_100", "online", "sex"},
		Filter:  "friends",
	})
	if err != nil {
		t.Fatalf("GetMembers() error = %v", err)
	}

	if gotPath != "/groups.getMembers" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("group_id") != "vk_test" {
		t.Fatalf("unexpected group_id: %q", gotForm.Get("group_id"))
	}
	if gotForm.Get("sort") != "id_asc" {
		t.Fatalf("unexpected sort: %q", gotForm.Get("sort"))
	}
	if gotForm.Get("offset") != "10" {
		t.Fatalf("unexpected offset: %q", gotForm.Get("offset"))
	}
	if gotForm.Get("count") != "2" {
		t.Fatalf("unexpected count: %q", gotForm.Get("count"))
	}
	if gotForm.Get("fields") != "photo_100,online,sex" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}
	if gotForm.Get("filter") != "friends" {
		t.Fatalf("unexpected filter: %q", gotForm.Get("filter"))
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Count != 2 {
		t.Fatalf("unexpected count: %d", resp.Count)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("unexpected items len: %d", len(resp.Items))
	}
	if resp.Items[0].FirstName != "Ivan" {
		t.Fatalf("unexpected first_name: %q", resp.Items[0].FirstName)
	}
	if resp.Items[0].Photo100 != "https://example.com/u1.jpg" {
		t.Fatalf("unexpected photo_100: %q", resp.Items[0].Photo100)
	}
	if resp.Items[1].Sex != 2 {
		t.Fatalf("unexpected sex: %d", resp.Items[1].Sex)
	}
}

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"count": 1,
				"items": [{"id": 1, "name": "Test Group", "type": "group"}]
			}
		}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	resp, err := Get(context.Background(), c, GetParams{
		UserID:   1,
		Extended: true,
		Filter:   []string{"admin"},
		Fields:   []string{"activity"},
		Offset:   0,
		Count:    10,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if resp == nil || resp.Count != 1 {
		t.Fatal("expected response with count=1")
	}
}

func TestGetAddresses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"count": 1,
				"items": [{
					"id": 1,
					"title": "Main Office",
					"address": "Street 1",
					"country_id": 1,
					"city_id": 1,
					"latitude": 55.75,
					"longitude": 37.61,
					"is_main_address": true
				}]
			}
		}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	resp, err := GetAddresses(context.Background(), c, GetAddressesParams{
		GroupID: 1,
		Count:   10,
		Fields:  []string{"city", "country"},
	})
	if err != nil {
		t.Fatalf("GetAddresses() error = %v", err)
	}
	if resp == nil || resp.Count != 1 {
		t.Fatal("expected response")
	}
}

func TestGetBanned(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"count": 1,
				"items": [{
					"type": "profile",
					"profile": {"id": 123, "first_name": "Banned"},
					"ban_info": {"admin_id": 1, "date": 1234567890, "reason": 1, "end_date": 0}
				}]
			}
		}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	resp, err := GetBanned(context.Background(), c, GetBannedParams{
		GroupID: 1,
		Count:   10,
	})
	if err != nil {
		t.Fatalf("GetBanned() error = %v", err)
	}
	if resp == nil || resp.Count != 1 {
		t.Fatal("expected response")
	}
}

func TestAddAddress(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"id": 1,
				"title": "New Office",
				"address": "Street 2",
				"country_id": 1,
				"city_id": 1,
				"is_main_address": false
			}
		}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	resp, err := AddAddress(context.Background(), c, AddAddressParams{
		GroupID:       1,
		Title:         "New Office",
		Address:       "Street 2",
		CountryID:     1,
		CityID:        1,
		IsMainAddress: false,
	})
	if err != nil {
		t.Fatalf("AddAddress() error = %v", err)
	}
	if resp == nil || resp.Title != "New Office" {
		t.Fatal("expected response")
	}
}

func TestEditAddress(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"id": 1,
				"title": "Edited Office",
				"address": "Street 3",
				"is_main_address": true
			}
		}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	resp, err := EditAddress(context.Background(), c, EditAddressParams{
		GroupID:   1,
		AddressID: 1,
		Title:     "Edited Office",
		Address:   "Street 3",
	})
	if err != nil {
		t.Fatalf("EditAddress() error = %v", err)
	}
	if resp == nil || resp.Title != "Edited Office" {
		t.Fatal("expected response")
	}
}

func TestDeleteAddress(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response": 1}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	ok, err := DeleteAddress(context.Background(), c, DeleteAddressParams{
		GroupID:   1,
		AddressID: 1,
	})
	if err != nil {
		t.Fatalf("DeleteAddress() error = %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
}

func TestAddCallbackServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response": {"server_id": 123}}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	resp, err := AddCallbackServer(context.Background(), c, AddCallbackServerParams{
		GroupID:   1,
		URL:       "https://example.com/callback",
		Title:     "Test Server",
		SecretKey: "secret",
	})
	if err != nil {
		t.Fatalf("AddCallbackServer() error = %v", err)
	}
	if resp == nil || resp.ServerID != 123 {
		t.Fatal("expected response with server_id=123")
	}
}

func TestEditCallbackServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response": 1}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	ok, err := EditCallbackServer(context.Background(), c, EditCallbackServerParams{
		GroupID:   1,
		ServerID:  123,
		URL:       "https://new.example.com",
		Title:     "New Title",
		SecretKey: "new_secret",
	})
	if err != nil {
		t.Fatalf("EditCallbackServer() error = %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
}

func TestDeleteCallbackServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response": 1}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	ok, err := DeleteCallbackServer(context.Background(), c, DeleteCallbackServerParams{
		GroupID:  1,
		ServerID: 123,
	})
	if err != nil {
		t.Fatalf("DeleteCallbackServer() error = %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
}

func TestDisableOnline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response": 1}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	ok, err := DisableOnline(context.Background(), c, DisableOnlineParams{GroupID: 1})
	if err != nil {
		t.Fatalf("DisableOnline() error = %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
}

func TestEnableOnline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response": 1}`))
	}))
	defer srv.Close()

	c := client.New(config.DefaultConfig(), client.WithBaseURL(srv.URL))

	ok, err := EnableOnline(context.Background(), c, EnableOnlineParams{GroupID: 1})
	if err != nil {
		t.Fatalf("EnableOnline() error = %v", err)
	}
	if !ok {
		t.Error("expected true")
	}
}

func TestValidateParams(t *testing.T) {
	t.Run("GetByIDParams valid", func(t *testing.T) {
		p := GetByIDParams{GroupIDs: []string{"1"}}
		if err := p.Validate(); err != nil {
			t.Errorf("Validate() error = %v", err)
		}
	})

	t.Run("GetByIDParams invalid", func(t *testing.T) {
		p := GetByIDParams{}
		if err := p.Validate(); err == nil {
			t.Error("Validate() should return error")
		}
	})

	t.Run("GetMembersParams valid", func(t *testing.T) {
		p := GetMembersParams{GroupID: "1", Count: 10}
		if err := p.Validate(); err != nil {
			t.Errorf("Validate() error = %v", err)
		}
	})

	t.Run("GetMembersParams invalid count", func(t *testing.T) {
		p := GetMembersParams{GroupID: "1", Count: -1}
		if err := p.Validate(); err == nil {
			t.Error("Validate() should return error for negative count")
		}
	})

	t.Run("AddAddressParams valid", func(t *testing.T) {
		p := AddAddressParams{GroupID: 1, Title: "Test"}
		if err := p.Validate(); err != nil {
			t.Errorf("Validate() error = %v", err)
		}
	})

	t.Run("AddAddressParams invalid", func(t *testing.T) {
		p := AddAddressParams{}
		if err := p.Validate(); err == nil {
			t.Error("Validate() should return error")
		}
	})

	t.Run("GetBannedParams valid", func(t *testing.T) {
		p := GetBannedParams{GroupID: 1, Count: 10}
		if err := p.Validate(); err != nil {
			t.Errorf("Validate() error = %v", err)
		}
	})

	t.Run("GetAddressesParams valid", func(t *testing.T) {
		p := GetAddressesParams{GroupID: 1, Count: 10}
		if err := p.Validate(); err != nil {
			t.Errorf("Validate() error = %v", err)
		}
	})
}
