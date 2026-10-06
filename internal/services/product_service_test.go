package services_test

import (
	"context"
	"testing"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/testhelper"
)

func TestProductService_CRUD_And_Listing(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	ctx := context.Background()

	desc := "A high quality developer tool"

	// 1. Creation success
	created, valErrors, err := productService.Create(ctx, services.CreateProductDTO{
		Name:        "Test Product A",
		Slug:        "test-product-a",
		Description: &desc,
		Price:       49.99,
		Status:      true,
	})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	if len(valErrors) > 0 {
		t.Fatalf("unexpected val errors on creation: %v", valErrors)
	}
	if created.ID == 0 || created.Name != "Test Product A" {
		t.Errorf("invalid created product: %v", created)
	}

	// 2. Creation validation failures
	_, valErrors, _ = productService.Create(ctx, services.CreateProductDTO{
		Name:   "P",
		Slug:   "INVALID SLUG!",
		Price:  -10.0,
		Status: true,
	})
	if valErrors["name"] == "" {
		t.Errorf("expected name error")
	}
	if valErrors["slug"] == "" {
		t.Errorf("expected slug format error")
	}
	if valErrors["price"] == "" {
		t.Errorf("expected price error")
	}

	// 3. Duplicate slug validation
	_, valErrors, _ = productService.Create(ctx, services.CreateProductDTO{
		Name:   "Another Product",
		Slug:   "test-product-a",
		Price:  19.99,
		Status: true,
	})
	if valErrors["slug"] == "" {
		t.Errorf("expected duplicate slug error")
	}

	// 4. Update
	newDesc := "Updated description"
	updated, valErrors, err := productService.Update(ctx, created.ID, services.UpdateProductDTO{
		Name:        "Test Product A Updated",
		Slug:        "test-product-a",
		Description: &newDesc,
		Price:       59.99,
		Status:      false,
	})
	if err != nil {
		t.Fatalf("update error: %v", err)
	}
	if len(valErrors) > 0 {
		t.Fatalf("unexpected update val errors: %v", valErrors)
	}
	if updated.Name != "Test Product A Updated" || updated.Price != 59.99 || updated.Status != false {
		t.Errorf("update values mismatch: %v", updated)
	}

	// Seed more products for search, filter & pagination
	p2, _, _ := productService.Create(ctx, services.CreateProductDTO{
		Name:   "Mobile Phone Deluxe",
		Slug:   "mobile-phone-deluxe",
		Price:  799.00,
		Status: true,
	})
	_, _, _ = productService.Create(ctx, services.CreateProductDTO{
		Name:   "Wireless Charger",
		Slug:   "wireless-charger",
		Price:  25.00,
		Status: true,
	})

	// 5. Search
	searchResults, err := productService.List(ctx, "phone", nil, 1, 10)
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if searchResults.Total != 1 || searchResults.Data[0].ID != p2.ID {
		t.Errorf("expected 1 search result with p2 ID, got %d", searchResults.Total)
	}

	// 6. Filtering by status
	activeStatus := true
	activeProducts, err := productService.List(ctx, "", &activeStatus, 1, 10)
	if err != nil {
		t.Fatalf("filter error: %v", err)
	}
	if activeProducts.Total != 2 {
		t.Errorf("expected 2 active products, got %d", activeProducts.Total)
	}

	inactiveStatus := false
	inactiveProducts, err := productService.List(ctx, "", &inactiveStatus, 1, 10)
	if err != nil {
		t.Fatalf("filter error: %v", err)
	}
	if inactiveProducts.Total != 1 {
		t.Errorf("expected 1 inactive product, got %d", inactiveProducts.Total)
	}

	// 7. Pagination
	paginated, err := productService.List(ctx, "", nil, 1, 2)
	if err != nil {
		t.Fatalf("pagination error: %v", err)
	}
	if len(paginated.Data) != 2 || paginated.Total != 3 || paginated.LastPage != 2 {
		t.Errorf("pagination mismatch: page size %d, total %d, lastPage %d", len(paginated.Data), paginated.Total, paginated.LastPage)
	}

	// 8. Deletion
	if err := productService.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	found, _ := productService.FindByID(ctx, created.ID)
	if found != nil {
		t.Errorf("expected deleted product not to be found")
	}
}
