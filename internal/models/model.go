package models

import "time"

// Cake model
type Cake struct {
	ID              int    `json:"id"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
	Weight          int    `json:"weight"`
	FullDescription string `json:"full_description"`
	Active          bool   `json:"active"`
	ImageURL        string `json:"image_url"`
}

// User model
type User struct {
	ID           int    `json:"id"`
	FullName     string `json:"fullname"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	PhoneNumber  string `json:"phone_number"`
}

// Delivery model
type Delivery struct {
	ID      int    `json:"id"`
	PointID int    `json:"point_id"`
	Cost    int    `json:"cost"`
	Status  string `json:"status"`
	Weight  int    `json:"weight"`
}

// DeliveryPoint model
type DeliveryPoint struct {
	ID           int    `json:"id"`
	Address      string `json:"address"`
	Rating       int    `json:"rating"`
	WorkingHours string `json:"status"`
	ContactPhone string `json:"contact_phone"`
}

// Order model
type Order struct {
	ID            int       `json:"id"`
	Time          time.Time `json:"time"`
	OrderStatus   string    `json:"order_status"`
	UserID        int       `json:"user_id"`
	PaymentMethod string    `json:"payment_method"`
	Cost          int       `json:"cost"`
}

// OrderCake model
type OrderCake struct {
	ID      int `json:"id"`
	OrderID int `json:"order_id"`
	CakeID  int `json:"cake_id"`
}

// Admin model
type Admin struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

// LogEntry model
type LogEntry struct {
	ID        int       `db:"id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

// InternOrder model
type InternOrder struct {
	Cakes []Cake
	Ord   Order
}

// GetOrdersResponse model
type GetOrdersResponse struct {
	Orders []InternOrder
}

// MakeOrderRequest model
type MakeOrderRequest struct {
	UserID        int      `json:"user_id"`
	Delivery      Delivery `json:"delivery"`
	Cakes         []Cake   `json:"cakes"`
	PaymentMethod string   `json:"payment_method"`
}

// MakeOrderResponse model
type MakeOrderResponse struct {
	OrderID      int    `json:"order_id"`
	DeliveryTime string `json:"delivery_time"`
}

// ViewOrdersRequest model
type ViewOrdersRequest struct{}

// ViewOrdersResponse model
type ViewOrdersResponse struct {
	Orders []struct {
		Cakes             []Cake    `json:"cakes"`
		OrderRegisterTime time.Time `json:"order_register_time"`
		PaymentMethod     string    `json:"payment_method"`
		OrderStatus       string    `json:"order_status"`
		DelivPoint        string    `json:"delivery_point"`
	} `json:"orders"`
}

// CancelOrderRequest model
type CancelOrderRequest struct {
	OrderID int `json:"order_id"`
}

type UpdateOrderRequest struct {
	OrderID       int    `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
}
