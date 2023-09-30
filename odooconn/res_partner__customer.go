package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) ResPartnerCustomerUnlink() {
	// customer
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerCustomerUnlink\n", umdl)

	stmt := ``
	type Partner struct {
		Ref   string `db:"kunnr"`
		Cname string `db:"cname"`
	}
	var dbrecs []Partner
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	partners, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"ref"})
	o.checkErr(err)
	ids := []int{}
	for _, or := range partners {
		for _, dr := range dbrecs {
			if or["ref"] == dr.Ref {
				ids = append(ids, int(or["id"].(float64)))
			}
		}
	}
	o.Unlink(umdl, ids)
}

// ResPartnerCustomer function
func (o *OdooConn) ResPartnerCustomer() {
	// customer
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerCustomer\n", umdl)

	stmt := ``
	type Partner struct {
		Name                     string  `db:"cname"`
		Ref                      string  `db:"kunnr"`
		ParentID                 string  `db:"parent"`
		Street                   string  `db:"street"`
		City                     string  `db:"city"`
		State                    string  `db:"state"`
		Country                  string  `db:"country"`
		Zip                      string  `db:"zip"`
		Salesregion              string  `db:"salesregion"`
		Type                     string  `db:"ctype"`
		IsCompany                bool    `db:"is_company"`
		CustomerNotes            string  `db:"customer_notes"`
		Customergroup            string  `db:"customergroup"`
		AltPricelist             string  `db:"alt_pricelist"`
		PropertyProductPricelist string  `db:"property_product_pricelist"`
		Phone                    string  `db:"phone"`
		Email                    string  `db:"email"`
		CreditLimit              float64 `db:"credit_limit"`
		DebitLimit               float64 `db:"debit_limit"`
		Taxjurcode               string  `db:"taxjurcode"`
		FiscalPosition           string  `db:"fiscal_position"`
		Taxkd                    string  `db:"taxkd"`
		VAT                      string  `db:"vat"`
		Payterms                 string  `db:"payterms"`
	}
	var rr []Partner

	custtype := "billto"
	switch custtype {
	case "billto":
		stmt = stmt + ` where c.kunnr in (select distinct billto from odoo.artg_customer_hierarchy) `
	case "shipto":
		stmt = stmt + ` where c.kunnr not in (select distinct billto from odoo.artg_customer_hierarchy) `
	}
	// stmt = stmt + ` and c.kunnr like '%1000018'`
	stmt = stmt + ` order by kunnr `
	// stmt = stmt + ` limit 10 `
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	payterms := o.AccountPaymentTermMap()

	// cids := o.ResCompanyMap()
	// cmap := []int{cids["A.R. Thomson Group"], cids["Groupe A.R. Thomson"]}
	plist := o.ProductPricelistMap()
	tagIDs, err := o.SearchRead("res.partner.category", oarg{}, 0, 0, []string{"name"})
	o.checkErr(err)
	tids := map[string]int{}
	for _, c := range tagIDs {
		tids[c["name"].(string)] = int(c["id"].(float64))
	}

	// tasker
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid, err := o.CountryID(v.Country)
			o.checkErr(err)
			sid, err := o.StateID(cid, v.State)
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			o.checkErr(err)

			altPriceListID := plist[v.AltPricelist]
			priceListID := plist[v.PropertyProductPricelist]

			payterm := payterms[v.Payterms]
			o.Log.Info("", "priceListID", priceListID, "altPriceListID", altPriceListID)

			if priceListID != altPriceListID {
				o.Log.Info("#### THEY ARE DIFFERENT ####")
			}

			ur := map[string]interface{}{
				"name":                       v.Name,
				"ref":                        v.Ref,
				"street":                     v.Street,
				"city":                       v.City,
				"country_id":                 cid,
				"zip":                        v.Zip,
				"phone":                      v.Phone,
				"email":                      v.Email,
				"is_company":                 v.IsCompany,
				"property_product_pricelist": priceListID,
				"property_payment_term_id":   payterm,
				"type":                       v.Type,
				"credit_limit":               v.CreditLimit,
			}

			// "customer_notes":             v.CustomerNotes, // requires model adjustment to add notes
			// "companies":                  cmap, // requires many-to-many adjustment to isolate delpro companies from artg companies

			pid, err := o.GetID(umdl, oarg{oarg{"ref", "=", v.ParentID}})
			o.checkErr(err)

			if v.Ref != v.ParentID {
				if pid != -1 {
					ur["parent_id"] = pid
				}
			}

			if sid != -1 {
				ur["state_id"] = sid
			}
			// Note: the expected format is 'CC##' (CC=Country Code, ##=VAT Number)
			// if v.VAT != "" {
			// 	ur["vat"] = v.VAT
			// }
			if v.FiscalPosition != "" {
				fiscalPositionID, err := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", v.FiscalPosition}, oarg{"country_id", "=", cid}})
				o.checkErr(err)
				ur["property_account_position_id"] = fiscalPositionID
			}

			tagMap := []int{tids["customer"]}
			tagMap = append(tagMap, tids[v.Salesregion])
			ur["category_id"] = tagMap

			o.Log.Info(umdl, "v", v, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ResPartnerCustomerLink function
func (o *OdooConn) ResPartnerCustomerLink() {
	// customer_link
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerCustomer\n", umdl)

	stmt := `
	select
	cname,kunnr,parent,street,city,state,country,zip,ctype,is_company
	,customer_notes,customergroup,alt_pricelist,property_product_pricelist
	,phone,email,credit_limit,debit_limit,taxjurcode,fiscal_position,taxkd,vat,payterms
	from odoo.artg_res_partner_customer c
	`
	type Partner struct {
		Name                     string  `db:"cname"`
		Ref                      string  `db:"kunnr"`
		ParentID                 string  `db:"parent"`
		Street                   string  `db:"street"`
		City                     string  `db:"city"`
		State                    string  `db:"state"`
		Country                  string  `db:"country"`
		Zip                      string  `db:"zip"`
		Type                     string  `db:"ctype"`
		IsCompany                bool    `db:"is_company"`
		CustomerNotes            string  `db:"customer_notes"`
		Customergroup            string  `db:"customergroup"`
		AltPricelist             string  `db:"alt_pricelist"`
		PropertyProductPricelist string  `db:"property_product_pricelist"`
		Phone                    string  `db:"phone"`
		Email                    string  `db:"email"`
		CreditLimit              float64 `db:"credit_limit"`
		DebitLimit               float64 `db:"debit_limit"`
		Taxjurcode               string  `db:"taxjurcode"`
		FiscalPosition           string  `db:"fiscal_position"`
		Taxkd                    string  `db:"taxkd"`
		VAT                      string  `db:"vat"`
		Payterms                 string  `db:"payterms"`
	}
	rr := []Partner{}

	stmt2 := stmt + `where c.kunnr <> c.parent`
	// fmt.Println(stmt2)
	err := o.DB.Select(&rr, stmt2)
	o.checkErr(err)
	recs := len(rr)
	// fmt.Println(recs)
	bar := progressbar.Default(int64(recs))

	payterms := o.AccountPaymentTermMap()
	// cids := o.ResCompanyMap()
	// cmap := []int{cids["A.R. Thomson Group"], cids["Groupe A.R. Thomson"]}
	plist := o.ProductPricelistMap()

	// tasker
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid, err := o.CountryID(v.Country)
			o.checkErr(err)
			sid, err := o.StateID(cid, v.State)
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			o.checkErr(err)
			pid, err := o.GetID(umdl, oarg{oarg{"ref", "=", v.ParentID}})
			o.checkErr(err)

			altPriceListID := plist[v.AltPricelist]
			priceListID, err := o.GetID("product.pricelist", oarg{oarg{"name", "=", v.PropertyProductPricelist}})
			o.checkErr(err)

			payterm := payterms[v.Payterms]
			o.Log.Info("", "priceListID", priceListID, "altPriceListID", altPriceListID)

			if priceListID != altPriceListID {
				o.Log.Info("#### THEY ARE DIFFERENT ####")
			}

			ur := map[string]interface{}{
				"name":                       v.Name,
				"ref":                        v.Ref,
				"street":                     v.Street,
				"city":                       v.City,
				"country_id":                 cid,
				"zip":                        v.Zip,
				"phone":                      v.Phone,
				"email":                      v.Email,
				"is_company":                 false,
				"property_product_pricelist": priceListID,
				"property_payment_term_id":   payterm,
				"type":                       v.Type,
				"credit_limit":               v.CreditLimit,
				// "customer":                   true, // requires model adjustment to add customer flag
				// "customer_notes":             v.CustomerNotes, // requires model adjustment to add notes
				// "companies":                  cmap, // requires many-to-many adjustment to isolate delpro companies from artg companies
				"parent_id": pid,
			}

			if sid != -1 {
				ur["state_id"] = sid
			}
			if v.VAT != "" {
				ur["vat"] = v.VAT
			}
			if v.FiscalPosition != "" {
				fiscalPositionID, err := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", v.FiscalPosition}, oarg{"country_id", "=", cid}})
				o.checkErr(err)
				ur["property_account_position_id"] = fiscalPositionID
			}
			o.Log.Info(umdl, "v", v, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
