package handler

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"inventoryService/model"
	pb "inventoryService/proto/inventory"
	"inventoryService/service"
	"log"
)

type InventoryHandler struct {
	pb.UnimplementedInventoryServiceServer
	productService       *service.ProductService
	categoryService      *service.CategoryService
	warehouseService     *service.WarehouseService
	inventoryItemService *service.InventoryItemService
	stockMovementService *service.StockMovementService
	supplierService      *service.SupplierService
}

func NewInventoryHandler(
	productService *service.ProductService,
	categoryService *service.CategoryService,
	warehouseService *service.WarehouseService,
	inventoryItemService *service.InventoryItemService,
	stockMovementService *service.StockMovementService,
	supplierService *service.SupplierService,
) *InventoryHandler {
	return &InventoryHandler{
		productService:       productService,
		categoryService:      categoryService,
		warehouseService:     warehouseService,
		inventoryItemService: inventoryItemService,
		stockMovementService: stockMovementService,
		supplierService:      supplierService,
	}
}

func (h *InventoryHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	internalProduct := convertPbToProductModel(req.Product)
	createdProduct, err := h.productService.CreateProduct(ctx, internalProduct)
	if err != nil {
		log.Printf("Error in CreateProduct: %v", err)
		return nil, err
	}
	return convertProductModelToPb(createdProduct), nil
}

func (h *InventoryHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	product, err := h.productService.GetProduct(ctx, id)
	if err != nil {
		log.Printf("Error in GetProduct: %v", err)
		return nil, err
	}
	return convertProductModelToPb(product), nil
}

func (h *InventoryHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	internalProduct := convertPbToProductModel(req.Product)
	err := h.productService.UpdateProduct(ctx, internalProduct)
	if err != nil {
		log.Printf("Error in UpdateProduct: %v", err)
		return nil, err
	}
	return convertProductModelToPb(internalProduct), nil
}

func (h *InventoryHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	err = h.productService.DeleteProduct(ctx, id)
	if err != nil {
		log.Printf("Error in DeleteProduct: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *InventoryHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.productService.ListProducts(ctx)
	if err != nil {
		log.Printf("Error in ListProducts: %v", err)
		return nil, err
	}

	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = convertProductModelToPb(product)
	}

	return &pb.ListProductsResponse{Products: pbProducts}, nil
}

func convertPbToProductModel(pbProduct *pb.Product) *model.Product {
	id, _ := uuid.Parse(pbProduct.Id)
	categoryID, _ := uuid.Parse(pbProduct.CategoryId)

	return &model.Product{
		ID:          id,
		Name:        pbProduct.Name,
		Description: pbProduct.Description,
		CategoryID:  categoryID,
		Price:       pbProduct.Price,
		SKU:         pbProduct.Sku,
	}
}

func convertProductModelToPb(product *model.Product) *pb.Product {
	return &pb.Product{
		Id:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		CategoryId:  product.CategoryID.String(),
		Price:       product.Price,
		Sku:         product.SKU,
	}
}

func (h *InventoryHandler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	internalCategory := convertPbToCategoryModel(req.Category)
	createdCategory, err := h.categoryService.CreateCategory(ctx, internalCategory)
	if err != nil {
		log.Printf("Error in CreateCategory: %v", err)
		return nil, err
	}
	return convertCategoryModelToPb(createdCategory), nil
}

func (h *InventoryHandler) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	category, err := h.categoryService.GetCategory(ctx, id)
	if err != nil {
		log.Printf("Error in GetCategory: %v", err)
		return nil, err
	}
	return convertCategoryModelToPb(category), nil
}

func (h *InventoryHandler) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
	internalCategory := convertPbToCategoryModel(req.Category)
	err := h.categoryService.UpdateCategory(ctx, internalCategory)
	if err != nil {
		log.Printf("Error in UpdateCategory: %v", err)
		return nil, err
	}
	return convertCategoryModelToPb(internalCategory), nil
}

