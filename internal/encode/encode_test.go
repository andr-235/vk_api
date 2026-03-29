package encode

import "testing"

type testParams struct {
	UserIDs []int    `url:"user_ids,comma,omitempty"`
	Fields  []string `url:"fields,comma,omitempty"`
	Name    string   `url:"name,omitempty"`
	Active  bool     `url:"active,omitempty"`
	Count   int      `url:"count,omitempty"`
}

func TestValues_Struct(t *testing.T) {
	v, err := Values(testParams{
		UserIDs: []int{1, 2, 3},
		Fields:  []string{"photo_100", "city"},
		Name:    "durov",
		Active:  true,
	})
	if err != nil {
		t.Fatalf("Values() error = %v", err)
	}

	if got := v.Get("user_ids"); got != "1,2,3" {
		t.Fatalf("unexpected user_ids: %q", got)
	}
	if got := v.Get("fields"); got != "photo_100,city" {
		t.Fatalf("unexpected fields: %q", got)
	}
	if got := v.Get("name"); got != "durov" {
		t.Fatalf("unexpected name: %q", got)
	}
	if got := v.Get("active"); got != "1" {
		t.Fatalf("unexpected active: %q", got)
	}
	if got := v.Get("count"); got != "" {
		t.Fatalf("expected omitted count, got %q", got)
	}
}
