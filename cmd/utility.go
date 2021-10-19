package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/rodaine/table"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)


// GetInventoryItems - Returns list of inventory items
func GetInventoryItems() (inventory []Inventory) {
	fileBytes, err := ioutil.ReadFile("../api/store-inventory.json")

	if err != nil {
		fmt.Println(err, fileBytes)
		panic(err)
	}

	err = json.Unmarshal(fileBytes, &inventory)

	if err != nil {
		ExitWithErrorMsg("Inventory Data is corrupted!! Please try again..")
		panic(err)
	}
	return inventory
}

// WriteInventoryDetails - Update the file with inventory data
func WriteInventoryDetails(inventory []Inventory) {
	inventoryBytes, err := json.Marshal(inventory)

	if err != nil {
		ExitWithErrorMsg(fmt.Sprintf("%v", err))
	}

	err = ioutil.WriteFile("../api/store-inventory.json", inventoryBytes, 0644)
	if err != nil {
		ExitWithErrorMsg(fmt.Sprintf("%v", err))
	}
}

// GetCustomerDetailsFromServer - Returns list of customers
func GetCustomerDetailsFromServer() (customers []CustomerInfo) {
	fileBytes, err := ioutil.ReadFile("../api/customer-records.json")

	if err != nil {
		fmt.Println(err, fileBytes)
		panic(err)
	}

	err = json.Unmarshal(fileBytes, &customers)

	if err != nil {
		ExitWithErrorMsg("Customer Data is corrupted!! Please try again..")
		panic(err)
	}
	return customers
}

// DisplayInventory - Show inventory list
func DisplayInventory() {
	items := GetInventoryItems()
	tbl := table.New("ID", "TITLE", "CATEGORY", "STOCK", "PRICE").WithPadding(3)

	DisplayMessage(fmt.Sprintf("%d items are available in our inventory", len(items)))
	for _, item := range items {
		tbl.AddRow(item.Id, item.Title, item.Category, item.Stock, item.Price)
	}

	tbl.Print()
}

// PrintBill - Show purchase bill
func PrintBill(items []Bills) {
	tbl := table.New("S.NO.", "ID", "TITLE", "CATEGORY", "PRICE", "QUANTITY", "AMOUNT").WithPadding(3)

	DisplayMessage(fmt.Sprintf("You have purchased %d item(s)", len(items)))
	for idx, item := range items {
		tbl.AddRow(idx+1, item.Id, item.Title, item.Category, item.Price, item.Quantity, item.Amount)
	}

	tbl.Print()
}

// WriteUserDetails - To Write customer details to json file
func WriteUserDetails(customers []CustomerInfo) {
	customersBytes, err := json.Marshal(customers)

	if err != nil {
		ExitWithErrorMsg(fmt.Sprintf("%v", err))
	}

	err = ioutil.WriteFile("../api/customer-records.json", customersBytes, 0644)
	if err != nil {
		ExitWithErrorMsg(fmt.Sprintf("%v", err))
	}
}

// AddNewUser - To add new customer to customers list
func AddNewUser(name, contact string) int {
	customers := GetCustomerDetailsFromServer()
	newCustomer := CustomerInfo{
		Id:      len(customers) + 1,
		Name:    name,
		Contact: contact,
	}

	customers = append(customers, newCustomer)
	WriteUserDetails(customers)
	DisplayMessage(fmt.Sprintf("Successfully added customer %s [ID : %d] ", newCustomer.Name, newCustomer.Id))
	return newCustomer.Id
}

// GetAdminDetails - Returns list of admins
func GetAdminDetails() (admin []Admin) {
	fileBytes, err := ioutil.ReadFile("../api/admin.json")

	if err != nil {
		fmt.Println(err, fileBytes)
		panic(err)
	}

	err = json.Unmarshal(fileBytes, &admin)

	if err != nil {
		ExitWithErrorMsg("Admin Data is corrupted!! Please try again..")
		panic(err)
	}
	return admin
}

