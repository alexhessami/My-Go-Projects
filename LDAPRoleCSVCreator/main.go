package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	//var password string
	fmt.Print("Please enter LDAP server password (Look up in 1Pass: 'LDAP - 198.54.96.66'): ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	passwordString := string(password)

	if err != nil {
		fmt.Println("\nError reading password:", err)
		return
	}

	//LDAP creds
	ldapURL := "198.54.96.66:389"
	bindDN := "cn=admin,dc=supplyframe,dc=com" // Your bind DN
	bindPassword := passwordString             // Your bind password

	//Connecting to server
	l, err := ldap.Dial("tcp", ldapURL)
	if err != nil {
		fmt.Println("LDAP connection error:", err)
		return
	}
	defer l.Close()

	err = l.Bind(bindDN, bindPassword)
	if err != nil {
		fmt.Println("LDAP bind error:", err)
		return
	}

	var answer string
	fmt.Println("\nWould you like to see a list of the roles (y or n): ")
	fmt.Scanln(&answer)

	switch answer {
	case "y":
		searchBase := "ou=roles,dc=supplyframe,dc=com" // Base DN for roles
		//fmt.Println("Searching for members in the role called", roleName, "in:", searchBase)
		searchFilter := "(objectClass=*)"

		searchRequest := ldap.NewSearchRequest(
			searchBase,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			searchFilter,
			[]string{"cn"},
			nil,
		)

		searchResult, err := l.Search(searchRequest)
		if err != nil {
			fmt.Println("LDAP search error:", err)
			return
		}

		fmt.Println("Role names:")
		for _, entry := range searchResult.Entries {
			roleName := entry.GetAttributeValue("cn")
			fmt.Println(roleName)
		}

	case "n":
	}

	//Ask for role name
	var roleName string
	fmt.Print("Enter the role name: ")
	fmt.Scanln(&roleName)

	//Search
	searchBase := "ou=roles,dc=supplyframe,dc=com" // Base DN for roles
	fmt.Println("Searching for members in the role called", roleName, "in:", searchBase)
	searchFilter := fmt.Sprintf("(cn=%s)", roleName)

	searchRequest := ldap.NewSearchRequest(
		searchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{"member"},
		nil,
	)

	searchResult, err := l.Search(searchRequest)
	if err != nil {
		fmt.Println("LDAP search error:", err)
		return
	}

	//Generate CSV
	var csvName string
	fmt.Print("Enter CSV file output name (.csv extension will be added automatically): ")
	fmt.Scanln(&csvName)
	csvFileName := csvName + ".csv"
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		fmt.Println("CSV file creation error:", err)
		return
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	//Write to CSV File
	for _, entry := range searchResult.Entries {
		for _, memberDN := range entry.GetAttributeValues("member") {
			//Remove extra text 
			memberDN = strings.Replace(memberDN, "cn=", "", -1)                         // Remove "cn="
			memberDN = strings.TrimSuffix(memberDN, ",ou=People,dc=supplyframe,dc=com") // Remove suffix

			csvWriter.Write([]string{memberDN})
		}
	}

	fmt.Println("CSV file created successfully.")
}
