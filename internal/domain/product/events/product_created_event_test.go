package product_events

import (
	"testing"
)

func TestNewProductCreatedEvent(t *testing.T) {
	tests := []struct {
		name       string
		prodName   string
		sku        int
		categories []string
		price      int
	}{
		{
			name:       "valid event",
			prodName:   "Notebook",
			sku:        12345,
			categories: []string{"Electronics", "Computers"},
			price:      3500,
		},
		{
			name:       "minimal event",
			prodName:   "Test",
			sku:        1,
			categories: []string{"Test"},
			price:      1,
		},
		{
			name:       "multiple categories",
			prodName:   "Gaming Setup",
			sku:        99999,
			categories: []string{"Gaming", "Electronics", "Computers", "Peripherals"},
			price:      15000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewProductCreatedEvent(tt.prodName, tt.sku, tt.categories, tt.price)

			if event == nil {
				t.Fatal("NewProductCreatedEvent() returned nil")
			}

			// Verificar campos
			if event.Name != tt.prodName {
				t.Errorf("Name = %v, want %v", event.Name, tt.prodName)
			}

			if event.Sku != tt.sku {
				t.Errorf("Sku = %v, want %v", event.Sku, tt.sku)
			}

			if len(event.Categories) != len(tt.categories) {
				t.Errorf("Categories length = %d, want %d", len(event.Categories), len(tt.categories))
			}

			for i, cat := range tt.categories {
				if event.Categories[i] != cat {
					t.Errorf("Categories[%d] = %v, want %v", i, event.Categories[i], cat)
				}
			}

			if event.Price != tt.price {
				t.Errorf("Price = %v, want %v", event.Price, tt.price)
			}
		})
	}
}

func TestProductCreatedEvent_EventName(t *testing.T) {
	event := NewProductCreatedEvent("Test", 123, []string{"Test"}, 100)

	expectedName := "product.created"
	if event.EventName() != expectedName {
		t.Errorf("EventName() = %v, want %v", event.EventName(), expectedName)
	}
}

func TestProductCreatedEvent_Fields(t *testing.T) {
	event := &ProductCreatedEvent{
		Name:       "Test Product",
		Sku:        999,
		Categories: []string{"Cat1", "Cat2", "Cat3"},
		Price:      500,
	}

	t.Run("Name", func(t *testing.T) {
		if event.Name != "Test Product" {
			t.Errorf("Name = %v, want %v", event.Name, "Test Product")
		}
	})

	t.Run("Sku", func(t *testing.T) {
		if event.Sku != 999 {
			t.Errorf("Sku = %v, want %v", event.Sku, 999)
		}
	})

	t.Run("Categories", func(t *testing.T) {
		if len(event.Categories) != 3 {
			t.Errorf("Categories length = %d, want 3", len(event.Categories))
		}
	})

	t.Run("Price", func(t *testing.T) {
		if event.Price != 500 {
			t.Errorf("Price = %v, want %v", event.Price, 500)
		}
	})
}

// Benchmark
func BenchmarkNewProductCreatedEvent(b *testing.B) {
	categories := []string{"Electronics", "Computers"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewProductCreatedEvent("Notebook", 12345, categories, 3500)
	}
}

func BenchmarkProductCreatedEvent_EventName(b *testing.B) {
	event := NewProductCreatedEvent("Test", 123, []string{"Test"}, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = event.EventName()
	}
}
