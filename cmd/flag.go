package main

import (
	"flag"
	"fmt"
	"os"
)

// ViewInventory - to view list of inventory items
func ViewInventory(viewInventoryCmd *flag.FlagSet, viewAll *bool, viewById *int, viewByCategory *string, viewAbovePrice *float64, viewBelowPrice *float64) {
	err := viewInventoryCmd.Parse(os.Args[2:])
	if err != nil {
		ExitWithErrorMsg(fmt.Sprintf("%s", err))
	}

	if !*viewAll && *viewById == -1 && *viewByCategory == "" && *viewAbovePrice == -1 && *viewBelowPrice == -1 {
		viewInventoryCmd.PrintDefaults()
		ExitWithErrorMsg("Invalid input. Require above mentioned data to execute")
	}

	if *viewAll {
		DisplayInventory()
		return
	}

	if *viewById != -1 {
		DisplayInventoryWithfilter(*viewById, "id")
		return
	}

	if *viewByCategory != "" {
		DisplayInventoryWithfilter(*viewByCategory, "category")
		return
	}

	if *viewAbovePrice != -1 {
		DisplayInventoryWithfilter(*viewAbovePrice, "price-above")
		return
	}

	if *viewBelowPrice != -1 {
		DisplayInventoryWithfilter(*viewBelowPrice, "price-below")
		return
	}

}

// RefillInventory - to refill inventory items
func RefillInventory(refillInventoryCmd *flag.FlagSet) {
	PrintGreetings("Welcome to Inventory Update!")
	adminName := AuthenticationForAdmin()
	StocksUpdateFromAdmin(adminName)
}

// CreateBill - bill items from store
func CreateBill(billItemCmd *flag.FlagSet) {
	id, name, contact := GetCustomerInfo()
	if id == -1 {
		ExitWithErrorMsg("Please enter valid data. Try again.")
	}
	PrintGreetings(fmt.Sprintf("Welcome %s [ID : %d, Contact : %s]",name, id, contact))
	amount, items, itemCount := PurchaseItems()
	if itemCount > 0 {
		WelcomeMessage(fmt.Sprintf("Your total amount is %f for %d item(s)", amount, itemCount))
		res := ConfirmPurchase()

		if res {
			UpdateInventoryStock(items)
			PrintBill(items)
			DisplayMessage(fmt.Sprintf("Total Amount is %f for %d item(s)", amount, itemCount))
			PrintGreetings("Thanks for purchasing here. Adios..See you soon!!!")
		} else {
			PrintGreetings("Thanks for visiting... Adios!")
		}
	} else {
		PrintGreetings("Thanks for visiting... Adios!")
	}

}

// showCustomers - show customers list
func showCustomers(customersCmd *flag.FlagSet) {
	ViewCustomers()
}

// showAdmins - show admins list
func showAdmins(adminCmd *flag.FlagSet) {
	ViewAdminDetails()
}