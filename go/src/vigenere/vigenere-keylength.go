// Author: Matt Gigliotti
// Fall 2018

package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"math"
)

func main() {

	filename := os.Args[1]
	ct_byte_array, err := ioutil.ReadFile(filename)
	check(err)

	ct_length := len(ct_byte_array)

	//used to track the closest
	var min_difference float64
	min_difference = 1

	//eventual selection
	chosen_keylength := 0

	for key_length := 3; key_length <= 20; key_length++ {

		//create cosets
		ct_cosets := make([][]int, key_length)
		for c, _ := range ct_cosets {
			ct_cosets[c] = make([]int, 26)
		}

		for i := 0; i < ct_length; i++ {
			letter := ct_byte_array[i]-65
			if letter >=0 && letter <= 25 {
				ct_cosets[i%key_length][letter]+=1
			}
			
		}

		var sum int
		var total_chars int
		var average float64
		average = 0
		//intermediate_ics := make([]float32, key_length)

		for c, _ := range ct_cosets {
			sum = 0
			total_chars = 0
			for i:=0; i<26; i++ {
				freq := int(ct_cosets[c][i])
				if freq >=1 {
					sum = sum + (freq * (freq - 1))
					total_chars+=freq 
				}
			}
			//intermediate_ics[c] = sum / (total_chars * (total_chars-1))
			average = average + (float64(sum) / float64(total_chars * (total_chars-1)))
		}

		average = average / float64(key_length)
		//fmt.Println("Keylength: ", key_length, "IC: ", average)
		english_ic := .0667

		if math.Abs(english_ic - average) < min_difference {
			min_difference = english_ic - average
			chosen_keylength = key_length
		}

	} 

	fmt.Println(chosen_keylength)
}

//function taken from gobyexample.com/reading-files
func check(e error) {
    if e != nil {
        panic(e)
    }
}