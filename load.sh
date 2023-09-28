#!/bin/bash
OHOST=localhost
ODB=modtest

# ,mrp_workcenter
case ${1} in
  "init")
    time go run . -j 8 -host $OHOST -d $ODB -schema http -port 8069 -f system_parameters,currency,company,partner_tags,approval_category,account_setup,account_payment_term,bank_table,uoms
    ;;
  "div")
    time go run . -j 8 -host $OHOST -d $ODB -schema http -port 8069 -f division,salesregion
    ;;
  "warehouse")
    time go run . -j 8 -host $OHOST -d $ODB -schema http -port 8069 -f plant,warehouse_init,warehouse_sequence_fix,warehouse,warehouse_location,branches
    ;;
  "users")
    time go run . -j 8 -host $OHOST -d $ODB -schema http -port 8069 -f users
    ;;
  "products")
    time go run . -j 8 -host $OHOST -d $ODB -schema http -port 8069 -f product_category,product_template_1000,pricelist_base,pricelist_customer,customer,vendor,mrp_workcenter_tag
    ;;
  *)
    shift
    time go run . -j 8 -host $OHOST -d $ODB -schema http -port 8069 -f $@
    ;;
esac