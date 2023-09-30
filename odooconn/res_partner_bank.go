package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ppreeper/pad"
	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) ResPartnerVendorsBankUnlink() {
	// vendor_bank_unlink
	mdl := "res_partner_bank"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerVendorsBank", umdl)

	ids, err := o.Search(umdl, oarg{})
	o.checkErr(err)

	pageSize := 100
	var ddList [][]int
	for i := 0; i <= (len(ids) / pageSize); i++ {
		ddList = append(ddList, ids[i*pageSize:(i+1)*pageSize])
	}

	bar := progressbar.Default(int64(len(ddList)))
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(len(ddList))
	for _, r := range ddList {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, r []int) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Unlink(umdl, r)

			<-sem
		}(sem, &wg, bar, r)
	}
	wg.Wait()
}

func (o *OdooConn) ResPartnerVendorsBank() {
	// vendor_bank
	mdl := "res_partner_bank"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerVendorsBank", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := ``

	limit := 2
	if limit > 0 {
		stmt += fmt.Sprintf(" limit %d", limit)
	}
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

	// As you can see below, the bank info is tagged onto the end of the bank name and the financial institution and transit number fields are left blank.  Canadian banks are a combination of 0 + institution # (3 digits) + transit # (5 digits)  - total of 9 digits which is also the number for the ABA/Routing number.  In the number below the bank info should be 0 004 15522 (no spaces)
	// US banks have 9 digits ABA/Routing numbers
	// Also the payment method hasn’t been transferred.  That needs to be transferred as well.  three options – C – Cheque, E – EFT or ACH, T – Wire transfer.  ODOO has Check, EFT, NACHA, and Manual.

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			// o.Log.Debug(umdl, "v", v)

			// companyID := cids[company]
			partnerID, err := o.GetID("res.partner", oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			o.checkErr(err)
			currencyID, err := o.GetID("res.currency", oarg{oarg{"name", "=", v.OrderCurrency}})
			o.checkErr(err)
			bankID, err := o.GetID("res.bank", oarg{oarg{"bic", "=", v.BIC}})
			o.checkErr(err)
			// bankCountry := o.SearchRead("res.bank", oarg{"id", "=", bankID}, 0, 0, []string{"country"})
			r, err := o.GetID(umdl, oarg{oarg{"partner_id", "=", partnerID}, oarg{"currency_id", "=", currencyID}, oarg{"bank_id", "=", bankID}, oarg{"acc_number", "=", v.BAN}})
			o.checkErr(err)

			ur := map[string]interface{}{
				"acc_number": v.BAN,
				// "sanitized_acc-number": v.Ref,
				"acc_holder_name": v.BANName,
				"partner_id":      partnerID,
				"bank_id":         bankID,
				"currency_id":     currencyID,
			}
			// if bank in US
			ur["aba_routing"] = v.BIC
			// if bank in CA
			ur["financial_institution_number"] = pad.RJustLen(pad.LJustLen(v.BIC, 4), 3)
			ur["bank_transit_number"] = pad.RJustLen(v.BIC, 5)

			o.Log.Debug(umdl, "record", ur, "r", r)
			if r == -1 {
				o.Record(umdl, r, ur)
			}

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

func (o *OdooConn) ResPartnerVendorsBankFix(limit int) {
	// vendor_bank_fix
	mdl := "res_partner_bank"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerVendorsBankFix", umdl)

	// banks := o.ModelMap("res.bank", "name")
	// o.Log.Debug(umdl, "banks", banks)

	bankAccounts, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{})
	o.checkErr(err)
	// bankAccounts := o.SearchRead(umdl, oarg{}, 0, 0, []string{})
	// o.Log.Debug(umdl, "bankAccounts", bankAccounts)
	recs := len(bankAccounts)
	bar := progressbar.Default(int64(recs))

	// sem := make(chan int, o.JobCount)
	// var wg sync.WaitGroup
	// wg.Add(recs)
	for _, b := range bankAccounts {
		// go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, b map[string]interface{}) {
		// 	defer bar.Add(1)
		// 	defer wg.Done()
		// 	sem <- 1
		r := int(b["id"].(float64))
		bankBic := b["bank_bic"].(string)

		ur := map[string]interface{}{}
		// if bank in US
		ur["aba_routing"] = bankBic
		// if bank in CA
		ur["financial_institution_number"] = pad.RJustLen(pad.LJustLen(bankBic, 4), 3)
		ur["bank_transit_number"] = pad.RJustLen(bankBic, 5)

		o.Log.Debug(umdl, "record", ur, "r", r)

		o.Record(umdl, r, ur)

		bar.Add(1)

		// 	<-sem
		// }(sem, &wg, bar, b)
	}
}
