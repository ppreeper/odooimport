package odooconn

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

// IRConfigParameter function
func (o *OdooConn) IRConfigParameter(key string, val interface{}) {
	mdl := "ir_config_parameter"
	umdl := "ir.config_parameter"
	fmt.Println(umdl, key)
	bar := progressbar.Default(int64(1))
	err := bar.Add(1)
	o.checkErr(err)

	r := o.GetID(umdl, oarg{oarg{"key", "=", key}})

	ur := map[string]interface{}{"key": key, "value": val}

	o.Log.Infow(mdl, "model", umdl, "record", ur)

	o.Record(umdl, r, ur)
}
