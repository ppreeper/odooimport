package odooconn

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

// DecimalPrecision function
func (o *OdooConn) DecimalPrecision(name string, digits int) {
	mdl := "decimal_precision"
	umdl := "decimal.precision"
	fmt.Println(umdl, name)
	bar := progressbar.Default(int64(1))
	err := bar.Add(1)
	o.checkErr(err)

	r, err := o.GetID(umdl, oarg{oarg{"name", "=", name}})
	o.checkErr(err)

	ur := map[string]interface{}{"name": name, "digits": digits}

	o.Log.Info(mdl, "model", umdl, "record", ur)

	o.Record(umdl, r, ur)
}
