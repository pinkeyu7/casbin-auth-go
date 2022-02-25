package driver

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

var (
	orm     *xorm.EngineGroup
	ormOnce sync.Once
)

// NewXorm return singleton xorm instance
func NewXorm() (*xorm.EngineGroup, error) {
	var err error
	ormOnce.Do(func() {
		err = newXorm()
	})
	return orm, err
}

func newXorm() error {
	var err error
	dsn := "%s:%s@(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	master, err := xorm.NewEngine("mysql", fmt.Sprintf(dsn, dbUser, dbPassword, dbHost, dbPort, dbName))

	if err != nil {
		log.WithFields(log.Fields{}).Error("[MySQL] Connect to MySQL Master error...", err)
		return err
	} else {
		log.WithFields(log.Fields{}).Info("[MySQL] Connected to MySQL Master")
	}

	dbSlave1User := os.Getenv("DB_SLAVE1_USER")
	dbSlave1Password := os.Getenv("DB_SLAVE1_PASSWORD")
	dbSlave1Host := os.Getenv("DB_SLAVE1_HOST")
	dbSlave1Port := os.Getenv("DB_SLAVE1_PORT")
	dbSlave1Name := os.Getenv("DB_SLAVE1_NAME")
	slave1, err := xorm.NewEngine("mysql", fmt.Sprintf(dsn, dbSlave1User, dbSlave1Password, dbSlave1Host, dbSlave1Port, dbSlave1Name))

	if err != nil {
		log.WithFields(log.Fields{}).Error("[MySQL] Connect to MySQL Slave1 error...", err)
		return err
	} else {
		log.WithFields(log.Fields{}).Info("[MySQL] Connected to MySQL Slave1")
	}

	master.TZLocation, _ = time.LoadLocation("UTC")
	master.DatabaseTZ, _ = time.LoadLocation("UTC")
	slave1.TZLocation, _ = time.LoadLocation("UTC")
	slave1.DatabaseTZ, _ = time.LoadLocation("UTC")

	slaves := []*xorm.Engine{slave1}
	orm, err = xorm.NewEngineGroup(master, slaves, xorm.LeastConnPolicy())

	if err != nil {
		log.Println(err)
		return err
	}

	if os.Getenv("XORM_MODE") == "debug" {
		log.WithFields(log.Fields{}).Info("[MySQL] XORM debug mode enabled")
		orm.ShowSQL(true)
		orm.Logger().SetLevel(xlog.LOG_DEBUG)
		pingErr := orm.Ping()
		if pingErr != nil {
			log.WithFields(log.Fields{}).Error(pingErr.Error())
		}
	}

	return err
}
