package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ProductCustomerPart function
func (o *OdooConn) ProductCustomerPart() {
	mdl := "product_customer_part"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductCustomerPart\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select distinct
	"right"(rk.kunnr, 7) AS kunnr,
	"right"(rk.matnr, 7) AS default_code,
	btrim(rk.kdmat) AS cust_matnr
	FROM sapdata.r_knmt rk
	left join odoo.mat_dnu dnu on rk.matnr = dnu.matnr 
	WHERE rk.vkorg <> '1010' AND btrim(rk.kdmat) <> ''
	and dnu.matnr is null 
	order by "right"(rk.kunnr, 7),"right"(rk.matnr, 7)
	-- limit 10
	`
	type Product struct {
		Customer    string `db:"kunnr"`
		DefaultCode string `db:"default_code"`
		Custmatnr   string `db:"cust_matnr"`
	}
	var rr []Product
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Product) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			partnerID, err := o.GetID("res.partner", oarg{oarg{"ref", "=", v.Customer}})
			o.checkErr(err)
			productID, err := o.GetID("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}})
			o.checkErr(err)

			r, err := o.GetID(umdl, oarg{oarg{"partner_id", "=", partnerID}, oarg{"product_tmpl_id", "=", productID}})
			o.checkErr(err)

			ur := map[string]interface{}{
				"partner_id":      partnerID,
				"product_tmpl_id": productID,
				"part_number":     v.Custmatnr,
			}

			o.Log.Info(umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
