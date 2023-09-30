package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// CrmTeam function
func (o *OdooConn) CrmTeam() {
	mdl := "crm_team"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	stmt := `
	select stl.team||' '||stl.division team,stl.leader,u.ename
	from ct.steamlead stl
	join sapdata.users u on stl.leader = u.username
	where stl.datestart::date <= current_date and current_date <= stl.dateend::date
	and stl.division <> 'ARTG'
	`
	type CrmTeam struct {
		Team   string `db:"team"`
		Leader string `db:"leader"`
		Ename  string `db:"ename"`
	}

	var rr []CrmTeam
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	wg.Add(recs)
	for _, v := range rr {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v CrmTeam) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Team}})
			o.checkErr(err)
			rp, err := o.GetID("res.partner", oarg{oarg{"name", "=", v.Ename}})
			o.checkErr(err)
			ru, err := o.GetID("res.users", oarg{oarg{"partner_id", "=", rp}})
			o.checkErr(err)

			ur := map[string]interface{}{
				"name":       v.Team,
				"company_id": 1,
				"user_id":    ru,
			}

			o.Log.Info(umdl, "ur", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}

// CrmTeamMembers function
func (o *OdooConn) CrmTeamMembers() {
	mdl := "crm_team"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v", umdl)

	stmt := `
	select sl.team||' '||sl.division team,sl.ename
	from ct.steammember_list sl
	where sl.datestart::date <= current_date and current_date <= sl.dateend::date
	and division <> 'ARTG'
	order by team,ename
	`
	type CrmTeam struct {
		Team  string `db:"team"`
		Ename string `db:"ename"`
	}

	var rr []CrmTeam
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)
	recs := len(rr)
	bar := progressbar.Default(int64(recs))

	// tasker
	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup
	wg.Add(recs)
	for _, v := range rr {
		// process
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, v CrmTeam) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			o.Log.Info("", "team", v)
			r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Team}})
			o.checkErr(err)
			rp, err := o.GetID("res.users", oarg{oarg{"name", "=", v.Ename}})
			o.checkErr(err)
			team, err := o.SearchRead(umdl, oarg{oarg{"name", "=", v.Team}}, 0, 0, []string{"member_ids"})
			o.checkErr(err)
			mids := team[0]["member_ids"]

			var mIDS []int
			if !valInSlice(rp, mIDS) {
				mIDS = append(mIDS, rp)
			}

			for _, v := range mids.([]interface{}) {
				mIDS = append(mIDS, int(v.(float64)))
			}

			ur := map[string]interface{}{
				"name":       v.Team,
				"member_ids": mIDS,
			}
			o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

			<-sem
		}(sem, &wg, bar, v)
	}
	wg.Wait()
}
