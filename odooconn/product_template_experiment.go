package odooconn

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) ProductTemplateUnlink2() {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplatePurge\n", umdl)

	stmt := `select distinct default_code from odoo.product_template order by default_code
	-- limit 1
	`

	type Product struct {
		DefaultCode string `db:"default_code"`
	}
	var dbrecs []Product
	if stmt == "" {
		return
	}
	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	odoorecs := o.SearchRead(umdl, oarg{}, 0, 0, []string{"default_code"})
	fmt.Println("products:", len(odoorecs))

	bar := progressbar.Default(int64(len(odoorecs)))
	// ids := []int{}
	for _, or := range odoorecs {
		for _, dr := range dbrecs {
			if or["default_code"] == dr.DefaultCode {
				// ids = append(ids, int(or["id"].(float64)))
				o.Unlink(umdl, []int{int(or["id"].(float64))})
			}
		}
		bar.Add(1)
	}
}

func (o *OdooConn) ProductTemplateUnlink3() {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplatePurge\n", umdl)

	type Product struct {
		DefaultCode string `db:"default_code"`
	}
	var dbrecs []Product

	stmt := `select distinct default_code from odoo.product_template2 order by default_code
	-- limit 1
	`
	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	fmt.Println("dbrecs:", len(dbrecs))

	odoorecs := o.SearchRead(umdl, oarg{}, 0, 0, []string{"default_code"})
	fmt.Println("odoorecs:", len(odoorecs))

	odooList := make(map[string]int)
	for _, r := range odoorecs {
		switch r["default_code"].(type) {
		case string:
			odooList[r["default_code"].(string)] = int(r["id"].(float64))
		}
	}

	var dList []int

	for _, dr := range dbrecs {
		if _, ok := odooList[dr.DefaultCode]; ok {
			dList = append(dList, odooList[dr.DefaultCode])
		}
	}

	bar := progressbar.Default(int64(len(dList)))
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(len(dList))
	for _, r := range dList {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, r int) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Unlink(umdl, []int{r})

			<-sem
		}(sem, &wg, bar, r)
	}
	wg.Wait()
}

// ProductTemplate function
func (o *OdooConn) ProductTemplate2() {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplate\n", umdl)

	stmt := `select "name"
	,default_code,default_code barcode,ptype detailed_type,list_price,standard_price
	,basic description,inspection description_pickingin,purchase description_purchase
	,category categ_id,matgrp
	from odoo.product_template2
	order by matgrp,default_code
	limit 9
	`

	type Product struct {
		Name                 string  `db:"name"`
		DefaultCode          string  `db:"default_code"`
		Barcode              string  `db:"barcode"`
		DetailedType         string  `db:"detailed_type"`
		ListPrice            float64 `db:"list_price"`
		StandardPrice        float64 `db:"standard_price"`
		DescriptionSale      string  `db:"description"`
		DescriptionPickingin string  `db:"description_pickingin"`
		DescriptionPurchase  string  `db:"description_purchase"`
		Category             string  `db:"categ_id"`
		Matgrp               string  `db:"matgrp"`
	}
	var rr []Product
	if stmt == "" {
		return
	}
	o.Log.Info(stmt)
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)

	taxSell := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for sales - 5%"}})
	taxPurchase := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for purchases - 5%"}})

	pgs := o.ProductCategoryMap()
	// uom := o.UomMapper()

	// tasker
	recs := len(rr)
	bar := progressbar.Default(int64(recs))
	for _, v := range rr {
		r := o.GetID(umdl, oarg{oarg{"default_code", "=", v.DefaultCode}})

		categID := -1
		if v.Matgrp != "" {
			categID = pgs[v.Matgrp]
		}

		ur := map[string]interface{}{
			"name":                  v.Name,
			"default_code":          v.DefaultCode,
			"barcode":               v.DefaultCode,
			"type":                  v.DetailedType,
			"list_price":            v.ListPrice,
			"standard_price":        v.StandardPrice,
			"description_sale":      v.DescriptionSale,
			"description_pickingin": v.DescriptionPickingin,
			"description_purchase":  v.DescriptionPurchase,
			"taxes_id":              []int{taxSell},
			"supplier_taxes_id":     []int{taxPurchase},
		}

		if categID != -1 {
			ur["categ_id"] = categID
		}

		o.Log.Infow(umdl, "ur", ur, "r", r)

		o.Record(umdl, r, ur)

		bar.Add(1)

	}
}

