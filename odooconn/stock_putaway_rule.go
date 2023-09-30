package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// StockPutawayRule Putaway Rules is what BIN location for product in Which Plant
func (o *OdooConn) StockPutawayRule() {
	mdl := "stock_putaway_rule"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v StockPutawayRule\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select distinct company,cityname,right(matnr,7) matnr,location,sbin from (
	select distinct 
	p.company 
	,p.cityname 
	,i.matnr 
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end as location
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end||'/'||i.sbin as sbin
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '1020'
	and p.werks not in ('1033','1013')
	and i.unrestricted <> '0'
	and trim(i.sbin) <> ''
	union
	select distinct 
	p.company 
	,p.cityname 
	,i.matnr 
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end as location
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end||'/'||i.sbin as sbin
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '1030'
	and p.werks not in ('1033','1013')
	and i.unrestricted <> '0'
	and trim(i.sbin) <> ''
	union
	select distinct 
	p.company 
	,p.cityname 
	,i.matnr 
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end as location
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end||'/'||i.sbin as sbin
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '1000'
	and p.werks not in ('1033','1013')
	and i.unrestricted <> '0'
	and trim(i.sbin) <> ''
	) l
	-- where l.company = $1
	order by company,location,sbin,matnr
	`
	type StockLocation struct {
		Company     string `db:"company"`
		City        string `db:"cityname"`
		DefaultCode string `db:"matnr"`
		Location    string `db:"location"`
		SBin        string `db:"sbin"`
	}
	var rr []StockLocation
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()

	// sls := o.StockLocationMap()

	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v StockLocation) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid := cids[v.Company]
			pTmplID, err := o.GetID("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}})
			o.checkErr(err)
			pID, err := o.GetID("product.product", oarg{oarg{"product_tmpl_id", "=", pTmplID}})
			o.checkErr(err)
			locInID, err := o.GetID("stock.location", oarg{oarg{"complete_name", "=", v.Location}})
			o.checkErr(err)
			locOutID, err := o.GetID("stock.location", oarg{oarg{"complete_name", "=", v.SBin}})
			o.checkErr(err)

			r, err := o.GetID(umdl,
				oarg{
					oarg{"product_id", "=", pID},
					oarg{"location_in_id", "=", locInID},
					oarg{"location_out_id", "=", locOutID},
					oarg{"company_id", "=", cid},
				},
			)
			o.checkErr(err)

			ur := map[string]interface{}{
				"company_id":      cid,
				"location_in_id":  locInID,
				"location_out_id": locOutID,
				"product_id":      pID,
			}

			o.Log.Info(mdl, "record", ur, "r", r, "v", v)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
