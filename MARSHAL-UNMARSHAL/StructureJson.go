package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type User struct {
	Userid  int     `json:"id"`
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

type Address struct {
	Area    string `json:"area"`
	Country string `json:"country"`
}

type TechDets struct {
	Technolgy  string  `json:"tech"`
	Experience float64 `json:"exp"`
}

type Tech struct {
	Userid   int        `json:"id"`
	TechDets []TechDets `json:"techDets"`
}

type ContactDets struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Contact struct {
	Userid      int         `json:"id"`
	ContactDets ContactDets `json:contactDets`
}

type MergeUser struct {
	Userid      int
	Name        string
	Address     Address
	TechDetails []TechDets
	Email       string
	Phone       string
}

var countrymap = map[string]string{
	"IND": "+91",
	"UK":  "+41",
}

func addContactPrefix(users []User, techDetails []Tech, contactDetails []Contact) map[int]MergeUser {

	merged := make([]MergeUser, len(users))
	for i, user := range users {

		mergeuser := MergeUser{}
		mergeuser.Userid = user.Userid
		mergeuser.Name = user.Name
		mergeuser.Address.Area = user.Address.Area
		mergeuser.Address.Country = user.Address.Country

		for _, tech := range techDetails {

			if user.Userid == tech.Userid {
				mergeuser.TechDetails = tech.TechDets
			}

		}

		for _, contact := range contactDetails {

			if user.Userid == contact.Userid {

				mergeuser.Email = contact.ContactDets.Email
				if valprefix, ok := countrymap[mergeuser.Address.Country]; ok {
					mergeuser.Phone = valprefix + " " + contact.ContactDets.Phone
				}

			}

		}

		merged[i] = mergeuser
	}
	file, _ := json.MarshalIndent(merged, "", " ")
	_ = ioutil.WriteFile("mergeuser.json", file, 0644)

	mapuser := make(map[int]MergeUser)
	for _, val := range merged {
		mapuser[val.Userid] = val
	}
	return mapuser
}

func main() {

	usercontent, err := ioutil.ReadFile("user.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	var users []User
	err2 := json.Unmarshal([]byte(usercontent), &users)
	if err2 != nil {
		fmt.Println("Error JSON Unmarshling for user file")
		fmt.Println(err2.Error())

	}

	techcontent, err3 := ioutil.ReadFile("tech.json")
	if err3 != nil {
		fmt.Println(err3.Error())
	}

	var techDetails []Tech
	err4 := json.Unmarshal([]byte(techcontent), &techDetails)
	if err4 != nil {
		fmt.Println("Error JSON Unmarshling for tech file")
		fmt.Println(err4.Error())

	}

	contactcontent, err5 := ioutil.ReadFile("contact.json")
	if err5 != nil {
		fmt.Println(err5.Error())
	}

	var contactDetails []Contact
	err6 := json.Unmarshal([]byte(contactcontent), &contactDetails)
	if err6 != nil {
		fmt.Println("Error JSON Unmarshling for contact file")
		fmt.Println(err6.Error())

	}

	mergeuser := addContactPrefix(users, techDetails, contactDetails)
	fmt.Println(mergeuser)
}
