// Author: Matt Gigliotti
// Fall 2018

package main

import (
	"os"
	"fmt"
	"io/ioutil"
)

func main() {

	//checkNumArgs(os.Args)
	enc_key := os.Args[1]
	filename := os.Args[2]

	key_length := checkKeyLength(enc_key)
	//TODO check key is all uppercase letters
	//TODO check file size

	//read file to byte array
	pt_byte_array, err := ioutil.ReadFile(filename)
	check(err)

	//create our ciphertext byte array
	ct_byte_array := make([]byte, len(pt_byte_array))

	//what index we are using in the key
	key_index := 0

	for pt_index := 0; pt_index < len(pt_byte_array); pt_index++ {

		//get index in alphabet of characters
		pt_char := pt_byte_array[pt_index] - 65
		key_char := enc_key[key_index] - 65

		//find the sum of the two characters
		//add 65 back to sum for ASCII value
		ct_char := ((pt_char + key_char) % 26) + 65

		//make substitution
		ct_byte_array[pt_index] = ct_char

		//increment our key index
		key_index = (key_index + 1) % key_length
	}

	fmt.Println(string(ct_byte_array))
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


