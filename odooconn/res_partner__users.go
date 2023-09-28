package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// ResPartnerUsers function
func (o *OdooConn) ResPartnerUsers() {
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v users\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	var uu []LDAPUser
	err := o.DB.Select(&uu, LDAPUserQuery)
	o.checkErr(err)
	recs := len(uu)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, u := range uu {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, u LDAPUser) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			name := u.CN
			r := o.GetID(umdl, oarg{oarg{"name", "=", name}})
			physOffice := u.PhysicalDeliveryOfficeName

			user := o.SearchRead("res.users", oarg{oarg{"name", "=", name}}, 0, 1, []string{})
			c := user[0]["company_id"].([]interface{})
			cid := int(c[0].(float64))
			company := o.SearchRead("res.company", oarg{oarg{"id", "=", cid}}, 0, 1, []string{"name"})
			company_id := company[0]["id"]
			company_name := company[0]["name"]
			parent_id := o.GetID("res.partner", oarg{oarg{"name", "=", company_name}})
			pid := o.GetID("res.partner", oarg{oarg{"name", "like", physOffice}, oarg{"parent_id", "=", parent_id}})

			o.Log.Infow("", "name", name, "company_id", cid, "physOffice", physOffice, "company", company, "company_id", company_id, "pid", pid)

			ur := map[string]interface{}{
				"name":      name,
				"parent_id": pid,
				"email":     u.Mail,
				"phone":     u.TelephoneNumber,
				"mobile":    u.Mobile,
				"function":  u.Title,
				"website":   u.WWWHomePage,
			}

			o.Log.Infow(umdl, "ur", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, u)
	}
	wg.Wait()
}
