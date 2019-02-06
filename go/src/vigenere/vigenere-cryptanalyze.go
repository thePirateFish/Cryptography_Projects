// Author: Matt Gigliotti
// Fall 2018

package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strconv"
)

func main() {

	filename := os.Args[1]
	key_length, _ := strconv.Atoi(os.Args[2])
	//key_length, _ = strconv.Atoi(key_length)
	ct_byte_array, err := ioutil.ReadFile(filename)
	check(err)

	ct_length := len(ct_byte_array)

	ct_cosets := make([][]int, key_length)
	coset_lengths := make([]int, key_length) //store size of each coset
	coset_frequencies := make([][]float64, key_length) //frequences i.e. 0.2

	english_frequencies := []float64{0.082, 0.014, 0.028, 0.038, 0.131, 0.029, 0.020, 
		0.053, 0.064, 0.001, 0.004, 0.034, 0.025, 0.071, 0.080, 0.020, 0.001, 0.068, 
		0.061, 0.105, 0.025, 0.009, 0.015, 0.002, 0.020, 0.001}

	//our candidates for each possible shift, for each coset
	chi_squared_values := make([][]float64, key_length)

	for c, _ := range ct_cosets {
		ct_cosets[c] = make([]int, 26)
		coset_frequencies[c] = make([]float64, 26)
		chi_squared_values[c] = make([]float64, 26)
	}
	for i := 0; i < ct_length; i++ {
		letter := ct_byte_array[i]-65
		if letter >=0 && letter <= 25 {
			ct_cosets[i%key_length][letter]+=1
			coset_lengths[i%key_length]+=1
		}
		
	}

	//calculate frequencies
	for c, _ := range ct_cosets {
		for i:=0; i<26; i++ {
			freq := int(ct_cosets[c][i])
			freq_percent := float64(freq) / float64(coset_lengths[c])
			coset_frequencies[c][i] = freq_percent
		}
	}

	var sum float64
	//calculate chi squared values
	for letter := 0; letter < 26; letter++ {
		
		for c, _ := range coset_frequencies {
			sum = 0
			for i:=0; i<26; i++ {
				sum = sum + ((coset_frequencies[c][(i+letter)%26] - english_frequencies[i]) * (coset_frequencies[c][(i+letter)%26] - english_frequencies[i])) / english_frequencies[i]
			}
			chi_squared_values[c][letter] = sum
			//fmt.Println("Coset: ", c, " Letter Offset: ", letter, "Chi sq: ", chi_squared_values[c][letter])
		}
	}

	keyword := make([]byte, key_length)

	for c, _ := range chi_squared_values {
		min := 100.0
		
		for i:=0; i<26; i++ {
			if chi_squared_values[c][i] < min {
				min = chi_squared_values[c][i]
				keyword[c] = byte(i+65)
			}
		}

	}

	fmt.Println(string(keyword))
}

//function taken from gobyexample.com/reading-files
func check(e error) {
    if e != nil {
        panic(e)
    }
}