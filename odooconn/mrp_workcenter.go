package odooconn

import (
	"fmt"
	"strings"
)

// MRPWorkcenter function
func (o *OdooConn) MRPWorkcenter() {
	mdl := "mrp_workcenter"
	umdl := strings.Replace(mdl, "_", ".", -1)
	fmt.Printf("\n%v MRPWorkcenter\n", umdl)
}
