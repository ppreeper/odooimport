package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// HREmployee function
func (o *OdooConn) HREmployee() {
	mdl := "hr_employee"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	var uu []LDAPUser
	err := o.DB.Select(&uu, LDAPUserQuery)
	o.checkErr(err)
	recs := len(uu)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()

	// tasker
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(len(uu))
	for _, u := range uu {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, u LDAPUser) {
			defer bar.Add(1)
			// bar.Add(1)
			defer wg.Done()
			sem <- 1
			// fmt.Println(u)

			o.Log.Info("", "user", u.CN, "office", u.PhysicalDeliveryOfficeName)
			cid := cids[u.Company]

			name := u.CN
			// office := u.PhysicalDeliveryOfficeName
			// fmt.Println(name, office, u.Department)

			resPartnerID, err := o.GetID("res.partner", oarg{oarg{"name", "=", name}})
			o.checkErr(err)
			resUserID, err := o.GetID("res.users", oarg{oarg{"partner_id", "=", resPartnerID}, oarg{"company_id", "=", cid}})
			o.checkErr(err)
			// fmt.Println(resPartnerID, resUserID)

			user, err := o.SearchRead("res.partner", oarg{oarg{"id", "=", resPartnerID}}, 0, 1, []string{"parent_id"})
			o.checkErr(err)
			// fmt.Println(user)

			p := user[0]["parent_id"].([]interface{})
			parentID := int(p[0].(float64))
			departmentID, err := o.GetID("hr.department", oarg{oarg{"name", "=", u.Department}, oarg{"company_id", "=", cid}})
			o.checkErr(err)

			// manager := ""
			// if u.Manager != "" {
			// 	mcn := strings.Split(u.Manager, ",")
			// 	m := strings.Split(mcn[0], "=")
			// 	manager = m[1]
			// }
			// managerID := o.GetID(umdl, oarg{oarg{"name", "=", manager}})

			r, err := o.GetID(umdl, oarg{oarg{"name", "=", name}, oarg{"company_id", "=", cid}})
			o.checkErr(err)

			ur := map[string]interface{}{
				"name":       name,
				"user_id":    resUserID,
				"company_id": cid,
				// "work_location": office,
				"address_id":    parentID,
				"department_id": departmentID,
				"job_title":     u.Title,
				"mobile_phone":  u.Mobile,
			}

			// if manager != "" {
			// 	ur["parent_id"] = managerID
			// }

			o.Log.Info(umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, u)
	}
	wg.Wait()
}

// HREmployeeManager function
func (o *OdooConn) HREmployeeManager() {
	mdl := "hr_employee"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	var uu []LDAPUser
	err := o.DB.Select(&uu, LDAPUserQuery)
	o.checkErr(err)
	bar := progressbar.Default(int64(len(uu)))

	cids := o.ResCompanyMap()

	// tasker
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(len(uu))
	for _, u := range uu {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, u LDAPUser) {
			defer bar.Add(1)
			// bar.Add(1)
			defer wg.Done()
			sem <- 1

			cid := cids[u.Company]
			name := u.CN
			r, err := o.GetID(umdl, oarg{oarg{"name", "=", name}, oarg{"company_id", "=", cid}})
			o.checkErr(err)

			manager := ""
			if u.Manager != "" {
				mcn := strings.Split(u.Manager, ",")
				m := strings.Split(mcn[0], "=")
				manager = m[1]
			}
			managerID, err := o.GetID(umdl, oarg{oarg{"name", "=", manager}})
			o.checkErr(err)

			ur := map[string]interface{}{}

			if manager != "" {
				ur["parent_id"] = managerID
				o.Log.Info(umdl, "record", ur, "r", r)
				o.Record(umdl, r, ur)
			}

			<-sem
		}(sem, &wg, bar, u)
	}
	wg.Wait()
}
