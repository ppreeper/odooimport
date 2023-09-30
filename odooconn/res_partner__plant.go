package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ResPartnerPlant function
func (o *OdooConn) ResPartnerPlant() {
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v plant\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select company,citycode,pprefix,werks,"name",pname,name1,name2,street,city,state,post_code,country,taxjurcode,txjcd,taxiw,tel_number,fax_number,time_zone,website
	from ct.company_plant_list
	order by company,werks
	`

	type Plant struct {
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
	var rr []Plant
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Plant) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Pname}})
			o.checkErr(err)
			companyID, err := o.CompanyID(v.Company)
			o.checkErr(err)
			parentID, err := o.PartnerID(v.Company)
			o.checkErr(err)
			countryID, err := o.CountryID(v.Country)
			o.checkErr(err)
			stateID, err := o.StateID(countryID, v.State)
			o.checkErr(err)
			propertyAccountPositionID, err := o.FiscalPosition(countryID, v.State)
			o.checkErr(err)

			ur := map[string]interface{}{
				"name":                         v.Pname,
				"company_id":                   companyID,
				"parent_id":                    parentID,
				"display_name":                 v.Pname,
				"type":                         "delivery",
				"ref":                          v.Werks,
				"website":                      v.Website,
				"street":                       v.Street,
				"zip":                          v.PostCode,
				"city":                         v.City,
				"state_id":                     stateID,
				"country_id":                   countryID,
				"phone":                        v.TelNumber,
				"is_company":                   true,
				"property_account_position_id": propertyAccountPositionID,
			}
			o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
