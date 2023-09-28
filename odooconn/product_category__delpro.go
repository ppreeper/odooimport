package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ProductCategoryDelpro function
func (o *OdooConn) ProductCategoryDelpro() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nproduct delpro\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	saleableID := o.GetID(umdl, oarg{oarg{"name", "=", "Saleable"}})

	//////////
	// not consumables
	stmt := `
	select distinct
	case
	when dpc.ptype = 'NONE' then 'MISCELLANEOUS PRODUCTS'
	when dpc.ptype = 'Regulators/Relief' then 'Regulators and Relief'
	when dpc.ptype = 'Tubing/Fittings' then 'Tubing and Fittings'
	else dpc.ptype end ptype
	from odoo.delpro_product_categories dpc
	order by case
	when dpc.ptype = 'NONE' then 'MISCELLANEOUS PRODUCTS'
	when dpc.ptype = 'Regulators/Relief' then 'Regulators and Relief'
	when dpc.ptype = 'Tubing/Fittings' then 'Tubing and Fittings'
	else dpc.ptype end
	`
	type MatCat struct {
		Ptype string `db:"ptype"`
	}

	var mgc []MatCat
	err := o.DB.Select(&mgc, stmt)
	o.checkErr(err)
	recs := len(mgc)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range mgc {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v MatCat) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			gid := o.GetID(umdl, oarg{oarg{"name", "=", v.Ptype}, oarg{"parent_id", "=", saleableID}})

			ur := map[string]interface{}{
				"name":                 v.Ptype,
				"parent_id":            saleableID,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}

			o.Log.Infow(mdl, "group", "matgrp3", "model", umdl, "record", ur, "gid", gid)

			o.Record(umdl, gid, ur)
			// if gid == -1 {
			// 	row, res, err := o.WriteRecord(umdl, gid, INSERT, ur)
			// 	o.ErrLog.Infow(umdl, "row", row, "res", res, "err", err)
			// } else {
			// 	if !o.NoUpdate {
			// 		row, res, err := o.WriteRecord(umdl, gid, UPDATE, ur)
			// 		o.ErrLog.Infow(umdl, "row", row, "res", res, "err", err)
			// 	}
			// }

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()

	stmt = `
	select distinct
	case
	when dpc.ptype = 'NONE' then 'MISCELLANEOUS PRODUCTS'
	when dpc.ptype = 'Regulators/Relief' then 'Regulators and Relief'
	when dpc.ptype = 'Tubing/Fittings' then 'Tubing and Fittings'
	else dpc.ptype end ptype
	,dpc.brand
	,'Kristina Jacob' buyer
	,'specialorder' day_review
	from odoo.delpro_product_categories dpc
	order by case
	when dpc.ptype = 'NONE' then 'MISCELLANEOUS PRODUCTS'
	when dpc.ptype = 'Regulators/Relief' then 'Regulators and Relief'
	when dpc.ptype = 'Tubing/Fittings' then 'Tubing and Fittings'
	else dpc.ptype end,
	dpc.brand
	`

	type MatCatGrp struct {
		Ptype     string `db:"ptype"`
		Brand     string `db:"brand"`
		Buyer     string `db:"buyer"`
		DayReview string `db:"day_review"`
	}

	var mgcg []MatCatGrp
	err = o.DB.Select(&mgcg, stmt)
	o.checkErr(err)
	recs = len(mgcg)
	bar = progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range mgcg {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v MatCatGrp) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			pid := o.GetID(umdl, oarg{oarg{"name", "=", v.Ptype}, oarg{"parent_id", "=", saleableID}})
			gid := o.GetID(umdl, oarg{oarg{"name", "=", v.Brand}, oarg{"parent_id", "=", pid}})

			bidPartner := o.GetID("res.partner", oarg{oarg{"name", "=", v.Buyer}})
			bid := o.GetID("res.users", oarg{oarg{"partner_id", "=", bidPartner}})

			ur := map[string]interface{}{
				"name":                 v.Brand,
				"parent_id":            pid,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}

			if v.DayReview != "" {
				ur["review_day"] = v.DayReview
			}

			if bid != -1 {
				ur["buyer_id"] = bid
			}

			o.Log.Infow(umdl, "record", ur, "parent_id", pid, "gid", gid)

			o.Record(umdl, gid, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
