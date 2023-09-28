package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// AccountPaymentTerm function
func (o *OdooConn) AccountPaymentTerm() {
	mdl := "account_payment_term"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	// AccountTerms struct
	type AccountTerms struct {
		Name string `db:"name"`
		Note string `db:"note"`
	}

	stmt := `select name,note from odoo.account_payment_term`

	var accountTerms []AccountTerms
	err := o.DB.Select(&accountTerms, stmt)
	o.checkErr(err)
	recs := len(accountTerms)
	bar := progressbar.Default(int64(recs))
	for _, v := range accountTerms {
		err := bar.Add(1)
		o.checkErr(err)

		// process
		r := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}})

		ur := map[string]interface{}{
			"name": v.Name,
			"note": v.Note,
		}

		o.Log.Infow(mdl, "model", umdl, "record", ur, "r", r)

		o.Record(umdl, r, ur)
	}
}

func (o *OdooConn) AccountPaymentTermMap() map[string]int {
	mdl := "account_payment_term"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc := o.SearchRead(umdl, oarg{}, 0, 0, []string{"name"})
	ids := map[string]int{}
	for _, c := range cc {
		ids[c["name"].(string)] = int(c["id"].(float64))
	}
	return ids
}

// AccountPaymentTermLine function
func (o *OdooConn) AccountPaymentTermLine() {
	mdl := "account_payment_term_line"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	// AccountTerms struct
	type AccountTerms struct {
		Name          string  `db:"name"`
		Sequence      int     `db:"sequence"`
		Value         string  `db:"value"`
		ValueAmount   float64 `db:"value_amount"`
		Days          int     `db:"days"`
		DayOfTheMonth int     `db:"day_of_the_month"`
		Option        string  `db:"option"`
	}

	stmt := `
	select 
	name
	,"sequence" 
	,value
	,value_amount
	,days
	,day_of_the_month
	,"option" 
	from odoo.account_payment_term_line
	order by name,"sequence"
	`
	var accountTerms []AccountTerms
	err := o.DB.Select(&accountTerms, stmt)
	o.checkErr(err)
	recs := len(accountTerms)
	bar := progressbar.Default(int64(recs))
	for _, v := range accountTerms {
		err := bar.Add(1)
		o.checkErr(err)

		// process
		r := o.GetID(umdl, oarg{oarg{"payment_id", "=", v.Name}})

		ur := map[string]interface{}{
			"sequence":         v.Sequence,
			"value":            v.Value,
			"value_amount":     v.ValueAmount,
			"days":             v.Days,
			"day_of_the_month": v.DayOfTheMonth,
			"option":           v.Option,
		}

		o.Log.Infow(mdl, "model", umdl, "record", ur, "r", r)

		o.Record(umdl, r, ur)
	}
}
