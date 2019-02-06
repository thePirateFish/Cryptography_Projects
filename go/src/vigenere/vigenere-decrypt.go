// Author: Matt Gigliotti
// Fall 2018

package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {

	//checkNumArgs(os.Args)
	dec_key := strings.TrimSpace(os.Args[1])
	filename := os.Args[2]


	key_length := checkKeyLength(dec_key)
	//TODO check key is all uppercase letters
	//TODO check file size

	//read file to byte array
	ct_byte_array, err := ioutil.ReadFile(filename)
	check(err)

	//create our plaintext byte array
	pt_byte_array := make([]byte, len(ct_byte_array))

	//what index we are using in the key
	key_index := 0

	for ct_index := 0; ct_index < len(ct_byte_array); ct_index++ {

		//get index in alphabet of characters
		ct_char := ct_byte_array[ct_index] - 65
		key_char := dec_key[key_index] - 65

		//Find difference 
		//add 65 back to sum for ASCII value

		pt_char := ((ct_char - key_char + 26) % 26) + 65

		//make substitution
		pt_byte_array[ct_index] = pt_char

		//increment our key index
		key_index = (key_index + 1) % key_length
	}

	fmt.Println(string(pt_byte_array))
}

//function taken from gobyexample.com/reading-files
func check(e error) {
    if e != nil {
        panic(e)
    }
}

/*func checkNumArgs(args []string) {
	if len(args) != 2 {
		panic ("Format for arguments: <key> <filename>")
	}
}*/

func checkKeyLength(key string) int {
	l := len(key)
	if l < 1 {
		panic("No key entered")
	}
	if l > 32 {
		panic("Key must be less than 32 characters")
	}
	return l
}