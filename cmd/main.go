package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Ramya")

	// view inventory
	viewInventoryCmd := flag.NewFlagSet("view", flag.ExitOnError)
	viewAll := viewInventoryCmd.Bool("all", false, "View entire inventory list")
	viewById := viewInventoryCmd.Int("id", -1, "ID in inventory list")
	viewByCategory := viewInventoryCmd.String("category", "", "Category of item")
	viewAbovePrice := viewInventoryCmd.Float64("price-above", -1, "Inventory Items above price")
	viewBelowPrice := viewInventoryCmd.Float64("price-below", -1, "Inventory Items below price")

	//	admin - refill
	refillInventoryCmd := flag.NewFlagSet("refill", flag.ExitOnError)
	//refillAdminAccess := refillInventoryCmd.String("admin", "", "Enter admin name")

	// billing
	billItemCmd := flag.NewFlagSet("bill", flag.ExitOnError)

	// display customer details
	customersCmd := flag.NewFlagSet("customer", flag.ExitOnError)

	// display admin details
	adminsCmd := flag.NewFlagSet("admin", flag.ExitOnError)

	if len(os.Args) < 2 {
		ExitWithErrorMsg("Insufficient arguments. Expected - view (or) refill (or) bill (or) customer (or) admin command for execution")
	}

	switch os.Args[1] {
	case "view":
		// to view list of inventory items
		ViewInventory(viewInventoryCmd, viewAll, viewById, viewByCategory, viewAbovePrice, viewBelowPrice)
	case "refill":
		// to refill inventory items
		RefillInventory(refillInventoryCmd)
	case "bill":
		// bill items from store
		CreateBill(billItemCmd)
	case "customer":
		// display customers list
		showCustomers(customersCmd)
	case "admin":
		// display admins list
		showAdmins(adminsCmd)

	default:
		ExitWithErrorMsg(fmt.Sprintf("Unknown arguments %s. Expected - view (or) refill (or) buy (or) customer (or) admin command for execution", os.Args[1]))
	}
}
