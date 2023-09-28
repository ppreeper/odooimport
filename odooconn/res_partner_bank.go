package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) ResPartnerVendorsBank() {
	// vendor_bank
	mdl := "res_partner_bank"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerVendorsBank", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select right(vb.vendorno,6) vendorno,vb."name"
	,v.order_currency
	,c."name" bcountry
	,vb.bic,vb.ban,vb.ban_name
	from sapdata.vendors v
	join odoo.artg_vendors_active ava on v.vendorno = ava.vendorno
	join sapdata.vendorsbanktemp vb on v.vendorno = vb.vendorno
	join ct.country c on vb.banks = c.country
	where v.vendortype <> 'YB05'
	and v."name" not like '%DO_NOT%'
	and v.vendorno = v.parent_id
	and vb.bic <> ''
	`
	type Partner struct {
		Ref           string `db:"vendorno"`
		Name          string `db:"name"`
		OrderCurrency string `db:"order_currency"`
		BCountry      string `db:"bcountry"`
		BIC           string `db:"bic"`
		BAN           string `db:"ban"`
		BANName       string `db:"ban_name"`
	}
	var rr []Partner
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// cids := o.ResCompanyMap()

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Log.Infow(umdl, "v", v)

			// companyID := cids[company]
			partnerID := o.GetID("res.partner", oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			currencyID := o.GetID("res.currency", oarg{oarg{"name", "=", v.OrderCurrency}})
			bankID := o.GetID("res.bank", oarg{oarg{"bic", "=", v.BIC}})
			r := o.GetID(umdl, oarg{oarg{"partner_id", "=", partnerID}, oarg{"currency_id", "=", currencyID}, oarg{"bank_id", "=", bankID}, oarg{"acc_number", "=", v.BAN}})

			ur := map[string]interface{}{
				"acc_number": v.BAN,
				// "sanitized_acc-number": v.Ref,
				"acc_holder_name": v.BANName,
				"partner_id":      partnerID,
				"bank_id":         bankID,
				"currency_id":     currencyID,
			}
			o.Log.Infow(umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
