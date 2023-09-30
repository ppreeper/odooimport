package main

import (
	"fmt"

	"github.com/ppreeper/odooimport/odooconn"
)

func execFlags(f string, o *odooconn.OdooConn, c Conf) {
	switch f {
	case "extemail":
		o.ResConfigSettings("external_email_server_default", true)
	case "decimal_precision":
		o.DecimalPrecision("Product Price", 6)
		o.DecimalPrecision("Discount", 6)
		o.DecimalPrecision("Stock Weight", 6)
		o.DecimalPrecision("Volume", 6)
		o.DecimalPrecision("Product Unit of Measure", 6)
		o.DecimalPrecision("Payment Terms", 6)
		o.DecimalPrecision("Quality Tests", 6)
	case "menu_sort":
		o.IrUiMenuSort()
	case "account_payment_term":
		o.AccountPaymentTerm()
		o.AccountPaymentTermLine()
	case "bank_table":
		o.ResPartnerBank()
	case "uom":
		o.UomCategory()
		o.UomUom()
	// partner categories
	case "partner_category":
		o.ResPartnerCategory()
	// add companies
	case "company":
		o.ResCompany()
	// partners
	// res_partner_companies
	// res_partner_cstype
	// product_template_companies
	// uom_output_text
	// plant
	case "plant":
		o.ResPartnerPlant()
	case "warehouse":
		o.StockWarehouse()
	case "warehouse_location":
		o.StockLocation()
	// users
	case "user":
		o.ResUsers()
		o.ResPartnerUsers()
		o.HRDepartment()
		o.HREmployee()
		o.HREmployeeManager()
	case "res_users":
		o.ResUsers()
	case "res_partner_users":
		o.ResPartnerUsers()
	case "hr_department":
		o.HRDepartment()
	case "hr_employee":
		o.HREmployee()
	case "hr_manager":
		o.HREmployeeManager()
	case "crm_team":
		o.CrmTeam()
	case "crm_team_members":
		o.CrmTeamMembers()
	// product categories
	// product_category_buyer_schedule
	// product_category_material
	case "pcc1":
		o.ProductCategoryConsumable1()
	case "pcc2":
		o.ProductCategoryConsumable2()
	case "pcc3":
		o.ProductCategoryConsumable3()
	case "pcp1":
		o.ProductCategoryProduct1()
	case "pcp2":
		o.ProductCategoryProduct2()
	case "pcp3":
		o.ProductCategoryProduct3()
	case "pcp4":
		o.ProductCategoryProduct4()
	case "pcp5":
		o.ProductCategoryProduct5()
	case "pcdelpro":
		o.ProductCategoryDelpro()
	// Pricegroup Lists
	case "pl_pg":
		o.PricelistPricegroup()
	case "pl_pg_mg":
		o.PricelistPricegroupMatGrpDiscounts()
	// Customer Specific Lists
	case "pl_cust":
		o.PricelistCustomer()
	case "pl_cust_def":
		o.PricelistCustomerDefault()
	case "pl_cust_mg":
		o.PricelistCustomerMatGrpDiscounts()
	case "pl_cust_no":
		o.PricelistCustomerNetoutItems()
	// Partners Customer
	// needs partner_long_notes
	// other partner_long_notes_report
	case "customer":
		o.ResPartnerCustomer()
	case "customer_link":
		o.ResPartnerCustomerLink()
	// Partners Vendor
	// partner_avl
	case "vendor":
		o.ResPartnerVendors()
	case "vendor_link":
		o.ResPartnerVendorsLink()
	case "vendor_bank":
		o.ResPartnerVendorsBank()
	case "vendor_delpro":
		o.ResPartnerVendorsDelpro()
	// products
	// product_template_inventory_type
	// product_template_mrp_type
	// product_template_sequence_number
	case "pnon":
		// nonvaluated materials
		o.ProductTemplate("UNBW")
	case "pservice":
		// artg service
		o.ProductTemplate("DIEN")
	case "pconsumable":
		// artg consumable products
		o.ProductTemplate("ZCON")
	case "praw":
		// artg raw materials
		o.ProductTemplate("ROH")
	case "psemi":
		// artg semifinished products
		o.ProductTemplate("HALB")
	case "pfini":
		// artg finished products
		o.ProductTemplate("FERT")
	case "product_customer":
		o.ProductCustomerPart()
	case "product_delpro":
		o.ProductTemplateDelpro() // delpro materials
	case "minmax":
		o.StockReorderPoint()
	case "putaway":
		o.StockPutawayRule()
	// Inventory
	case "inventory":
		o.StockInventoryInitial("Groupe A.R. Thomson")
	// Vendor Pricelist
	case "vend_price":
		o.ProductSupplierinfo("ROH")
		o.ProductSupplierinfo("HALB")
		o.ProductSupplierinfo("FERT")
	case "mrp_bom":
		o.MRPBom()
	case "mrp_bom_item":
		o.MRPBomLine("A.R. Thomson Group Manufacturing")
	case "mrp_bom_op":
		o.MRPBomOP("A.R. Thomson Group Manufacturing")
	case "mrp_work_center":
		o.MRPWorkcenter()
	default:
		fmt.Println("selected: ", f)
	}

	ilog.Info("odooimport end")
}
