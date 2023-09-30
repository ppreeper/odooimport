package odooconn

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/ppreeper/odooimport/database"
	"github.com/ppreeper/odoojrpc"
)

var (
	// INSERT data
	INSERT = true
	// UPDATE data
	UPDATE = false
)

// type alias to reduce typing
type oarg = odoojrpc.FilterArg

// OdooConn structure to provide basic connection
type OdooConn struct {
	Hostname  string
	Port      int
	Database  string
	Username  string
	Password  string
	Schema    string
	JobCount  int
	NoUpdate  bool
	BatchSize int
	DB        *database.Database
	Log       *slog.Logger
	*odoojrpc.Odoo
}

// WriteRecord function
func (o *OdooConn) WriteRecord(umdl string, r int, mode bool, ur map[string]interface{}) (row int, res bool, err error) {
	if mode == INSERT {
		o.Log.Info("INSERT", "model", umdl, "record", ur)
		row, err = o.Create(umdl, ur)
		o.checkErr(err)
		o.Log.Info("INSERT", "model", umdl, "record", r, "err", err)
	} else if mode == UPDATE {
		o.Log.Info("UPDATE", "model", umdl, "record", ur)
		res, err = o.Update(umdl, r, ur)
		o.checkErr(err)
		o.Log.Info("UPDATE", "model", umdl, "record", r, "err", err)
	}
	return
}

func (o *OdooConn) Record(umdl string, r int, ur map[string]interface{}) {
	if r == -1 {
		row, res, err := o.WriteRecord(umdl, r, INSERT, ur)
		if err != nil {
			o.Log.Info(umdl, "row", row, "res", res, "err", err)
		}
	} else {
		if !o.NoUpdate {
			row, res, err := o.WriteRecord(umdl, r, UPDATE, ur)
			if err != nil {
				o.Log.Info(umdl, "row", row, "res", res, "err", err)
			}
		}
	}
}

// NewOdooConn initializer
func NewOdooConn(oc OdooConn) *OdooConn {
	oc.Odoo = &odoojrpc.Odoo{
		Hostname: oc.Hostname,
		Port:     oc.Port,
		Username: oc.Username,
		Password: oc.Password,
		Schema:   oc.Schema,
		Database: oc.Database,
	}
	return &oc
}

// Utils
// checkErr function
func (o *OdooConn) checkErr(err error) {
	if err != nil {
		o.Log.Error(err.Error())
	}
}

func CheckErr[T any](val T, err error) T {
	if err != nil {
		fmt.Println()
		return val
	}
	return val
}

func pager(batch, total int) {
	batchTotal := math.Ceil(float64(total) / float64(batch))
	fmt.Println("batchTotal", batchTotal)
}

func getPages(vallist []int, pagesize int) [][]int {
	// getPage size and return list subset
	// return slice of int slices
	return [][]int{}
}

func removeDuplicate(vals []int) []int {
	// remove duplicate values in slice
	return []int{}
}

func valInSlice[T comparable](val T, vals []T) bool {
	for _, b := range vals {
		if b == val {
			return true
		}
	}
	return false
}
