package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ppreeper/odooimport/database"
	"github.com/ppreeper/odooimport/odooconn"
	"go.uber.org/zap"
)

var (
	log      *zap.SugaredLogger
	errorLog *zap.SugaredLogger
)

func main() {
	log = setupLogging("import.log")
	errorLog = setupLogging("error.log")
	log.Info("odooimport start")

	// Flags
	// Config File
	var c Conf
	c.getConf("config.yml")

	var hostname, dbase, flags string
	jobCount := 8
	var noUpdate bool

	flag.StringVar(&hostname, "host", "odoocrm.arthomson.com", "odoo database")
	flag.StringVar(&dbase, "d", "odoocrm", "odoo database")
	flag.StringVar(&flags, "f", "", "flags")
	flag.IntVar(&jobCount, "j", 8, "job count")
	flag.BoolVar(&noUpdate, "n", false, "no update")
	flag.Parse()
	fmt.Println(dbase)
	if dbase == "" {
		fmt.Println("no database specified")
		os.Exit(1)
	}
	ff := strings.Split(flags, ",")

	// connect to source database

	// open database connection
	sdb, err := database.OpenDatabase(database.Database{
		Name:     c.Source.Database,
		Driver:   c.Source.Driver,
		Host:     c.Source.Host,
		Port:     c.Source.Port,
		Database: c.Source.Database,
		Username: c.Source.Username,
		Password: c.Source.Password,
	})
	checkErr(err)
	defer func() {
		if err := sdb.Close(); err != nil {
			checkErr(err)
		}
	}()
	checkErr(err)

	// odoo connection

	o := odooconn.NewOdooConn(odooconn.OdooConn{
		Hostname: hostname,
		Port:     c.Port,
		Username: c.Login,
		Password: c.Password,
		Schema:   c.Schema,
		Database: dbase,
		DB:       sdb,
		NoUpdate: noUpdate,
		Log:      log,
		ErrLog:   errorLog,
		JobCount: jobCount,
	})
	// odoo login to get uid
	err = o.Login()
	checkErr(err)

	if flags != "" {
		for _, v := range ff {
			execFlags(v, o, c)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		errorLog.Errorw(err.Error())
	}
}

func fatalErr(err error) {
	if err != nil {
		errorLog.Fatalw(err.Error())
	}
}

func setupLogging(logName string) (log *zap.SugaredLogger) {
	_, err := os.Stat(logName)
	if os.IsNotExist(err) {
		file, err := os.Create(logName)
		fatalErr(err)
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Printf("Error closing file: %s\n", err)
			}
		}()
	}
	err = os.Truncate(logName, 0)
	checkErr(err)
	// Logger Setup
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logName}
	logger, _ := cfg.Build()
	log = logger.Sugar()
	return
}
