package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// StockWarehouse function
func (o *OdooConn) StockWarehouse() {
	mdl := "stock_warehouse"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Println(umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select company,citycode,pprefix,werks,"name",pname,name1,name2,street,city,state,post_code,country,taxjurcode,txjcd,taxiw,tel_number,fax_number,time_zone,website
	from ct.company_plant_list
	order by company,werks
	`
	type Warehouse struct {
		Company    string `db:"company"`
		Citycode   string `db:"citycode"`
		Pprefix    string `db:"pprefix"`
		Werks      string `db:"werks"`
		Name       string `db:"name"`
		Pname      string `db:"pname"`
		Name1      string `db:"name1"`
		Name2      string `db:"name2"`
		Street     string `db:"street"`
		City       string `db:"city"`
		State      string `db:"state"`
		PostCode   string `db:"post_code"`
		Country    string `db:"country"`
		Taxjurcode string `db:"taxjurcode"`
		Txjcd      string `db:"txjcd"`
		Taxiw      string `db:"taxiw"`
		TelNumber  string `db:"tel_number"`
		FaxNumber  string `db:"fax_number"`
		TimeZone   string `db:"time_zone"`
		Website    string `db:"website"`
	}
	var rr []Warehouse
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Warehouse) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid := o.GetID("res.company", oarg{oarg{"name", "=", v.Company}, oarg{"city", "=", v.City}})
			pid := o.GetID("res.partner", oarg{oarg{"name", "=", v.Pname}, oarg{"city", "=", v.City}})
			r := o.GetID(umdl, oarg{oarg{"company_id", "=", cid}, oarg{"partner_id", "=", pid}})

			switch v.Pprefix {
			case "ARSUR":
				r = 1
			case "ARMTL":
				r = o.GetID(umdl, oarg{oarg{"name", "=", v.Company}})
				if r == -1 {
					r = o.GetID(umdl, oarg{oarg{"code", "=", v.Pprefix}})
				}
			case "DASUR":
				r = o.GetID(umdl, oarg{oarg{"name", "=", v.Company}})
				if r == -1 {
					r = o.GetID(umdl, oarg{oarg{"code", "=", v.Pprefix}})
				}
			case "DTSUR":
				r = o.GetID(umdl, oarg{oarg{"name", "=", v.Company}})
				if r == -1 {
					r = o.GetID(umdl, oarg{oarg{"code", "=", v.Pprefix}})
				}
			}

			ur := map[string]interface{}{
				"name":       v.Pname,
				"code":       v.Pprefix,
				"company_id": cid,
				"partner_id": pid,
			}
			o.Log.Infow(mdl, "##### model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
