package outbound

import (
	"database/sql"

	"github.com/LeonLow97/internal/pkg/kafkago"
	"github.com/LeonLow97/internal/ports"
)

type repository struct {
	db                *sql.DB
	segmentioInstance *kafkago.Segmentio
}

func NewRepository(db *sql.DB, segmentioInstance *kafkago.Segmentio) ports.Repository {
	return &repository{
		db:                db,
		segmentioInstance: segmentioInstance,
	}
}
