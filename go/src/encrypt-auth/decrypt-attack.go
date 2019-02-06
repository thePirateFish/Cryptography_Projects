// Author: Matt Gigliotti
// Fall 2018

package main

import (
	"flag"
	"os/exec"
	"io/ioutil"
	"fmt"
	"log"
)

func main() {

	var ciphertext_file string
	var decrypt_this_file string = "decrypt_this.txt"
	flag.StringVar(&ciphertext_file, "i", "cipher.txt", "ciphertext file")
	flag.Parse()

	ct_byte_array, err := ioutil.ReadFile(ciphertext_file)
	_=ct_byte_array
	check(err)

	current_ct_block := make([]byte, 16)

	// -16 because of the IV in the ct array
	decrypted_text := make([]byte, len(ct_byte_array)-16)

	//iterating backwards over blocks over the ciphertext
	// >=16 so that we stop after doing everything up until the IV
	for ct_block_index := len(ct_byte_array)-16; ct_block_index >=16; ct_block_index-=16 {

		// current_ct_block moves backward, byte by byte
		copy(current_ct_block, ct_byte_array[ct_block_index:ct_block_index+16])
		intermediate_bytes := make([]byte, 16)
		iv := make([]byte, 16)

		//individual byte in block
		for b := 1; b <= 16; b++ {
			iv_index := 16-b
			actual_xor_ct_byte := ct_byte_array[ct_block_index-b]
			//trying values until padding is correct
			for i := 0; i < 256; i++ {
				
				iv[iv_index] = byte(i)

				// send for decryption blocks C_1, C_2 
				two_blocks_to_decrypt := append(iv, current_ct_block...)
				err = ioutil.WriteFile(decrypt_this_file, two_blocks_to_decrypt, 0644)
				check(err)

				out, err := exec.Command("./decrypt-test", "-i", decrypt_this_file).Output()
			    if err != nil {
			        log.Fatal(err)
			    }

				return_message := string(out)
				if (return_message != "INVALID PADDING") {
					intermediate_bytes[iv_index] = iv[iv_index] ^ byte(b)
					pt_byte := intermediate_bytes[iv_index] ^ actual_xor_ct_byte
					//decrypted_text = append([]byte{pt_byte}, decrypted_text...)
					copy(decrypted_text[(ct_block_index-16)+iv_index:], ([]byte{pt_byte}))	//copy in decrypted pt

					for j := iv_index; j < 16 && b < 16; j++ {
						// set iv so that padding for next round is correct in rightmost bytes
						iv[j] = byte(b+1) ^ intermediate_bytes[j]
					}

					break
				}	
			}
		}
	}

	without_padding := remove_padding(decrypted_text)
	l := len(without_padding)
	m, tag := without_padding[:l-32], without_padding[l-32:]
	_=tag
	fmt.Println(string(m))
}

func remove_padding(m_double_prime [] byte) [] byte {

	last_index := len(m_double_prime)-1
	last_byte := m_double_prime[last_index]
	pad_num := int(last_byte)
	m_prime := m_double_prime[:last_index-pad_num+1]

	return m_prime
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}