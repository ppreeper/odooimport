package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ProductCategoryConsumable1 function
func (o *OdooConn) ProductCategoryConsumable1() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nconsumables1\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	allID, err := o.GetID(umdl, oarg{oarg{"name", "=", "All"}})
	o.checkErr(err)

	type Group1 struct {
		Group1 string `db:"group1"`
	}

	// consumables group1
	stmt := `
	select distinct
	trim(ma.group1) group1
	from ct.matgrp_assign ma
	left join odoo.matkl_dnu md on ma.matkl = md.matkl
	where md.matkl is null
	and trim(ma.group1) = 'CONSUMABLES'
	AND trim(ma.group1) <> ''
	order by trim(ma.group1)
	`

	var ccg1 []Group1
	err = o.DB.Select(&ccg1, stmt)
	o.checkErr(err)
	recs := len(ccg1)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range ccg1 {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Group1) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			gid, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}})
			o.checkErr(err)
			gval, err := o.SearchRead(umdl, oarg{oarg{"name", "=", v.Group1}}, 0, 0, []string{"property_cost_method", "property_valuation"})
			o.checkErr(err)
			o.Log.Info("", "gval", gval)

			// Costing Method::property_cost_method
			// Standard Price: standard
			// First In First Out (FIFO): fifo
			// Average Cost (AVCO): average

			// Inventory Valuation::property_valuation
			// Manual: manual_periodic
			// Automated: real_time

			ur := map[string]interface{}{
				"name":                 v.Group1,
				"parent_id":            allID,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}

			o.Log.Info(mdl, "group", "consumables1", "model", umdl, "record", ur, "gid", gid)

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
}

// ProductCategoryConsumable2 function
func (o *OdooConn) ProductCategoryConsumable2() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nconsumables2\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	allID, err := o.GetID(umdl, oarg{oarg{"name", "=", "All"}})
	o.checkErr(err)

	type Group2 struct {
		Group1 string `db:"group1"`
		Group2 string `db:"group2"`
	}

	//////////
	// consumables group2
	stmt := `
	select distinct
	trim(ma.group1) group1
	,trim(ma.group2) group2
	from ct.matgrp_assign ma
	left join odoo.matkl_dnu md on ma.matkl = md.matkl
	where md.matkl is null
	and trim(ma.group1) = 'CONSUMABLES'
	AND trim(ma.group1) <> ''
	AND trim(ma.group2) <> ''
	order by trim(ma.group1),trim(ma.group2)
	`

	var ccg2 []Group2
	err = o.DB.Select(&ccg2, stmt)
	o.checkErr(err)
	recs := len(ccg2)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range ccg2 {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Group2) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			pid, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}, oarg{"parent_id", "=", allID}})
			o.checkErr(err)
			gid := -1
			if pid != -1 {
				gid, err = o.GetID(umdl, oarg{oarg{"name", "=", v.Group2}, oarg{"parent_id", "=", pid}})
				o.checkErr(err)
			}

			// Costing Method::property_cost_method
			// Standard Price: standard
			// First In First Out (FIFO): fifo
			// Average Cost (AVCO): average

			// Inventory Valuation::property_valuation
			// Manual: manual_periodic
			// Automated: real_time

			ur := map[string]interface{}{
				"name":                 v.Group2,
				"parent_id":            pid,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}
			o.Log.Info(mdl, "group", "consumables2", "model", umdl, "record", ur, "gid", gid)

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
}

// ProductCategoryConsumable3 function
func (o *OdooConn) ProductCategoryConsumable3() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nconsumables3\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	allID, err := o.GetID(umdl, oarg{oarg{"name", "=", "All"}})
	o.checkErr(err)

	type MatGrp2 struct {
		Matgrp    string `db:"matgrp"`
		Matkl     string `db:"matkl"`
		Group1    string `db:"group1"`
		Group2    string `db:"group2"`
		Buyer     string `db:"buyer"`
		DayReview string `db:"day_review"`
	}

	//////////
	// consumables matgrp2
	stmt := `select distinct
	trim(mga.matgrp) matgrp
	,trim(mga.matkl) matkl
	,trim(mga.group1) group1
	,trim(mga.group2) group2
	,trim(case when b.buyer is null then '' else b.buyer end) buyer
	,lower(replace(trim(case when b.day_review is null then 'specialorder' else b.day_review end),' ','')) day_review
	from ct.matgrp_assign mga
	left join odoo.matkl_dnu md on mga.matkl = md.matkl
	left join ct.buyer_schedule b on mga.matkl = b.matkl
	where md.matkl is null
	and trim(mga.group1) = 'CONSUMABLES'
	AND trim(mga.group2) <> ''
	AND trim(mga.group3) = ''
	AND trim(mga.matgrp) <> ''
	and trim(mga.matkl) like '____'
	order by trim(mga.group1),trim(mga.group2),trim(mga.matgrp)
	`

	var ccmg2 []MatGrp2
	err = o.DB.Select(&ccmg2, stmt)
	o.checkErr(err)
	recs := len(ccmg2)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range ccmg2 {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v MatGrp2) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			pid, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}, oarg{"parent_id", "=", allID}})
			o.checkErr(err)
			sid := -1
			if pid != -1 {
				sid, err = o.GetID(umdl, oarg{oarg{"name", "=", v.Group2}, oarg{"parent_id", "=", pid}})
				o.checkErr(err)
			}
			gid := -1
			if sid != -1 {
				gid, err = o.GetID(umdl, oarg{oarg{"name", "=", v.Matgrp}, oarg{"parent_id", "=", sid}})
				o.checkErr(err)
			}

			// Costing Method::property_cost_method
			// Standard Price: standard
			// First In First Out (FIFO): fifo
			// Average Cost (AVCO): average

			// Inventory Valuation::property_valuation
			// Manual: manual_periodic
			// Automated: real_time

			// with buyer id setting
			// bidPartner := o.GetID("res.partner", oarg{oarg{"name", "=", v.Buyer}})
			// bid := o.GetID("res.users", oarg{oarg{"partner_id", "=", bidPartner}})

			ur := map[string]interface{}{
				"name":                 v.Matgrp,
				"parent_id":            sid,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}

			// with buyer id setting
			// if v.Matkl != "" {
			// 	ur["material_group"] = v.Matkl
			// }

			// if v.DayReview != "" {
			// 	ur["review_day"] = v.DayReview
			// }

			// if bid != -1 {
			// 	ur["buyer_id"] = bid
			// }

			o.Log.Info(mdl, "group", "consumables matgrp2", "model", umdl, "record", ur, "gid", gid)

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
}
