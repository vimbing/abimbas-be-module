package adidas_backend

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	BlankCosmetics = Cosmetics{
		Image: "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ee/Logo_brand_Adidas.png/608px-Logo_brand_Adidas.png",
		Name:  "NOT FOUND",
		Price: "NOT FOUND",
		Size:  "NOT FOUND",
	}
)

func GetCosmetics(checkoutIdResponse *CheckoutIdResponse) Cosmetics {
	var cosmetics Cosmetics

	if len(checkoutIdResponse.Items) < 1 {
		return BlankCosmetics
	}

	product := checkoutIdResponse.Items[0]

	cosmetics.Image = product.Product.Thumbnail
	cosmetics.Name = product.Product.Name
	cosmetics.Price = fmt.Sprint(product.Prices.Total.Current)
	cosmetics.Size = product.Size.Name

	return cosmetics
}

func GetOrderNumber(orderResponse *OrderResponse) string {
	if len(orderResponse.ProductGroups) < 1 {
		return uuid.NewString()
	}

	if len(orderResponse.ProductGroups[0].OrderItems) < 1 {
		return uuid.NewString()
	}

	return orderResponse.OrderNumber
}
