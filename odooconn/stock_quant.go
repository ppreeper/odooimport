package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) StockQuantUnlink(pageSize int) {
	// inventory_unlink
	mdl := "stock_quant"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v StockQuant %v\n", umdl, pageSize)
	odoorecs := o.Search(umdl, oarg{})
	ddList := getPages(odoorecs, pageSize)

	bar := progressbar.Default(int64(len(ddList)))
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(len(ddList))
	for _, r := range ddList {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, r []int) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Unlink(umdl, r)

			// bar.Add(1)
			<-sem
		}(sem, &wg, bar, r)
	}
	wg.Wait()
}

// StockInventory function to create an inventory adjustment
// you can find the inventory adjustments in
// the Inventory > Operations > Inventory Adjustments
func (o *OdooConn) StockQuant(company string) {
	mdl := "stock_quant"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v StockQuant %v\n", umdl, company)

	salesorg := "1000"

	stmt := ``

	type Line struct {
		Company     string  `db:"company"`
		Pprefix     string  `db:"pprefix"`
		Location    string  `db:"location"`
		Sbin        string  `db:"sbin"`
		DefaultCode string  `db:"default_code"`
		Matkl       string  `db:"matkl"`
		Quantity    float64 `db:"quantity"`
		CategoryID  string  `db:"category_id"`
		UomId       string  `db:"uom_id"`
	}
	rr := []Line{}
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)

	// pgs := o.ProductCategoryMap()
	// uom := o.UomMapper()
	cids := o.ModelMap("res.company", "name")

	recs := len(rr)
	bar := progressbar.Default(int64(recs))
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Line) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			productTmplId := o.GetID("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}})
			productId := o.GetID("product.product", oarg{oarg{"product_tmpl_id", "=", productTmplId}, oarg{"default_code", "=", v.DefaultCode}})

			companyId := cids[v.Company]

			// categID := -1
			// if v.Matkl != "" {
			// 	categID = pgs[v.Matkl]
			// }

			// catID := uom[v.CategoryID].ID
			// uomID := uom[v.CategoryID].Units[v.UomId]

			location := v.Pprefix + "/" + v.Location
			if v.Sbin != "" {
				location += "/" + v.Sbin
			}

			locationId := o.GetID("stock.location", oarg{oarg{"complete_name", "=", location}})

			r := o.GetID(umdl, oarg{oarg{"product_id", "=", productId}, oarg{"company_id", "=", companyId}, oarg{"location_id", "=", locationId}})

			ur := map[string]interface{}{
				"product_id":             productId,
				"company_id":             companyId,
				"location_id":            locationId,
				"inventory_quantity":     v.Quantity,
				"inventory_quantity_set": true,
			}
			o.Log.Debugw(umdl, "r", r, "ur", ur)

			row := o.Record(umdl, r, ur)
			if row == -1 {
				o.Log.Errorw(umdl, "v", v, "ur", ur)
			}

			// err := bar.Add(1)
			// o.checkErr(err)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

func (o *OdooConn) StockQuantConsignment() {
	mdl := "stock_quant"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v StockQuantConsignment\n", umdl)

	stmt := ``

	type Line struct {
		Company     string  `db:"company"`
		Pprefix     string  `db:"pprefix"`
		Location    string  `db:"location"`
		Sbin        string  `db:"sbin"`
		DefaultCode string  `db:"default_code"`
		Matkl       string  `db:"matkl"`
		Quantity    float64 `db:"quantity"`
		CategoryID  string  `db:"category_id"`
		UomId       string  `db:"uom_id"`
		Shipto      string  `db:"shipto"`
	}
	rr := []Line{}
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)

	cids := o.ModelMap("res.company", "name")

	recs := len(rr)
	bar := progressbar.Default(int64(recs))
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(recs)
	for _, v := range rr {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v Line) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			productTmplId := o.GetID("product.template", oarg{oarg{"default_code", "=", v.DefaultCode}})
			productId := o.GetID("product.product", oarg{oarg{"product_tmpl_id", "=", productTmplId}, oarg{"default_code", "=", v.DefaultCode}})

			ownerID := o.GetID("res.partner", oarg{oarg{"ref", "=", v.Shipto}})

			companyId := cids[v.Company]

			location := v.Location
			// Consignment Location cannot be a global location, we need to have a consignment location under the branches
			// location := v.Pprefix + "/" + v.Location
			// if v.Sbin != "" {
			// 	location += "/" + v.Sbin
			// }

			locationId := o.GetID("stock.location", oarg{oarg{"complete_name", "=", location}})

			r := o.GetID(umdl, oarg{oarg{"product_id", "=", productId}, oarg{"company_id", "=", companyId}, oarg{"location_id", "=", locationId}})

			ur := map[string]interface{}{
				"product_id":             productId,
				"company_id":             companyId,
				"location_id":            locationId,
				"inventory_quantity":     v.Quantity,
				"inventory_quantity_set": true,
				"owner_id":               ownerID,
			}

			o.Log.Debugw(umdl, "r", r, "ur", ur)
			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
