package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ResPartnerCategory function
func (o *OdooConn) ResPartnerCategory() {
	mdl := "res_partner_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `select region,bztxt rname from sapdata.salesregions`

	type PartnerCategory struct {
		Region string `db:"region"`
		Rname  string `db:"rname"`
	}
	var rr []PartnerCategory
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v PartnerCategory) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1
			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Rname}})
			o.checkErr(err)

			ur := map[string]interface{}{"name": v.Rname}
			o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
