package ldap

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"
)

// LDAPconf config structure
type LDAPconf struct {
	Host     string `json:"host,omitempty"`
	BindDN   string `json:"binddn,omitempty"`
	BindPass string `json:"bindpass,omitempty"`
	Base     string `json:"base,omitempty"`
	Filter   string `json:"filter,omitempty"`
}

// LDAPServer struct
type LDAPServer struct {
	URL      string
	BindDN   string
	BindPass string
	Base     string
	Conn     *ldap.Conn
	Log      *zap.SugaredLogger
}

// Dial connect to ldap server
func (s *LDAPServer) Dial() error {
	l, err := ldap.DialURL(s.URL)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	s.Conn = l
	return err
}

// Bind to LDAP server
func (s *LDAPServer) Bind(user, pass string) error {
	err := s.Conn.Bind(user, pass)
	return err
}

// Search the ldap server
// func (s *LDAPServer) Search(filter string) []map[string]interface{} {
// 	err := s.Dial()
// 	s.checkErr(err)
// 	defer s.Conn.Close()
// 	err = s.Bind(s.BindDN, s.BindPass)
// 	s.checkErr(err)
// 	filter = strings.Replace(filter, "\r\n", "", -1)
// 	filter = strings.Replace(filter, "\n", "", -1)
// 	attributes := []string{"*"}
// 	sreq := ldap.NewSearchRequest(
// 		s.Base,
// 		ldap.ScopeWholeSubtree,
// 		ldap.NeverDerefAliases,
// 		0, 0, false,
// 		filter, attributes, nil,
// 	)
// 	sres, err := s.Conn.Search(sreq)
// 	s.checkErr(err)
// 	fmt.Println(len(sres.Entries))
// 	var uu []map[string]interface{}
// 	for _, v := range sres.Entries {
// 		if len(v.Attributes) > 0 {
// 			m := make(map[string]interface{})
// 			for _, attr := range v.Attributes {
// 				if len(attr.Values) == 1 {
// 					m[attr.Name] = attr.Values[0]
// 				} else if len(attr.Values) > 1 {
// 					m[attr.Name] = attr.Values
// 				}
// 			}
// 			uu = append(uu, m)
// 		}
// 	}
// 	return uu
// }

type User struct {
	C                          string `json:"c,omitempty"`
	CN                         string `json:"cn,omitempty"`
	CO                         string `json:"co,omitempty"`
	Company                    string `json:"company,omitempty"`
	Department                 string `json:"department,omitempty"`
	DisplayName                string `json:"display_name,omitempty"`
	DistinguishedName          string `json:"distinguished_name,omitempty"`
	FacsimileTelephoneNumber   string `json:"facsimile_telephone_number,omitempty"`
	GivenName                  string `json:"given_name,omitempty"`
	L                          string `json:"l,omitempty"`
	Mail                       string `json:"mail,omitempty"`
	MailNickname               string `json:"mail_nickname,omitempty"`
	Manager                    string `json:"manager,omitempty"`
	MiddleName                 string `json:"middle_name,omitempty"`
	Mobile                     string `json:"mobile,omitempty"`
	Name                       string `json:"name,omitempty"`
	Pager                      string `json:"pager,omitempty"`
	PhysicalDeliveryOfficeName string `json:"physical_delivery_office_name,omitempty"`
	PostalCode                 string `json:"postal_code,omitempty"`
	SAMAccountName             string `json:"sam_account_name,omitempty"`
	SN                         string `json:"sn,omitempty"`
	ST                         string `json:"st,omitempty"`
	StreetAddress              string `json:"street_address,omitempty"`
	TelephoneNumber            string `json:"telephone_number,omitempty"`
	Title                      string `json:"title,omitempty"`
	UserPrincipalName          string `json:"user_principal_name,omitempty"`
	WWWHomePage                string `json:"www_home_page,omitempty"`
}

