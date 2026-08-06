package domain

type ProductType string

const (
	ProductTypeUnspecified  ProductType = "UNSPECIFIED"
	ProductTypeSubscription ProductType = "SUBSCRIPTION"
	ProductTypeBundle       ProductType = "BUNDLE"
)

type Product struct {
	ID            string
	Type          ProductType
	Name          string
	Description   string
	DurationDays  int
	PriceAmount   int64
	PriceCurrency string
	BundleItems   []BundleItem
}

type BundleItem struct {
	ProductID string
	Name      string
}
