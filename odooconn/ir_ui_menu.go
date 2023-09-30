package odooconn

import (
	"fmt"
	"sort"
	"strings"
)

// IrUiMenuSort function
func (o *OdooConn) IrUiMenuSort() {
	mdl := "ir_ui_menu"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v IrUiMenuSort\n", umdl)

	menuList := []string{}
	mm, err := o.SearchRead(umdl, oarg{oarg{"parent_id", "=", false}}, 0, 0, []string{"name", "sequence"})
	o.checkErr(err)

	for _, m := range mm {
		menuList = append(menuList, m["name"].(string))
	}
	sort.Slice(menuList, func(i, j int) bool { return strings.ToLower(menuList[i]) < strings.ToLower(menuList[j]) })
	for k, m := range menuList {

		r := -1
		for _, i := range mm {
			if i["name"].(string) == m {
				r = int(i["id"].(float64))
			}
		}

		ur := map[string]interface{}{
			"sequence": k + 1,
		}
		o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

		o.Record(umdl, r, ur)
	}
}
