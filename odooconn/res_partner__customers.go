package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ResPartnerCustomer function
func (o *OdooConn) ResPartnerCustomer() {
	// customer
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
	var rr []Partner

	stmt1 := stmt + `where c.kunnr = c.parent`
	err := o.DB.Select(&rr, stmt1)
	o.checkErr(err)
	recs := len(rr)
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

			cid := o.CountryID(v.Country)
			sid := o.StateID(cid, v.State)
			r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})

			altPriceListID := plist[v.AltPricelist]
			priceListID := plist[v.PropertyProductPricelist]

			payterm := payterms[v.Payterms]
			o.Log.Infow("", "priceListID", priceListID, "altPriceListID", altPriceListID)

			if priceListID != altPriceListID {
				o.Log.Infow("#### THEY ARE DIFFERENT ####")
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
				// "customer":                   true, // requires model adjustment to add customer flag
				"type":         v.Type,
				"credit_limit": v.CreditLimit,
				// "customer_notes":             v.CustomerNotes, // requires model adjustment to add notes
				// "companies":                  cmap, // requires many-to-many adjustment to isolate delpro companies from artg companies
			}
			if sid != -1 {
				ur["state_id"] = sid
			}
			if v.VAT != "" {
				ur["vat"] = v.VAT
			}
			if v.FiscalPosition != "" {
				fiscalPositionID := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", v.FiscalPosition}, oarg{"country_id", "=", cid}})
				ur["property_account_position_id"] = fiscalPositionID
			}
			o.Log.Infow(umdl, "v", v, "model", umdl, "record", ur, "r", r)

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

			cid := o.CountryID(v.Country)
			sid := o.StateID(cid, v.State)
			r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			pid := o.GetID(umdl, oarg{oarg{"ref", "=", v.ParentID}})

			altPriceListID := plist[v.AltPricelist]
			priceListID := o.GetID("product.pricelist", oarg{oarg{"name", "=", v.PropertyProductPricelist}})

			payterm := payterms[v.Payterms]
			o.Log.Infow("", "priceListID", priceListID, "altPriceListID", altPriceListID)

			if priceListID != altPriceListID {
				o.Log.Infow("#### THEY ARE DIFFERENT ####")
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
				fiscalPositionID := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", v.FiscalPosition}, oarg{"country_id", "=", cid}})
				ur["property_account_position_id"] = fiscalPositionID
			}
			o.Log.Infow(umdl, "v", v, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ResPartnerCustomerSROdooCRM function
func (o *OdooConn) ResPartnerCustomerSROdooCRM() {
	// customer_sr_crm
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerCustomerSROdooCRM\n", umdl)

	stmt := `
	select
	c.cname,c.kunnr,cc.username,cc.ename
	from odoo.artg_res_partner_customer c
	join (
	select distinct cs.kunnr,cs.username,u.ename
	from ct.customer_sr cs
	join sapdata.users u on cs.username = u.username
	where cs.salesorg <> '1010' and cs.datestart::date <= current_date and current_date <= cs.dateend::date
	and cs.division = 'Control'
	) cc on c.kunnr = right(cc.kunnr,7)
	order by c.cname,c.kunnr
	`
	type PartnerUser struct {
		Name     string `db:"cname"`
		Ref      string `db:"kunnr"`
		Username string `db:"username"`
		Ename    string `db:"ename"`
	}
	rr := []PartnerUser{}

	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	// fmt.Println(recs)
	bar := progressbar.Default(int64(recs))

	// tasker
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v PartnerUser) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"ref", "=", v.Ref}})
			userID := o.GetID("res.users", oarg{oarg{"name", "=", v.Ename}})

			ur := map[string]interface{}{
				"user_id": userID,
			}

			// o.Log.Infow(umdl, "v", v, "model", umdl, "record", ur, "r", r)
			o.Log.Infow(umdl, "r", r, "ur", ur)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ResPartnerCustomerDelpro function
