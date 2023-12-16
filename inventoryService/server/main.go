package main

import (
	"github.com/gocql/gocql"
	"google.golang.org/grpc"
	"inventoryService/handler"
	pb "inventoryService/proto/inventory"
	"inventoryService/repository"
	"inventoryService/service"
	"log"
	"net"
)

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "inventory"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	cqlStatement := `CREATE INDEX IF NOT EXISTS ON categories (name)`
	err = session.Query(cqlStatement).Exec()
	if err != nil {
		log.Fatalf("Failed to create index on 'name': %v", err)
	}

	cqlStatement = `CREATE INDEX IF NOT EXISTS ON products (sku)`
	err = session.Query(cqlStatement).Exec()
	if err != nil {
		log.Fatalf("Failed to create index on 'sku': %v", err)
	}

	cqlStatement = `CREATE INDEX IF NOT EXISTS ON products (name)`
	err = session.Query(cqlStatement).Exec()
	if err != nil {
		log.Fatalf("Failed to create index on 'name': %v", err)
	}

	cqlStatement = `CREATE INDEX IF NOT EXISTS ON warehouses (name)`
	err = session.Query(cqlStatement).Exec()
	if err != nil {
		log.Fatalf("Failed to create index on 'name': %v", err)
	}

	cqlStatement = `CREATE INDEX IF NOT EXISTS ON suppliers (name)`
	err = session.Query(cqlStatement).Exec()
	if err != nil {
		log.Fatalf("Failed to create index on 'name': %v", err)
	}

	productRepo := repository.NewProductRepository(session)
	categoryRepo := repository.NewCategoryRepository(session)
	warehouseRepo := repository.NewWarehouseRepository(session)
	inventoryItemRepo := repository.NewInventoryItemRepository(session)
	stockMovementRepo := repository.NewStockMovementRepository(session)
	supplierRepo := repository.NewSupplierRepository(session)

	productService := service.NewProductService(productRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	inventoryItemService := service.NewInventoryItemService(inventoryItemRepo, productRepo, warehouseRepo)
	stockMovementService := service.NewStockMovementService(stockMovementRepo)
	supplierService := service.NewSupplierService(supplierRepo)

	inventoryHandler := handler.NewInventoryHandler(
		productService, categoryService, warehouseService,
		inventoryItemService, stockMovementService, supplierService,
	)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, inventoryHandler)

	log.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
