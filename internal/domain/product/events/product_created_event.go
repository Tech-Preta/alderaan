package product_events

type ProductCreatedEvent struct {
	Name       string
	Sku        int
	Categories []string
	Price      int
}

func NewProductCreatedEvent(name string, sku int, categories []string, price int) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		Name:       name,
		Sku:        sku,
		Categories: categories,
		Price:      price,
	}
}

func (e *ProductCreatedEvent) EventName() string {
	return "product.created"
}