// Search the ldap server
func (s *LDAPServer) Search(filter string) []User {
	err := s.Dial()
	s.checkErr(err)
	defer s.Conn.Close()
	err = s.Bind(s.BindDN, s.BindPass)
	s.checkErr(err)
	filter = strings.Replace(filter, "\r\n", "", -1)
	filter = strings.Replace(filter, "\n", "", -1)
	attributes := []string{"*"}
	sreq := ldap.NewSearchRequest(
		s.Base,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter, attributes, nil,
	)
	sres, err := s.Conn.Search(sreq)
	s.checkErr(err)
	var uu []User
	for _, v := range sres.Entries {
		if len(v.Attributes) > 0 {
			var m User
			for _, attr := range v.Attributes {
				switch attr.Name {
				case "c":
					m.C = attr.Values[0]
				case "cn":
					m.CN = attr.Values[0]
				case "co":
					m.CO = attr.Values[0]
				case "company":
					m.Company = attr.Values[0]
				case "department":
					m.Department = attr.Values[0]
				case "displayName":
					m.DisplayName = attr.Values[0]
				case "distinguishedName":
					m.DistinguishedName = attr.Values[0]
				case "facsimileTelephoneNumber":
					m.FacsimileTelephoneNumber = attr.Values[0]
				case "givenName":
					m.GivenName = attr.Values[0]
				case "l":
					m.L = attr.Values[0]
				case "mail":
					m.Mail = attr.Values[0]
				case "mailNickname":
					m.MailNickname = attr.Values[0]
				case "manager":
					m.Manager = attr.Values[0]
				case "middleName":
					m.MiddleName = attr.Values[0]
				case "mobile":
					m.Mobile = attr.Values[0]
				case "name":
					m.Name = attr.Values[0]
				case "pager":
					m.Pager = attr.Values[0]
				case "physicalDeliveryOfficeName":
					m.PhysicalDeliveryOfficeName = attr.Values[0]
				case "postalCode":
					m.PostalCode = attr.Values[0]
				case "sAMAccountName":
					m.SAMAccountName = attr.Values[0]
				case "sn":
					m.SN = attr.Values[0]
				case "st":
					m.ST = attr.Values[0]
				case "streetAddress":
					m.StreetAddress = attr.Values[0]
				case "telephoneNumber":
					m.TelephoneNumber = attr.Values[0]
				case "title":
					m.Title = attr.Values[0]
				case "userPrincipalName":
					m.UserPrincipalName = attr.Values[0]
				case "wWWHomePage":
					m.WWWHomePage = attr.Values[0]
				}
			}
			uu = append(uu, m)
		}
	}
	return uu
}

// Filter results based on parameter containing value to filter
func (s *LDAPServer) Filter(dd []map[string]interface{}, p string, v string) (uu []map[string]interface{}) {
	for _, u := range dd {
		filt := strings.Contains(strings.ToLower(u[p].(string)), v)
		if !filt {
			uu = append(uu, u)
		}
	}
	return uu
}

// FilterSAMAccountName results based on parameter containing value to filter
func (s *LDAPServer) FilterSAMAccountName(dd []User, v string) (uu []User) {
	for _, u := range dd {
		// fmt.Println("user", u)
		filt := strings.Contains(strings.ToLower(u.SAMAccountName), v)
		// fmt.Println(filt, u.Name)
		if !filt {
			uu = append(uu, u)
		}
		// filt := strings.Contains(strings.ToLower(u[p].(string)), v)
		// if !filt {
		// 	uu = append(uu, u)
		// }
	}
	return
}

// LikeSAMAccountName results based on parameter containing value to select
func (s *LDAPServer) LikeSAMAccountName(dd []User, v string) (uu []User) {
	for _, u := range dd {
		filt := strings.Contains(strings.ToLower(u.SAMAccountName), v)
		if filt {
			uu = append(uu, u)
		}
	}
	return
}

// checkErr function
func (s *LDAPServer) checkErr(err error) {
	if err != nil {
		s.Log.Errorw(err.Error())
	}
}
