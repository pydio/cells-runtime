package sql

import (
	"context"

	"gorm.io/gorm"

	"github.com/pydio/cells-runtime/errors"
)

// NewAbstract initializes a new Abstract to be used by DAO implementation
func NewAbstract(db *gorm.DB) *Abstract {
	return &Abstract{DB: db}
}

type Abstract struct {
	*gorm.DB
	migrateModels func() []any
}

// WithModels sets the migrateModels factory for auto migration
func (a *Abstract) WithModels(factory func() []any) *Abstract {
	a.migrateModels = factory
	return a
}

// Session creates a gorm session with given context and default config (SkipDefaultTransaction: true)
func (a *Abstract) Session(ctx context.Context) *gorm.DB {
	return a.DB.Session(&gorm.Session{SkipDefaultTransaction: true}).WithContext(ctx)
}

func (a *Abstract) Migrate(ctx context.Context) error {
	if a.migrateModels == nil {
		return errors.New("no models factoring given, implement your own Migrate")
	}
	return a.Session(ctx).AutoMigrate(a.migrateModels()...)
}

// MigrateWithCtx triggers AutoMigrate in context
func (a *Abstract) MigrateWithCtx(ctx context.Context, dst ...any) error {
	return a.Session(ctx).AutoMigrate(dst...)
}
