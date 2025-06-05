package main

type Account struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Orders []*Order `json:"orders"`
}
