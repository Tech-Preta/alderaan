package product_entity

import (
	"testing"

	shared_events "github.com/williamkoller/golang-domain-driven-design/internal/shared/domain/events"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name           string
		productName    string
		sku            int
		categories     []string
		price          int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:        "valid product",
			productName: "Notebook",
			sku:         12345,
			categories:  []string{"Electronics", "Computers"},
			price:       3500,
			wantErr:     false,
		},
		{
			name:           "empty name",
			productName:    "",
			sku:            12345,
			categories:     []string{"Electronics"},
			price:          3500,
			wantErr:        true,
			expectedErrMsg: "name is required",
		},
		{
			name:           "zero sku",
			productName:    "Notebook",
			sku:            0,
			categories:     []string{"Electronics"},
			price:          3500,
			wantErr:        true,
			expectedErrMsg: "sku is required",
		},
		{
			name:           "negative sku",
			productName:    "Notebook",
			sku:            -1,
			categories:     []string{"Electronics"},
			price:          3500,
			wantErr:        true,
			expectedErrMsg: "sku is required",
		},
		{
			name:           "empty categories",
			productName:    "Notebook",
			sku:            12345,
			categories:     []string{},
			price:          3500,
			wantErr:        true,
			expectedErrMsg: "categories is required",
		},
		{
			name:           "nil categories",
			productName:    "Notebook",
			sku:            12345,
			categories:     nil,
			price:          3500,
			wantErr:        true,
			expectedErrMsg: "categories is required",
		},
		{
			name:           "zero price",
			productName:    "Notebook",
			sku:            12345,
			categories:     []string{"Electronics"},
			price:          0,
			wantErr:        true,
			expectedErrMsg: "price is required",
		},
		{
			name:           "negative price",
			productName:    "Notebook",
			sku:            12345,
			categories:     []string{"Electronics"},
			price:          -100,
			wantErr:        true,
			expectedErrMsg: "price is required",
		},
		{
			name:        "multiple categories",
			productName: "Gaming Mouse",
			sku:         99999,
			categories:  []string{"Electronics", "Gaming", "Accessories"},
			price:       299,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatcher := shared_events.NewEventDispatcher()

			product, event, err := NewProduct(tt.productName, tt.sku, tt.categories, tt.price, dispatcher)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewProduct() expected error, got nil")
					return
				}
				if err.Error() != tt.expectedErrMsg {
					t.Errorf("NewProduct() error = %v, want %v", err.Error(), tt.expectedErrMsg)
				}
				if product != nil {
					t.Errorf("NewProduct() expected nil product, got %v", product)
				}
				if event != nil {
					t.Errorf("NewProduct() expected nil event, got %v", event)
				}
			} else {
				if err != nil {
					t.Errorf("NewProduct() unexpected error = %v", err)
					return
				}
				if product == nil {
					t.Error("NewProduct() expected product, got nil")
					return
				}
				if event == nil {
					t.Error("NewProduct() expected event, got nil")
					return
				}

				// Verificar campos do produto
				if product.GetName() != tt.productName {
					t.Errorf("Product.Name = %v, want %v", product.GetName(), tt.productName)
				}
				if product.GetSku() != tt.sku {
					t.Errorf("Product.Sku = %v, want %v", product.GetSku(), tt.sku)
				}
				if product.GetPrice() != tt.price {
					t.Errorf("Product.Price = %v, want %v", product.GetPrice(), tt.price)
				}
				if len(product.GetCategories()) != len(tt.categories) {
					t.Errorf("Product.Categories length = %v, want %v", len(product.GetCategories()), len(tt.categories))
				}
			}
		})
	}
}

func TestNewProduct_WithNilDispatcher(t *testing.T) {
	product, event, err := NewProduct("Test Product", 123, []string{"Test"}, 100, nil)

	if err != nil {
		t.Errorf("NewProduct() with nil dispatcher unexpected error = %v", err)
	}
	if product == nil {
		t.Error("NewProduct() with nil dispatcher expected product, got nil")
	}
	if event == nil {
		t.Error("NewProduct() with nil dispatcher expected event, got nil")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		inputName  string
		sku        int
		categories []string
		price      int
		wantValid  bool
		wantErr    string
	}{
		{
			name:       "all valid",
			inputName:  "Product",
			sku:        100,
			categories: []string{"Category"},
			price:      50,
			wantValid:  true,
			wantErr:    "",
		},
		{
			name:       "invalid name",
			inputName:  "",
			sku:        100,
			categories: []string{"Category"},
			price:      50,
			wantValid:  false,
			wantErr:    "name is required",
		},
		{
			name:       "invalid sku",
			inputName:  "Product",
			sku:        0,
			categories: []string{"Category"},
			price:      50,
			wantValid:  false,
			wantErr:    "sku is required",
		},
		{
			name:       "invalid categories",
			inputName:  "Product",
			sku:        100,
			categories: []string{},
			price:      50,
			wantValid:  false,
			wantErr:    "categories is required",
		},
		{
			name:       "invalid price",
			inputName:  "Product",
			sku:        100,
			categories: []string{"Category"},
			price:      0,
			wantValid:  false,
			wantErr:    "price is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := Validate(tt.inputName, tt.sku, tt.categories, tt.price)

			if valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v", valid, tt.wantValid)
			}

			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("Validate() expected error %v, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestProduct_Getters(t *testing.T) {
	product := &Product{
		Name:       "Test Product",
		Sku:        12345,
		Categories: []string{"Cat1", "Cat2"},
		Price:      999,
	}

	t.Run("GetName", func(t *testing.T) {
		if got := product.GetName(); got != "Test Product" {
			t.Errorf("GetName() = %v, want %v", got, "Test Product")
		}
	})

	t.Run("GetSku", func(t *testing.T) {
		if got := product.GetSku(); got != 12345 {
			t.Errorf("GetSku() = %v, want %v", got, 12345)
		}
	})

	t.Run("GetCategories", func(t *testing.T) {
		got := product.GetCategories()
		if len(got) != 2 {
			t.Errorf("GetCategories() length = %v, want %v", len(got), 2)
		}
		if got[0] != "Cat1" || got[1] != "Cat2" {
			t.Errorf("GetCategories() = %v, want [Cat1 Cat2]", got)
		}
	})

	t.Run("GetPrice", func(t *testing.T) {
		if got := product.GetPrice(); got != 999 {
			t.Errorf("GetPrice() = %v, want %v", got, 999)
		}
	})
}

// Benchmark tests
func BenchmarkNewProduct(b *testing.B) {
	dispatcher := shared_events.NewEventDispatcher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = NewProduct("Test Product", 123, []string{"Category"}, 100, dispatcher)
	}
}

func BenchmarkValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Validate("Product", 100, []string{"Category"}, 50)
	}
}
