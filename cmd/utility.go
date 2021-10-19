package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Inventory struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type CustomerInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type Bills struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Quantity int  `json:"quantity"`
	Amount float64 `json:"amount"`
}

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


func DisplayInventory() {
	items := GetInventoryItems()
	tbl := table.New("ID", "TITLE", "CATEGORY", "STOCK", "PRICE").WithPadding(3)

	DisplayMessage(fmt.Sprintf("%d items are available in our inventory", len(items)))
	for _, item := range items {
		tbl.AddRow(item.Id, item.Title, item.Category, item.Stock, item.Price)
	}

	tbl.Print()
}

func PrintBill(items []Bills) {
	tbl := table.New("S.NO.", "ID", "TITLE", "CATEGORY", "PRICE", "QUANTITY", "AMOUNT").WithPadding(3)

	DisplayMessage(fmt.Sprintf("You have purchased %d item(s)", len(items)))
	for idx, item := range items {
		tbl.AddRow(idx+1, item.Id, item.Title, item.Category, item.Price, item.Quantity, item.Amount)
	}

	tbl.Print()
}

//func ConvertStructToMap(item Inventory) map[string]interface{} {
//	m, _ := json.Marshal(item)
//	var x map[string]interface{}
//	_ = json.Unmarshal(m, &x)
//
//	return x
//}

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

func AddNewUser(name, contact string) int{
	customers := GetCustomerDetailsFromServer()
	newCustomer := CustomerInfo {
		Id: len(customers)+1,
		Name : name,
		Contact: contact,
	}

	customers = append(customers, newCustomer)
	WriteUserDetails(customers)
	DisplayMessage(fmt.Sprintf("Successfully added customer %s [ID : %d] ", newCustomer.Name, newCustomer.Id))
	return newCustomer.Id
}

func GetAdminAccess() {
//
}

func GetCustomerInfo() (int, string, string) {
	customerType := ""
	prompt := &survey.Select{
		Message: "Choose customer relation :",
		Options: []string{"New Customer", "Existing Customer"},
	}
	survey.AskOne(prompt, &customerType)

	if(customerType == "New Customer") {
		var contact string
		name := StringPrompt("Enter Customer Name : ")
		fmt.Print("Enter Contact Number : ")
		_, err := fmt.Scanln(&contact)
		if err != nil {
			ExitWithErrorMsg(fmt.Sprintf("%v", err))
		}
		id := AddNewUser(name, contact)
		return id,  name, contact
	} else {
		customers := GetCustomerDetailsFromServer()
		id := StringPrompt("Enter Customer Id : ")
		userId, err := strconv.Atoi(id)
		if err != nil {
			ExitWithErrorMsg(fmt.Sprintf("%d is not a valid id. Please try again", userId))
		}
		for _, customer := range customers {
			if customer.Id ==  userId{
				return customer.Id, customer.Name, customer.Contact
			}
		}
		return -1,  "", ""
	}
}

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
			msg = fmt.Sprintf("%d",value.(int))
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
			msg = fmt.Sprintf("%f",value.(float64))
		}
	case "price-below":
		{
			for _, item := range items {
				if value.(float64) >= item.Price {
					count++
					tbl.AddRow(item.Id, item.Title, item.Category, item.Stock, item.Category, item.Price)
				}
			}
			msg = fmt.Sprintf("%f",value.(float64))
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

func PurchaseItems() (float64, []Bills, int) {
	amount := 0.0
	purchaseList := []Bills{}
	inventoryItems := GetInventoryItems()
	fmt.Println()
	DisplayMessage("Enter -1 as item id to finish billing")
	fmt.Println()

	for {
		item := StringPrompt("Enter inventory item id [1-48] : ")
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
				Id: itemId,
				Title: purchaseItem.Title,
				Category: purchaseItem.Category,
				Price: purchaseItem.Price,
				Quantity: itemQuantity,
				Amount: purchaseItem.Price * float64(itemQuantity),
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

func ConfirmPurchase() bool{
	confirmation := false
	prompt := &survey.Confirm{
		Message: "Do you like to proceed with payment? ",
	}
	survey.AskOne(prompt, &confirmation)

	return confirmation
}

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

func WelcomeMessage(msg string) {
	magenta := color.New(color.FgCyan).Add(color.Bold)
	magenta.Println(msg)
}


func PrintGreetings(msg string) {
	green := color.New(color.FgGreen).Add(color.Bold)
	green.Println(msg)
}