// ProductTemplate function
func (o *OdooConn) ProductTemplate3() {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplate\n", umdl)

	stmt := `select "name"
	,default_code,default_code barcode,ptype detailed_type,list_price,standard_price
	,basic description,inspection description_pickingin,purchase description_purchase
	,category categ_id,matgrp
	from odoo.product_template2
	order by matgrp,default_code
	limit 10000
	`

	type Product struct {
		Name                 string  `db:"name"`
		DefaultCode          string  `db:"default_code"`
		Barcode              string  `db:"barcode"`
		DetailedType         string  `db:"detailed_type"`
		ListPrice            float64 `db:"list_price"`
		StandardPrice        float64 `db:"standard_price"`
		DescriptionSale      string  `db:"description"`
		DescriptionPickingin string  `db:"description_pickingin"`
		DescriptionPurchase  string  `db:"description_purchase"`
		Category             string  `db:"categ_id"`
		Matgrp               string  `db:"matgrp"`
	}
	var rr []Product
	if stmt == "" {
		return
	}
	o.Log.Info(stmt)
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)

	taxSell := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for sales - 5%"}})
	taxPurchase := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for purchases - 5%"}})

	pgs := o.ProductCategoryMap()
	// uom := o.UomMapper()

	// tasker
	recs := len(rr)
	bar := progressbar.Default(int64(recs))
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Product) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			r := o.GetID(umdl, oarg{oarg{"default_code", "=", v.DefaultCode}})

			categID := -1
			if v.Matgrp != "" {
				categID = pgs[v.Matgrp]
			}

			ur := map[string]interface{}{
				"name":                  v.Name,
				"default_code":          v.DefaultCode,
				"barcode":               v.DefaultCode,
				"type":                  v.DetailedType,
				"list_price":            v.ListPrice,
				"standard_price":        v.StandardPrice,
				"description_sale":      v.DescriptionSale,
				"description_pickingin": v.DescriptionPickingin,
				"description_purchase":  v.DescriptionPurchase,
				"taxes_id":              []int{taxSell},
				"supplier_taxes_id":     []int{taxPurchase},
			}

			if categID != -1 {
				ur["categ_id"] = categID
			}

			o.Log.Infow(umdl, "ur", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

func (o *OdooConn) ProductTemplate4() {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplate\n", umdl)

	type Product struct {
		Name                 string  `db:"name"`
		DefaultCode          string  `db:"default_code"`
		Barcode              string  `db:"barcode"`
		DetailedType         string  `db:"detailed_type"`
		ListPrice            float64 `db:"list_price"`
		StandardPrice        float64 `db:"standard_price"`
		DescriptionSale      string  `db:"description"`
		DescriptionPickingin string  `db:"description_pickingin"`
		DescriptionPurchase  string  `db:"description_purchase"`
		Category             string  `db:"categ_id"`
		Matgrp               string  `db:"matgrp"`
	}
	var dbrecs []Product
	stmt := `select distinct
	"name",default_code,default_code barcode,ptype detailed_type
	,list_price , standard_price
	,basic description,inspection description_pickingin,purchase description_purchase
	,category categ_id,matgrp
	from odoo.product_template2
	order by matgrp,default_code
	limit 10000
	`
	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	odoorecs := o.SearchRead(umdl, oarg{}, 0, 0, []string{"default_code"})

	odooList := make(map[string]int)
	for _, r := range odoorecs {
		switch r["default_code"].(type) {
		case string:
			odooList[r["default_code"].(string)] = int(r["id"].(float64))
		}
	}

	uList := make(map[string]int)
	lList := make(map[string]struct{})

	for _, dr := range dbrecs {
		if _, ok := odooList[dr.DefaultCode]; ok {
			uList[dr.DefaultCode] = odooList[dr.DefaultCode]
		} else {
			lList[dr.DefaultCode] = struct{}{}
		}
	}

	// fmt.Println("uList", len(uList), uList)
	// fmt.Println("lList", len(lList))

	// recs := len(loadList)
	// bar := progressbar.Default(int64(recs))

	// load records
	header := []string{
		"name",
		"default_code",
		"barcode",
		"detailed_type",
		"list_price",
		"standard_price",
		"description_sale",
		"description_pickingin",
		"description_purchase",
		"taxes_id",
		"supplier_taxes_id",
	}
	fmt.Println(header)

	recs := make([]interface{}, 0, o.BatchSize)
	// fmt.Println("recs", len(recs), cap(recs), recs)

	i := 0
	start := 0
	end := 0
	// load records
	fmt.Println("load records", len(lList))
	bar := progressbar.Default(int64(len(lList)))
	for k, dr := range dbrecs {
		if _, ok := lList[dr.DefaultCode]; ok {
			rec := []interface{}{
				dr.Name,
				dr.DefaultCode,
				dr.DefaultCode,
				dr.DetailedType,
				dr.ListPrice,
				dr.StandardPrice,
				dr.DescriptionSale,
				dr.DescriptionPickingin,
				dr.DescriptionPurchase,
				"GST for sales - 5%",
				"GST for purchases - 5%",
			}
			recs = append(recs, rec)
			i++
			if i == o.BatchSize || k == len(dbrecs)-1 {
				end = start + i
				// fmt.Println("flush and clear", start, end)
				if len(recs) > 0 {
					// o.Log.Infow(umdl, "start", start, "end", end, "diff", start-end)
					out, err := o.Load(umdl, header, recs)
					if err != nil {
						o.Log.Errorw(umdl, "out", out, "err", err)
					}
					bar.Add(end - start)
				}
				start += i
				recs = []interface{}{}
				i = 0
			}
		}
	}

	// update records
	fmt.Println("update records", len(uList))
	// taxSell := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for sales - 5%"}})
	// taxPurchase := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for purchases - 5%"}})

	// pgs := o.ProductCategoryMap()
	// uom := o.UomMapper()

	// bar = progressbar.Default(int64(len(uList)))
	// for _, dr := range dbrecs {
	// 	if _, ok := uList[dr.DefaultCode]; ok {
	// 		r := uList[dr.DefaultCode]

	// 		categID := -1
	// 		if dr.Matgrp != "" {
	// 			categID = pgs[dr.Matgrp]
	// 		}

	// 		var ur = map[string]interface{}{
	// 			"name":                  dr.Name,
	// 			"default_code":          dr.DefaultCode,
	// 			"barcode":               dr.DefaultCode,
	// 			"type":                  dr.DetailedType,
	// 			"list_price":            dr.ListPrice,
	// 			"standard_price":        dr.StandardPrice,
	// 			"description_sale":      dr.DescriptionSale,
	// 			"description_pickingin": dr.DescriptionPickingin,
	// 			"description_purchase":  dr.DescriptionPurchase,
	// 			"taxes_id":              []int{taxSell},
	// 			"supplier_taxes_id":     []int{taxPurchase},
	// 		}

	// 		if categID != -1 {
	// 			ur["categ_id"] = categID
	// 		}

	// 		o.Log.Infow(umdl, "ur", ur, "r", r)

	// 		o.Record(umdl, r, ur)

	// 		bar.Add(1)
	// 	}
	// }
}

func (o *OdooConn) ProductTemplate5() {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplate\n", umdl)

	type Product struct {
		Name                 string  `db:"name"`
		DefaultCode          string  `db:"default_code"`
		Barcode              string  `db:"barcode"`
		DetailedType         string  `db:"detailed_type"`
		ListPrice            float64 `db:"list_price"`
		StandardPrice        float64 `db:"standard_price"`
		DescriptionSale      string  `db:"description"`
		DescriptionPickingin string  `db:"description_pickingin"`
		DescriptionPurchase  string  `db:"description_purchase"`
		Category             string  `db:"categ_id"`
		Matgrp               string  `db:"matgrp"`
	}
	var dbrecs []Product
	stmt := `select distinct
	"name",default_code,default_code barcode,ptype detailed_type
	,list_price , standard_price
	,basic description,inspection description_pickingin,purchase description_purchase
	,category categ_id,matgrp
	from odoo.product_template2
	order by matgrp,default_code
	limit 10000
	`
	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	odoorecs := o.SearchRead(umdl, oarg{}, 0, 0, []string{"default_code"})

	odooList := make(map[string]int)
	for _, r := range odoorecs {
		switch r["default_code"].(type) {
		case string:
			odooList[r["default_code"].(string)] = int(r["id"].(float64))
		}
	}

	uList := make(map[string]int)
	lList := make(map[string]struct{})

	for _, dr := range dbrecs {
		if _, ok := odooList[dr.DefaultCode]; ok {
			uList[dr.DefaultCode] = odooList[dr.DefaultCode]
		} else {
			lList[dr.DefaultCode] = struct{}{}
		}
	}

	// fmt.Println("uList", len(uList), uList)
	// fmt.Println("lList", len(lList))

	// recs := len(loadList)
	// bar := progressbar.Default(int64(recs))

	// load records
	header := []string{
		"name",
		"default_code",
		"barcode",
		"detailed_type",
		"list_price",
		"standard_price",
		"description_sale",
		"description_pickingin",
		"description_purchase",
		"taxes_id",
		"supplier_taxes_id",
	}
	o.Log.Infow(umdl, "header", header)

	// var recs = make([]interface{}, 0, o.BatchSize)
	// fmt.Println("recs", len(recs), cap(recs), recs)

	taxSell := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for sales - 5%"}})
	taxPurchase := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for purchases - 5%"}})

	pgs := o.ProductCategoryMap()
	// uom := o.UomMapper()

	// load records
	fmt.Println("load records", len(lList))
	bar := progressbar.Default(int64(len(lList)))

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(len(lList))
	for _, dr := range dbrecs {
		if _, ok := lList[dr.DefaultCode]; ok {
			go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, dr Product) {
				defer bar.Add(1)
				defer wg.Done()
				sem <- 1

				categID := -1
				if dr.Matgrp != "" {
					categID = pgs[dr.Matgrp]
				}

				ur := map[string]interface{}{
					"name":                  dr.Name,
					"default_code":          dr.DefaultCode,
					"barcode":               dr.DefaultCode,
					"type":                  dr.DetailedType,
					"list_price":            dr.ListPrice,
					"standard_price":        dr.StandardPrice,
					"description_sale":      dr.DescriptionSale,
					"description_pickingin": dr.DescriptionPickingin,
					"description_purchase":  dr.DescriptionPurchase,
					"taxes_id":              []int{taxSell},
					"supplier_taxes_id":     []int{taxPurchase},
				}

				if categID != -1 {
					ur["categ_id"] = categID
				}

				out, err := o.Create(umdl, ur)
				if err != nil {
					o.Log.Errorw(umdl, "out", out, "err", err)
				}

				<-sem
			}(sem, &wg, bar, dr)
		}
	}
	wg.Wait()

	// update records
	fmt.Println("update records", len(uList))

	bar = progressbar.Default(int64(len(uList)))
	wg.Add(len(uList))
	for _, dr := range dbrecs {
		if _, ok := uList[dr.DefaultCode]; ok {
			go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, dr Product) {
				defer bar.Add(1)
				defer wg.Done()
				sem <- 1

				r := uList[dr.DefaultCode]

				categID := -1
				if dr.Matgrp != "" {
					categID = pgs[dr.Matgrp]
				}

				ur := map[string]interface{}{
					"name":                  dr.Name,
					"default_code":          dr.DefaultCode,
					"barcode":               dr.DefaultCode,
					"type":                  dr.DetailedType,
					"list_price":            dr.ListPrice,
					"standard_price":        dr.StandardPrice,
					"description_sale":      dr.DescriptionSale,
					"description_pickingin": dr.DescriptionPickingin,
					"description_purchase":  dr.DescriptionPurchase,
					"taxes_id":              []int{taxSell},
					"supplier_taxes_id":     []int{taxPurchase},
				}

				if categID != -1 {
					ur["categ_id"] = categID
				}

				o.Log.Infow(umdl, "ur", ur, "r", r)

				out, err := o.Update(umdl, r, ur)
				if err != nil {
					o.Log.Errorw(umdl, "out", out, "err", err)
				}

				<-sem
			}(sem, &wg, bar, dr)
		}
	}
	wg.Wait()
}

