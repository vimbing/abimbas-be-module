package adidas_backend

type CheckoutIdPayload struct {
	Items []Items `json:"items"`
}

type Items struct {
	ProductID          string `json:"product_id"`
	VariationProductID string `json:"variation_product_id"`
	Quantity           int    `json:"quantity"`
}

type PaymentPayload struct {
	ID string `json:"id"`
}

type AddressFillPayload struct {
	Phone       string `json:"phone"`
	City        string `json:"city"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	PostalCode  string `json:"postal_code"`
	LastName    string `json:"last_name"`
	CountryCode string `json:"country_code"`
	Type        string `json:"type"`
	FirstName   string `json:"first_name"`
}

type AddressPayload struct {
	LocationID string `json:"location_id"`
}

type OrderPayload struct {
	CheckoutID string `json:"checkout_id"`
}

type DeliveryPayload struct {
	ShippingMethodID string `json:"shipping_method_id"`
}

type PickupPayload struct {
	LocationID  string      `json:"location_id"`
	GeoLocation GeoLocation `json:"geo_location"`
}

type GeoLocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type DiscountPayload struct {
	VoucherCode string `json:"voucher_code"`
}

type RegisterPayload struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	MembershipConsent bool   `json:"membership_consent"`
	DormantPeriod     string `json:"dormant_period"`
}

type AddPaypalLinkToApiPayload struct {
	Body string `json:"body"`
}

type PaymentCCPayload struct {
	CheckoutID string  `json:"checkout_id"`
	NewCard    NewCard `json:"new_card"`
}
type NewCard struct {
	Type       string `json:"type"`
	Encryption string `json:"encryption"`
}
