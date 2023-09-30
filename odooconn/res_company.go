package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// ResCompany function
func (o *OdooConn) ResCompany() {
	mdl := "res_company"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	stmt := `
	select
	"id"
	,"name"
	,"parent_company"
	,"street"
	,"city"
	,"state"
	,"country"
	,"zip"
	,"phone"
	,"email"
	,"website"
	,"vat"
	,"paperformat_id"
	,"chart_template_id"
	,"fiscalyear_last_month"
	,"fiscalyear_last_day"
	from ct.company_list
	where name not like '%Manufacturing%'
	order by parent_company,id
	`
	type Company struct {
		ID                  int    `db:"id"`
		Name                string `db:"name"`
		ParentCompany       string `db:"parent_company"`
		Street              string `db:"street"`
		City                string `db:"city"`
		State               string `db:"state"`
		Country             string `db:"country"`
		Zip                 string `db:"zip"`
		Phone               string `db:"phone"`
		Email               string `db:"email"`
		Website             string `db:"website"`
		VAT                 string `db:"vat"`
		PaperformatID       int    `db:"paperformat_id"`
		ChartTemplateID     int    `db:"chart_template_id"`
		FiscalyearLastMonth int    `db:"fiscalyear_last_month"`
		FiscalyearLastDay   int    `db:"fiscalyear_last_day"`
	}

	var rr []Company
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	for _, v := range rr {
		err := bar.Add(1)
		o.checkErr(err)

		r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}})
		o.checkErr(err)
		pid, err := o.GetID(umdl, oarg{oarg{"name", "=", v.ParentCompany}})
		o.checkErr(err)
		// chartID := o.GetID("account.chart.template", oarg{oarg{"name", "like", "Canada"}})
		paperID, err := o.GetID("report.paperformat", oarg{oarg{"name", "=", "US Letter"}})
		o.checkErr(err)
		// accountJournalID := o.GetID("account.journal", oarg{oarg{"name", "=", "Miscellaneous Operations"}})
		salesTaxID, err := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for sales - 5%"}})
		o.checkErr(err)
		purchaseTaxID, err := o.GetID("account.tax", oarg{oarg{"name", "=", "GST for purchases - 5%"}})
		o.checkErr(err)
		// currencyExchangeJournalID := o.GetID("account.journal", oarg{oarg{"name", "=", "Exchange Difference"}})
		// incomeCurrencyExchangeAccountID := o.GetID("account.account", oarg{oarg{"name", "like", "420000"}})
		// expenseCurrencyExchangeAccountID := o.GetID("account.account", oarg{oarg{"name", "like", "550000"}})
		cid, err := o.CountryID(v.Country)
		o.checkErr(err)
		sid, err := o.StateID(cid, v.State)
		o.checkErr(err)

		ur := map[string]interface{}{
			"name":       v.Name,
			"street":     v.Street,
			"city":       v.City,
			"state_id":   sid,
			"country_id": cid,
			"zip":        v.Zip,
			"phone":      v.Phone,
			"email":      v.Email,
			"website":    v.Website,
			"vat":        v.VAT,
			// "chart_template_id":                  chartID,
			"paperformat_id": paperID,
			// "account_tax_periodicity_journal_id": accountJournalID,
			// "currency_exchange_journal_id":       currencyExchangeJournalID,
			// "income_currency_exchange_account_id":  incomeCurrencyExchangeAccountID,
			// "expense_currency_exchange_account_id": expenseCurrencyExchangeAccountID,
			"account_sale_tax_id":     salesTaxID,
			"account_purchase_tax_id": purchaseTaxID,
			"fiscalyear_last_month":   fmt.Sprintf("%d", v.FiscalyearLastMonth),
			"fiscalyear_last_day":     v.FiscalyearLastDay,
		}
		o.Log.Info(umdl, "ur", ur, "pid", pid, "r", r, "v.ID", v.ID)

		if v.ID == 1 {
			if r == -1 {
				row, res, err := o.WriteRecord(umdl, v.ID, UPDATE, ur)
				if err != nil {
					o.Log.Info(umdl, "row", row, "res", res, "err", err)
				}
				o.ResCompanyLDAP(v.ID)
			} else {
				if !o.NoUpdate {
					row, res, err := o.WriteRecord(umdl, r, UPDATE, ur)
					if err != nil {
						o.Log.Info(umdl, "row", row, "res", res, "err", err)
					}
					o.ResCompanyLDAP(r)
				}
			}
		} else {
			ur["parent_id"] = pid
			if r == -1 {
				row, res, err := o.WriteRecord(umdl, r, INSERT, ur)
				if err != nil {
					o.Log.Info(umdl, "row", row, "res", res, "err", err)
				}
				o.ResCompanyLDAP(r)
			} else {
				if !o.NoUpdate {
					row, res, err := o.WriteRecord(umdl, r, UPDATE, ur)
					if err != nil {
						o.Log.Info(umdl, "row", row, "res", res, "err", err)
					}
					o.ResCompanyLDAP(r)
				}
			}
		}
	}
}

func (o *OdooConn) ResCompanyMap() map[string]int {
	mdl := "res_company"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"name"})
	o.checkErr(err)
	cids := map[string]int{}
	for _, c := range cc {
		cids[c["name"].(string)] = int(c["id"].(float64))
	}
	return cids
}
