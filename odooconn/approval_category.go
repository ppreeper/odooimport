package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func (o *OdooConn) ApprovalCategory() {
	mdl := "approval_category"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v Approvals\n", umdl)

	type Record struct {
		Name            string `json:"name,omitempty" db:"name"`
		SequenceCode    string `json:"sequence_code,omitempty" db:"sequence_code"`
		ApprovalMinimum int    `json:"approval_minimum,omitempty" db:"approval_minimum"`
		Company         string `json:"company,omitempty" db:"company"`
	}
	var dbrecs []Record

	stmt := ``

	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	cids, err := o.ModelMap("res.company", "name")
	o.checkErr(err)

	// tasker
	recs := len(dbrecs)
	bar := progressbar.Default(int64(recs))
	for _, v := range dbrecs {
		cid := cids[v.Company]

		r, err := o.GetID(umdl, oarg{oarg{"name", "=", v.Name}, oarg{"company_id", "=", cid}})
		o.checkErr(err)

		ur := map[string]interface{}{
			"name":             v.Name,
			"sequence_code":    v.SequenceCode,
			"approval_minimum": v.ApprovalMinimum,
			"company_id":       cid,
		}

		o.Log.Debug(umdl, "ur", ur, "r", r)

		o.Record(umdl, r, ur)

		bar.Add(1)
	}
}
