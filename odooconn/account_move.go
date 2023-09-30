package odooconn

import (
	"fmt"
	"strings"
)

// AccountMoveLine function
func (o *OdooConn) AccountMove() {
	mdl := "account_move"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	// cids := o.ResCompanyMap()

	// sem := make(chan int, o.JobCount)
	// var wg sync.WaitGroup

	// bar := progressbar.Default(int64(len(dd)))
	// wg.Add(len(dd))
	// for _, u := range dd {
	// 	go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, u HRDept) {
	// 		defer bar.Add(1)
	// 		// bar.Add(1)
	// 		defer wg.Done()
	// 		sem <- 1

	// 		company_id := cids[u.Company]
	// 		r := o.GetID(umdl, [][]interface{}{{"name", "=", u.Department}, {"company_id", "=", company_id}})

	// 		var ur = map[string]interface{}{
	// 			"name":       u.Department,
	// 			"company_id": company_id,
	// 		}
	// 		o.Log.Info(umdl, "u", u, "record", ur, "r", r)

	// o.Record(umdl, r, ur)

	// 		<-sem
	// 	}(sem, &wg, bar, u)
	// }
	// wg.Wait()
}

// AccountMoveLine function
func (o *OdooConn) AccountMoveLine() {
	mdl := "account_move_line"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)
	// cids := o.ResCompanyMap()
}
