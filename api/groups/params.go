package groups

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
