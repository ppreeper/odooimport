package odooconn

import (
	"fmt"
	"strings"
)

// ResConfigSettings function
func (o *OdooConn) ResConfigSettings(setting string, nval interface{}) {
	mdl := "res_config_settings"
	umdl := strings.Replace(mdl, "_", ".", -1)

	// id
	// ,create_uid,create_date,write_uid,write_date
	// ,company_id
	// ,user_default_rights,external_email_server_default
	// ,module_base_import
	// ,module_google_calendar,module_microsoft_calendar
	// ,module_google_drive,module_google_spreadsheet
	// ,module_auth_oauth,module_auth_ldap,module_base_gengo
	// ,module_account_inter_company_rules,module_pad,module_voip
	// ,module_web_unsplash,module_partner_autocomplete
	// ,module_base_geolocalize,module_google_recaptcha
	// ,group_multi_currency,show_effect,fail_counter
	// ,alias_domain,map_box_token,unsplash_access_key
	// ,auth_signup_reset_password,auth_signup_uninvited
	// ,auth_signup_template_user_id,group_discount_per_so_line
	// ,group_uom,group_product_variant,module_sale_product_configurator
	// ,module_sale_product_matrix,group_stock_packaging
	// ,group_product_pricelist,group_sale_pricelist
	// ,product_pricelist_setting,product_weight_in_lbs
	// ,product_volume_volume_in_cubic_feet
	// ,disable_redirect_firebase_dynamic_link,enable_ocn
	// ,digest_emails,digest_id,chart_template_id
	// ,module_account_accountant,group_analytic_accounting
	// ,group_analytic_tags,group_warning_account,group_cash_rounding
	// ,group_show_line_subtotals_tax_excluded
	// ,group_show_line_subtotals_tax_included,group_show_sale_receipts
	// ,group_show_purchase_receipts,show_line_subtotals_tax_selection
	// ,module_account_budget,module_account_payment
	// ,module_account_reports,module_account_check_printing
	// ,module_account_batch_payment,module_account_sepa
	// ,module_account_sepa_direct_debit,module_account_plaid
	// ,module_account_yodlee,module_account_bank_statement_import_qif
	// ,module_account_bank_statement_import_ofx
	// ,module_account_bank_statement_import_csv
	// ,module_account_bank_statement_import_camt,module_currency_rate_live
	// ,module_account_intrastat,module_product_margin,module_l10n_eu_service
	// ,module_account_taxcloud,module_account_invoice_extract
	// ,module_snailmail_account,use_invoice_terms,group_auto_done_setting
	// ,module_sale_margin,use_quotation_validity_days
	// ,group_warning_sale,group_sale_delivery_address,group_proforma_sales
	// ,default_invoice_policy,deposit_default_product_id
	// ,module_delivery,module_delivery_dhl,module_delivery_fedex
	// ,module_delivery_ups,module_delivery_usps,module_delivery_bpost
	// ,module_delivery_easypost,module_product_email_template
	// ,module_sale_coupon,module_sale_amazon,automatic_invoice
	// ,template_id,confirmation_template_id,group_sale_order_template
	// ,module_sale_quotation_builder,crm_alias_prefix
	// ,generate_lead_from_alias,group_use_lead,group_use_recurring_revenues
	// ,module_crm_iap_lead,module_crm_iap_lead_website
	// ,module_crm_iap_lead_enrich,module_mail_client_extension
	// ,lead_enrich_auto,lead_mining_in_pipeline
	// predictive_lead_scoring_start_date_str
	// ,predictive_lead_scoring_fields_str,module_project_forecast
	// ,module_hr_timesheet,group_subtask_project,group_project_rating
	// ,group_project_recurring_tasks,geoloc_provider_id
	// ,geoloc_provider_googlemap_key,module_hr_presence,module_hr_skills
	// ,hr_presence_control_o.Log.n,hr_presence_control_email
	// ,hr_presence_control_ip,module_hr_attendance,hr_employee_self_edit
	// ,module_project_timesheet_synchro,module_project_timesheet_holidays
	// ,timesheet_min_duration,timesheet_rounding,module_industry_fsm_report
	// ,module_industry_fsm_sale,invoiced_timesheet,group_industry_fsm_quotations
	// ,module_account_predictive_bills,group_fiscal_year
	// ,module_procurement_jit,module_product_expiry,group_stock_production_lot
	// ,group_lot_on_delivery_slip,group_stock_tracking_lot
	// ,group_stock_tracking_owner,group_stock_adv_location
	// ,group_warning_stock,group_stock_sign_delivery,module_stock_picking_batch
	// ,module_stock_barcode,module_stock_sms,group_stock_multi_locations
	// ,module_stock_landed_costs,group_display_incoterm,group_lot_on_invoice
	// ,use_security_lead,default_picking_policy,lock_confirmed_po
	// ,po_order_approval,default_purchase_method
	// ,group_warning_purchase,module_account_3way_match
	// ,module_purchase_requisition,module_purchase_product_matrix
	// ,use_po_lead,group_send_reminder,module_stock_dropshipping
	// ,is_installed_sale,use_manufacturing_lead,group_mrp_byproducts
	// ,module_mrp_mps,module_mrp_plm,module_mrp_workorder
	// ,module_quality_control,module_mrp_subcontracting
	// ,group_mrp_routings,group_locked_by_default,website_id
	// ,group_multi_website,recaptcha_public_key,recaptcha_private_key
	// ,recaptcha_min_score,module_website_sale_delivery
	// ,sale_delivery_settings,group_delivery_invoice_address
	// ,module_website_sale_digital,module_website_sale_wishlist
	// ,module_website_sale_comparison,module_website_sale_stock
	// ,module_account,inventory_availability,available_threshold
	// ,module_website_sale_slides,module_website_slides_forum
	// ,module_website_slides_survey,module_mass_mailing_slides
	// ,group_attendance_use_pin
	// rr := o.SearchRead(umdl, oarg{{"company_id", "=", 1}}, 0, 0, []string{"company_id", f})
	// rr := o.SearchRead(umdl, oarg{{"company_id", "=", 1}}, 0, 0, []string{})
	// rr := o.Search(umdl, oarg{{"company_id", "=", 1}})
	rr, err := o.SearchRead("ir.config_parameter", oarg{oarg{"company_id", "=", 1}}, 0, 0, []string{})
	o.checkErr(err)
	fmt.Println(rr)
	// rid := rr[len(rr)-1]
	// fmt.Println(rr, rid)
	// r := o.Read(umdl, []int{rid}, []string{})
	// ur := r[0]
	// fmt.Println(ur[f])
	ur := map[string]interface{}{}
	// for k, v := range ur {
	// 	if k != "id" {
	// 		nr[k] = v
	// 		// fmt.Println(k)
	// 	}
	// }
	// r := rr[len(rr)-1][f]
	// fmt.Println(r)
	// var ur = map[string]interface{}{}
	// switch nvalType := nval.(type) {
	// case int:
	// 	o.Log.Info("int", "value", nvalType)
	// 	ur[setting] = nvalType
	// case float64:
	// 	o.Log.Info("float64", "value", nvalType)
	// 	ur[setting] = nvalType
	// case string:
	// 	o.Log.Info("string", "value", nvalType)
	// 	ur[setting] = nvalType
	// case bool:
	// 	o.Log.Info("bool", "value", nvalType)
	// 	ur[setting] = nvalType
	// default:
	// 	o.Log.Info("no valid type")
	// }
	o.Log.Info(mdl, "model", umdl, "record", ur, "setting", setting)

	// o.WriteRecord(umdl, -1, INSERT, ur)
}
