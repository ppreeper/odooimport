package odooconn

import (
	"strings"
)

func (o *OdooConn) ResPartnerMap() map[string]int {
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"name"})
	o.checkErr(err)
	cids := map[string]int{}
	for _, c := range cc {
		cids[c["name"].(string)] = int(c["id"].(float64))
	}
	return cids
}

func (o *OdooConn) ResPartnerRefMap() map[string]int {
	mdl := "res_partner"
	umdl := strings.Replace(mdl, "_", ".", -1)
	cc, err := o.SearchRead(umdl, oarg{}, 0, 0, []string{"ref"})
	o.checkErr(err)
	cids := map[string]int{}
	for _, c := range cc {
		cids[c["ref"].(string)] = int(c["id"].(float64))
	}
	return cids
}
