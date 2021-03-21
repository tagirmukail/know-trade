package types

type Order struct {
	ID     string
	Symbol string
	Side   string
	Status string
	Price  float64
	Size   float64
	Other  map[string]interface{}
}
