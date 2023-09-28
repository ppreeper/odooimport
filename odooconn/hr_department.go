package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// HRDepartment function
func (o *OdooConn) HRDepartment() {
	mdl := "hr_department"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select distinct
	case when u.physicaldeliveryofficename = 'Montreal' then 'Groupe A.R. Thomson' when u.company = 'Delpro' then 'DelPro Automation Inc.' else u.company end company
	,u.department
	from artg.users u
	where u.company like any(array['A.R%','Del%'])
	and u.msexchresourcedisplay <> 'Room'
	and u.physicaldeliveryofficename not in ('','Datacenter','Consultant','Inactive')
	and u.department not in ('','Datacenter','Inactive')
	and u.name not like '%Superlok'
	and u.name not like '%Delpro'
	order by company,u.department
	`
	type HRDept struct {
		Company    string `db:"company"`
		Department string `db:"department"`
	}
	dd := []HRDept{}

	err := o.DB.Select(&dd, stmt)
	o.checkErr(err)
	recs := len(dd)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()

	// tasker
	wg.Add(len(dd))
	for _, u := range dd {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, u HRDept) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			company_id := cids[u.Company]
			r := o.GetID(umdl, oarg{oarg{"name", "=", u.Department}, oarg{"company_id", "=", company_id}})

			ur := map[string]interface{}{
				"name":       u.Department,
				"company_id": company_id,
			}
			o.Log.Infow(umdl, "u", u, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, u)
	}
	wg.Wait()
}
