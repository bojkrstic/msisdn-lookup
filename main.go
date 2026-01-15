package main

import (
	"fmt"
	"lookup/lookup"
	"lookup/web"
	"net/http"
)

func main() {

	http.HandleFunc("/", web.IndexHandler)
	http.HandleFunc("/lookup", lookup.Handler)
	http.HandleFunc("/lookup-view", web.LookupViewHandler)
	http.HandleFunc("/batch", lookup.BatchHandler)

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)

	fmt.Println(lookup.Country("+393383260866")) // Italy
	fmt.Println(lookup.Country("+38164123456"))  // Serbia
	fmt.Println(lookup.Country("+38598123456"))  // Croatia

	fmt.Println(lookup.NumberType("+393383260866")) // mobile
	fmt.Println(lookup.NumberType("+390636918899")) // fixed
	fmt.Println(lookup.NumberType("+38164123456"))  // mobile
	fmt.Println(lookup.NumberType("+38111345678"))  // fixed

	fmt.Println(lookup.IsValidLength("+393383260866")) // true  (Italy)
	fmt.Println(lookup.IsValidLength("+390636918899")) // true  (Italy fixed)

	fmt.Println(lookup.IsValidLength("+38164123456")) // true  (Serbia mobile)
	fmt.Println(lookup.IsValidLength("+38111345678")) // true  (Serbia fixed)

	fmt.Println(lookup.IsValidLength("+393"))  // false
	fmt.Println(lookup.IsValidLength("+3816")) // false

}
