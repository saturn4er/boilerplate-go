package example

import (
	"context"
	"fmt"
	"log"

	"github.com/go-pnp/go-pnp/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/saturn4er/boilerplate-go/example/auth/authservice"
	"github.com/saturn4er/boilerplate-go/example/auth/authstorage"
	"github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	"github.com/saturn4er/boilerplate-go/example/segmentation/segmentationstorage"
	"github.com/saturn4er/boilerplate-go/lib/txoutbox"
)

func Run() {
	dsn := "host=localhost user=postgres password=postgres dbname=example port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	var logger *logging.Logger // nil is safe — Logger methods are nil-receiver tolerant

	var processors []txoutbox.MessageProcessor

	// Initialize storages
	authStorage := authstorage.NewStorages(db, logger, processors)
	_ = segmentationstorage.NewStorages(db, logger, processors)

	// Initialize services
	authSvc := &authservice.Service{Storage: authStorage}
	_ = &segmentationservice.Service{} // would be wired to a message consumer

	// Demonstrate registration flow: creates a user and sends a SetUserTagCommand
	ctx := context.Background()

	user, err := authSvc.Register(ctx, "alice@example.com", "Alice", authservice.RoleUser)
	if err != nil {
		log.Fatalf("register user: %v", err)
	}

	fmt.Printf("Registered user: %s (ID: %s)\n", user.Name, user.ID)
}
