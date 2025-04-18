package schema

import "time"

type OrderWine struct {
	WineID   int `json:"wine_id"`
	Quantity int `json:"quantity"`
}

type CreatePaymentIntentRequest struct {
	Amount     int64       `json:"amount"`
	Currency   string      `json:"currency"`
	UserID     int         `json:"user_id"`
	OrderWines []OrderWine `json:"order_wines"`
}

type CreatePaymentIntentResponse struct {
	ClientSecret string `json:"client_secret"`
}

type OrderWineResponse struct {
	WineID   int  `json:"wine_id"`
	Wine     Wine `json:"wine"`
	Quantity int  `json:"quantity"`
}

type GetOrderResponse struct {
	ID               string              `json:"id"`
	UserID           int                 `json:"user_id"`
	StripePurchaseID string              `json:"stripe_purchase_id"`
	Amount           int64               `json:"amount"`
	Status           string              `json:"status"`
	OrderWines       []OrderWineResponse `json:"order_wines"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}
