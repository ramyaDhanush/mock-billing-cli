package main

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
	Id       int     `json:"id"`
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Amount   float64 `json:"amount"`
}

type Admin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
