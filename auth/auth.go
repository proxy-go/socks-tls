package auth

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
)

func LdapVerify(ldapAddr string, ldapBaseDN string, username string, password string) bool {
	// Connect to LDAP server
	l, err := ldap.Dial("tcp", ldapAddr)
	if err != nil {
		log.Printf("Connect failed: %s", err)
		return false
	}
	defer l.Close()

	// Bind with service account
	err = l.Bind(fmt.Sprintf("uid=%s,ou=people,%s", username, ldapBaseDN), password)
	if err != nil {
		log.Printf("Bind failed: %s", err)
		return false
	}

	// Search for user
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=people,%s", ldapBaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=person)(uid=%s))", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Printf("Search failed: %s", err)
		return false
	}

	if len(sr.Entries) != 1 {
		log.Printf("User %s not found or too many entries returned", username)
		return false
	}

	// Bind with user's DN and password
	userDN := sr.Entries[0].DN
	err = l.Bind(userDN, password)
	if err != nil {
		log.Printf("Bind failed: %s", err)
		return false
	}

	// Authentication successful
	log.Printf("%s authentication successful", username)
	return true
}
