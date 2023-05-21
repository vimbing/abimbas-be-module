package adidas_backend

import "time"

type LoginResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type CheckoutIdResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		CheckoutCreation struct {
			Href string `json:"href"`
		} `json:"checkout_creation"`
		CheckoutAbandon struct {
			Href string `json:"href"`
		} `json:"checkout_abandon"`
		DeliveryOptionLookup struct {
			Href string `json:"href"`
		} `json:"delivery_option_lookup"`
		PaymentOptionLookup struct {
			Href string `json:"href"`
		} `json:"payment_option_lookup"`
		BillingAddressUpdate struct {
			Href string `json:"href"`
		} `json:"billing_address_update"`
		VoucherApplication struct {
			Href string `json:"href"`
		} `json:"voucher_application"`
		OrderCreation struct {
			Name string `json:"name"`
			Href string `json:"href"`
		} `json:"order_creation"`
	} `json:"_links"`
	ID       string `json:"id"`
	Currency string `json:"currency"`
	Items    []struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Thumbnail struct {
				Href string `json:"href"`
			} `json:"thumbnail"`
		} `json:"_links"`
		Product struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Color     string `json:"color"`
		} `json:"product"`
		Size struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"size"`
		Quantity int `json:"quantity"`
		Prices   struct {
			Total struct {
				Current  int `json:"current"`
				Original int `json:"original"`
			} `json:"total"`
			Unit struct {
				Current  int `json:"current"`
				Original int `json:"original"`
			} `json:"unit"`
			Tax float64 `json:"tax"`
		} `json:"prices"`
	} `json:"items"`
	Selected struct {
		Delivery struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Location struct {
				ID          string `json:"id"`
				FirstName   string `json:"first_name"`
				LastName    string `json:"last_name"`
				FullName    string `json:"full_name"`
				Address1    string `json:"address1"`
				Phone       string `json:"phone"`
				City        string `json:"city"`
				CountryCode string `json:"country_code"`
				PostalCode  string `json:"postal_code"`
			} `json:"location"`
			Lines []struct {
				ID             string `json:"id"`
				ShippingMethod struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Cost struct {
						DisplayValue string `json:"display_value"`
						Value        int    `json:"value"`
					} `json:"cost"`
					DeliveryTime struct {
						Date struct {
							From string `json:"from"`
							To   string `json:"to"`
						} `json:"date"`
						Time struct {
							From string `json:"from"`
							To   string `json:"to"`
						} `json:"time"`
						Slot struct {
							From time.Time `json:"from"`
							To   time.Time `json:"to"`
						} `json:"slot"`
					} `json:"delivery_time"`
				} `json:"shipping_method"`
			} `json:"lines"`
		} `json:"delivery"`
		PaymentMethod struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Config struct {
				Type string `json:"type"`
				Icon string `json:"icon"`
			} `json:"config"`
		} `json:"payment_method"`
		BillingAddress struct {
			ID          string `json:"id"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			FullName    string `json:"full_name"`
			Address1    string `json:"address1"`
			Phone       string `json:"phone"`
			City        string `json:"city"`
			CountryCode string `json:"country_code"`
			PostalCode  string `json:"postal_code"`
		} `json:"billing_address"`
	} `json:"selected"`
	PaymentSummary struct {
		IsPaymentRequired bool `json:"is_payment_required"`
		Amounts           []struct {
			Payment struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Config struct {
					Type string `json:"type"`
					Icon string `json:"icon"`
				} `json:"config"`
			} `json:"payment"`
			Amount struct {
				DisplayValue string `json:"display_value"`
				Value        int    `json:"value"`
			} `json:"amount"`
		} `json:"amounts"`
	} `json:"payment_summary"`
	OrderSummary struct {
		Lines []struct {
			Type         string `json:"type"`
			Title        string `json:"title"`
			DisplayValue string `json:"display_value"`
			Value        int    `json:"value"`
		} `json:"lines"`
		TotalPrice struct {
			Type         string `json:"type"`
			Title        string `json:"title"`
			DisplayValue string `json:"display_value"`
			Value        int    `json:"value"`
		} `json:"total_price"`
		MembershipPoints struct {
			Type         string `json:"type"`
			Title        string `json:"title"`
			DisplayValue string `json:"display_value"`
			Value        int    `json:"value"`
		} `json:"membership_points"`
	} `json:"order_summary"`
	LegalLinks []struct {
		Type string `json:"type"`
		Body string `json:"body,omitempty"`
		Name string `json:"name"`
		Link string `json:"link"`
	} `json:"legal_links"`
	Terms []struct {
		ID          string `json:"id"`
		Text        string `json:"text"`
		Disclaimers []struct {
			Title       string `json:"title"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"disclaimers"`
	} `json:"terms"`
}

type GetItemDataResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	ProductID       string `json:"product_id"`
	ModelNumber     string `json:"model_number"`
	OriginalPrice   int    `json:"original_price"`
	DisplayCurrency string `json:"display_currency"`
	Orderable       bool   `json:"orderable"`
	BadgeText       string `json:"badge_text"`
	BadgeColor      string `json:"badge_color"`
	Embedded        struct {
		Variations []struct {
			Size               string  `json:"size"`
			TechnicalSize      string  `json:"technical_size"`
			Orderable          bool    `json:"orderable"`
			AbsoluteSize       float64 `json:"absolute_size"`
			VariationProductID string  `json:"variation_product_id"`
			StockLevel         int     `json:"stock_level"`
			Links              struct {
				SimilarProducts struct {
					Href string `json:"href"`
				} `json:"similar_products"`
			} `json:"_links,omitempty"`
			LowOnStockMessage string `json:"low_on_stock_message,omitempty"`
		} `json:"variations"`
	} `json:"_embedded"`
	IsHype bool `json:"is_hype"`
}

type OrderResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		VisualImage struct {
			Href string `json:"href"`
		} `json:"visual_image"`
		Help struct {
			Href string `json:"href"`
		} `json:"help"`
		OrderConfirmationMessage struct {
			Href string `json:"href"`
		} `json:"order_confirmation_message"`
	} `json:"_links"`
	ThirdPartyPayment struct {
		PaymentMethodID   string `json:"payment_method_id"`
		PaymentMethodName string `json:"payment_method_name"`
		Base64HTML        string `json:"base64_html"`
	} `json:"third_party_payment"`
	OrderType            string    `json:"order_type"`
	OrderNumber          string    `json:"order_number"`
	TotalPrice           int       `json:"total_price"`
	OrderItemsSubtotal   int       `json:"order_items_subtotal"`
	ShippingSubtotal     int       `json:"shipping_subtotal"`
	PromotionsSubtotal   int       `json:"promotions_subtotal"`
	TaxTotal             float64   `json:"tax_total"`
	Currency             string    `json:"currency"`
	OrderDate            time.Time `json:"order_date"`
	OrderStatus          string    `json:"order_status"`
	OrderStatusLocalized string    `json:"order_status_localized"`
	BillingInfo          struct {
		ID          string `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		FullName    string `json:"full_name"`
		Address1    string `json:"address1"`
		Phone       string `json:"phone"`
		City        string `json:"city"`
		CountryCode string `json:"country_code"`
		PostalCode  string `json:"postal_code"`
	} `json:"billing_info"`
	ShippingInfo struct {
		ID             string `json:"id"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		FullName       string `json:"full_name"`
		Address1       string `json:"address1"`
		Phone          string `json:"phone"`
		City           string `json:"city"`
		CountryCode    string `json:"country_code"`
		PostalCode     string `json:"postal_code"`
		ShippingMethod string `json:"shipping_method"`
	} `json:"shipping_info"`
	CustomerInfo struct {
		CustomerEmail string `json:"customer_email"`
	} `json:"customer_info"`
	PaymentInfo struct {
		Amount            int    `json:"amount"`
		Currency          string `json:"currency"`
		PaymentMethodID   string `json:"payment_method_id"`
		PaymentMethodIcon string `json:"payment_method_icon"`
		PaymentMethodName string `json:"payment_method_name"`
	} `json:"payment_info"`
	TotalView []struct {
		Title        string `json:"title"`
		Value        int    `json:"value"`
		Type         string `json:"type"`
		DisplayValue string `json:"display_value"`
	} `json:"total_view"`
	ProductGroups []struct {
		ID                      string    `json:"id"`
		ExpectedDelivery        time.Time `json:"expected_delivery"`
		DeliveryStatus          string    `json:"delivery_status"`
		DeliveryStatusLocalized string    `json:"delivery_status_localized"`
		ShippingName            string    `json:"shipping_name"`
		OrderItems              []struct {
			Links struct {
				Thumbnail struct {
					Href string `json:"href"`
				} `json:"thumbnail"`
				Product struct {
					Href string `json:"href"`
				} `json:"product"`
				Variations struct {
					Href string `json:"href"`
				} `json:"variations"`
			} `json:"_links"`
			ItemID             string `json:"item_id"`
			OrderLineID        string `json:"order_line_id"`
			ProductID          string `json:"product_id"`
			ProductName        string `json:"product_name"`
			ProductType        string `json:"product_type"`
			OriginalPrice      int    `json:"original_price"`
			CurrentPrice       int    `json:"current_price"`
			Quantity           int    `json:"quantity"`
			Size               string `json:"size"`
			VariationProductID string `json:"variation_product_id"`
			ReturnInfo         struct {
				Status string `json:"status"`
			} `json:"return_info"`
			Mtbr      bool   `json:"mtbr"`
			ColorName string `json:"color_name"`
		} `json:"order_items"`
	} `json:"product_groups"`
	Cancelable           bool          `json:"cancelable"`
	Returnable           bool          `json:"returnable"`
	AvailableCancelTypes []interface{} `json:"available_cancel_types"`
	Embedded             struct {
		RenderedContent []struct {
			Href           string `json:"href"`
			Type           string `json:"type"`
			Block          string `json:"block"`
			AdditionalInfo struct {
				PromotionID   string `json:"promotionId"`
				PromotionName string `json:"promotionName"`
			} `json:"additional_info"`
		} `json:"rendered_content"`
	} `json:"_embedded"`
}

type DiscountedOrderResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		VisualImage struct {
			Href string `json:"href"`
		} `json:"visual_image"`
		Help struct {
			Href string `json:"href"`
		} `json:"help"`
	} `json:"_links"`
	ThirdPartyPayment struct {
		PaymentMethodID   string `json:"payment_method_id"`
		PaymentMethodName string `json:"payment_method_name"`
		Base64HTML        string `json:"base64_html"`
	} `json:"third_party_payment"`
	OrderType            string    `json:"order_type"`
	OrderNumber          string    `json:"order_number"`
	TotalPrice           float64   `json:"total_price"`
	OrderItemsSubtotal   float64   `json:"order_items_subtotal"`
	ShippingSubtotal     int       `json:"shipping_subtotal"`
	PromotionsSubtotal   float64   `json:"promotions_subtotal"`
	TaxTotal             float64   `json:"tax_total"`
	Currency             string    `json:"currency"`
	OrderDate            time.Time `json:"order_date"`
	OrderStatus          string    `json:"order_status"`
	OrderStatusLocalized string    `json:"order_status_localized"`
	BillingInfo          struct {
		ID          string `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		FullName    string `json:"full_name"`
		Address1    string `json:"address1"`
		Phone       string `json:"phone"`
		City        string `json:"city"`
		CountryCode string `json:"country_code"`
		PostalCode  string `json:"postal_code"`
	} `json:"billing_info"`
	ShippingInfo struct {
		ID             string `json:"id"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		FullName       string `json:"full_name"`
		Address1       string `json:"address1"`
		Phone          string `json:"phone"`
		City           string `json:"city"`
		CountryCode    string `json:"country_code"`
		PostalCode     string `json:"postal_code"`
		ShippingMethod string `json:"shipping_method"`
	} `json:"shipping_info"`
	CustomerInfo struct {
		CustomerEmail string `json:"customer_email"`
	} `json:"customer_info"`
	PaymentInfo struct {
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"`
		PaymentMethodID   string  `json:"payment_method_id"`
		PaymentMethodIcon string  `json:"payment_method_icon"`
		PaymentMethodName string  `json:"payment_method_name"`
	} `json:"payment_info"`
	TotalView []struct {
		Title        string `json:"title"`
		Value        int    `json:"value"`
		Type         string `json:"type"`
		DisplayValue string `json:"display_value"`
	} `json:"total_view"`
	ProductGroups []struct {
		ID                      string    `json:"id"`
		ExpectedDelivery        time.Time `json:"expected_delivery"`
		DeliveryStatus          string    `json:"delivery_status"`
		DeliveryStatusLocalized string    `json:"delivery_status_localized"`
		ShippingName            string    `json:"shipping_name"`
		OrderItems              []struct {
			Links struct {
				Thumbnail struct {
					Href string `json:"href"`
				} `json:"thumbnail"`
				Product struct {
					Href string `json:"href"`
				} `json:"product"`
				Variations struct {
					Href string `json:"href"`
				} `json:"variations"`
			} `json:"_links"`
			ItemID                 string  `json:"item_id"`
			OrderLineID            string  `json:"order_line_id"`
			ProductID              string  `json:"product_id"`
			ProductName            string  `json:"product_name"`
			ProductType            string  `json:"product_type"`
			OriginalPrice          int     `json:"original_price"`
			CurrentPrice           float64 `json:"current_price"`
			DiscountPercentageText string  `json:"discount_percentage_text"`
			Quantity               int     `json:"quantity"`
			Size                   string  `json:"size"`
			VariationProductID     string  `json:"variation_product_id"`
			ReturnInfo             struct {
				Status string `json:"status"`
			} `json:"return_info"`
			Mtbr      bool   `json:"mtbr"`
			ColorName string `json:"color_name"`
		} `json:"order_items"`
	} `json:"product_groups"`
	Cancelable           bool          `json:"cancelable"`
	Returnable           bool          `json:"returnable"`
	AvailableCancelTypes []interface{} `json:"available_cancel_types"`
	Embedded             struct {
		RenderedContent []interface{} `json:"rendered_content"`
	} `json:"_embedded"`
}

type GetAddressesResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		User struct {
			Href string `json:"href"`
		} `json:"user"`
	} `json:"_links"`
	Addresses []struct {
		ID              string `json:"id"`
		AddressType     string `json:"address_type"`
		Phone           string `json:"phone"`
		Address1        string `json:"address1"`
		City            string `json:"city"`
		PostalCode      string `json:"postal_code"`
		Country         string `json:"country"`
		CountryCode     string `json:"country_code"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Type            string `json:"type"`
		UsedForBilling  bool   `json:"used_for_billing"`
		UsedForShipping bool   `json:"used_for_shipping"`
	} `json:"addresses"`
}

type RegisterResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Login struct {
			Href string `json:"href"`
		} `json:"login"`
		Logout struct {
			Href string `json:"href"`
		} `json:"logout"`
		User struct {
			Href string `json:"href"`
		} `json:"user"`
	} `json:"_links"`
	ID           string `json:"id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type UmbrellaApiLinkAddResponse struct {
	URL string `json:"url"`
}

type ItemDataFrontendResponse struct {
	ID                 string `json:"id"`
	AvailabilityStatus string `json:"availability_status"`
	VariationList      []struct {
		Sku                string `json:"sku"`
		Size               string `json:"size"`
		Availability       int    `json:"availability"`
		AvailabilityStatus string `json:"availability_status"`
	} `json:"variation_list"`
}
