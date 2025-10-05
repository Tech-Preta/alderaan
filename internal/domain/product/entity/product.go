package product_entity

import (
	"errors"

	product_events "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/events"
	shared_events "github.com/williamkoller/golang-domain-driven-design/internal/shared/domain/events"
)

type Product struct {
	Name       string
	Sku        int
	Categories []string
	Price      int
	*product_events.ProductCreatedEvent
}

func NewProduct(name string, sku int, categories []string, price int, dispatcher *shared_events.EventDispatcher) (*Product, *product_events.ProductCreatedEvent, error) {
	ok, err := Validate(name, sku, categories, price)

	if !ok {
		return nil, nil, err
	}

	p := &Product{Name: name, Sku: sku, Categories: categories, Price: price}

	event := product_events.NewProductCreatedEvent(p.GetName(), p.GetSku(), p.GetCategories(), p.GetPrice())
	if dispatcher != nil {
		dispatcher.Dispatch(event.EventName(), event)
	}

	return p, event, nil

}

func Validate(name string, sku int, categories []string, price int) (bool, error) {
	if name == "" {
		return false, errors.New("name is required")
	}

	if sku <= 0 {
		return false, errors.New("sku is required")
	}

	if len(categories) == 0 {
		return false, errors.New("categories is required")
	}

	if price <= 0 {
		return false, errors.New("price is required")
	}

	return true, nil
}

func (p *Product) GetName() string {
	return p.Name
}

func (p *Product) GetSku() int {
	return p.Sku
}

func (p *Product) GetCategories() []string {
	return p.Categories
}

func (p *Product) GetPrice() int {
	return p.Price
}
