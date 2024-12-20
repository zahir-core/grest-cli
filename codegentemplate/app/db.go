package app

import (
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"grest.dev/grest"
)

// DB returns a pointer to the dbUtil instance (db).
// If db is not initialized, it creates a new dbUtil instance, configures it, and assigns it to db.
// It ensures that only one instance of dbUtil is created and reused.
func DB() *dbUtil {
	if db == nil {
		db = &dbUtil{}
		db.configure()
	}
	return db
}

// db is a pointer to a dbUtil instance.
// It is used to store and access the singleton instance of dbUtil.
var db *dbUtil

// dbUtil represents a db utility.
// It embeds grest.DB, indicating that dbUtil inherits from grest.DB.
type dbUtil struct {
	grest.DB
}

// configure configures the db utility instance.
// It connect to main db corresponding environment variables.
// You can configure here to connect to multiple db based on connection name if needed.
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
		Logger().Fatal("Failed to connect to main DB",
			slog.Any("err", err),
			slog.String("driver", c.Driver),
			slog.String("host", c.Host),
			slog.Int("port", c.Port),
			slog.String("user", c.User),
			slog.String("password", c.Password),
			slog.String("db_name", c.DbName),
		)
	}
	return d
}

// Connect connect to the db and store to config based on connName key.
func (d *dbUtil) Connect(connName string, c grest.DBConfig) error {
	dbLogLevel := gormlogger.Error
	if DB_IS_DEBUG {
		dbLogLevel = gormlogger.Info
	}
	dialector := postgres.Open(c.DSN())
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  dbLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
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

// setupReplicas setup replica to automatic read and write connection switching.
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
