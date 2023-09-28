package odooconn

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// AccountAccount function
func (o *OdooConn) AccountGroup() {
	mdl := "account_group"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v ProductTemplate\n", umdl)

	stmt := `select code_prefix_start,code_prefix_end,name from odoo.account_group order by code_prefix_start`

	type Group struct {
		CodePrefixStart string `json:"code_prefix_start,omitempty" db:"code_prefix_start"`
		CodePrefixEnd   string `json:"code_prefix_end,omitempty" db:"code_prefix_end"`
		Name            string `json:"name,omitempty" db:"name"`
	}
	var dbrecs []Group
	if stmt == "" {
		return
	}
	o.Log.Info(stmt)
	err := o.DB.Select(&dbrecs, stmt)
	o.checkErr(err)

	// tasker
	recs := len(dbrecs)
	bar := progressbar.Default(int64(recs))
	for _, v := range dbrecs {
		r := o.GetID(umdl, oarg{oarg{"code_prefix_start", "=", v.CodePrefixStart}})

		ur := map[string]interface{}{
			"code_prefix_start": v.CodePrefixStart,
			"code_prefix_end":   v.CodePrefixEnd,
			"name":              v.Name,
		}

		o.Log.Infow(umdl, "ur", ur, "r", r)

		o.Record(umdl, r, ur)

		bar.Add(1)

	}
}
