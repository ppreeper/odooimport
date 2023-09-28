package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// MRPBomLine function
func (o *OdooConn) MRPBomLine(c string) {
	mdl := "mrp_bom_line"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v MRPBomLine\n", umdl)

	stmt := `
	select default_code,company,plant,code,qty,uom,kit from odoo.artg_bom_list_import 
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
	}
	var rr []BOM
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()

	bomlinestmt := `
	select default_code,item_code,stlkn,item_qty,item_uom from odoo.artg_bom_item_import where default_code = $1
	`
	type BOMLine struct {
		DefaultCode string  `db:"default_code"`
		ItemCode    string  `db:"item_code"`
		Line        float64 `db:"stlkn"`
		ItemQty     float64 `db:"item_qty"`
		ItemUOM     string  `db:"item_uom"`
	}

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v BOM) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			// o.Log.Infow(umdl, "v", v)
			companyID := cids[v.Company]
			productID := o.GetID("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}})
			bomID := o.GetID("mrp.bom", oarg{oarg{"product_tmpl_id", "=", productID}, oarg{"company_id", "=", companyID}})

			var bb []BOMLine
			err := o.DB.Select(&bb, bomlinestmt, v.DefaultCode)
			o.checkErr(err)
			// brecs := len(rr)
			o.Log.Infow(umdl, "bb", bb)
			for _, b := range bb {
				uomID := o.GetID("uom.uom", oarg{oarg{"name", "=", b.ItemUOM}})
				productItemID := o.GetID("product.product", oarg{oarg{"default_code", "=", b.ItemCode}})
				r := o.GetID(umdl, oarg{
					oarg{"product_id", "=", productItemID},
					oarg{"company_id", "=", companyID},
					oarg{"product_uom_id", "=", uomID},
					oarg{"bom_id", "=", bomID},
				})
				// o.Log.Infow(umdl, "b", b, "productItemID", productItemID, "r", r)

				ur := map[string]interface{}{
					"product_id":     productItemID,
					"company_id":     companyID,
					"product_qty":    b.ItemQty,
					"product_uom_id": uomID,
					"bom_id":         bomID,
				}

				o.Log.Infow(umdl, "ur", ur, "r", r)

				o.Record(umdl, r, ur)
			}

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
