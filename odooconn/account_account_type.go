package odooconn

import (
	"strings"
)

func (o *OdooConn) AccountAccountTypeMap() map[string]int {
	mdl := "account_account_type"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"name"})
	o.checkErr(err)

	cids := map[string]int{}
	for _, c := range cc {
		cids[c["name"].(string)] = int(c["id"].(float64))
	}
	return cids
}
