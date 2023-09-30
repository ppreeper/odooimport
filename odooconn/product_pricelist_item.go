package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) PricelistPricegroupMatGrpDiscounts() {
	// pl_pg_mg
	mdl := "product_pricelist_item"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select
	left(c.vakey,4) salesorg
	,right(left(c.vakey,6),2) saleschan
	,right(c.vakey,4) matkl
	,"right"("left"(c.vakey::text, 8), 2) as pricegroup
	,"right"("left"(c.vakey::text, 8), 2)||' '||btrim((replace(t188t.vtext::text, 'ARTG-'::text, ''::text) || ' '::text) || CASE WHEN c.kwaeh::text <> 'CAD'::text THEN c.kwaeh ELSE ''::character varying END::text) as "name"
	,kbetr/-10.0 as discount
	,datab::date::varchar date_start
	,datbi::date::varchar date_end
	,kwaeh as currency
	from sapdata.raw_conditions c
	JOIN sap.t188t t188t ON "right"("left"(c.vakey::text, 8), 2) = t188t.konda::text
	WHERE c.kschl::text = 'ZCMG'::text AND "left"(c.vakey::text, 4) <> '1010'::text;
	`
	type Pricelist struct {
		Salesorg   string  `db:"salesorg"`
		Saleschan  string  `db:"saleschan"`
		Matkl      string  `db:"matkl"`
		Pricegroup string  `db:"pricegroup"`
		Name       string  `db:"name"`
		Discount   float64 `db:"discount"`
		DateStart  string  `db:"date_start"`
		DateEnd    string  `db:"date_end"`
		Currency   string  `db:"currency"`
	}
	var pp []Pricelist
	err := o.DB.Select(&pp, stmt)
	o.checkErr(err)
	recs := len(pp)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, p := range pp {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, p Pricelist) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			curid, err := o.GetID("res.currency", oarg{oarg{"name", "=", p.Currency}})
			o.checkErr(err)
			listID, err := o.GetID("product.pricelist", oarg{oarg{"name", "=", p.Name}, oarg{"currency_id", "=", curid}})
			o.checkErr(err)
			categID, err := o.GetID("product.category", oarg{oarg{"name", "like", p.Matkl}})
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"categ_id", "=", categID}, oarg{"pricelist_id", "=", listID}, oarg{"currency_id", "=", curid}})
			o.checkErr(err)
			ur := map[string]interface{}{
				"categ_id":      categID,
				"applied_on":    "2_product_category",
				"base":          "list_price",
				"compute_price": "percentage",
				"pricelist_id":  listID,
				"percent_price": p.Discount,
				"date_start":    p.DateStart,
				"date_end":      p.DateEnd,
			}
			// o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

			o.Log.Info(umdl, "ur", ur, "rid", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, p)
	}
	wg.Wait()
}

func (o *OdooConn) PricelistCustomerDefault() {
	// pl_cust_def
	mdl := "product_pricelist_item"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v PricelistCustomerDefault\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	stmt := `
	select distinct
	c.kunnr,c.currency
	,cust.cname,cust.pricinggroup
	,app.pricelist
	from (
	-- kunnr in both 1000 1020
	select distinct a.kunnr,a.currency from (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	,c.kwaeh as currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1000'
	)a
	join (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	,c.kwaeh as currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1020'
	) g on a.kunnr = g.kunnr
	union
	-- kunnr exclusive 1000
	select distinct a.kunnr,a.currency from (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	,c.kwaeh as currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1000'
	)a
	left join (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	,c.kwaeh as currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1020' ) g on a.kunnr = g.kunnr
	where g.kunnr is null
	union
	-- kunnr exclusive 1020
	select distinct g.kunnr,g.currency from (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	,c.kwaeh as currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1020'
	)g
	left join (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	,c.kwaeh as currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1000' ) a on g.kunnr = a.kunnr
	where a.kunnr is null
	) c
	JOIN (
	SELECT DISTINCT c_1.kunnr,pricinggroup,trim(trim(name1)||' '||trim(name2))||' '||right(c_1.kunnr, 7) AS cname
	FROM sapdata.customers c_1
	LEFT JOIN odoo.artg_customers_dnu dnu ON c_1.kunnr = dnu.kunnr
	WHERE dnu.kunnr IS null
	) cust ON c.kunnr = cust.kunnr
	join (
	select distinct
	"right"("left"(c.vakey, 8), 2) as pricegroup
	,"right"("left"(c.vakey, 8), 2)||' '||btrim((replace(t188t.vtext, 'ARTG-', '') || ' ') || CASE WHEN c.kwaeh <> 'CAD' THEN c.kwaeh ELSE ''::character varying END) as pricelist
	,c.kwaeh as currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	JOIN sap.t188t t188t ON "right"("left"(c.vakey, 8), 2) = t188t.konda
	join (select distinct pricinggroup from sapdata.customerstemp c where trim(pricinggroup) <> '') cp on "right"("left"(c.vakey, 8), 2) = cp.pricinggroup
	WHERE c.kschl = 'ZCMG' AND "left"(c.vakey, 4) <> '1010'
	) app on cust.pricinggroup = app.pricegroup and c.currency = app.currency
	`

	type Pricelist struct {
		Ref              string `db:"kunnr"`
		Currency         string `db:"currency"`
		Name             string `db:"cname"`
		PricingGroup     string `db:"pricinggroup"`
		DefaultPricelist string `db:"pricelist"`
	}
	var pp []Pricelist
	err := o.DB.Select(&pp, stmt)
	o.checkErr(err)
	recs := len(pp)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, p := range pp {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, p Pricelist) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			curid, err := o.GetID("res.currency", oarg{oarg{"name", "=", p.Currency}})
			o.checkErr(err)
			listID, err := o.GetID("product.pricelist", oarg{oarg{"name", "=", p.Name}, oarg{"currency_id", "=", curid}})
			o.checkErr(err)
			defaultListID, err := o.GetID("product.pricelist", oarg{oarg{"name", "=", p.DefaultPricelist}, oarg{"currency_id", "=", curid}})
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"base_pricelist_id", "=", defaultListID}, oarg{"pricelist_id", "=", listID}, oarg{"currency_id", "=", curid}})
			o.checkErr(err)

			ur := map[string]interface{}{
				"applied_on":        "3_global",
				"base":              "pricelist",
				"base_pricelist_id": defaultListID,
				"pricelist_id":      listID,
				"compute_price":     "formula",
			}
			o.Log.Info(umdl, "ur", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, p)
	}
	wg.Wait()
}

func (o *OdooConn) PricelistCustomerMatGrpDiscounts() {
	// pl_cust_mg
	mdl := "product_pricelist_item"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v PricelistCustomerMatGrpDiscounts\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select distinct
	cust.cname,
	"right"(c.vakey::text, 4) AS matkl,
	right("left"("right"(c.vakey::text, 14), 10),7) as kunnr,
	CASE WHEN c.konwa::text = '%'::text THEN c.kbetr / -10.0::numeric ELSE c.kbetr END AS discount,
	c.datab::date::varchar AS date_start,
	c.datbi::date::varchar AS date_end,
	c.kwaeh AS currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt::text = cm.mandt::text AND c.vakey::text = cm.vakey::text AND c.datab::text = cm.datab::text AND c.datbi::text = cm.datbi::text
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt::text = cm2.mandt::text AND c.vakey::text = cm2.vakey::text AND c.datab::text = cm2.datab::text AND c.datbi::text = cm2.datbi::text AND c.knumh::text = cm2.knumh::text
	JOIN (
	SELECT DISTINCT c_1.kunnr,
	trim(trim(name1)||' '||trim(name2))||' '||right(c_1.kunnr, 7) AS cname
	FROM sapdata.customers c_1
	LEFT JOIN odoo.artg_customers_dnu dnu ON c_1.kunnr::text = dnu.kunnr::text
	WHERE dnu.kunnr IS null
	) cust ON "left"("right"(c.vakey::text, 14), 10) = cust.kunnr
	WHERE c.kschl::text = 'ZMGC'::text AND "left"(c.vakey::text, 4) = '1000'::text
	order by cust.cname,"right"(c.vakey::text, 4)
	`

	type Pricelist struct {
		Name      string  `db:"cname"`
		CustNo    string  `db:"kunnr"`
		Matkl     string  `db:"matkl"`
		Discount  float64 `db:"discount"`
		DateStart string  `db:"date_start"`
		DateEnd   string  `db:"date_end"`
		Currency  string  `db:"currency"`
	}
	var pp []Pricelist
	err := o.DB.Select(&pp, stmt)
	o.checkErr(err)
	recs := len(pp)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, p := range pp {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, p Pricelist) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			curid, err := o.GetID("res.currency", oarg{oarg{"name", "=", p.Currency}})
			o.checkErr(err)
			listID, err := o.GetID("product.pricelist", oarg{oarg{"name", "=", p.Name}, oarg{"currency_id", "=", curid}})
			o.checkErr(err)
			categID, err := o.GetID("product.category", oarg{oarg{"name", "like", p.Matkl}})
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"categ_id", "=", categID}, oarg{"pricelist_id", "=", listID}, oarg{"currency_id", "=", curid}})
			o.checkErr(err)
			o.Log.Info("", "name", p.Name, "matkl", p.Matkl, "discount", p.Discount)
			ur := map[string]interface{}{
				"categ_id":      categID,
				"applied_on":    "2_product_category",
				"base":          "list_price",
				"compute_price": "percentage",
				"pricelist_id":  listID,
				"percent_price": p.Discount,
				"date_start":    p.DateStart,
				"date_end":      p.DateEnd,
			}
			// o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)
			o.Log.Info(umdl, "ur", ur, "rid", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, p)
	}
	wg.Wait()
}

