package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// MRPBom function
func (o *OdooConn) MRPBom() {
	mdl := "mrp_bom"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v MRPBom\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select default_code,company,plant,code,qty,uom,kit,bomtext 
	from odoo.artg_bom_list_import
	-- where bomtext <> ''
	-- limit 10
	`
	// stmt = stmt + ` limit 10`
	type BOM struct {
		DefaultCode string  `db:"default_code"`
		Company     string  `db:"company"`
		Plant       string  `db:"plant"`
		Code        string  `db:"code"`
		Qty         float64 `db:"qty"`
		UOM         string  `db:"uom"`
		Kit         float64 `db:"kit"`
		BOMText     string  `db:"bomtext"`
	}
	var rr []BOM
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()
	uuom := o.UomUomMap()

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v BOM) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid := cids[v.Company]
			productID, err := o.GetID("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}})
			o.checkErr(err)
			// uomID := o.GetID("uom.uom", oarg{oarg{"name", "=", v.UOM}})
			uomID := uuom[v.UOM]

			r, err := o.GetID(umdl, oarg{oarg{"product_tmpl_id", "=", productID}, oarg{"company_id", "=", cid}})
			o.checkErr(err)

			ur := map[string]interface{}{
				"product_tmpl_id":       productID,
				"product_qty":           v.Qty,
				"product_uom_id":        uomID,
				"company_id":            cid,
				"code":                  v.Code,
				"production_order_text": v.BOMText,
			}
			if v.Kit == 1 {
				ur["type"] = "phantom"
			} else {
				pickType, err := o.GetID("stock.picking.type", oarg{oarg{"company_id", "=", cid}, oarg{"sequence_code", "=", "MO"}, oarg{"default_location_src_id", "like", v.Code}})
				o.checkErr(err)
				ur["type"] = "normal"
				ur["picking_type_id"] = pickType
			}

			o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
