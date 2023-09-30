package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// ResCompany function
func (o *OdooConn) ResCompanyBranch() {
	// branches
	mdl := "res_company_branch"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	stmt := ``

	type Warehouse struct {
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
	var rr []Warehouse
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)

	cids, err := o.ModelMap("res.company", "name")
	o.checkErr(err)
	
	recs := len(rr)
	bar := progressbar.Default(int64(recs))
	for _, v := range rr {
		r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Pname}})
		o.checkErr(err)
		cid := cids[v.Company]
		o.checkErr(err)
		pid, err := o.GetID("res.partner", oarg{oarg{"name", "=", v.Pname}, oarg{"parent_id", "=", v.Company}})
		o.checkErr(err)
		partner, err := o.SearchRead("res.partner", oarg{oarg{"name", "=", v.Pname}, oarg{"parent_id", "=", v.Company}}, 0, 0, []string{"name", "parent_id", "phone"})
		o.checkErr(err)

		defaultDeliveryRouteId, err := o.GetID("stock.location.route", oarg{oarg{"name", "=", "From " + v.Pprefix}, oarg{"company_id", "=", cid}})
		o.checkErr(err)
		accountAnalyticId, err := o.GetID("account.analytic.account", oarg{oarg{"name", "=", v.Pname}, oarg{"company_id", "=", cid}})
		o.checkErr(err)
		soRequestApprovalId, err := o.GetID("approval.category", oarg{oarg{"name", "=", "Sale Approval - " + v.Pprefix}, oarg{"company_id", "=", cid}})
		o.checkErr(err)

		// TODO: needs default incoming routes

		ur := map[string]interface{}{
			"name":       v.Pname,
			"code":       v.Pprefix,
			"company_id": cid,
			"partner_id": pid,
			"phone":      partner[0]["phone"].(string),
		}

		if defaultDeliveryRouteId != -1 {
			ur["default_delivery_route_id"] = defaultDeliveryRouteId
		}

		if accountAnalyticId != -1 {
			ur["account_analytic_id"] = accountAnalyticId
		}

		if soRequestApprovalId != -1 {
			ur["so_request_approval_id"] = soRequestApprovalId
		}

		o.Log.Debug(umdl, "r", r, "ur", ur)
		o.Record(umdl, r, ur)

		bar.Add(1)
	}
}
