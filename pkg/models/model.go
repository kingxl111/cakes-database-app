package models

// Database model
type Cake struct {
	ID 			int		`json:"id"`
	Description string	`json:"description"`	
	Price		int		`json:"price"`
	Weight		int		`json:"weight"`
}

type User struct {
	ID 			 int	`json:"id"`
	FullName 	 string	`json:"fullname"`
	Username 	 string `json:"username"`
	Email 		 string	`json:"email"`
	PasswordHash string	`json:"password_hash"`
	PhoneNumber  string	`json:"phone_number"`
}

type Delivery struct {
	ID 		int			`json:"id"`
	PointID int			`json:"point_id"`
	Cost 	int			`json:"cost"`
	Status 	string 		`json:"status"`
	Weight 	int			`json:"weight"`
}

type DeliveryPoint struct {
	ID 				int			`json:"id"`
	Address 		string		`json:"address"`
	Rating 			int			`json:"rating"`
	WorkingHours 	string 		`json:"status"`
	ContactPhone 	int			`json:"weight"`
}

type Order struct {
	ID 				int 	`json:"id"`
	Time 			string	`json:"time"`
	OrderStatus 	string	`json:"order_status"`
	UserID 			int		`json:"user_id"`
	PaymentMethod 	string 	`json:"payment_method"`
}

type OrderCake struct {
	ID 		int `json:"id"`
	OrderID int `json:"order_id"`
	CakeID 	int `json:"cake_id"`
}

type Admin struct {
	ID 				int 	`json:"id"`
	Username 		string 	`json:"username"`
	PasswordHash 	string 	`json:"password_hash"`
}

// RequestModels
type MakeOrderRequest struct {
	UserID		 int 		  `json:"userID"`
    Delivery     Delivery     `json:"delivery"`     
    Cakes        []Cake       `json:"cakes"`         
    PaymentMethod string      `json:"paymentMethod"`
}