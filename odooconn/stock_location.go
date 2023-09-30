package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) StockLocation() {
	mdl := "stock_location"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v StockLocation\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	companyIDs := o.ResCompanyMap()

	stmt := `
	select distinct company,cityname,warehouse,location from (
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,prefix||citycode warehouse
	,case when i.sloc = '0001' then 'Stock' else i.sloc end as location
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '1020'
	and p.werks not in ('1033','1013')
	and trim(sbin) = ''
	and i.unrestricted <> '0'
	union
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,prefix||citycode warehouse
	,case when i.sloc = '0001' then 'Stock' else i.sloc end as location
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '1000'
	and p.werks not in ('1033','1013')
	and trim(sbin) = ''
	and i.unrestricted <> '0'
	union 
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,prefix||citycode warehouse
	,case when i.sloc = '0001' then 'Stock' else i.sloc end as location
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '2000'
	and p.werks not in ('1033','1013')
	and trim(sbin) = ''
	and i.unrestricted <> '0'
	union
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,prefix||citycode warehouse
	,case when i.sloc = '0001' then 'Stock' else i.sloc end as location
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '2010'
	and p.werks not in ('1033','1013')
	and trim(sbin) = ''
	and i.unrestricted <> '0'
	) l
	-- where l.company = $1
	order by company,warehouse,location
	`
	type WarehouseLocation struct {
		Company   string `db:"company"`
		City      string `db:"cityname"`
		Warehouse string `db:"warehouse"`
		Location  string `db:"location"`
	}
	var ww []WarehouseLocation
	err := o.DB.Select(&ww, stmt)
	o.checkErr(err)
	recs := len(ww)
	bar := progressbar.Default(int64(recs))

	for _, v := range ww {
		err := bar.Add(1)
		o.checkErr(err)
		cid := companyIDs[v.Company]
		name := v.Location
		pid, err := o.GetID(umdl, oarg{oarg{"complete_name", "=", v.Warehouse}, oarg{"company_id", "=", cid}})
		o.checkErr(err)
		r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Location}, oarg{"location_id", "=", pid}, oarg{"company_id", "=", cid}})
		o.checkErr(err)
		ur := map[string]interface{}{
			"name":        name,
			"location_id": pid,
			"company_id":  cid,
		}
		o.Log.Info(mdl, "record", ur, "r", r, "v", v)
	}

	stmt = `
	select distinct company,cityname,location,sbin from (
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end as location
	,i.sbin 
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '1020'
	and p.werks not in ('1033','1013')
	and i.unrestricted <> '0'
	and trim(i.sbin) <> ''
	union
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end as location
	,i.sbin 
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '1030'
	and p.werks not in ('1033','1013')
	and i.unrestricted <> '0'
	and trim(i.sbin) <> ''
	union
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end as location
	,i.sbin 
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '1000'
	and p.werks not in ('1033','1013')
	and i.unrestricted <> '0'
	and trim(i.sbin) <> ''
	union	
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end as location
	,i.sbin 
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '2000'
	and p.werks not in ('1033','1013')
	and i.unrestricted <> '0'
	and trim(i.sbin) <> ''
	union	
	select distinct 
	i.werks
	,p.company 
	,p.cityname 
	,case when i.sloc = '0001' then prefix||citycode||'/Stock' else prefix||citycode||'/'||i.sloc end as location
	,i.sbin 
	from sapdata.inventorytemp i
	join odoo.plantstemp p on i.werks = p.werks
	where p.salesorg = '2010'
	and p.werks not in ('1033','1013')
	and i.unrestricted <> '0'
	and trim(i.sbin) <> ''
	) l
	-- where l.company = $1
	order by company,location,sbin
	`
	type StockLocation struct {
		Company  string `db:"company"`
		City     string `db:"cityname"`
		Location string `db:"location"`
		SBin     string `db:"sbin"`
	}
	var rr []StockLocation
	err = o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs = len(rr)
	bar = progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v StockLocation) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid := companyIDs[v.Company]
			name := v.SBin
			pid, err := o.GetID(umdl, oarg{oarg{"complete_name", "=", v.Location}, oarg{"company_id", "=", cid}})
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.SBin}, oarg{"location_id", "=", pid}, oarg{"company_id", "=", cid}})
			o.checkErr(err)

			ur := map[string]interface{}{
				"name":        name,
				"location_id": pid,
				"company_id":  cid,
			}

			o.Log.Info(mdl, "record", ur, "r", r, "v", v)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

func (o *OdooConn) StockLocationMap() map[string]int {
	mdl := "stock_location"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"name"})
	o.checkErr(err)
	cids := map[string]int{}
	for _, c := range cc {
		cids[c["complete_name"].(string)] = int(c["id"].(float64))
	}
	return cids
}