func (o *OdooConn) ProductTemplate6() {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplate\n", umdl)

	type Product struct {
		Name                 string  `db:"name"`
		DefaultCode          string  `db:"default_code"`
		Barcode              string  `db:"barcode"`
		DetailedType         string  `db:"detailed_type"`
		ListPrice            float64 `db:"list_price"`
		StandardPrice        float64 `db:"standard_price"`
		DescriptionSale      string  `db:"description"`
		DescriptionPickingin string  `db:"description_pickingin"`
		DescriptionPurchase  string  `db:"description_purchase"`
		Category             string  `db:"categ_id"`
		Matgrp               string  `db:"matgrp"`
	}
	var dbrecs []Product
	stmt := `select distinct
	"name",default_code,default_code barcode,ptype detailed_type
	,list_price , standard_price
	,basic description,inspection description_pickingin,purchase description_purchase
	,category categ_id,matgrp
	from odoo.product_template2
	order by matgrp,default_code
	limit 10000
	`
	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	odoorecs := o.SearchRead(umdl, oarg{}, 0, 0, []string{"default_code"})

	odooList := make(map[string]int)
	for _, r := range odoorecs {
		switch r["default_code"].(type) {
		case string:
			odooList[r["default_code"].(string)] = int(r["id"].(float64))
		}
	}

	uList := make(map[string]int)
	lList := make(map[string]struct{})

	for _, dr := range dbrecs {
		if _, ok := odooList[dr.DefaultCode]; ok {
			uList[dr.DefaultCode] = odooList[dr.DefaultCode]
		} else {
			lList[dr.DefaultCode] = struct{}{}
		}
	}

	taxSell := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for sales - 5%"}})
	taxPurchase := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for purchases - 5%"}})

	pgs := o.ProductCategoryMap()
	// uom := o.UomMapper()

	// load records
	fmt.Println("load records", len(lList))

	type PRecord struct {
		ID     int
		Record map[string]interface{}
	}

	out1 := make(chan PRecord)
	out2 := make(chan PRecord)

	go func() {
		defer close(out1)
		for _, dr := range dbrecs {
			if _, ok := lList[dr.DefaultCode]; ok {
				categID := -1
				if dr.Matgrp != "" {
					categID = pgs[dr.Matgrp]
				}
				rec := PRecord{
					ID: -1,
				}

				ur := map[string]interface{}{
					"name":                  dr.Name,
					"default_code":          dr.DefaultCode,
					"barcode":               dr.DefaultCode,
					"type":                  dr.DetailedType,
					"list_price":            dr.ListPrice,
					"standard_price":        dr.StandardPrice,
					"description_sale":      dr.DescriptionSale,
					"description_pickingin": dr.DescriptionPickingin,
					"description_purchase":  dr.DescriptionPurchase,
					"taxes_id":              []int{taxSell},
					"supplier_taxes_id":     []int{taxPurchase},
				}

				if categID != -1 {
					ur["categ_id"] = categID
				}

				rec.Record = ur
				out1 <- rec
			}
		}
	}()

	go func() {
		defer close(out2)
		rec := PRecord{
			ID: -1,
		}
		out2 <- rec
	}()

	// for {
	// 	select {
	// 	case msg := <-out1:
	// 		fmt.Println("out1", msg)
	// 	case msg := <-out2:
	// 		fmt.Println("out2", msg)
	// 	}
	// }
	li := 0
	for msg := range out1 {
		li++
		fmt.Println("out1", li, msg)
		o.Create(umdl, msg.Record)
	}
	for msg := range out2 {
		fmt.Println("out2", msg)
	}

	// update records

	// bar = progressbar.Default(int64(len(uList)))
	// for _, dr := range dbrecs {
	// 	if _, ok := uList[dr.DefaultCode]; ok {
	// 		r := uList[dr.DefaultCode]

	// 		categID := -1
	// 		if dr.Matgrp != "" {
	// 			categID = pgs[dr.Matgrp]
	// 		}

	// 		var ur = map[string]interface{}{
	// 			"name":                  dr.Name,
	// 			"default_code":          dr.DefaultCode,
	// 			"barcode":               dr.DefaultCode,
	// 			"type":                  dr.DetailedType,
	// 			"list_price":            dr.ListPrice,
	// 			"standard_price":        dr.StandardPrice,
	// 			"description_sale":      dr.DescriptionSale,
	// 			"description_pickingin": dr.DescriptionPickingin,
	// 			"description_purchase":  dr.DescriptionPurchase,
	// 			"taxes_id":              []int{taxSell},
	// 			"supplier_taxes_id":     []int{taxPurchase},
	// 		}

	// 		if categID != -1 {
	// 			ur["categ_id"] = categID
	// 		}

	// 		o.Log.Infow(umdl, "ur", ur, "r", r)

	// 		o.Record(umdl, r, ur)

	// 		bar.Add(1)
	// 	}
	// }
}

// ProductTemplate function

func pager(batch, total int) {
	batchTotal := math.Ceil(float64(total) / float64(batch))
	fmt.Println("batchTotal", batchTotal)
}

// Create record
// func (o *Odoo) Load(model string, header []string, records []interface{}) (out int, err error) {
// 	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "load", header, records)
// 	if err != nil {
// 		return -1, err
// 	}
// 	// fmt.Printf("\n\n Create: %v", v)
// 	switch v := v.(type) {
// 	case float64:
// 		out = int(v)
// 	case interface{}:
// 		// ids := v.(map[string]interface{})["ids"].([]interface{})
// 		// message := v.(map[string]interface{})["message"].([]interface{})
// 		// nextrow := v.(map[string]interface{})["nextrow"].(int)
// 		// if len(message) != 0 {
// 		// 	err = fmt.Errorf("create record error model: %s message: %s record: %v", model, message, ids)
// 		// }
// 		// fmt.Println(ids, message)
// 		out = 0
// 	default:
// 		out = -1
// 	}
// 	return
// }
