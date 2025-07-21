package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"register/pkg/logger"
)

type QueryBuilder struct {
	db *gorm.DB
}

func NewQueryBuilder(databaseURL string, logger logger.Logger) *QueryBuilder {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})

	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	return &QueryBuilder{db: db}
}

func (qb *QueryBuilder) AutoMigrate(dst ...interface{}) error {
	return qb.db.AutoMigrate(dst...)
}

func (qb *QueryBuilder) Create(value interface{}) *QueryBuilder {
	return &QueryBuilder{db: qb.db.Create(value)}
}

func (qb *QueryBuilder) First(dest interface{}) *QueryBuilder {
	return &QueryBuilder{db: qb.db.First(dest)}
}

func (qb *QueryBuilder) Find(dest interface{}) *QueryBuilder {
	return &QueryBuilder{db: qb.db.Find(dest)}
}

func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{db: qb.db.Where(query, args...)}
}

func (qb *QueryBuilder) Delete(value interface{}) *QueryBuilder {
	return &QueryBuilder{db: qb.db.Delete(value)}
}

func (qb *QueryBuilder) Model(value interface{}) *QueryBuilder {
	return &QueryBuilder{db: qb.db.Model(value)}
}

func (qb *QueryBuilder) Update(column string, value interface{}) *QueryBuilder {
	return &QueryBuilder{db: qb.db.Update(column, value)}
}

func (qb *QueryBuilder) Updates(values interface{}) *QueryBuilder {
	return &QueryBuilder{db: qb.db.Updates(values)}
}

func (qb *QueryBuilder) Error() error {
	return qb.db.Error
}

func (qb *QueryBuilder) RowsAffected() int64 {
	return qb.db.RowsAffected
}
