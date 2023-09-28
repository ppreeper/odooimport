package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// UomCategory function
func (o *OdooConn) UomCategory() {
	mdl := "uom_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Println(umdl)

	uomCat := []struct {
		Name string
	}{{"Area"}}

	bar := progressbar.Default(int64(len(uomCat)))
	for _, v := range uomCat {
		err := bar.Add(1)
		o.checkErr(err)
		r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}})
		ur := map[string]interface{}{
			"name": v.Name,
		}
		o.Log.Infow(mdl, "model", umdl, "record", ur, "r", r)

		o.Record(umdl, r, ur)
	}
}
