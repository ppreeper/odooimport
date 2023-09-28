package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

////////////////////
// Pricelists
////////////////////

// PricelistPricegroup Pricelist Pricegroup
func (o *OdooConn) PricelistPricegroup() {
	// pl_pg
	mdl := "product_pricelist"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select distinct
	"right"("left"(c.vakey::text, 8), 2) as pricegroup
	,"right"("left"(c.vakey::text, 8), 2)||' '||btrim((replace(t188t.vtext::text, 'ARTG-'::text, ''::text) || ' '::text) || CASE WHEN c.kwaeh::text <> 'CAD'::text THEN c.kwaeh ELSE ''::character varying END::text) as "name"
	,c.kwaeh as currency
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt::text = cm.mandt::text AND c.vakey::text = cm.vakey::text AND c.datab::text = cm.datab::text AND c.datbi::text = cm.datbi::text
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt::text = cm2.mandt::text AND c.vakey::text = cm2.vakey::text AND c.datab::text = cm2.datab::text AND c.datbi::text = cm2.datbi::text AND c.knumh::text = cm2.knumh::text
	JOIN sap.t188t t188t ON "right"("left"(c.vakey::text, 8), 2) = t188t.konda::text
	join (select distinct pricinggroup from sapdata.customerstemp c where trim(pricinggroup) <> '') cp on "right"("left"(c.vakey::text, 8), 2) = cp.pricinggroup
	WHERE c.kschl::text = 'ZCMG'::text AND "left"(c.vakey::text, 4) <> '1010'::text;
	`
	type Pricelist struct {
		Pricegroup string `db:"pricegroup"`
		Name       string `db:"name"`
		Currency   string `db:"currency"`
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

			curid := o.GetID("res.currency", oarg{oarg{"name", "=", p.Currency}})
			r := o.GetID(umdl, oarg{oarg{"name", "=", p.Name}, oarg{"currency_id", "=", curid}})

			ur := map[string]interface{}{
				"name":        p.Name,
				"currency_id": curid,
			}

			o.Log.Infow(umdl, "ur", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, p)
	}
	wg.Wait()
}

// //////////////////
// Pricelist Customer Specific
// //////////////////
func (o *OdooConn) PricelistCustomer() {
	// pl_cust
	mdl := "product_pricelist"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v PricelistCustomer\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select distinct
	c.kunnr
	,cust.cname
	from (
	-- kunnr in both 1000 1020
	select distinct a.kunnr from (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1000'
	)a
	join (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1020'
	) g on a.kunnr = g.kunnr
	union
	-- kunnr exclusive 1000
	select distinct a.kunnr from (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1000'
	)a
	left join (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1020' ) g on a.kunnr = g.kunnr
	where g.kunnr is null
	union
	-- kunnr exclusive 1020
	select distinct g.kunnr from (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1020'
	)g
	left join (
	select distinct
	"left"("right"(c.vakey, 14), 10) AS kunnr
	FROM sapdata.raw_conditions c
	JOIN sapdata.conditionsmaxdate cm ON c.mandt = cm.mandt AND c.vakey = cm.vakey AND c.datab = cm.datab AND c.datbi = cm.datbi
	JOIN sapdata.conditionsmaxdocno cm2 ON c.mandt = cm2.mandt AND c.vakey = cm2.vakey AND c.datab = cm2.datab AND c.datbi = cm2.datbi AND c.knumh = cm2.knumh
	WHERE c.kschl = 'ZMGC' AND "left"(c.vakey, 4) = '1000' ) a on g.kunnr = a.kunnr
	where a.kunnr is null
	) c
	JOIN (
	SELECT DISTINCT c_1.kunnr
	,trim(trim(name1)||' '||trim(name2))||' '||right(c_1.kunnr, 7) AS cname
	FROM sapdata.customers c_1
	LEFT JOIN odoo.artg_customers_dnu dnu ON c_1.kunnr = dnu.kunnr
	WHERE dnu.kunnr IS null
	) cust ON c.kunnr = cust.kunnr
	order by name,kunnr
	`

	type Pricelist struct {
		CustNo string `db:"kunnr"`
		Name   string `db:"cname"`
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

			r := o.GetID(umdl, oarg{oarg{"name", "like", p.Name}})

			ur := map[string]interface{}{
				"name": p.Name,
			}

			o.Log.Infow(umdl, "ur", ur, "rid", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, p)
	}
	wg.Wait()
}

// ProductPricelist function
func (o *OdooConn) ProductPricelist() {
	mdl := "product_pricelist"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select salesorg,company,pricegroup code,name,currency from odoo.product_pricelist
	`
	type Pricelist struct {
		Salesorg string `db:"salesorg"`
		Company  string `db:"company"`
		Code     string `db:"code"`
		Name     string `db:"name"`
		Currency string `db:"currency"`
	}
	var rr []Pricelist
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	curr := o.ResCurrencyMap()

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Pricelist) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1
			cid := o.CompanyID(v.Company)
			// curid := o.GetID("res.currency", oarg{oarg{"name", "=", v.Currency}})
			curid := curr[v.Currency]
			r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"currency_id", "=", curid}, oarg{"company_id", "=", cid}})
			ur := map[string]interface{}{
				"name":        v.Name,
				"currency_id": curid,
				"company_id":  cid,
			}
			o.Log.Infow(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

func (o *OdooConn) ProductPricelistMap() map[string]int {
	mdl := "product_pricelist"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc := o.SearchRead(umdl, oarg{}, 0, 0, []string{"name"})
	cids := map[string]int{}
	for _, c := range cc {
		cids[c["name"].(string)] = int(c["id"].(float64))
	}
	return cids
}
