package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// UomUom function
func (o *OdooConn) UomUom() {
	mdl := "uom_uom"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Println(umdl)

	stmt := `
	select "name",category_id,factor,rounding,uom_type,measure_type
	from odoo.uom_uom uu
	order by category_id
	,case uom_type
	when 'reference' then 1
	when 'smaller' then 2
	when 'bigger' then 3
	else 0 end
	,"name"
	`
	uoms := []struct {
		Name        string  `db:"name"`
		CategoryID  string  `db:"category_id"`
		Factor      float64 `db:"factor"`
		Rounding    float64 `db:"rounding"`
		UomType     string  `db:"uom_type"`
		MeasureType string  `db:"measure_type"`
	}{}
	err := o.DB.Select(&uoms, stmt)
	o.checkErr(err)

	bar := progressbar.Default(int64(len(uoms)))
	for _, v := range uoms {
		err := bar.Add(1)
		o.checkErr(err)
		r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}})
		pr := o.GetID("uom.category", oarg{oarg{"name", "=", v.CategoryID}})
		ur := map[string]interface{}{
			"name":        v.Name,
			"category_id": pr,
			"factor":      v.Factor,
			"uom_type":    v.UomType,
			"rounding":    v.Rounding,
		}
		o.Log.Infow(mdl, "model", umdl, "record", ur, "r", r)

		o.Record(umdl, r, ur)
	}
}

func (o *OdooConn) UomUomMap() map[string]int {
	mdl := "uom_uom"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc := o.SearchRead(umdl, oarg{}, 0, 0, []string{"name"})
	cids := map[string]int{}
	for _, c := range cc {
		cids[c["name"].(string)] = int(c["id"].(float64))
	}
	return cids
}
