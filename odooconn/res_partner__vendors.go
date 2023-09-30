package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ResPartnerVendors function
func (o *OdooConn) ResPartnerVendors() {
	// vendor
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerVendors", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select right(v.vendorno,6) vendorno,right(v.parent_id,6) parent_id
	,v."name",v.name1,v.name2,v.name3,v.name4,
	v.street,v.city,v.region,v.country,v.postal,v.vendortype,v.telephone,v.email,
	v.is_approved_vendor,vendor_notes,v.gl_reconcillation,
	v.pay_method,v.pay_method_desc,v.pay_term,v.order_currency,v.min_order_value
	from sapdata.vendors v
	join odoo.artg_vendors_active ava on v.vendorno = ava.vendorno
	where v.vendortype <> 'YB05'
	`
	type Partner struct {
		Ref              string `db:"vendorno"`
		VParent          string `db:"parent_id"`
		Name             string `db:"name"`
		Name1            string `db:"name1"`
		Name2            string `db:"name2"`
		Name3            string `db:"name3"`
		Name4            string `db:"name4"`
		Street           string `db:"street"`
		City             string `db:"city"`
		State            string `db:"region"`
		Country          string `db:"country"`
		Zip              string `db:"postal"`
		Vendortype       string `db:"vendortype"`
		Telephone        string `db:"telephone"`
		Email            string `db:"email"`
		IsApprovedVendor bool   `db:"is_approved_vendor"`
		VendorNotes      string `db:"vendor_notes"`
		GLReconcillation string `db:"gl_reconcillation"`
		PayMethod        string `db:"pay_method"`
		PayMethodDesc    string `db:"pay_method_desc"`
		PayTerm          string `db:"pay_term"`
		OrderCurrency    string `db:"order_currency"`
		MinOrder         string `db:"min_order_value"`
	}
	var rr []Partner
	stmt1 := stmt + `and v.vendorno = v.parent_id order by v.vendorno`
	o.Log.Info(stmt1)
	err := o.DB.Select(&rr, stmt1)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	payterms := o.AccountPaymentTermMap()
	cids := o.ResCompanyMap()
	cmap := []int{cids["A.R. Thomson Group"]}
	currs := o.ResCurrencyMap()

	// parent vendors
	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Log.Info(umdl, "v", v)

			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			o.checkErr(err)
			payterm := payterms[v.PayTerm]

			cid, err := o.CountryID(v.Country)
			o.checkErr(err)
			sid, err := o.StateID(cid, v.State)
			o.checkErr(err)
			// currencyID := o.GetID("res.currency", oarg{oarg{"name", "=", v.OrderCurrency}})
			currencyID := currs[v.OrderCurrency]

			ur := map[string]interface{}{
				"name":                              v.Name,
				"ref":                               v.Ref,
				"street":                            v.Street,
				"city":                              v.City,
				"country_id":                        cid,
				"zip":                               v.Zip,
				"phone":                             v.Telephone,
				"email":                             v.Email,
				"is_company":                        true,
				"supplier":                          true,
				"property_supplier_payment_term_id": payterm,
				"property_purchase_currency_id":     currencyID,
				"is_approved_vendor":                v.IsApprovedVendor,
				"vendor_notes":                      v.VendorNotes,
				"companies":                         cmap,
			}

			if sid != -1 {
				ur["state_id"] = sid
			}
			if v.Country == "Canada" {
				fiscalPositionID, err := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", v.State}, oarg{"country_id", "=", cid}})
				o.checkErr(err)
				ur["property_account_position_id"] = fiscalPositionID
			} else {
				fiscalPositionID, err := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", "International"}})
				o.checkErr(err)
				ur["property_account_position_id"] = fiscalPositionID
			}

			o.Log.Info(umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ResPartnerVendorsLink function
func (o *OdooConn) ResPartnerVendorsLink() {
	// vendor_link
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerVendors", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select right(v.vendorno,6) vendorno,right(v.parent_id,6) parent_id
	,v."name",v.name1,v.name2,v.name3,v.name4,
	v.street,v.city,v.region,v.country,v.postal,v.vendortype,v.telephone,v.email,
	v.is_approved_vendor,vendor_notes,v.gl_reconcillation,
	v.pay_method,v.pay_method_desc,v.pay_term,v.order_currency,v.min_order_value
	from sapdata.vendors v
	join odoo.artg_vendors_active ava on v.vendorno = ava.vendorno
	where v.vendortype <> 'YB05'
	`
	type Partner struct {
		Ref              string `db:"vendorno"`
		VParent          string `db:"parent_id"`
		Name             string `db:"name"`
		Name1            string `db:"name1"`
		Name2            string `db:"name2"`
		Name3            string `db:"name3"`
		Name4            string `db:"name4"`
		Street           string `db:"street"`
		City             string `db:"city"`
		State            string `db:"region"`
		Country          string `db:"country"`
		Zip              string `db:"postal"`
		Vendortype       string `db:"vendortype"`
		Telephone        string `db:"telephone"`
		Email            string `db:"email"`
		IsApprovedVendor bool   `db:"is_approved_vendor"`
		VendorNotes      string `db:"vendor_notes"`
		GLReconcillation string `db:"gl_reconcillation"`
		PayMethod        string `db:"pay_method"`
		PayMethodDesc    string `db:"pay_method_desc"`
		PayTerm          string `db:"pay_term"`
		OrderCurrency    string `db:"order_currency"`
		MinOrder         string `db:"min_order_value"`
	}
	rr := []Partner{}

	stmt2 := stmt + `and v.vendorno <> v.parent_id order by v.vendorno`
	o.Log.Info(stmt2)
	err := o.DB.Select(&rr, stmt2)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	payterms := o.AccountPaymentTermMap()
	// cids := o.ResCompanyMap()
	currs := o.ResCurrencyMap()

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Log.Info(umdl, "v", v)

			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			o.checkErr(err)
			payterm := payterms[v.PayTerm]
			// currencyID := o.GetID("res.currency", oarg{oarg{"name", "=", v.OrderCurrency}})
			currencyID := currs[v.OrderCurrency]

			pid, err := o.GetID(umdl, oarg{oarg{"ref", "=", v.VParent}})
			o.checkErr(err)
			cid, err := o.CountryID(v.Country)
			o.checkErr(err)
			sid, err := o.StateID(cid, v.State)
			o.checkErr(err)

			ur := map[string]interface{}{
				"name":                              v.Name,
				"ref":                               v.Ref,
				"street":                            v.Street,
				"city":                              v.City,
				"country_id":                        cid,
				"zip":                               v.Zip,
				"phone":                             v.Telephone,
				"email":                             v.Email,
				"is_company":                        false,
				"supplier":                          true,
				"property_supplier_payment_term_id": payterm,
				"property_purchase_currency_id":     currencyID,
				"is_approved_vendor":                v.IsApprovedVendor,
				"vendor_notes":                      v.VendorNotes,
				"parent_id":                         pid,
				"type":                              "delivery",
			}
			if sid != -1 {
				ur["state_id"] = sid
			}
			if v.Country == "Canada" {
				fiscalPositionID, err := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", v.State}, oarg{"country_id", "=", cid}})
				o.checkErr(err)
				ur["property_account_position_id"] = fiscalPositionID
			} else {
				fiscalPositionID, err := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", "International"}})
				o.checkErr(err)
				ur["property_account_position_id"] = fiscalPositionID
			}

			o.Log.Info(umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ResPartnerVendorsDelpro function
func (o *OdooConn) ResPartnerVendorsDelpro() {
	// vendor_delpro
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerVendorsDelpro", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select "name",right("ref",6) "ref",street ,city ,state ,country ,postal ,telephone ,lower(email) email
	from odoo.artg_vendors av
	join odoo.artg_vendors_active ava on av."ref" = ava.vendorno
	union
	select distinct "name",'' "ref", trim(street1 || ' ' ||street2) street,city ,province state,country,postalcode postal, phone1 telephone,email
	from odoo.delpro_vendors
	`
	type Partner struct {
		Name      string `db:"name"`
		Ref       string `db:"ref"`
		Street    string `db:"street"`
		City      string `db:"city"`
		State     string `db:"state"`
		Country   string `db:"country"`
		Postal    string `db:"postal"`
		Telephone string `db:"telephone"`
		Email     string `db:"email"`
	}
	var rr []Partner
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			o.checkErr(err)
			cid, err := o.CountryID(v.Country)
			o.checkErr(err)
			sid, err := o.StateID(cid, v.State)
			o.checkErr(err)
			ur := map[string]interface{}{
				"name":       v.Name,
				"ref":        v.Ref,
				"street":     v.Street,
				"city":       v.City,
				"country_id": cid,
				"zip":        v.Postal,
				"phone":      v.Telephone,
				"email":      v.Email,
				"is_company": true,
			}
			if sid != -1 {
				ur["state_id"] = sid
			}
			o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
