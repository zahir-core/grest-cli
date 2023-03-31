package app

import (
	"errors"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"grest.dev/grest"
)

func DB() DBInterface {
	if db == nil {
		db = &dbUtil{}
		db.configure()
	}
	return db
}

type DBInterface interface {
	RegisterConn(connName string, conn *gorm.DB)
	Conn(connName string) (*gorm.DB, error)
	Close()
	RegisterTable(connName string, t grest.Table) error
	MigrateTable(tx *gorm.DB, connName string, mTable grest.MigrationTable) error
	RegisterSeeder(connName, seederKey string, seederHandler grest.SeederHandler) error
	RunSeeder(tx *gorm.DB, connName string, seedTable grest.SeederTable) error
	Connect(connName string, c grest.DBConfig) error
	IsNotFoundError(err error) bool
}

var db *dbUtil

type dbUtil struct {
	grest.DB
}

func (d *dbUtil) configure() *dbUtil {
	c := grest.DBConfig{}
	c.Driver = DB_DRIVER
	c.Host = DB_HOST
	c.Port = DB_PORT
	c.User = DB_USERNAME
	c.Password = DB_PASSWORD
	c.DbName = DB_DATABASE
	err := d.Connect("main", c)
	if err != nil {
		Logger().Fatal().
			Err(err).
			Str("driver", c.Driver).
			Str("host", c.Host).
			Int("port", c.Port).
			Str("user", c.User).
			Str("password", c.Password).
			Str("db_name", c.DbName).
			Msg("Failed to connect to main DB")
	}
	return d
}

func (d *dbUtil) Connect(connName string, c grest.DBConfig) error {
	dialector := postgres.Open(c.DSN())
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	if DB_IS_DEBUG {
		gormDB = gormDB.Debug()
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(DB_MAX_OPEN_CONNS)
	sqlDB.SetMaxIdleConns(DB_MAX_OPEN_CONNS)
	sqlDB.SetConnMaxLifetime(DB_CONN_MAX_LIFETIME)

	d.RegisterConn(connName, gormDB)
	d.setupReplicas(gormDB, c)
	return nil
}

// Automatic read and write connection switching
func (d *dbUtil) setupReplicas(db *gorm.DB, c grest.DBConfig) {
	if DB_HOST_READ != "" {
		dialector := postgres.Open(c.DSN())
		sourcesDialector := []gorm.Dialector{dialector}
		replicasDialector := []gorm.Dialector{}
		replicas := strings.Split(DB_HOST_READ, ",")
		for _, replica := range replicas {
			c.Host = replica
			dialector := postgres.Open(c.DSN())
			replicasDialector = append(replicasDialector, dialector)
		}
		if len(replicasDialector) == 0 {
			replicasDialector = sourcesDialector
		}
		db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  sourcesDialector,
			Replicas: replicasDialector,
			Policy:   dbresolver.RandomPolicy{},
		}))
	}
}

func (*dbUtil) IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
