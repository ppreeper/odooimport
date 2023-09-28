package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// StockInventory function to create an inventory adjustment
// you can find the inventory adjustments in
// the Inventory > Operations > Inventory Adjustments
func (o *OdooConn) StockInventoryInitial(company string) {
	mdli := "stock_inventory"
	umdli := strings.Replace(mdli, "_", ".", -1)
	fmt.Printf("\n%v StockInventoryInitial %v\n", umdli, company)

	company_id := o.CompanyID(company)
	scName := "Initial Count " + company

	invCountID := o.GetID(umdli, oarg{oarg{"name", "=", scName}, oarg{"state", ""}, oarg{"company_id", "=", company_id}})
	ur := map[string]interface{}{
		"name":                     scName,
		"company_id":               company_id,
		"prefill_counted_quantity": "zero",
	}

	o.Log.Infow(mdli, "model", umdli, "ur", ur, "r", invCountID)

	o.Record(umdli, invCountID, ur)

	// q := `select id,name,state,company_id,prefill_counted_quantity,exhausted from stock_inventory;`
	// id |   name    |  state  | company_id | prefill_counted_quantity | exhausted
	// ---+-----------+---------+------------+--------------------------+-----------
	//  1 | Inventory | confirm |          1 | counted                  | f

	mdll := "stock_inventory_line"
	umdll := strings.Replace(mdll, "_", ".", -1)

	fmt.Printf("\n%v StockInventoryInitial\n", umdll)

	stmtGART := `
	select
	'G'||right(i.matnr,7) as default_code
	,prefix||citycode||'/Stock' as location_id
	,i.unrestricted as product_qty
	,case when uu.name is null then 'ea' else uu.name end as product_uom
	,i.werks
	,p.company
	from sapdata.inventorytemp i
	join sapdata.materials m on i.mandt = m.mandt and i.matnr = m.matnr
	left join odoo.uom_uom uu on m.meins = uu.uom
	join odoo.plantstemp p on i.werks = p.werks
	where i.unrestricted <> 0
	and i.werks = '1024' and p.company = 'Groupe A.R. Thomson'
	union
	select
	right(i.matnr,7) as default_code
	,prefix||citycode||'/Stock' as location_id
	,i.unrestricted as product_qty
	,case when uu.name is null then 'ea' else uu.name end as product_uom
	,i.werks
	,p.company
	from sapdata.inventorytemp i
	join sapdata.materials m on i.mandt = m.mandt and i.matnr = m.matnr
	left join odoo.uom_uom uu on m.meins = uu.uom
	join odoo.plantstemp p on i.werks = p.werks
	where i.unrestricted <> 0
	and i.werks not in ('1004','1013','1024') and p.company = 'A.R. Thomson Group'
	order by location_id,default_code
	`
	type Line struct {
		DefaultCode string  `db:"default_code"`
		Location    string  `db:"location_id"`
		ProductQty  float64 `db:"product_qty"`
		ProductUOM  string  `db:"product_uom"`
		Plant       string  `db:"werks"`
		Company     string  `db:"company"`
	}
	rr := []Line{}
	err := o.DB.Select(&rr, stmtGART)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	for _, v := range rr {
		err := bar.Add(1)
		o.checkErr(err)

		product := o.SearchRead("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}, oarg{"company_id", "=", company_id}}, 0, 0, []string{"categ_id"})
		productID := -1
		categID := -1
		if len(product) == 1 {
			productID = int(product[0]["id"].(float64))
			categID = int(product[0]["categ_id"].([]interface{})[0].(float64))
		}
		location_id := o.GetID("stock.location", oarg{oarg{"complete_name", "=", v.Location}})
		uomID := o.GetID("uom.uom", oarg{oarg{"name", "=", v.ProductUOM}})

		r := o.GetID(umdll, oarg{oarg{"inventory_id", "=", invCountID}, oarg{"product_id", "=", productID}, oarg{"location_id", "=", location_id}, oarg{"company_id", "=", company_id}})

		ur := map[string]interface{}{
			"inventory_id": invCountID,
			"product_id":   productID,
			"location_id":  location_id,
			"company_id":   company_id,
			"product_qty":  v.ProductQty,
		}
		if categID != -1 {
			ur["categ_id"] = categID
		}
		if uomID != -1 {
			ur["product_uom_id"] = uomID
		}
		o.Log.Infow(umdll, "record", ur, "r", r)

		// o.Record(umdl, r, ur)
	}

	// q := `select id,name,state,company_id,prefill_counted_quantity,exhausted from stock_inventory;`
	// id |   name    |  state  | company_id | prefill_counted_quantity | exhausted
	// ---+-----------+---------+------------+--------------------------+-----------
	//  1 | Inventory | confirm |          1 | counted                  | f

	// q := "select id,is_editable,inventory_id,partner_id,product_id,product_uom_id,product_qty,categ_id,location_id,package_id,prod_lot_id,company_id,theoretical_qty from stock_inventory_line;"
	//  id | is_editable | inventory_id | partner_id | product_id | product_uom_id | product_qty | categ_id | location_id | package_id | prod_lot_id | company_id | theoretical_qty
	// ----+-------------+--------------+------------+------------+----------------+-------------+----------+-------------+------------+-------------+------------+-----------------
	//   1 | t           |            1 |          1 |          4 |              1 |   10.000000 |        1 |           8 |            |             |          1 |        0.000000

	// 	select distinct
	// '[G'||right(i.matnr,7)||'] '||m.maktx as product
	// ,case when i.sloc <> '0001' then trim(prefix||citycode||'/'||i.sloc)
	// else trim(prefix||citycode||'/Stock'||case when trim(i.sbin)<> '' then '/'||i.sbin else '' end)
	// end as location_id
	// ,i.unrestricted as product_qty
	// ,p.company
	// from sapdata.inventorytemp i
	// join sapdata.materials m on i.mandt = m.mandt and i.matnr = m.matnr
	// left join odoo.mat_dnutemp dnu on i.matnr = dnu.matnr
	// left join odoo.uom_uom uu on m.meins = uu.uom
	// join odoo.plantstemp p on i.werks = p.werks
	// where i.unrestricted <> 0
	// and i.werks = '1024' and p.company = 'Groupe A.R. Thomson'
	// and dnu.matnr is null
	// order by location_id,product

	// select distinct
	// '['||right(i.matnr,7)||'] '||m.maktx as product
	// ,case when i.sloc <> '0001' then trim(prefix||citycode||'/'||i.sloc)
	// else trim(prefix||citycode||'/Stock'||case when trim(i.sbin)<> '' then '/'||i.sbin else '' end)
	// end as location_id
	// ,i.unrestricted as product_qty
	// ,p.company
	// from sapdata.inventorytemp i
	// join sapdata.materials m on i.mandt = m.mandt and i.matnr = m.matnr
	// left join odoo.mat_dnutemp dnu on i.matnr = dnu.matnr
	// left join odoo.uom_uom uu on m.meins = uu.uom
	// join odoo.plantstemp p on i.werks = p.werks
	// where i.unrestricted <> 0
	// and i.werks not in ('1004','1013','1024') and p.company = 'A.R. Thomson Group'
	// and dnu.matnr is null
	// and m.mtart not in ('ROH','HALB')
	// order by location_id,product

	// select distinct
	// '[M'||right(i.matnr,7)||'] '||m.maktx as product
	// ,case when i.sloc <> '0001' then trim(prefix||citycode||'/'||i.sloc)
	// else trim(prefix||citycode||'/Stock'||case when trim(i.sbin)<> '' then '/'||i.sbin else '' end)
	// end as location_id
	// ,i.unrestricted as product_qty
	// ,p.company
	// from sapdata.inventorytemp i
	// join sapdata.materials m on i.mandt = m.mandt and i.matnr = m.matnr
	// left join odoo.mat_dnutemp dnu on i.matnr = dnu.matnr
	// left join odoo.uom_uom uu on m.meins = uu.uom
	// join odoo.plantstemp p on i.werks = p.werks
	// where i.unrestricted <> 0
	// and i.werks not in ('1004','1013','1024') and p.company = 'A.R. Thomson Group Manufacturing'
	// and dnu.matnr is null
	// and m.mtart in ('ROH','HALB')
	// order by location_id,product
}
