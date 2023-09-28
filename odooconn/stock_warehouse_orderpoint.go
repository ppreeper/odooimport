package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) StockReorderPoint() {
	mdl := "stock_warehouse_orderpoint"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v StockReorderPoint\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select default_code,company,warehouse,location,stock_min,stock_max,lot_size_min,uom,reorder_rule,reorder_rule_description,mtart
	from (
	select distinct
	right(m.matnr,7) as default_code
	,p.company
	,p.cityname||' '||p.ccode as warehouse
	,case when trim(i.sbin) = '' then trim(p.prefix||p.citycode||'/Stock') 
	when i.sbin is null then trim(p.prefix||p.citycode||'/Stock')
	else trim(p.prefix||p.citycode||'/Stock/'||i.sbin) end as location
	,i.stock_min,i.stock_max,i.lot_size_min
	,uu."name" as uom
	,i.reorder_rule,i.reorder_rule_description
	,m.mtart
	from sapdata.materials m
	join (
	select
	mandt,matnr,werks,'1030' salesorg,name1,sloc,sbin,unrestricted,returnqty,inqualinsp,"blocked",suom,suom_per_x,suom_per_y,perunit
	,unitcost,inventorycost,stock_min,stock_max,stock_rounding,lot_size_min,mrp_type,mrp_type_description,reorder_rule,reorder_rule_description
	from sapdata.inventorytemp i where salesorg = '1000'
	) i on m.mandt = i.mandt and m.matnr = i.matnr
	join odoo.uom_uom uu on m.meins = uu.uom
	left join odoo.mat_dnu dnu on m.matnr = dnu.matnr
	left join odoo.plantstemp p on i.salesorg = p.salesorg and i.werks = p.werks
	where dnu.matnr is null
	and i.werks not in ('1006','1013','1024')
	and p.salesorg = '1030'
	and i.stock_min <> 0
	and m.mtart in ('ROH','HALB')
	union 
	select distinct
	right(m.matnr,7) as default_code
	,p.company
	,p.cityname||' '||p.ccode as warehouse
	,case when trim(i.sbin) = '' then trim(p.prefix||p.citycode||'/Stock') 
	when i.sbin is null then trim(p.prefix||p.citycode||'/Stock')
	else trim(p.prefix||p.citycode||'/Stock/'||i.sbin) end as location
	,i.stock_min,i.stock_max,i.lot_size_min
	,uu."name" as uom
	,i.reorder_rule,i.reorder_rule_description
	,m.mtart
	from sapdata.materials m
	join (
	select
	mandt,matnr,werks,salesorg,name1,sloc,sbin,unrestricted,returnqty,inqualinsp,"blocked",suom,suom_per_x,suom_per_y,perunit
	,unitcost,inventorycost,stock_min,stock_max,stock_rounding,lot_size_min,mrp_type,mrp_type_description,reorder_rule,reorder_rule_description
	from sapdata.inventorytemp i where salesorg = '1000'
	) i on m.mandt = i.mandt and m.matnr = i.matnr
	join odoo.uom_uom uu on m.meins = uu.uom
	left join odoo.mat_dnu dnu on m.matnr = dnu.matnr
	left join odoo.plantstemp p on i.salesorg = p.salesorg and i.werks = p.werks
	where dnu.matnr is null
	and i.werks not in ('1006','1013','1024')
	and p.salesorg = '1000'
	and i.stock_min <> 0
	and m.mtart not in ('DIEN','ROH','HALB')
	union		
	select distinct
	right(m.matnr,7) as default_code
	,p.company
	,p.cityname||' '||p.ccode as warehouse
	,case when trim(i.sbin) = '' then trim(p.prefix||p.citycode||'/Stock') 
	when i.sbin is null then trim(p.prefix||p.citycode||'/Stock')
	else trim(p.prefix||p.citycode||'/Stock/'||i.sbin) end as location
	,i.stock_min,i.stock_max,i.lot_size_min
	,uu."name" as uom
	,i.reorder_rule,i.reorder_rule_description
	,m.mtart
	from sapdata.materials m
	join (
	select
	mandt,matnr,werks,salesorg,name1,sloc,sbin,unrestricted,returnqty,inqualinsp,"blocked",suom,suom_per_x,suom_per_y,perunit
	,unitcost,inventorycost,stock_min,stock_max,stock_rounding,lot_size_min,mrp_type,mrp_type_description,reorder_rule,reorder_rule_description
	from sapdata.inventorytemp i where salesorg = '1020'
	) i on m.mandt = i.mandt and m.matnr = i.matnr
	join odoo.uom_uom uu on m.meins = uu.uom
	left join odoo.mat_dnu dnu on m.matnr = dnu.matnr
	left join odoo.plantstemp p on i.salesorg = p.salesorg and i.werks = p.werks
	where dnu.matnr is null
	and i.werks in ('1024')
	and p.salesorg = '1020'
	and i.stock_min <> 0
	and m.mtart in ('FERT')
	) sor
	order by default_code,company
	`
	type Reorderpoint struct {
		DefaultCode            string  `db:"default_code"`
		Company                string  `db:"company"`
		Warehouse              string  `db:"warehouse"`
		Location               string  `db:"location"`
		StockMin               float64 `db:"stock_min"`
		StockMax               float64 `db:"stock_max"`
		LotSize                float64 `db:"lot_size_min"`
		UOM                    string  `db:"uom"`
		ReorderRule            string  `db:"reorder_rule"`
		ReorderRuleDescription string  `db:"reorder_rule_description"`
		Mtart                  string  `db:"mtart"`
	}
	var rr []Reorderpoint
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	companyIDs := o.ResCompanyMap()
	routeID := o.GetID("stock.location.route", oarg{oarg{"name", "=", "Buy"}})

	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Reorderpoint) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			companyID := companyIDs[v.Company]

			warehouseID := o.GetID("stock.warehouse", oarg{oarg{"name", "=", v.Warehouse}, oarg{"company_id", "=", companyID}})
			locationID := o.GetID("stock.location", oarg{oarg{"complete_name", "=", v.Location}, oarg{"company_id", "=", companyID}})
			productTmplID := o.GetID("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}})
			productID := o.GetID("product.product", oarg{oarg{"product_tmpl_id", "=", productTmplID}})

			r := o.GetID(umdl, oarg{
				oarg{"warehouse_id", "=", warehouseID},
				oarg{"location_id", "=", locationID},
				oarg{"product_id", "=", productID},
				oarg{"company_id", "=", companyID},
			})

			ur := map[string]interface{}{
				"product_id":      productID,
				"warehouse_id":    warehouseID,
				"location_id":     locationID,
				"company_id":      companyID,
				"route_id":        routeID,
				"product_min_qty": v.StockMin,
				"product_max_qty": v.StockMax,
				"qty_multiple":    v.LotSize,
				"trigger":         "manual",
			}
			o.Log.Infow(umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// StockWarehouseOrderpoint function
func (o *OdooConn) StockWarehouseOrderpoint() {
	mdl := "stock_warehouse_orderpoint"
	umdl := strings.Replace(mdl, "_", ".", -1)
	stmt := `
    select 
    "name",product_id,location_id
    ,product_min_qty,product_max_qty
    from odoo.artg_stock_warehouse_orderpointtemp 
    where product_min_qty <> 0 and product_max_qty <> 0
    order by location_id,product_id`
	rr := []struct {
		Name          string  `db:"name"`
		ProductID     string  `db:"product_id"`
		LocationID    string  `db:"location_id"`
		ProductMinQty float64 `db:"product_min_qty"`
		ProductMaxQty float64 `db:"product_max_qty"`
	}{}
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	o.Log.Infow(mdl, "model", umdl, "record", recs)
	// id   |         name         | trigger | active | snoozed_until | warehouse_id | location_id | product_id | product_category_id | product_min_qty | product_max_qty | qty_multiple | group_id | company_id | route_id | qty_to_order | create_uid |        create_date         | write_uid |         write_date         | bom_id | supplier_id
}
