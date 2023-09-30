package odooconn

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) StockPickingTypeUnlink() {
	mdl := "stock_picking_type"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v StockPickingTypeUnlink\n", umdl)

	stmt := ``
	type Rec struct {
		Pprefix string `db:"pprefix"`
		Pname   string `db:"pname"`
		Company string `db:"company"`
	}
	var dbrecs []Rec
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	odoorecs, err := o.SearchRead(umdl, oarg{"active", "=", true}, 0, 0, []string{"warehouse_id", "company_id", "active"})
	o.checkErr(err)
	fmt.Println("odoorecs", len(odoorecs))
	type ORec struct {
		Row       int
		Warehouse string
		Company   string
		Active    bool
	}

	var orecs []ORec
	for _, v := range odoorecs {
		rid := int(v["id"].(float64))

		var warehouse_id string
		switch vv := v["warehouse_id"].(type) {
		case string:
			warehouse_id = vv
		case bool:
			warehouse_id = strconv.FormatBool(vv)
		case interface{}:
			warehouse_id = vv.([]interface{})[1].(string)
		default:
			warehouse_id = ""
		}

		var company_id string
		switch vv := v["company_id"].(type) {
		case string:
			company_id = vv
		case bool:
			company_id = strconv.FormatBool(vv)
		case interface{}:
			company_id = vv.([]interface{})[1].(string)
		default:
			company_id = ""
		}

		var active bool
		switch vv := v["active"].(type) {
		case bool:
			active = vv
		default:
			active = false
		}

		orecs = append(orecs, ORec{Row: rid, Warehouse: warehouse_id, Company: company_id, Active: active})
	}

	odoorecs, err = o.SearchRead(umdl, oarg{"active", "=", false}, 0, 0, []string{"warehouse_id", "company_id", "active"})
	o.checkErr(err)
	fmt.Println("odoorecs", len(odoorecs))
	for _, v := range odoorecs {
		rid := int(v["id"].(float64))

		var warehouse_id string
		switch vv := v["warehouse_id"].(type) {
		case string:
			warehouse_id = vv
		case bool:
			warehouse_id = strconv.FormatBool(vv)
		case interface{}:
			warehouse_id = vv.([]interface{})[1].(string)
		default:
			warehouse_id = ""
		}

		var company_id string
		switch vv := v["company_id"].(type) {
		case string:
			company_id = vv
		case bool:
			company_id = strconv.FormatBool(vv)
		case interface{}:
			company_id = vv.([]interface{})[1].(string)
		default:
			company_id = ""
		}

		var active bool
		switch vv := v["active"].(type) {
		case bool:
			active = vv
		default:
			active = false
		}

		orecs = append(orecs, ORec{Row: rid, Warehouse: warehouse_id, Company: company_id, Active: active})
	}
	fmt.Println("orecs", len(orecs))

	var dlist []int
	for _, d := range dbrecs {
		for _, o := range orecs {
			if strings.Contains(o.Warehouse, d.Pname) && strings.Contains(o.Company, d.Company) {
				dlist = append(dlist, o.Row)
				fmt.Println(o)
			}
		}
	}
	sort.Ints(dlist)
	dlist = removeDuplicate(dlist)
	fmt.Println(dlist)
	fmt.Println(umdl, "remove", len(dlist), "records")
	for _, v := range dlist {
		o.Unlink(umdl, []int{v})
	}
}

func (o *OdooConn) StockPickingType() {
	mdl := "stock_picking_type"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v StockPickingType\n", umdl)

	stmt := ``

	type StockPickingType struct {
		Name                  string `db:"name"`
		SequenceId            string `db:"sequence_id"`
		SequenceCode          string `db:"sequence_code"`
		DefaultLocationSrcId  string `db:"default_location_src_id"`
		DefaultLocationDestId string `db:"default_location_dest_id"`
		Code                  string `db:"code"`
		ReturnPickingTypeId   string `db:"return_picking_type_id"`
		WarehouseId           string `db:"warehouse_id"`
		PrintLabel            bool   `db:"print_label"`
		ShowOperations        bool   `db:"show_operations"`
		ShowReserved          bool   `db:"show_reserved"`
		ReservationMethod     string `db:"reservation_method"`
		CompanyId             string `db:"company_id"`
		CompanyBranchId       string `db:"company_branch_id"`
	}
	var dbrecs []StockPickingType
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	fmt.Println("dbrecs", len(dbrecs))

	companyIDs, err := o.ModelMap("res.company", "name")
	o.checkErr(err)

	// tasker
	recs := len(dbrecs)
	bar := progressbar.Default(int64(recs))
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(recs)
	for _, v := range dbrecs {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v StockPickingType) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid := companyIDs[v.CompanyId]
			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"company_id", "=", cid}})
			o.checkErr(err)
			sequenceID, err := o.GetID("ir.sequence", oarg{oarg{"name", "like", v.WarehouseId}, oarg{"company_id", "=", cid}, oarg{"prefix", "like", v.SequenceCode}})
			o.checkErr(err)
			if v.SequenceCode == "DS" {
				sequenceID, err = o.GetID("ir.sequence", oarg{oarg{"name", "like", v.SequenceId}, oarg{"company_id", "=", cid}, oarg{"prefix", "like", v.SequenceCode}})
				o.checkErr(err)
			}
			warehouseID, err := o.GetID("stock.warehouse", oarg{oarg{"name", "=", v.WarehouseId}, oarg{"company_id", "=", cid}})
			o.checkErr(err)
			companyBranchId, err := o.GetID("res.company.branch", oarg{oarg{"name", "=", v.WarehouseId}, oarg{"company_id", "=", cid}})
			o.checkErr(err)

			defaultLocationSrcId, err := o.GetID("stock.location", oarg{oarg{"complete_name", "=", v.DefaultLocationSrcId}, oarg{"company_id", "=", cid}})
			o.checkErr(err)
			defaultLocationDestId, err := o.GetID("stock.location", oarg{oarg{"complete_name", "=", v.DefaultLocationDestId}, oarg{"company_id", "=", cid}})
			o.checkErr(err)

			o.Log.Error(v.DefaultLocationSrcId, defaultLocationSrcId, v.DefaultLocationDestId, defaultLocationDestId)

			ur := map[string]interface{}{
				"name":          v.Name,
				"code":          v.Code,
				"sequence_id":   sequenceID,
				"sequence_code": v.SequenceCode,
				"company_id":    cid,
			}

			if warehouseID != -1 {
				ur["warehouse_id"] = warehouseID
			}

			if defaultLocationSrcId != -1 {
				ur["default_location_src_id"] = defaultLocationSrcId
			}

			if defaultLocationDestId != -1 {
				ur["default_location_dest_id"] = defaultLocationDestId
			}
			if companyBranchId != -1 {
				ur["company_branch_id"] = companyBranchId
			}

			o.Log.Debug(umdl, "r", r, "ur", ur)

			o.Record(umdl, r, ur)

			// bar.Add(1)
			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
