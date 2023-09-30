package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/ppreeper/odooimport/database"
	"github.com/ppreeper/odooimport/odooconn"
)

var ilog *slog.Logger

func main() {
	ilog = setupLogging("import.log")
	ilog.Info("odooimport start")

	// Flags
	var config, flags string
	jobCount := 8
	var noUpdate bool

	flag.StringVar(&config, "c", "config.yml", "config file")
	flag.StringVar(&flags, "f", "", "flags")
	flag.IntVar(&jobCount, "j", 8, "job count")
	flag.BoolVar(&noUpdate, "n", false, "no update")
	flag.Parse()

	// Config File
	var c Conf
	c.getConf(config)
	fmt.Println("config", c)

	ff := strings.Split(flags, ",")
	fmt.Println("flags", ff)

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
		Log:      ilog,
		JobCount: jobCount,
	})
	// odoo login to get uid
	checkErr(o.Login())

	if flags != "" {
		for _, v := range ff {
			execFlags(v, o, c)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		ilog.Error(err.Error())
	}
}

func fatalErr(err error) {
	if err != nil {
		ilog.Error(err.Error())
		os.Exit(1)
	}
}

func setupLogging(logName string) *slog.Logger {
	// check for file existence
	_, err := os.Stat(logName)
	if os.IsNotExist(err) {
		file, err := os.Create(logName)
		fatalErr(err)
		defer file.Close()
	}
	// if exists truncate file
	err = os.Truncate(logName, 0)
	checkErr(err)

	f, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	fatalErr(err)
	logwriter := io.Writer(f)
	return slog.New(slog.NewTextHandler(logwriter, nil))
}
