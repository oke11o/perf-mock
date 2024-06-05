package handler

type AuthResponse struct {
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type ListResponse struct {
	Result []*ListItem `json:"result,omitempty"`
}

type ListItem struct {
	ItemId int64 ` json:"item_id,omitempty"`
}

type OrderResponse struct {
	OrderId int64 `json:"order_id,omitempty"`
}
