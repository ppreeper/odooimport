package odooconn

import "strings"

// ResCompanyLDAP function
func (o *OdooConn) ResCompanyLDAP(cid int) {
	mdl := "res_company_ldap"
	umdl := strings.Replace(mdl, "_", ".", -1)

	stmt := `
	select
	"company_id"
	,"sequence"
	,"ldap_server"
	,"ldap_server_port"
	,"ldap_binddn"
	,"ldap_password"
	,"ldap_filter"
	,"ldap_base"
	,"create_user"
	,"ldap_tls"
	from odoo.company_list_ldap
	`
	type LdapConn struct {
		CompanyID      int    `db:"company_id"`
		Sequence       string `db:"sequence"`
		LdapServer     string `db:"ldap_server"`
		LdapServerPort int    `db:"ldap_server_port"`
		LdapBinddn     string `db:"ldap_binddn"`
		LdapPassword   string `db:"ldap_password"`
		LdapFilter     string `db:"ldap_filter"`
		LdapBase       string `db:"ldap_base"`
		CreateUser     bool   `db:"create_user"`
		LdapTLS        bool   `db:"ldap_tls"`
	}
	var rr []LdapConn
	err := o.DB.Select(&rr, stmt)
	o.checkErr(err)

	for _, v := range rr {
		if v.CompanyID == cid {
			r, err := o.GetID(umdl, oarg{oarg{"company", "=", cid}, oarg{"sequence", "=", v.Sequence}})
			o.checkErr(err)
			ur := map[string]interface{}{
				"sequence":         v.Sequence,
				"company":          cid,
				"ldap_server":      v.LdapServer,
				"ldap_server_port": v.LdapServerPort,
				"ldap_binddn":      v.LdapBinddn,
				"ldap_password":    v.LdapPassword,
				"ldap_filter":      v.LdapFilter,
				"ldap_base":        v.LdapBase,
				"ldap_tls":         v.LdapTLS,
				"create_user":      v.CreateUser,
			}
			o.Log.Info(mdl, "model", umdl, "record", ur, "r", r)

			o.Record(umdl, r, ur)

		}
	}
}