// StocksUpdateFromAdmin - To update stocks by admin
func StocksUpdateFromAdmin(admin string) {
	inventoryItems := GetInventoryItems()
	updatedStocks := []int{}
	fmt.Println()
	DisplayMessage("Enter -1 as item id to finish stock updates")
	fmt.Println()

	for {
		item := StringPrompt(fmt.Sprintf("Enter inventory item id [1-%d]", len(inventoryItems)))
		itemId, err := strconv.Atoi(item)
		if itemId == -1 { // terminate the list
			break
		}
		quantity := StringPrompt("Enter stock quantity : ")
		itemQuantity, err := strconv.Atoi(quantity)
		if err != nil {
			ExitWithErrorMsg(fmt.Sprintf("%v", err))
		}
		switch {
		case itemId >= 1 && itemId <= len(inventoryItems):
			stockItem := inventoryItems[itemId-1]
			stockItem.Stock = itemQuantity
			updatedStocks = append(updatedStocks, stockItem.Id)

			fmt.Printf("ID : %d, Item : %s, Category : %s, Price : %f, Stock : %d\n",
				stockItem.Id, stockItem.Title, stockItem.Category, stockItem.Price, stockItem.Stock)
		default:
			DisplayMessage(fmt.Sprintf("ID %d is not available. Choose between 1 & %d", itemId, len(inventoryItems)))
		}
	}
	if len(updatedStocks) > 0 {
		WriteInventoryDetails(inventoryItems)
		WriteInventoryStockLog(admin, updatedStocks)
	}
}

// WriteInventoryStockLog - Write a log file with admin action
func WriteInventoryStockLog(admin string, stocks []int) {
	// convert int slice to a comma separated string
	dataString := strings.Trim(strings.Replace(fmt.Sprint(stocks), " ", ",", -1), "[]")
	//to insert new line while writing to file
	dataString += "\n"
	data := []byte(fmt.Sprintf("%s updated stock following item(s) :  %s", admin, dataString))

	f, err := os.OpenFile("../logs/admin-stock-updates.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		ExitWithErrorMsg(fmt.Sprintf("%v", err))
	}

	PrintGreetings("Successfully updated the stocks :) ")

}

// AuthenticationForAdmin -Authenticate the admin login process
func AuthenticationForAdmin() string {
	name := StringPrompt("Enter your name : ")

	password := ""
	prompt2 := &survey.Password{
		Message: "Enter your password : ",
	}
	survey.AskOne(prompt2, &password)

	DisplayMessage("Authentication is in progress.... ")
	fmt.Println(".................................")
	adminDetails := GetAdminDetails()
	for _, admin := range adminDetails {
		if admin.Name == name && admin.Password == password {
			PrintGreetings(fmt.Sprintf("Log in successful for %s", name))
			return name
		}
	}
	ExitWithErrorMsg("Authentication Failed. Try again.")
	return ""
}

// GetCustomerInfo - To receive customer info from userinput
func GetCustomerInfo() (int, string, string) {
	customerType := ""
	prompt := &survey.Select{
		Message: "Choose customer relation :",
		Options: []string{"New Customer", "Existing Customer"},
	}
	survey.AskOne(prompt, &customerType)

	if customerType == "New Customer" {
		var contact string
		name := StringPrompt("Enter Customer Name : ")
		fmt.Print("Enter Contact Number : ")
		_, err := fmt.Scanln(&contact)
		if err != nil {
			ExitWithErrorMsg(fmt.Sprintf("%v", err))
		}
		id := AddNewUser(name, contact)
		return id, name, contact
	} else {
		customers := GetCustomerDetailsFromServer()
		id := StringPrompt("Enter Customer Id : ")
		userId, err := strconv.Atoi(id)
		if err != nil {
			ExitWithErrorMsg(fmt.Sprintf("%d is not a valid id. Please try again", userId))
		}
		for _, customer := range customers {
			if customer.Id == userId {
				return customer.Id, customer.Name, customer.Contact
			}
		}
		return -1, "", ""
	}
}

// DisplayInventoryWithfilter - To display inventory list with filter from input
func DisplayInventoryWithfilter(value interface{}, filter string) {
	items := GetInventoryItems()

	tbl := table.New("ID", "TITLE", "CATEGORY", "STOCK", "PRICE").WithPadding(3)

	count := 0
	var msg string

	switch filter {
	case "id":
		{
			for _, item := range items {
				if item.Id == value.(int) {
					count++
					tbl.AddRow(item.Id, item.Title, item.Category, item.Stock, item.Price)
				}
			}
			msg = fmt.Sprintf("%d", value.(int))
		}
	case "category":
		{
			for _, item := range items {
				category := strings.ToLower(item.Category)
				if category == value.(string) {
					count++
					tbl.AddRow(item.Id, item.Title, item.Category, item.Stock, item.Category, item.Price)
				}
			}
			msg = value.(string)
		}
	case "price-above":
		{
			for _, item := range items {
				if value.(float64) <= item.Price {
					count++
					tbl.AddRow(item.Id, item.Title, item.Category, item.Stock, item.Category, item.Price)
				}
			}
			msg = fmt.Sprintf("%f", value.(float64))
		}
	case "price-below":
		{
			for _, item := range items {
				if value.(float64) >= item.Price {
					count++
					tbl.AddRow(item.Id, item.Title, item.Category, item.Stock, item.Category, item.Price)
				}
			}
			msg = fmt.Sprintf("%f", value.(float64))
		}
	default:
		DisplayMessage(fmt.Sprintf("Invalid entry  %s", filter))
	}

	if count > 0 {
		DisplayMessage(fmt.Sprintf("%d Result(s) for %s : %s", count, filter, msg))
		tbl.Print()
	} else {
		DisplayMessage(fmt.Sprintf("No item found with %s", filter))
	}
}