func (o *OdooConn) PricelistCustomerNetoutItems() {
	// pl_cust_no
	mdl := "product_pricelist_item"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v PricelistCustomerNetoutItems\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	SELECT distinct
	cust.cname,
	"right"("left"(c.vakey::text, 16), 7) AS kunnr,
	"right"(c.vakey::text, 7) AS matnr,
	c.kbetr AS net_price,
	c.kpein AS per,
	c.datab::date::varchar AS date_start,
	c.datbi::date::varchar AS date_end,
	c.kwaeh AS currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt::text = cm.mandt::text AND c.vakey::text = cm.vakey::text AND c.datab::text = cm.datab::text AND c.datbi::text = cm.datbi::text
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt::text = cm2.mandt::text AND c.vakey::text = cm2.vakey::text AND c.datab::text = cm2.datab::text AND c.datbi::text = cm2.datbi::text AND c.knumh::text = cm2.knumh::text
	JOIN (
	SELECT DISTINCT c_1.kunnr,
	trim(trim(name1)||' '||trim(name2))||' '||right(c_1.kunnr, 7) AS cname
	FROM sapdata.customers c_1
	LEFT JOIN odoo.artg_customers_dnu dnu ON c_1.kunnr::text = dnu.kunnr::text
	WHERE dnu.kunnr IS null
	) cust ON "right"("left"(c.vakey::text, 16), 10) = cust.kunnr
	WHERE c.kotabnr::text = '005'::text AND c.kschl::text = 'ZPR0'::text AND "left"(c.vakey::text, 4) = '1000'::text
	order by cust.cname,"right"(c.vakey::text, 7)
	`

	type Pricelist struct {
		Name      string  `db:"cname"`
		CustNo    string  `db:"kunnr"`
		Matnr     string  `db:"matnr"`
		NetPrice  float64 `db:"net_price"`
		Per       string  `db:"per"`
		DateStart string  `db:"date_start"`
		DateEnd   string  `db:"date_end"`
		Currency  string  `db:"currency"`
	}
	var pp []Pricelist
	err := o.DB.Select(&pp, stmt)
	o.checkErr(err)
	recs := len(pp)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, p := range pp {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, p Pricelist) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			curID, err := o.GetID("res.currency", oarg{oarg{"name", "=", p.Currency}})
			o.checkErr(err)
			listID, err := o.GetID("product.pricelist", oarg{oarg{"name", "=", p.Name}, oarg{"currency_id", "=", curID}})
			o.checkErr(err)
			pid, err := o.GetID("product.template", oarg{oarg{"default_code", "=", p.Matnr}})
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"product_tmpl_id", "=", pid}, oarg{"pricelist_id", "=", listID}, oarg{"currency_id", "=", curID}})
			o.checkErr(err)
			ur := map[string]interface{}{
				"product_tmpl_id": pid,
				"applied_on":      "1_product",
				"pricelist_id":    listID,
				"currency_id":     curID,
				"date_start":      p.DateStart,
				"date_end":        p.DateEnd,
				"compute_price":   "fixed",
				"fixed_price":     p.NetPrice,
			}
			// o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)
			o.Log.Info(umdl, "ur", ur, "rid", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, p)
	}
	wg.Wait()
}

// ProductPricelistItemMatgroup function
func (o *OdooConn) ProductPricelistItemMatgroup() {
	mdl := "product_pricelist_item"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Println(umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select salesorg,company,pricegroup
	,"name",matkl,discount
	,date_start::date::varchar date_start
	,date_end::date::varchar date_end
	,currency
	from odoo.product_pricelist_item_matgroup
	order by company,pricegroup,matkl
	`
	type Pricelist struct {
		Salesorg   string  `db:"salesorg"`
		Company    string  `db:"company"`
		Pricegroup string  `db:"pricegroup"`
		Name       string  `db:"name"`
		Matkl      string  `db:"matkl"`
		Discount   float64 `db:"discount"`
		DateStart  string  `db:"date_start"`
		DateEnd    string  `db:"date_end"`
		Currency   string  `db:"currency"`
	}
	var rr []Pricelist
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Pricelist) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1
			// s := "select id,categ_id,applied_on,base,base_pricelist_id,pricelist_id,company_id,currency_id,date_start,date_end,percent_price from product_pricelist_item;"
			cid, err := o.CompanyID(v.Company)
			o.checkErr(err)
			curid, err := o.GetID("res.currency", oarg{oarg{"name", "=", v.Currency}})
			o.checkErr(err)
			listID, err := o.GetID("product.pricelist", oarg{oarg{"name", "=", v.Name}, oarg{"currency_id", "=", curid}, oarg{"company_id", "=", cid}})
			o.checkErr(err)
			categID, err := o.GetID("product.category", oarg{oarg{"name", "like", v.Matkl}})
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"categ_id", "=", categID}, oarg{"pricelist_id", "=", listID}, oarg{"currency_id", "=", curid}, oarg{"company_id", "=", cid}})
			o.checkErr(err)
			ur := map[string]interface{}{
				"categ_id":      categID,
				"applied_on":    "2_product_category",
				"base":          "list_price",
				"compute_price": "percentage",
				"pricelist_id":  listID,
				"company_id":    cid,
				"percent_price": v.Discount,
				"date_start":    v.DateStart,
				"date_end":      v.DateEnd,
			}
			o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
