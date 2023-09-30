package odooconn

import (
	"fmt"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type LDAPUser struct {
	C                          string `db:"c"`
	CN                         string `db:"cn"`
	CO                         string `db:"co"`
	Company                    string `db:"company"`
	Department                 string `db:"department"`
	DisplayName                string `db:"displayname"`
	DistinguishedName          string `db:"distinguishedname"`
	FacsimileTelephoneNumber   string `db:"facsimiletelephonenumber"`
	GivenName                  string `db:"givenname"`
	L                          string `db:"l"`
	Mail                       string `db:"mail"`
	MailNickname               string `db:"mailnickname"`
	Manager                    string `db:"manager"`
	MiddleName                 string `db:"middlename"`
	Mobile                     string `db:"mobile"`
	Name                       string `db:"name"`
	Pager                      string `db:"pager"`
	PhysicalDeliveryOfficeName string `db:"physicaldeliveryofficename"`
	PostalCode                 string `db:"postalcode"`
	SAMAccountName             string `db:"samaccountname"`
	SN                         string `db:"sn"`
	ST                         string `db:"st"`
	StreetAddress              string `db:"streetaddress"`
	TelephoneNumber            string `db:"telephonenumber"`
	Title                      string `db:"title"`
	UserPrincipalName          string `db:"userprincipalname"`
	WWWHomePage                string `db:"wwwhomepage"`
	PropertyWarehouseID        string `db:"property_warehouse_id,omitempty"`
	TZ                         string `db:"tz,omitempty"`
}

var LDAPUserQuery = `
	select c,cn,co
	,case when u.physicaldeliveryofficename = 'Montreal' then 'Groupe A.R. Thomson' when u.company = 'Delpro' then 'DelPro Automation Inc.' else u.company end company
	,department,displayname,distinguishedname
	,facsimiletelephonenumber,givenname,l,mail,mailnickname,manager
	,middlename,mobile,u.name,pager
	,case
		when u.physicaldeliveryofficename = 'Concepcion' then 'Surrey'
		when u.physicaldeliveryofficename = 'Clairmont' then 'Grande Prairie'
		else physicaldeliveryofficename end physicaldeliveryofficename
	,postalcode,samaccountname,sn,st,streetaddress,telephonenumber
	,title,userprincipalname,wwwhomepage
	,'' property_warehouse_id
	,case when u.physicaldeliveryofficename = 'Concepcion' then 'America/Santiago' else cpl.tz end tz
	from artg.users u
	left join (select distinct name,tz from ct.company_plant_list) cpl on u.physicaldeliveryofficename = cpl.name
	where u.company like any(array['A.R%','Del%'])
	and u.msexchresourcedisplay <> 'Room'
	and u.extensionattribute15 not like '%positional%'
	and u.physicaldeliveryofficename not in ('','Datacenter','Consultant','Inactive')
	and u.department not in ('','Datacenter','Inactive')
	and u.name not like '%Superlok'
	and u.name not like '%Delpro'
	order by company,cn
	-- limit 5
	`

// ResUsers function run the ldap2sql update first
func (o *OdooConn) ResUsers() {
	mdl := "res_users"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v\n", umdl)

	sem := make(chan int, o.JobCount)
	var wg sync.WaitGroup

	var uu []LDAPUser
	err := o.DB.Select(&uu, LDAPUserQuery)
	o.checkErr(err)
	recs := len(uu)
	bar := progressbar.Default(int64(recs))

	cids := o.ResCompanyMap()
	c1 := cids["A.R. Thomson Group"]
	c2 := cids["Groupe A.R. Thomson"]
	c4 := cids["DelPro Automation Inc."]
	c5 := cids["DelPro Technical Inc."]

	// tasker
	wg.Add(recs)
	for _, u := range uu {
		go func(sem chan int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, u LDAPUser) {
			defer bar.Add(1)
			defer wg.Done()
			sem <- 1

			name := u.CN
			upn := u.UserPrincipalName
			warehouseID, err := o.GetID("stock.warehouse", oarg{oarg{"name", "=", u.PropertyWarehouseID}})
			o.checkErr(err)
			var cid int
			var ccids []int
			switch u.Company {
			case "A.R. Thomson Group":
				cid = c1
				ccids = []int{c1, c2}
				switch u.Department {
				case "Accounting":
					ccids = []int{c1, c2, c4, c5}
				case "IT":
					ccids = []int{c1, c2, c4, c5}
				}
			case "Groupe A.R. Thomson":
				cid = c1
				ccids = []int{c1, c2}
			case "DelPro Automation Inc.":
				cid = c1
				ccids = []int{c4, c5, c1}
			default:
				cid = -1
			}

			r, err := o.GetID(umdl, oarg{oarg{"name", "=", name}})
			o.checkErr(err)
	
			ur := map[string]interface{}{
				"name":                  name,
				"login":                 upn,
				"company_id":            cid,
				"company_ids":           ccids,
				"property_warehouse_id": warehouseID,
				"tz":                    u.TZ,
			}

			o.Log.Info(mdl, "model", umdl, "record", ur, "rid", r)

			if r != 2 {
				o.Record(umdl, r, ur)
			}

			<-sem
		}(sem, &wg, bar, u)
	}
	wg.Wait()
}

func (o *OdooConn) ResUsersMap() map[string]int {
	mdl := "res_users"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"name"})
	o.checkErr(err)
	cids := map[string]int{}
	for _, c := range cc {
		cids[c["name"].(string)] = int(c["id"].(float64))
	}
	return cids
}