// UpdateInventoryStock - Update the stock value based on user action
func UpdateInventoryStock(bills []Bills) {
	inventoryItems := GetInventoryItems()

	for _, bill := range bills {
		inventoryItems[bill.Id-1].Stock -= bill.Quantity
		if inventoryItems[bill.Id-1].Stock <= 0 {
			inventoryItems[bill.Id-1].Stock = 100
		}
	}

	WriteInventoryDetails(inventoryItems)
}

// PurchaseItems - Process of purchasing items from inventory
func PurchaseItems() (float64, []Bills, int) {
	amount := 0.0
	purchaseList := []Bills{}
	inventoryItems := GetInventoryItems()
	fmt.Println()
	DisplayMessage("Enter -1 as item id to finish billing")
	fmt.Println()

	for {
		item := StringPrompt(fmt.Sprintf("Enter inventory item id [1-%d]", len(inventoryItems)))
		itemId, err := strconv.Atoi(item)
		if itemId == -1 { // terminate the list
			break
		}
		quantity := StringPrompt("Enter Quantity : ")
		itemQuantity, err := strconv.Atoi(quantity)
		if err != nil {
			ExitWithErrorMsg(fmt.Sprintf("%v", err))
		}
		switch {
		case itemId >= 1 && itemId <= len(inventoryItems):
			purchaseItem := inventoryItems[itemId-1]
			newPurchase := Bills{
				Id:       itemId,
				Title:    purchaseItem.Title,
				Category: purchaseItem.Category,
				Price:    purchaseItem.Price,
				Quantity: itemQuantity,
				Amount:   purchaseItem.Price * float64(itemQuantity),
			}

			fmt.Printf("ID : %d, Item : %s, Category : %s, Price : %f, Quantity : %d | Amount = %f\n",
				newPurchase.Id, newPurchase.Title, newPurchase.Category, newPurchase.Price, newPurchase.Quantity, newPurchase.Amount)

			purchaseList = append(purchaseList, newPurchase)
			amount += newPurchase.Amount
		default:
			DisplayMessage(fmt.Sprintf("ID %d is not available. Choose between 1 & %d", itemId, len(inventoryItems)))
		}
	}

	return amount, purchaseList, len(purchaseList)
}

// ConfirmPurchase -To confirm the purchase
func ConfirmPurchase() bool {
	confirmation := false
	prompt := &survey.Confirm{
		Message: "Do you like to proceed with payment? ",
	}
	survey.AskOne(prompt, &confirmation)

	return confirmation
}

// ExitWithErrorMsg -  Display error with formatting
func ExitWithErrorMsg(msg string) {
	color.Set(color.FgHiRed)
	fmt.Println(msg)
	color.Unset()
	os.Exit(1)
}

// DisplayMessage - Display a message
func DisplayMessage(msg string) {
	magenta := color.New(color.FgMagenta).Add(color.Bold)
	magenta.Println(msg)
}

// StringPrompt - To prompt user with a question & return response
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

// WelcomeMessage - Display message with formatting
func WelcomeMessage(msg string) {
	magenta := color.New(color.FgCyan).Add(color.Bold)
	magenta.Println(msg)
}

// PrintGreetings - Display message with formatting
func PrintGreetings(msg string) {
	green := color.New(color.FgGreen).Add(color.Bold)
	green.Println(msg)
}

// ViewCustomers - Table of all customers
func ViewCustomers() {
	customers:= GetCustomerDetailsFromServer()

	if len(customers) > 0 {
		WelcomeMessage("Customer Details")
		tbl := table.New("ID", "CUSTOMER NAME", "CONTACT NO")

		for _, customer := range customers {
			tbl.AddRow(customer.Id, customer.Name, customer.Contact)
		}

		tbl.Print()
	} else {
		DisplayMessage("No customer records found")
	}
}

// ViewAdminDetails - Table of all admins
func ViewAdminDetails() {
	admins := GetAdminDetails()
	if len(admins) > 0 {
		WelcomeMessage("Admin Details")

		tbl := table.New("NAME", "PASSWORD")

		for _, admin := range admins {
			tbl.AddRow(admin.Name, admin.Password)
		}

		tbl.Print()
	} else {
		DisplayMessage("No admin records found")
	}
}