func (h *InventoryHandler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	err = h.categoryService.DeleteCategory(ctx, id)
	if err != nil {
		log.Printf("Error in DeleteCategory: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *InventoryHandler) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	categories, err := h.categoryService.ListCategories(ctx)
	if err != nil {
		log.Printf("Error in ListCategories: %v", err)
		return nil, err
	}

	pbCategories := make([]*pb.Category, len(categories))
	for i, category := range categories {
		pbCategories[i] = convertCategoryModelToPb(category)
	}

	return &pb.ListCategoriesResponse{Categories: pbCategories}, nil
}

func convertPbToCategoryModel(pbCategory *pb.Category) *model.Category {
	id, _ := uuid.Parse(pbCategory.Id)
	return &model.Category{
		ID:          id,
		Name:        pbCategory.Name,
		Description: pbCategory.Description,
	}
}

func convertCategoryModelToPb(category *model.Category) *pb.Category {
	return &pb.Category{
		Id:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
	}
}

func (h *InventoryHandler) CreateInventoryItem(ctx context.Context, req *pb.CreateInventoryItemRequest) (*pb.InventoryItem, error) {
	internalItem := convertPbToInventoryItemModel(req.InventoryItem)
	createdItem, err := h.inventoryItemService.CreateInventoryItem(ctx, internalItem)
	if err != nil {
		log.Printf("Error in CreateInventoryItem: %v", err)
		return nil, err
	}
	return convertInventoryItemModelToPb(createdItem), nil
}

func (h *InventoryHandler) GetInventoryItem(ctx context.Context, req *pb.GetInventoryItemRequest) (*pb.InventoryItem, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	item, err := h.inventoryItemService.GetInventoryItem(ctx, id)
	if err != nil {
		log.Printf("Error in GetInventoryItem: %v", err)
		return nil, err
	}
	return convertInventoryItemModelToPb(item), nil
}

func (h *InventoryHandler) UpdateInventoryItem(ctx context.Context, req *pb.UpdateInventoryItemRequest) (*pb.InventoryItem, error) {
	internalItem := convertPbToInventoryItemModel(req.InventoryItem)
	err := h.inventoryItemService.UpdateInventoryItem(ctx, internalItem)
	if err != nil {
		log.Printf("Error in UpdateInventoryItem: %v", err)
		return nil, err
	}
	return convertInventoryItemModelToPb(internalItem), nil
}

func (h *InventoryHandler) DeleteInventoryItem(ctx context.Context, req *pb.DeleteInventoryItemRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	err = h.inventoryItemService.DeleteInventoryItem(ctx, id)
	if err != nil {
		log.Printf("Error in DeleteInventoryItem: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *InventoryHandler) ListInventoryItems(ctx context.Context, req *pb.ListInventoryItemsRequest) (*pb.ListInventoryItemsResponse, error) {
	items, err := h.inventoryItemService.ListInventoryItems(ctx)
	if err != nil {
		log.Printf("Error in ListInventoryItems: %v", err)
		return nil, err
	}

	pbItems := make([]*pb.InventoryItem, len(items))
	for i, item := range items {
		pbItems[i] = convertInventoryItemModelToPb(item)
	}

	return &pb.ListInventoryItemsResponse{InventoryItems: pbItems}, nil
}

func convertPbToInventoryItemModel(pbItem *pb.InventoryItem) *model.InventoryItem {
	id, _ := uuid.Parse(pbItem.Id)
	productID, _ := uuid.Parse(pbItem.ProductId)
	warehouseID, _ := uuid.Parse(pbItem.WarehouseId)

	return &model.InventoryItem{
		ID:              id,
		ProductID:       productID,
		WarehouseID:     warehouseID,
		Quantity:        int(pbItem.Quantity),
		ReorderLevel:    int(pbItem.ReorderLevel),
		ReorderQuantity: int(pbItem.ReorderQuantity),
	}
}

func convertInventoryItemModelToPb(item *model.InventoryItem) *pb.InventoryItem {
	return &pb.InventoryItem{
		Id:              item.ID.String(),
		ProductId:       item.ProductID.String(),
		WarehouseId:     item.WarehouseID.String(),
		Quantity:        int32(item.Quantity),
		ReorderLevel:    int32(item.ReorderLevel),
		ReorderQuantity: int32(item.ReorderQuantity),
	}
}

func (h *InventoryHandler) CreateStockMovement(ctx context.Context, req *pb.CreateStockMovementRequest) (*pb.StockMovement, error) {
	internalMovement := convertPbToStockMovementModel(req.StockMovement)
	createdMovement, err := h.stockMovementService.CreateStockMovement(ctx, internalMovement)
	if err != nil {
		log.Printf("Error in CreateStockMovement: %v", err)
		return nil, err
	}
	return convertStockMovementModelToPb(createdMovement), nil
}

func (h *InventoryHandler) GetStockMovement(ctx context.Context, req *pb.GetStockMovementRequest) (*pb.StockMovement, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	movement, err := h.stockMovementService.GetStockMovement(ctx, id)
	if err != nil {
		log.Printf("Error in GetStockMovement: %v", err)
		return nil, err
	}
	return convertStockMovementModelToPb(movement), nil
}

func (h *InventoryHandler) UpdateStockMovement(ctx context.Context, req *pb.UpdateStockMovementRequest) (*pb.StockMovement, error) {
	internalMovement := convertPbToStockMovementModel(req.StockMovement)
	err := h.stockMovementService.UpdateStockMovement(ctx, internalMovement)
	if err != nil {
		log.Printf("Error in UpdateStockMovement: %v", err)
		return nil, err
	}
	return convertStockMovementModelToPb(internalMovement), nil
}

func (h *InventoryHandler) DeleteStockMovement(ctx context.Context, req *pb.DeleteStockMovementRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	err = h.stockMovementService.DeleteStockMovement(ctx, id)
	if err != nil {
		log.Printf("Error in DeleteStockMovement: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *InventoryHandler) ListStockMovements(ctx context.Context, req *pb.ListStockMovementsRequest) (*pb.ListStockMovementsResponse, error) {
	movements, err := h.stockMovementService.ListStockMovements(ctx)
	if err != nil {
		log.Printf("Error in ListStockMovements: %v", err)
		return nil, err
	}

	pbMovements := make([]*pb.StockMovement, len(movements))
	for i, movement := range movements {
		pbMovements[i] = convertStockMovementModelToPb(movement)
	}

	return &pb.ListStockMovementsResponse{StockMovements: pbMovements}, nil
}

func convertPbToStockMovementModel(pbMovement *pb.StockMovement) *model.StockMovement {
	id, _ := uuid.Parse(pbMovement.Id)
	inventoryItemID, _ := uuid.Parse(pbMovement.InventoryItemId)
	sourceWarehouseID, _ := uuid.Parse(pbMovement.SourceWarehouseId)
	destinationWarehouseID, _ := uuid.Parse(pbMovement.DestinationWarehouseId)
	return &model.StockMovement{
		ID:                     id,
		InventoryItemID:        inventoryItemID,
		Type:                   model.StockMovementType(pbMovement.Type),
		Quantity:               int(pbMovement.Quantity),
		Date:                   pbMovement.Date.AsTime(),
		SourceWarehouseID:      sourceWarehouseID,
		DestinationWarehouseID: destinationWarehouseID,
	}
}

func convertStockMovementModelToPb(movement *model.StockMovement) *pb.StockMovement {
	return &pb.StockMovement{
		Id:                     movement.ID.String(),
		InventoryItemId:        movement.InventoryItemID.String(),
		Type:                   pb.StockMovementType(movement.Type),
		Quantity:               int32(movement.Quantity),
		Date:                   timestamppb.New(movement.Date),
		SourceWarehouseId:      movement.SourceWarehouseID.String(),
		DestinationWarehouseId: movement.DestinationWarehouseID.String(),
	}
}

func (h *InventoryHandler) CreateSupplier(ctx context.Context, req *pb.CreateSupplierRequest) (*pb.Supplier, error) {
	internalSupplier := convertPbToSupplierModel(req.Supplier)
	createdSupplier, err := h.supplierService.CreateSupplier(ctx, internalSupplier)
	if err != nil {
		log.Printf("Error in CreateSupplier: %v", err)
		return nil, err
	}
	return convertSupplierModelToPb(createdSupplier), nil
}

func (h *InventoryHandler) GetSupplier(ctx context.Context, req *pb.GetSupplierRequest) (*pb.Supplier, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	supplier, err := h.supplierService.GetSupplier(ctx, id)
	if err != nil {
		log.Printf("Error in GetSupplier: %v", err)
		return nil, err
	}
	return convertSupplierModelToPb(supplier), nil
}

func (h *InventoryHandler) UpdateSupplier(ctx context.Context, req *pb.UpdateSupplierRequest) (*pb.Supplier, error) {
	internalSupplier := convertPbToSupplierModel(req.Supplier)
	err := h.supplierService.UpdateSupplier(ctx, internalSupplier)
	if err != nil {
		log.Printf("Error in UpdateSupplier: %v", err)
		return nil, err
	}
	return convertSupplierModelToPb(internalSupplier), nil
}

func (h *InventoryHandler) DeleteSupplier(ctx context.Context, req *pb.DeleteSupplierRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	err = h.supplierService.DeleteSupplier(ctx, id)
	if err != nil {
		log.Printf("Error in DeleteSupplier: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *InventoryHandler) ListSuppliers(ctx context.Context, req *pb.ListSuppliersRequest) (*pb.ListSuppliersResponse, error) {
	suppliers, err := h.supplierService.ListSuppliers(ctx)
	if err != nil {
		log.Printf("Error in ListSuppliers: %v", err)
		return nil, err
	}

	pbSuppliers := make([]*pb.Supplier, len(suppliers))
	for i, supplier := range suppliers {
		pbSuppliers[i] = convertSupplierModelToPb(supplier)
	}

	return &pb.ListSuppliersResponse{Suppliers: pbSuppliers}, nil
}

func convertPbToSupplierModel(pbSupplier *pb.Supplier) *model.Supplier {
	id, _ := uuid.Parse(pbSupplier.Id)
	return &model.Supplier{
		ID:          id,
		Name:        pbSupplier.Name,
		ContactInfo: pbSupplier.ContactInfo,
	}
}

func convertSupplierModelToPb(supplier *model.Supplier) *pb.Supplier {
	return &pb.Supplier{
		Id:          supplier.ID.String(),
		Name:        supplier.Name,
		ContactInfo: supplier.ContactInfo,
	}
}

func (h *InventoryHandler) CreateWarehouse(ctx context.Context, req *pb.CreateWarehouseRequest) (*pb.Warehouse, error) {
	internalWarehouse := convertPbToWarehouseModel(req.Warehouse)
	createdWarehouse, err := h.warehouseService.CreateWarehouse(ctx, internalWarehouse)
	if err != nil {
		log.Printf("Error in CreateWarehouse: %v", err)
		return nil, err
	}
	return convertWarehouseModelToPb(createdWarehouse), nil
}

func (h *InventoryHandler) GetWarehouse(ctx context.Context, req *pb.GetWarehouseRequest) (*pb.Warehouse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	warehouse, err := h.warehouseService.GetWarehouse(ctx, id)
	if err != nil {
		log.Printf("Error in GetWarehouse: %v", err)
		return nil, err
	}
	return convertWarehouseModelToPb(warehouse), nil
}

func (h *InventoryHandler) UpdateWarehouse(ctx context.Context, req *pb.UpdateWarehouseRequest) (*pb.Warehouse, error) {
	internalWarehouse := convertPbToWarehouseModel(req.Warehouse)
	err := h.warehouseService.UpdateWarehouse(ctx, internalWarehouse)
	if err != nil {
		log.Printf("Error in UpdateWarehouse: %v", err)
		return nil, err
	}
	return convertWarehouseModelToPb(internalWarehouse), nil
}

func (h *InventoryHandler) DeleteWarehouse(ctx context.Context, req *pb.DeleteWarehouseRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	err = h.warehouseService.DeleteWarehouse(ctx, id)
	if err != nil {
		log.Printf("Error in DeleteWarehouse: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *InventoryHandler) ListWarehouses(ctx context.Context, req *pb.ListWarehousesRequest) (*pb.ListWarehousesResponse, error) {
	warehouses, err := h.warehouseService.ListWarehouses(ctx)
	if err != nil {
		log.Printf("Error in ListWarehouses: %v", err)
		return nil, err
	}

	pbWarehouses := make([]*pb.Warehouse, len(warehouses))
	for i, warehouse := range warehouses {
		pbWarehouses[i] = convertWarehouseModelToPb(warehouse)
	}

	return &pb.ListWarehousesResponse{Warehouses: pbWarehouses}, nil
}

func convertPbToWarehouseModel(pbWarehouse *pb.Warehouse) *model.Warehouse {
	id, _ := uuid.Parse(pbWarehouse.Id)
	return &model.Warehouse{
		ID:       id,
		Name:     pbWarehouse.Name,
		Location: pbWarehouse.Location,
	}
}

func convertWarehouseModelToPb(warehouse *model.Warehouse) *pb.Warehouse {
	return &pb.Warehouse{
		Id:       warehouse.ID.String(),
		Name:     warehouse.Name,
		Location: warehouse.Location,
	}
}
