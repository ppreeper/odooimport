package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// MRPBomLine function
func (o *OdooConn) MRPBomOP(c string) {
	mdl := "mrp_routing_workcenter"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v MRPBomOP\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select default_code,company,plant,code,qty,uom,kit 
	from odoo.artg_bom_list_import
	-- where default_code = '3000745'
	-- limit 10	
	`
	// stmt = stmt + ` limit 10`
	type BOM struct {
		DefaultCode string  `db:"default_code"`
		Company     string  `db:"company"`
		Plant       string  `db:"plant"`
		Code        string  `db:"code"`
		Qty         float64 `db:"qty"`
		UOM         string  `db:"uom"`
		Kit         float64 `db:"kit"`
	}
	var rr []BOM
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()
	companyID := cids[c]

	bomopstmt := `
	select * from (
	select 
	right(abo.matnr,7) matnr 
	,abo.werks,abo.plnnr,abo.plnkn
	,'100' as seq
	,'setup' as name
	,case abo.werks
	when '1001' then workcenter||'-EDM'
	when '1003' then workcenter||'-SUR'
	when '1004' then workcenter||'-SAR'
	when '1014' then workcenter||'-FTM'
	when '1016' then workcenter||'-RED'
	else '' end workcenter_name
	,ltxa1,ltxa2 
	,abo.verwe duration_per_qty
	,abo.plnme
	,abo.lar02 machine
	,case abo.setup_uom
	when 'STD' then 60*abo.setup_time
	else abo.setup_time end as time_cycle_manual
	,case abo.setup_uom 
	when 'STD' then 'MIN'
	else abo.setup_uom end as op_uom
	from odoo.artg_bom_operationstemp abo 
	join (
	select abr.* 
	from odoo.artg_boms_reviewed abr
	join odoo.artg_parts_reviewed apr on abr.matnr = apr.matnr 
	left join odoo.artg_bom_kits abk on abr.matnr = right(abk.matnr,7)
	where trim(abr.plant_wins) = '1'
	and trim(apr.mfg_item) = '1'
	and abk.matnr is null
	) abr on right(abo.matnr,7) = abr.matnr and abo.werks = abr.plant
	where workcenter <> 'FINAL'
	union all
	select 
	right(abo.matnr,7) matnr 
	,abo.werks,abo.plnnr,abo.plnkn
	,'101' as seq
	,'machine' as name
	,case abo.werks
	when '1001' then workcenter||'-EDM'
	when '1003' then workcenter||'-SUR'
	when '1004' then workcenter||'-SAR'
	when '1014' then workcenter||'-FTM'
	when '1016' then workcenter||'-RED'
	else '' end workcenter_name
	,ltxa1,ltxa2 
	,abo.verwe duration_per_qty
	,abo.plnme
	,abo.lar01 machine
	,case abo.machine_uom
	when 'STD' then abo.machine_time*60
	when 'H' then abo.machine_time*60
	else abo.machine_time end as time_cycle_manual
	,case abo.machine_uom 
	when 'STD' then 'MIN'
	when 'H' then 'MIN'
	when 'S' then 'MIN'
	else abo.machine_uom end as op_uom
	from odoo.artg_bom_operationstemp abo 
	join (
	select abr.* 
	from odoo.artg_boms_reviewed abr
	join odoo.artg_parts_reviewed apr on abr.matnr = apr.matnr 
	left join odoo.artg_bom_kits abk on abr.matnr = right(abk.matnr,7)
	where trim(abr.plant_wins) = '1'
	and trim(apr.mfg_item) = '1'
	and abk.matnr is null
	) abr on right(abo.matnr,7) = abr.matnr and abo.werks = abr.plant
	where workcenter <> 'FINAL'
	) ops
	where matnr = $1
	and time_cycle_manual <> 0
	order by matnr,werks,plnnr,plnkn,seq
	`

	type BOMOP struct {
		Matnr           string `db:"matnr"`
		Werks           string `db:"werks"`
		Plnnr           string `db:"plnnr"`
		Plnkn           string `db:"plnkn"`
		Seq             string `db:"seq"`
		Name            string `db:"name"`
		WorkcenterName  string `db:"workcenter_name"`
		Ltxa1           string `db:"ltxa1"`
		Ltxa2           string `db:"ltxa2"`
		DurationPerQty  string `db:"duration_per_qty"`
		Plnme           string `db:"plnme"`
		Machine         string `db:"machine"`
		TimeCycleManual string `db:"time_cycle_manual"`
		OpUOM           string `db:"op_uom"`
	}

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v BOM) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			// o.Log.Info(umdl, "v", v, "companyID", companyID)
			productID, err := o.GetID("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}})
			o.checkErr(err)
			bomID, err := o.GetID("mrp.bom", oarg{oarg{"product_tmpl_id", "=", productID}, oarg{"company_id", "=", companyID}})
			o.checkErr(err)

			var bb []BOMOP
			err = o.DB.Select(&bb, bomopstmt, v.DefaultCode)
			// err := o.DB.Select(&bb, bomopstmt)
			o.checkErr(err)
			o.Log.Info(umdl, "bb", bb)
			for _, b := range bb {
				// o.Log.Info(umdl, "b", b)

				workcenterID, err := o.GetID("mrp.workcenter", oarg{
					oarg{"company_id", "=", companyID},
					oarg{"name", "=", b.WorkcenterName},
				})
				o.checkErr(err)

				r, err := o.GetID(umdl, oarg{
					oarg{"name", "=", b.Name},
					oarg{"workcenter_id", "=", workcenterID},
					oarg{"bom_id", "=", bomID},
					oarg{"company_id", "=", companyID},
				})
				o.checkErr(err)
				ur := map[string]interface{}{
					"name":              b.Name,
					"workcenter_id":     workcenterID,
					"bom_id":            bomID,
					"company_id":        companyID,
					"time_cycle_manual": b.TimeCycleManual,
					"duration_per_qty":  b.DurationPerQty,
				}

				o.Log.Info(umdl, "record", ur, "r", r)

				o.Record(umdl, r, ur)
			}

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
