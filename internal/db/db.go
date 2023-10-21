package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/synctv-org/synctv/internal/conf"
	"github.com/synctv-org/synctv/internal/model"
	_ "github.com/synctv-org/synctv/utils/fastJSONSerializer"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(d *gorm.DB) error {
	db = d
	return AutoMigrate(new(model.Movie), new(model.Room), new(model.User), new(model.RoomUserRelation), new(model.UserProvider))
}

func AutoMigrate(dst ...any) error {
	var err error
	switch conf.Conf.Database.Type {
	case conf.DatabaseTypeMysql:
		err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(dst...)
	case conf.DatabaseTypeSqlite3, conf.DatabaseTypePostgres:
		err = db.AutoMigrate(dst...)
	default:
		log.Fatalf("unknown database type: %s", conf.Conf.Database.Type)
	}
	return err
}

func DB() *gorm.DB {
	return db
}

func Close() {
	log.Info("closing db")
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("failed to get db: %s", err.Error())
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Errorf("failed to close db: %s", err.Error())
		return
	}
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func OrderByAsc(column string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(column + " asc")
	}
}

func OrderByDesc(column string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(column + " desc")
	}
}

func OrderByCreatedAtAsc(db *gorm.DB) *gorm.DB {
	return db.Order("created_at asc")
}

func OrderByCreatedAtDesc(db *gorm.DB) *gorm.DB {
	return db.Order("created_at desc")
}

func OrderByIDAsc(db *gorm.DB) *gorm.DB {
	return db.Order("id asc")
}

func OrderByIDDesc(db *gorm.DB) *gorm.DB {
	return db.Order("id desc")
}
