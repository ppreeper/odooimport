package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) AccountAccountUnlink() {
	mdl := "account_account"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v AccountAccountUnlink\n", umdl)

	type Account struct {
		Code      string `json:"code,omitempty" db:"code"`
		Aname     string `json:"aname,omitempty" db:"aname"`
		Atype     string `json:"atype,omitempty" db:"atype"`
		Reconcile bool   `json:"reconcile,omitempty" db:"reconcile"`
	}
	var dbrecs []Account

	stmt := ``

	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	odoorecs, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"code"})
	o.checkErr(err)

	fmt.Println("accounts:", len(odoorecs))

	// ids := []int{}
	for _, or := range odoorecs {
		for _, dr := range dbrecs {
			if or["code"] == dr.Code {
				// ids = append(ids, int(or["id"].(float64)))
				fmt.Println("unlink", or["id"], dr.Code, dr.Aname)
				o.Unlink(umdl, []int{int(or["id"].(float64))})
			}
		}
	}
}

// AccountAccount function
func (o *OdooConn) AccountAccount() {
	mdl := "account_account"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v AccountAccount\n", umdl)

	type Account struct {
		Code      string `json:"code,omitempty" db:"code"`
		Aname     string `json:"aname,omitempty" db:"aname"`
		Atype     string `json:"atype,omitempty" db:"atype"`
		Reconcile bool   `json:"reconcile,omitempty" db:"reconcile"`
	}
	var dbrecs []Account

	stmt := ``

	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	ats := o.AccountAccountTypeMap()

	// tasker
	recs := len(dbrecs)
	bar := progressbar.Default(int64(recs))
	for _, v := range dbrecs {
		r, err := o.GetID(umdl, oarg{oarg{"code", "=", v.Code}})
		o.checkErr(err)

		typeID := ats[v.Atype]

		ur := map[string]interface{}{
			"code":         v.Code,
			"name":         v.Aname,
			"user_type_id": typeID,
			"reconcile":    v.Reconcile,
		}

		o.Log.Debug(umdl, "ur", ur, "r", r)

		o.Record(umdl, r, ur)

		bar.Add(1)

	}
}
