package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) ResPartnerBank() {
	mdl := "res_bank"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ResPartnerBank\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select distinct
	upper(trim(bnka.banka)) as "name"
	,upper(trim(bnka.stras)) as street
	,upper(trim(bnka.ort01)) as city
	,trim(s."name") as state
	,trim(c."name") as country
	,trim(bnka.bankl) as bic
	,trim(bnka.swift) as swift
	from sap.bnka bnka
	join ct.state s on bnka.provz = s.state
	join ct.country c on bnka.banks = c.country
	order by trim(bnka.bankl)
	`
	type Bank struct {
		Name    string `db:"name"`
		Street  string `db:"street"`
		City    string `db:"city"`
		State   string `db:"state"`
		Country string `db:"country"`
		BIC     string `db:"bic"`
		SWIFT   string `db:"swift"`
	}
	var rr []Bank
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(len(rr))
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Bank) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid, err := o.CountryID(v.Country)
			o.checkErr(err)
			sid, err := o.StateID(cid, v.State)
			o.checkErr(err)
			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"bic", "=", v.BIC}})
			o.checkErr(err)

			ur := map[string]interface{}{
				"name":    v.Name,
				"street":  v.Street,
				"city":    v.City,
				"state":   sid,
				"country": cid,
				"bic":     v.BIC,
			}
			o.Log.Info(umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
