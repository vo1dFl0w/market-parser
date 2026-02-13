package chromium

type ProductsResponse struct {
	Prods []Product `json:"products"`
}

type Product struct {
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	CanonicalURL string  `json:"canonical_url"`
}