func (o *OdooConn) ResPartnerCustomerDelpro() {
	// customer_delpro
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerCustomerDelpro\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select
	company
	,dc."name"
	-- ,"type"
	,cid
	,parent_cid
	,contact
	,street1||' '||street2 as street
	,city
	-- ,province as state
	,case when s."name" is null then dc.province else s."name" end state
	,country
	,postalcode as zip
	,'contact' as type
	,pricelist as property_product_pricelist
	,salesperson as user_id
	,phone1 as phone
	--,phone2
	--,faxnumber
	,email
	,website
	,taxcode
	,case paymentterms
	when '0.00 % discount if paid within 0 days. Net due within 0 days.' then 'Immediate Payment'
	when '0.00 % discount if paid within 0 days. Net due within 30 days.' then '30 Days'
	when '0.00 % discount if paid within 0 days. Net due within 45 days.' then '45 Days'
	when '0.00 % discount if paid within 0 days. Net due within 60 days.' then 'NT60'
	when '0.00 % discount if paid within 0 days. Net due within 40 days.' then '45 Days'
	when '0.00 % discount if paid within 0 days. Net due within 52 days.' then 'NT60'
	when '0.00 % discount if paid within 0 days. Net due within 55 days.' then 'NT60'
	when '0.00 % discount if paid within 0 days. Net due within 90 days.' then 'NET90'
	else paymentterms end payterms
	,case country
	when 'Canada' then 'CAD'
	when 'United States' then 'USD'
	else 'CAD' end currency
	,trim(trim(field1)||' '||trim(field2)||' '||trim(field3)||' '||trim(field4)||' '||trim(field5)) fieldinfo
	,shipinfo
	from odoo.delpro_customerlist dc
	left join ct.state s on dc.province = s.state
	where type = 'B'
	order by company,name
	`
	type Partner struct {
		Company   string `db:"company"`
		Name      string `db:"name"`
		CID       string `db:"cid"`
		ParentCID string `db:"parent_cid"`
		Contact   string `db:"contact"`
		Street    string `db:"street"`
		City      string `db:"city"`
		State     string `db:"state"`
		Country   string `db:"country"`
		Zip       string `db:"zip"`
		Type      string `db:"type"`
		Pricelist string `db:"property_product_pricelist"`
		UserID    string `db:"user_id"`
		Phone     string `db:"phone"`
		Email     string `db:"email"`
		Website   string `db:"website"`
		Taxcode   string `db:"taxcode"`
		Payterms  string `db:"payterms"`
		Currency  string `db:"currency"`
		Fieldinfo string `db:"fieldinfo"`
		Shipinfo  string `db:"shipinfo"`
	}
	var rr []Partner
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	payterms := o.AccountPaymentTermMap()
	cids := o.ResCompanyMap()

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Log.Infow("", "v", v)
			companyID := cids[v.Company]
			cid := o.CountryID(v.Country)
			sid := o.StateID(cid, v.State)
			r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"company_id", "=", companyID}})
			userID := o.GetID("res.users", oarg{oarg{"name", "=", v.UserID}, oarg{"company_id", "=", companyID}})
			// priceListID := o.GetID("product.pricelist", oarg{{"name", "=", v.Pricelist}, {"company_id", "=", companyID}})
			payterm := payterms[v.Payterms]
			fiscalPositionID := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", v.State}, oarg{"country_id", "=", cid}, oarg{"company_id", "=", companyID}})

			// o.Log.Infow("", "priceListID", priceListID)

			ur := map[string]interface{}{
				"name":       v.Name,
				"company_id": companyID,
				"ref":        v.CID,
				"street":     v.Street,
				"city":       v.City,
				"country_id": cid,
				"zip":        v.Zip,
				"phone":      v.Phone,
				"email":      v.Email,
				"is_company": true,
				"type":       v.Type,
				// "property_product_pricelist": priceListID,
				"property_payment_term_id": payterm,
				"customer":                 true,
			}
			if userID != -1 {
				ur["user_id"] = userID
			}
			if sid != -1 {
				ur["state_id"] = sid
			}
			if fiscalPositionID != -1 {
				ur["property_account_position_id"] = fiscalPositionID
			}
			o.Log.Infow(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// ResPartnerCustomerLinkDelpro function
func (o *OdooConn) ResPartnerCustomerLinkDelpro() {
	// customer_delpro_link
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerCustomerLinkDelpro\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select
	company
	,dc."name"
	-- ,"type"
	,cid
	,parent_cid
	,contact
	,street1||' '||street2 as street
	,city
	-- ,province as state
	,case when s."name" is null then dc.province else s."name" end state
	,country
	,postalcode as zip
	,'delivery' as type
	,pricelist as property_product_pricelist
	,salesperson as user_id
	,phone1 as phone
	--,phone2
	--,faxnumber
	,email
	,website
	,taxcode
	,case paymentterms
	when '0.00 % discount if paid within 0 days. Net due within 0 days.' then 'Immediate Payment'
	when '0.00 % discount if paid within 0 days. Net due within 30 days.' then '30 Days'
	when '0.00 % discount if paid within 0 days. Net due within 45 days.' then '45 Days'
	when '0.00 % discount if paid within 0 days. Net due within 60 days.' then 'NT60'
	when '0.00 % discount if paid within 0 days. Net due within 40 days.' then '45 Days'
	when '0.00 % discount if paid within 0 days. Net due within 52 days.' then 'NT60'
	when '0.00 % discount if paid within 0 days. Net due within 55 days.' then 'NT60'
	when '0.00 % discount if paid within 0 days. Net due within 90 days.' then 'NET90'
	else paymentterms end payterms
	,case country
	when 'Canada' then 'CAD'
	when 'United States' then 'USD'
	else 'CAD' end currency
	,trim(trim(field1)||' '||trim(field2)||' '||trim(field3)||' '||trim(field4)||' '||trim(field5)) fieldinfo
	,shipinfo
	from odoo.delpro_customerlist dc
	left join ct.state s on dc.province = s.state
	where type = 'S' and cid <> parent_cid
	order by company,name
	`
	type Partner struct {
		Company   string `db:"company"`
		Name      string `db:"name"`
		CID       string `db:"cid"`
		ParentCID string `db:"parent_cid"`
		Contact   string `db:"contact"`
		Street    string `db:"street"`
		City      string `db:"city"`
		State     string `db:"state"`
		Country   string `db:"country"`
		Zip       string `db:"zip"`
		Type      string `db:"type"`
		Pricelist string `db:"property_product_pricelist"`
		UserID    string `db:"user_id"`
		Phone     string `db:"phone"`
		Email     string `db:"email"`
		Website   string `db:"website"`
		Taxcode   string `db:"taxcode"`
		Payterms  string `db:"payterms"`
		Currency  string `db:"currency"`
		Fieldinfo string `db:"fieldinfo"`
		Shipinfo  string `db:"shipinfo"`
	}
	var rr []Partner
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	payterms := o.AccountPaymentTermMap()
	cids := o.ResCompanyMap()

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Partner) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Log.Infow("", "v", v)
			companyID := cids[v.Company]
			cid := o.CountryID(v.Country)
			sid := o.StateID(cid, v.State)
			r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"company_id", "=", companyID}})
			pid := o.GetID(umdl, oarg{oarg{"ref", "=", v.ParentCID}, oarg{"company_id", "=", companyID}})
			userID := o.GetID("res.users", oarg{oarg{"name", "=", v.UserID}, oarg{"company_id", "=", companyID}})
			// priceListID := o.GetID("product.pricelist", oarg{{"name", "=", v.Pricelist}, {"company_id", "=", companyID}})
			payterm := payterms[v.Payterms]
			fiscalPositionID := o.GetID("account.fiscal.position", oarg{oarg{"name", "like", v.State}, oarg{"country_id", "=", cid}, oarg{"company_id", "=", companyID}})

			// o.Log.Infow("", "priceListID", priceListID)

			ur := map[string]interface{}{
				"name":       v.Name,
				"company_id": companyID,
				"ref":        v.CID,
				"parent_id":  pid,
				"street":     v.Street,
				"city":       v.City,
				"country_id": cid,
				"zip":        v.Zip,
				"phone":      v.Phone,
				"email":      v.Email,
				"is_company": false,
				"type":       v.Type,
				// "property_product_pricelist": priceListID,
				"property_payment_term_id": payterm,
			}
			if userID != -1 {
				ur["user_id"] = userID
			}
			if sid != -1 {
				ur["state_id"] = sid
			}
			if fiscalPositionID != -1 {
				ur["property_account_position_id"] = fiscalPositionID
			}
			o.Log.Infow(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
