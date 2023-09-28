package odooconn

import (
	"github.com/ppreeper/odooimport/database"
	"github.com/ppreeper/odoojrpc"
	"go.uber.org/zap"
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
	Hostname string
	Port     int
	Database string
	Username string
	Password string
	Schema   string
	JobCount int
	NoUpdate bool
	DB       *database.Database
	Log      *zap.SugaredLogger
	ErrLog   *zap.SugaredLogger
	*odoojrpc.Odoo
}

// WriteRecord function
func (o *OdooConn) WriteRecord(umdl string, r int, mode bool, ur map[string]interface{}) (row int, res bool, err error) {
	if mode == INSERT {
		o.Log.Infow("INSERT", "model", umdl, "record", ur)
		row, err = o.Create(umdl, ur)
		o.checkErr(err)
		o.Log.Infow("INSERT", "model", umdl, "record", r, "err", err)
	} else if mode == UPDATE {
		o.Log.Infow("UPDATE", "model", umdl, "record", ur)
		res, err = o.Update(umdl, r, ur)
		o.checkErr(err)
		o.Log.Infow("UPDATE", "model", umdl, "record", r, "err", err)
	}
	return
}

func (o *OdooConn) Record(umdl string, r int, ur map[string]interface{}) {
	if r == -1 {
		row, res, err := o.WriteRecord(umdl, r, INSERT, ur)
		if err != nil {
			o.ErrLog.Infow(umdl, "row", row, "res", res, "err", err)
		}
	} else {
		if !o.NoUpdate {
			row, res, err := o.WriteRecord(umdl, r, UPDATE, ur)
			if err != nil {
				o.ErrLog.Infow(umdl, "row", row, "res", res, "err", err)
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

// checkErr function
func (o *OdooConn) checkErr(err error) {
	if err != nil {
		o.ErrLog.Errorw(err.Error())
	}
}
