package outbound

import (
	"github.com/LeonLow97/internal/ports"
	"github.com/LeonLow97/internal/pkg/kafkago"
	pb "github.com/LeonLow97/proto"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type Repository struct {
	db                *sqlx.DB
	grpcConn          pb.InventoryServiceClient
	segmentioInstance *kafkago.Segmentio
}

func NewRepository(db *sqlx.DB, conn *grpc.ClientConn, segmentioInstance *kafkago.Segmentio) ports.Repository {
	return &Repository{
		db:                db,
		grpcConn:          pb.NewInventoryServiceClient(conn),
		segmentioInstance: segmentioInstance,
	}
}
