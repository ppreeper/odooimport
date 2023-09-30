package odooconn

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// ProductSupplierinfo function
func (o *OdooConn) ProductSupplierinfo(ptfilt string) {
	mdl := "product_supplierinfo"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Println(umdl)

	stmt := `
	select
	right(avpl.vendorno,6) as vendorno
	,right(avpl.matnr,7) matnr
	,avpl.min_qty 
	,avpl.per
	,avpl.puom
	,avpl.price
	,avpl.currency
	,avpl.date_start
	,avpl.date_end
	,avpl.delay
	from sapdata.materials m
	left join odoo.mat_dnu dnu on m.matnr = dnu.matnr 
	join odoo.artg_vendor_pricelisttemp avpl on m.matnr = avpl.matnr
	where dnu.matnr is null and mtart = $1
	order by vendorno,matnr,date_start
	-- limit 10
	`
	type VPL struct {
		Vendorno  string  `db:"vendorno"`
		Matnr     string  `db:"matnr"`
		MinQty    float64 `db:"min_qty"`
		Per       float64 `db:"per"`
		UOM       string  `db:"puom"`
		Price     float64 `db:"price"`
		Currency  string  `db:"currency"`
		DateStart string  `db:"date_start"`
		DateEnd   string  `db:"date_end"`
		Delay     int     `db:"delay"`
	}
	var rr []VPL
	err := o.DB.Select(&rr, stmt, ptfilt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()
	companyID := -1
	if ptfilt == "ROH" {
		companyID = cids["A.R. Thomson Group Manufacturing"]
	} else if ptfilt == "HALB" {
		companyID = cids["A.R. Thomson Group Manufacturing"]
	} else {
		companyID = cids["A.R. Thomson Group"]
	}
	o.Log.Info("", "companyID", companyID, "recs", recs)

	curr := o.ResCurrencyMap()

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v VPL) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Log.Info(umdl, "v", v)

			// currencyID := o.GetID("res.currency", oarg{oarg{"name", "=", v.Currency}})
			currencyID := curr[v.Currency]

			t, err := time.Parse("20060102", v.DateStart)
			o.checkErr(err)
			dateStart := t.Format("2006-01-02")
			t, err = time.Parse("20060102", v.DateEnd)
			o.checkErr(err)
			dateEnd := t.Format("2006-01-02")

			vendorID, err := o.GetID("res.partner", oarg{oarg{"ref", "=", v.Vendorno}})
			o.checkErr(err)

			// vendor := o.SearchRead("res.partner", oarg{oarg{"ref", "=", v.Vendorno}}, 0, 0, []string{"name"})
			// vendorName := vendor[0]["name"].(string)
			// // fmt.Println(vendorName)
			// o.Log.Info("vendorName:", vendorName)

			prodTmplID, err := o.GetID("product.template", oarg{oarg{"default_code", "=", v.Matnr}})
			o.checkErr(err)

			pUOM, err := o.GetID("uom.uom", oarg{oarg{"name", "=", v.UOM}})
			o.checkErr(err)

			r, err := o.GetID(umdl,
				oarg{
					oarg{"currency_id", "=", currencyID},
					oarg{"date_end", "=", dateEnd},
					oarg{"date_start", "=", dateStart},
					oarg{"name", "=", vendorID},
					oarg{"product_tmpl_id", "=", prodTmplID},
					// oarg{"product_uom", "=", pUOM},
				},
			)
			o.checkErr(err)

			ur := map[string]interface{}{
				"company_id":      companyID,
				"currency_id":     currencyID,
				"date_end":        dateEnd,
				"date_start":      dateStart,
				"delay":           v.Delay,
				"min_qty":         v.MinQty,
				"name":            vendorID,
				"price":           v.Price,
				"product_tmpl_id": prodTmplID,
				"product_uom":     pUOM,
			}
			o.Log.Info(umdl, "ur", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ProductSupplierinfoTopParts function
func (o *OdooConn) ProductSupplierinfoTopParts() {
	mdl := "product_supplierinfo"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Println(umdl)

	// id,name,product_name,product_code,sequence,min_qty,price
	// ,company_id,currency_id,date_start,date_end,product_id,product_tmpl_id,delay
	// ,create_uid,create_date,write_uid,write_date
	// Purchase Info Records ?? links product to supplier
	stmt := `
	select
	right(avpl.vendorno,6) vendorno
	,avpl.min_qty
	,avpl.price
	,avpl.currency
	,avpl.date_start
	,avpl.date_end
	,avpl.matnr
	,avpl.delay
	from odoo.artg_vendor_pricelist avpl
	join odoo.artg_vendors_active ava on avpl.vendorno = ava.vendorno
	join odoo.artg_topparts tp on avpl.matnr = tp."ref"
	`
	rr := []struct {
		Vendorno  string  `db:"vendorno"`
		MinQty    float64 `db:"min_qty"`
		Price     float64 `db:"price"`
		Currency  string  `db:"currency"`
		DateStart string  `db:"date_start"`
		DateEnd   string  `db:"date_end"`
		Matnr     string  `db:"matnr"`
		Delay     int     `db:"delay"`
	}{}
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))
	for _, v := range rr {
		err := bar.Add(1)
		o.checkErr(err)

		vendorID, err := o.GetID("res.partner", oarg{oarg{"ref", "=", v.Vendorno}})
		o.checkErr(err)
		currencyID, err := o.GetID("res.currency", oarg{oarg{"name", "=", v.Currency}})
		o.checkErr(err)
		prodTmplID, err := o.GetID("product.template", oarg{oarg{"default_code", "=", v.Matnr}})
		o.checkErr(err)
		r, err := o.GetID(umdl, oarg{oarg{"name", "=", vendorID}, oarg{"currency_id", "=", currencyID}, oarg{"product_tmpl_id", "=", prodTmplID}})
		o.checkErr(err)
		t, err := time.Parse("20060102", v.DateStart)
		o.checkErr(err)
		dateStart := t.Format("2006-01-02")
		t, err = time.Parse("20060102", v.DateEnd)
		o.checkErr(err)
		dateEnd := t.Format("2006-01-02")

		ur := map[string]interface{}{
			"name":            vendorID,
			"min_qty":         v.MinQty,
			"date_start":      dateStart,
			"date_end":        dateEnd,
			"price":           v.Price,
			"currency_id":     currencyID,
			"product_tmpl_id": prodTmplID,
			"delay":           v.Delay,
		}
		o.Log.Info(mdl, "model", umdl, "record", ur)

		o.Record(umdl, r, ur)
	}
}
