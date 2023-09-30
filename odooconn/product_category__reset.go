package odooconn

import (
	"fmt"
	"strings"
)

// ProductCategoryConsumable3 function
func (o *OdooConn) ProductCategoryReset() {
	mdl := "product_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\nconsumables3\n", umdl)

	odoorecs, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"display_name"})
	o.checkErr(err)
	fmt.Println("product.categories:", len(odoorecs))
	odooIDs, err := o.Search(umdl, oarg{})
	o.checkErr(err)
	fmt.Println("product.categories:", len(odooIDs))

	// bar := progressbar.Default(int64(recs))

	incCode, err := o.GetID("account.account", oarg{oarg{"code", "=", "420000"}})
	o.checkErr(err)
	expCode, err := o.GetID("account.account", oarg{oarg{"code", "=", "511100"}})
	o.checkErr(err)

	// tasker
	// wg.Add(recs)
	for _, r := range odooIDs {
		// pid := o.GetID(umdl, oarg{oarg{"name", "=", v.Group1}, oarg{"parent_id", "=", allID}})
		// sid := -1
		// if pid != -1 {
		// 	sid = o.GetID(umdl, oarg{oarg{"name", "=", v.Group2}, oarg{"parent_id", "=", pid}})
		// }
		// gid := -1
		// if sid != -1 {
		// 	gid = o.GetID(umdl, oarg{oarg{"name", "=", v.Matgrp}, oarg{"parent_id", "=", sid}})

		// with buyer id setting
		// bidPartner := o.GetID("res.partner", oarg{oarg{"name", "=", v.Buyer}})
		// bid := o.GetID("res.users", oarg{oarg{"partner_id", "=", bidPartner}})

		ur := map[string]interface{}{
			// 	"name":                              v.Matgrp,
			// 	"parent_id":                         sid,
			// 	"property_cost_method":              v.PropertyCostMethod,
			// 	"property_valuation":                "real_time",
			"property_account_income_categ_id":  incCode,
			"property_account_expense_categ_id": expCode,
		}

		// with buyer id setting
		// if v.Matkl != "" {
		// 	ur["material_group"] = v.Matkl
		// }

		// if v.DayReview != "" {
		// 	ur["review_day"] = v.DayReview
		// }

		// if bid != -1 {
		// 	ur["buyer_id"] = bid
		// }

		o.Log.Info(umdl, "record", ur, "r", r)

		o.Record(umdl, r, ur)

	}
}
