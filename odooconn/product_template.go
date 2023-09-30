package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ProductTemplate function
func (o *OdooConn) ProductTemplate(ptfilt string) {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplate3 %v\n", umdl, ptfilt)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `select
	default_code,"name",sequence_number,category
	,group1,group2,group3,matgrp,basic,inspection,purchase
	,mtart,ptype,itype,mrp_rule,uom,industry_description
	,weight_gross,weight_net,weight_uom,standard_price,list_price,sell_per
	from odoo.product_template where mtart = $1
	`

	type Product struct {
		DefaultCode          string  `db:"default_code"`
		Name                 string  `db:"name"`
		SortCode             string  `db:"sequence_number"`
		Category             string  `db:"category"`
		Group1               string  `db:"group1"`
		Group2               string  `db:"group2"`
		Group3               string  `db:"group3"`
		Matgrp               string  `db:"matgrp"`
		DescriptionSale      string  `db:"basic"`
		DescriptionPickingin string  `db:"inspection"`
		DescriptionPurchase  string  `db:"purchase"`
		Mtart                string  `db:"mtart"`
		Ptype                string  `db:"ptype"`
		Itype                string  `db:"itype"`
		MRPRule              string  `db:"mrp_rule"`
		UOM                  string  `db:"uom"`
		IndustryDescription  string  `db:"industry_description"`
		WeightGross          float64 `db:"weight_gross"`
		WeightNet            string  `db:"weight_net"`
		WeightUOM            string  `db:"weight_uom"`
		StandardPrice        float64 `db:"standard_price"`
		ListPrice            float64 `db:"list_price"`
		SellPer              float64 `db:"sell_per"`
	}
	var rr []Product
	if stmt == "" {
		return
	}
	o.Log.Info(stmt)
	err := o.DB.Select(&rr, stmt, ptfilt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	taxSell, err := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for sales - 5%"}})
	o.checkErr(err)
	taxPurchase, err := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for purchases - 5%"}})
	o.checkErr(err)
	
	pgs := o.ProductCategoryMap()
	uuom := o.UomUomMap()

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Product) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			r, err := o.GetID(umdl, oarg{oarg{"default_code", "=", v.DefaultCode}})
			o.checkErr(err)

			categID := -1
			if v.Matgrp != "" {
				categID = pgs[v.Matgrp]
			}

			uomID := uuom[v.UOM]
			// v.StandardPrice

			ur := map[string]interface{}{
				"name":                  v.Name,
				"default_code":          v.DefaultCode,
				"barcode":               v.DefaultCode,
				"type":                  v.Ptype,
				"list_price":            v.ListPrice,
				"standard_price":        v.StandardPrice,
				"description_sale":      v.DescriptionSale,
				"description_pickingin": v.DescriptionPickingin,
				"description_purchase":  v.DescriptionPurchase,
				"uom_id":                uomID,
				"uom_po_id":             uomID,
				"taxes_id":              []int{taxSell},
				"supplier_taxes_id":     []int{taxPurchase},
			}

			if categID != -1 {
				ur["categ_id"] = categID
			}
			if v.Itype != "" {
				ur["inventory_type"] = v.Itype
			}
			if v.MRPRule != "" {
				ur["mrp_type"] = v.MRPRule
			}
			if v.SortCode != "" {
				ur["sequence_number"] = v.SortCode
			}
			if uomID == -2 {
				o.Log.Info(umdl)
			}

			o.Log.Info(umdl, "ur", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ProductTemplateDelpro function
func (o *OdooConn) ProductTemplateDelpro() {
	mdl := "product_template"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplateDelpro\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	saleableID, err := o.GetID("product.category", oarg{oarg{"name", "=", "Saleable"}})
	o.checkErr(err)
	
	// saleable
	stmt := `
	select distinct 
	dit.company,
	case when dit.company = 'DelPro Technical Inc.' then 'DT' else 'DA' end prefix,
	dit.name,
	dit.description ,
	trim(trim(dit.description)||' '||trim(dit.description_sale)) description_sale,
	dit.sell_uom,
	dit.category,
	dit.standard_price,
	case 
	when dpc.ptype = 'NONE' then 'MISCELLANEOUS PRODUCTS' 
	when dpc.ptype = 'Regulators/Relief' then 'Regulators and Relief'
	when dpc.ptype = 'Tubing/Fittings' then 'Tubing and Fittings'
	when dpc.ptype is null then ''
	else dpc.ptype end ptype
	,case when dpc.brand is null then '' else dpc.brand end brand
	from odoo.delpro_items_trimmed dit 
	left join odoo.delpro_product_categories dpc on dit.description = dpc.description
	order by dit.company,dit.name
	limit 10
	`
	type Product struct {
		Company         string  `db:"company"`
		Prefix          string  `db:"prefix"`
		Name            string  `db:"name"`
		Description     string  `db:"description"`
		DescriptionSale string  `db:"description_sale"`
		UOM             string  `db:"sell_uom"`
		Category        string  `db:"category"`
		StandardPrice   float64 `db:"standard_price"`
		Ptype           string  `db:"ptype"`
		Brand           string  `db:"brand"`
	}
	var rr []Product
	err = o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Product) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			companyID := cids[v.Company]
			o.Log.Info("company", "v.Company", v.Company, "companyID", companyID)

			pid, err := o.GetID("product.category", oarg{oarg{"name", "=", v.Ptype}, oarg{"parent_id", "=", saleableID}})
			o.checkErr(err)
			gid, err := o.GetID("product.category", oarg{oarg{"name", "=", v.Brand}, oarg{"parent_id", "=", pid}})
			o.checkErr(err)
			if v.Ptype == "" {
				gid = saleableID
			}

			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"default_code", "like", v.Name}, oarg{"company_id", "=", companyID}})
			o.checkErr(err)
			uomID, err := o.GetID("uom.uom", oarg{oarg{"name", "=", "ea"}})
			o.checkErr(err)
			ptype := "product"
			// v.Prefix + v.Name

			ur := map[string]interface{}{
				"name":             v.Name,
				"company_id":       companyID,
				"default_code":     "",
				"barcode":          "",
				"type":             ptype,
				"standard_price":   v.StandardPrice,
				"description":      v.Description,
				"description_sale": v.DescriptionSale,
				"uom_id":           uomID,
				"uom_po_id":        uomID,
				"categ_id":         gid,
				"inventory_type":   "finishedgoods",
				"mrp_type":         "ondemand",
			}

			// if categID != -1 {
			// 	ur["categ_id"] = categID
			// }
			// if v.Mattype != "" {
			// 	ur["inventory_type"] = v.Mattype
			// }
			// if v.MRPRule != "" {
			// 	ur["mrp_type"] = v.MRPRule
			// }

			o.Log.Info(umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
