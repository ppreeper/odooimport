package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ProductCategoryProduct1 function
func (o *OdooConn) ProductCategoryProduct1() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nproduct group1\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	saleableID := o.GetID(umdl, oarg{oarg{"name", "=", "Saleable"}})

	type Group1 struct {
		Group1 string `db:"group1"`
	}

	// not consumables group1
	stmt := `select distinct
	trim(ma.group1) group1
	from ct.matgrp_assign ma
	left join odoo.matkl_dnu md on ma.matkl = md.matkl
	where md.matkl is null
	and trim(ma.group1) <> 'CONSUMABLES'
	AND trim(ma.group1) <> ''
	order by trim(ma.group1)
	`

	var rrg1 []Group1
	err := o.DB.Select(&rrg1, stmt)
	o.checkErr(err)
	recs := len(rrg1)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rrg1 {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Group1) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			gid := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}})

			// Costing Method::property_cost_method
			// Standard Price: standard
			// First In First Out (FIFO): fifo
			// Average Cost (AVCO): average

			// Inventory Valuation::property_valuation
			// Manual: manual_periodic
			// Automated: real_time

			ur := map[string]interface{}{
				"name":                 v.Group1,
				"parent_id":            saleableID,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}
			o.Log.Infow(mdl, "group", "group1", "model", umdl, "record", ur, "gid", gid)

			o.Record(umdl, gid, ur)
			// if gid == -1 {
			// 	row, res, err := o.WriteRecord(umdl, gid, INSERT, ur)
			// 	if err != nil {
			// 		o.ErrLog.Infow(umdl, "row", row, "res", res, "err", err)
			// 	}
			// } else {
			// 	if !o.NoUpdate {
			// 		row, res, err := o.WriteRecord(umdl, gid, UPDATE, ur)
			// 		if err != nil {
			// 			o.ErrLog.Infow(umdl, "row", row, "res", res, "err", err)
			// 		}
			// 	}
			// }

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ProductCategoryProduct2 function
func (o *OdooConn) ProductCategoryProduct2() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nproduct group2\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	saleableID := o.GetID(umdl, oarg{oarg{"name", "=", "Saleable"}})

	type Group2 struct {
		Group1 string `db:"group1"`
		Group2 string `db:"group2"`
	}

	//////////
	// not consumables group2
	stmt := `select distinct
	trim(ma.group1) group1
	,trim(group2) group2
	from ct.matgrp_assign ma
	left join odoo.matkl_dnu md on ma.matkl = md.matkl
	where md.matkl is null
	and trim(ma.group1) <> 'CONSUMABLES'
	AND trim(ma.group1) <> ''
	AND trim(ma.group2) <> ''
	order by trim(ma.group1),trim(ma.group2)
	`

	var rrg2 []Group2
	err := o.DB.Select(&rrg2, stmt)
	o.checkErr(err)
	recs := len(rrg2)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rrg2 {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Group2) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			pid := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}, oarg{"parent_id", "=", saleableID}})
			gid := -1
			if pid != -1 {
				gid = o.GetID(umdl, oarg{oarg{"name", "=", v.Group2}, oarg{"parent_id", "=", pid}})
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
			o.Log.Infow(mdl, "group", "group2", "model", umdl, "record", ur, "gid", gid)

			o.Record(umdl, gid, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ProductCategoryProduct3 function
func (o *OdooConn) ProductCategoryProduct3() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nproduct group3\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	saleableID := o.GetID(umdl, oarg{oarg{"name", "=", "Saleable"}})

	type Group3 struct {
		Group1 string `db:"group1"`
		Group2 string `db:"group2"`
		Group3 string `db:"group3"`
	}

	//////////
	// not consumables group3
	stmt := `select distinct
	trim(ma.group1) group1
	,trim(group2) group2
	,trim(group3) group3
	from ct.matgrp_assign ma
	left join odoo.matkl_dnu md on ma.matkl = md.matkl
	where md.matkl is null
	and trim(ma.group1) <> 'CONSUMABLES'
	AND trim(ma.group1) <> ''
	AND trim(ma.group2) <> ''
	AND trim(ma.group3) <> ''
	order by trim(ma.group1),trim(ma.group2),trim(ma.group3)
	`

	var rrg3 []Group3
	err := o.DB.Select(&rrg3, stmt)
	o.checkErr(err)
	recs := len(rrg3)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rrg3 {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Group3) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			pid := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}, oarg{"parent_id", "=", saleableID}})
			sid := -1
			if pid != -1 {
				sid = o.GetID(umdl, oarg{oarg{"name", "=", v.Group2}, oarg{"parent_id", "=", pid}})
			}
			gid := -1
			if sid != -1 {
				gid = o.GetID(umdl, oarg{oarg{"name", "=", v.Group3}, oarg{"parent_id", "=", sid}})
			}

			// Costing Method::property_cost_method
			// Standard Price: standard
			// First In First Out (FIFO): fifo
			// Average Cost (AVCO): average

			// Inventory Valuation::property_valuation
			// Manual: manual_periodic
			// Automated: real_time

			ur := map[string]interface{}{
				"name":                 v.Group3,
				"parent_id":            sid,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}
			o.Log.Infow(mdl, "group", "group3", "model", umdl, "record", ur, "gid", gid)

			o.Record(umdl, gid, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ProductCategoryProduct4 function
func (o *OdooConn) ProductCategoryProduct4() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nproduct group4\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	saleableID := o.GetID(umdl, oarg{oarg{"name", "=", "Saleable"}})

	type MatGrp2 struct {
		Matgrp    string `db:"matgrp"`
		Matkl     string `db:"matkl"`
		Group1    string `db:"group1"`
		Group2    string `db:"group2"`
		Buyer     string `db:"buyer"`
		DayReview string `db:"day_review"`
	}

	//////////
	// not consumables matgrp2
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
	and trim(mga.group1) <> 'CONSUMABLES'
	AND trim(mga.group1) <> ''
	AND trim(mga.group2) <> ''
	AND trim(mga.group3) = ''
	AND trim(mga.matgrp) <> ''
	and trim(mga.matkl) like '____'
	order by trim(mga.group1),trim(mga.group2),trim(mga.matgrp)
	`

	var mg2 []MatGrp2
	err := o.DB.Select(&mg2, stmt)
	o.checkErr(err)
	recs := len(mg2)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range mg2 {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v MatGrp2) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			pid := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}, oarg{"parent_id", "=", saleableID}})
			sid := -1
			if pid != -1 {
				sid = o.GetID(umdl, oarg{oarg{"name", "=", v.Group2}, oarg{"parent_id", "=", pid}})
			}
			gid := -1
			if sid != -1 {
				gid = o.GetID(umdl, oarg{oarg{"name", "=", v.Matgrp}, oarg{"parent_id", "=", sid}})
			}

			// with buyerID modification
			// bidPartner := o.GetID("res.partner", oarg{oarg{"name", "=", v.Buyer}})
			// bid := o.GetID("res.users", oarg{oarg{"partner_id", "=", bidPartner}})

			// Costing Method::property_cost_method
			// Standard Price: standard
			// First In First Out (FIFO): fifo
			// Average Cost (AVCO): average

			// Inventory Valuation::property_valuation
			// Manual: manual_periodic
			// Automated: real_time

			ur := map[string]interface{}{
				"name":                 v.Matgrp,
				"parent_id":            sid,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}
			// with buyerID modification
			// if v.Matkl != "" {
			// 	ur["material_group"] = v.Matkl
			// }

			// if v.DayReview != "" {
			// 	ur["review_day"] = v.DayReview
			// }

			// if bid != -1 {
			// 	ur["buyer_id"] = bid
			// }

			o.Log.Infow(mdl, "group", "matgrp2", "model", umdl, "record", ur, "gid", gid)

			o.Record(umdl, gid, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ProductCategoryProduct5 function
func (o *OdooConn) ProductCategoryProduct5() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nproduct group5\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	saleableID := o.GetID(umdl, oarg{oarg{"name", "=", "Saleable"}})

	type MatGrp3 struct {
		Matgrp    string `db:"matgrp"`
		Matkl     string `db:"matkl"`
		Group1    string `db:"group1"`
		Group2    string `db:"group2"`
		Group3    string `db:"group3"`
		Buyer     string `db:"buyer"`
		DayReview string `db:"day_review"`
	}

	//////////
	// not consumables matgrp3
	stmt := `
	select distinct
	trim(mga.matgrp) matgrp
	,trim(mga.matkl) matkl
	,trim(mga.group1) group1
	,trim(mga.group2) group2
	,trim(mga.group3) group3
	,trim(case when b.buyer is null then '' else b.buyer end) buyer
	,lower(replace(trim(case when b.day_review is null then 'specialorder' else b.day_review end),' ','')) day_review
	from ct.matgrp_assign mga
	left join odoo.matkl_dnu md on mga.matkl = md.matkl
	left join ct.buyer_schedule b on mga.matkl = b.matkl
	where md.matkl is null
	and trim(mga.group1) <> 'CONSUMABLES'
	AND trim(mga.group1) <> ''
	AND trim(mga.group2) <> ''
	AND trim(mga.group3) <> ''
	AND trim(mga.matgrp) <> ''
	and trim(mga.matkl) like '____'
	order by trim(mga.group1),trim(mga.group2),trim(mga.group3),trim(mga.matgrp)
	`

	var mg3 []MatGrp3
	err := o.DB.Select(&mg3, stmt)
	o.checkErr(err)
	recs := len(mg3)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range mg3 {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v MatGrp3) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			pid := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}, oarg{"parent_id", "=", saleableID}})
			sid := -1
			if pid != -1 {
				sid = o.GetID(umdl, oarg{oarg{"name", "=", v.Group2}, oarg{"parent_id", "=", pid}})
			}
			tid := -1
			if sid != -1 {
				tid = o.GetID(umdl, oarg{oarg{"name", "=", v.Group3}, oarg{"parent_id", "=", sid}})
			}
			gid := -1
			if tid != -1 {
				gid = o.GetID(umdl, oarg{oarg{"name", "=", v.Matgrp}, oarg{"parent_id", "=", tid}})
			}

			// with buyerID modification
			// bidPartner := o.GetID("res.partner", oarg{oarg{"name", "=", v.Buyer}})
			// bid := o.GetID("res.users", oarg{oarg{"partner_id", "=", bidPartner}})

			// Costing Method::property_cost_method
			// Standard Price: standard
			// First In First Out (FIFO): fifo
			// Average Cost (AVCO): average

			// Inventory Valuation::property_valuation
			// Manual: manual_periodic
			// Automated: real_time

			ur := map[string]interface{}{
				"name":                 v.Matgrp,
				"parent_id":            tid,
				"property_cost_method": "average",
				"property_valuation":   "real_time",
			}

			// with buyerID modification
			// if v.Matkl != "" {
			// 	ur["material_group"] = v.Matkl
			// }

			// if v.DayReview != "" {
			// 	ur["review_day"] = v.DayReview
			// }

			// if bid != -1 {
			// 	ur["buyer_id"] = bid
			// }

			o.Log.Infow(mdl, "group", "matgrp3", "model", umdl, "record", ur, "gid", gid)

			o.Record(umdl, gid, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
