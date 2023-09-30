package odooconn

import "strings"

// StockInventoryLine function
func (o *OdooConn) StockInventoryLine() {
	mdl := "stock_inventory_line"
	umdl := strings.Replace(mdl, "_", ".", -1)

	stmt := `
	select "werks","company","parent","type","website","vtweg","spart",
	"kunnr","lifnr","fabkl","name1","street","city","region","post_code",
	"country","taxjurcode","taxiw","tel_number","fax_number","time_zone",
	"citycode","cityname","prefix","ccode","is_company"
	from odoo.plants
	`
	rr := []struct{}{}
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)

	o.Log.Info(mdl, "model", umdl, "record", recs)
}
