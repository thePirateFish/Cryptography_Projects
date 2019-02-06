// Author: Matt Gigliotti
// Fall 2018

package main

import (
	"os"
	"flag"
	"fmt"
	"io/ioutil"
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strings"
	"errors"
)

//const BLOCK_SIZE int = 64

func main() {
	var mode string
	var k_enc [] byte
	var k_mac [] byte
	var arg_key_string string
	var input_file_name string
	var output_file_name string
	mode = os.Args[1]

	fs := flag.NewFlagSet("fs", flag.ContinueOnError)

	fs.StringVar(&arg_key_string, "k", "0", "32-byte key in hexadecimal")
	fs.StringVar(&input_file_name, "i", "input.txt", "input file")
	fs.StringVar(&output_file_name, "o", "output.txt", "output file")

	fs.Parse(os.Args[2:])
	arg_key_string = strings.Trim(arg_key_string, " ")


	data, err := hex.DecodeString(arg_key_string)
	if err != nil {
    	panic(err)
	}

	k_enc, k_mac = data[:16], data[16:]

	if (mode == "encrypt") {
		pt_byte_array, err := ioutil.ReadFile(input_file_name)
		check(err)
		c := encrypt(k_enc, k_mac, pt_byte_array)
		err = ioutil.WriteFile(output_file_name, c, 0644)
		check(err)
	}

	if (mode == "decrypt") {
		ct_byte_array, err := ioutil.ReadFile(input_file_name)
		check(err)
		p, err := decrypt(k_enc, k_mac, ct_byte_array)
		check(err)
		err = ioutil.WriteFile(output_file_name, p, 0644)
		check(err)
	}
	

}

func encrypt(k_enc, k_mac, m [] byte) [] byte {

	tag := hmac_sha256(k_mac, m)
	m_prime := append(m, tag...) //variadic function, interesting
	m_double_prime := apply_pkcs_5_padding(m_prime)
	iv := generate_iv(16)
	var iv_orig = make([]byte, 16)
	copy(iv_orig, iv)
	c_prime := aes_cbc_enc(k_enc, iv, m_double_prime)
	c := append(iv_orig, c_prime...)

	return c
}

func decrypt(k_enc, k_mac, c [] byte) ([] byte, error) {

	iv, c_prime := c[:16], c[16:]
	m_double_prime := aes_cbc_dec(k_enc, iv, c_prime)
	m_prime, err := validate_padding(m_double_prime)
	check(err)
	l := len(m_prime)
	if (l < 32) {
		fmt.Print("INVALID MAC")
		return m_prime, errors.New("INVALID MAC")
	}
	m, tag := m_prime[:l-32], m_prime[l-32:]
	plaintext, err2 := validate_mac(k_mac, m, tag)
	check(err2)


	return plaintext, nil

}

func aes_cbc_dec(k_enc, iv, ciphertext [] byte) [] byte {

	aes_cipher, error := aes.NewCipher(k_enc)
	check(error)
	block_size := aes_cipher.BlockSize() //should be 16 bytes
	num_blocks := len(ciphertext) / block_size

	plaintext := make([]byte, len(ciphertext))
	ct_last := make([]byte, block_size)
	ct_block := make([]byte, block_size)
	decrypted_block := make([]byte, block_size)
	ct_last = iv //our XOR vector starts out as IV

	for i := 0; i < num_blocks; i++ {
		block_offset := i*block_size
		copy(ct_block, ciphertext[block_offset:block_offset+block_size])
		aes_cipher.Decrypt(decrypted_block, ct_block)
		
		for j:= 0; j < block_size; j++ {
			decrypted_block[j] ^= ct_last[j] 
		}
		//set up ct_last to be the ciphertext block we just decrypted
		copy(ct_last, ct_block)
		copy(plaintext[block_offset:], decrypted_block)
	}

	return plaintext


}

func aes_cbc_enc(k_enc, iv, plaintext [] byte) [] byte {

	aes_cipher, error := aes.NewCipher(k_enc)
	check(error)
	block_size := aes_cipher.BlockSize() //should be 16 bytes
	num_blocks := len(plaintext) / block_size

	ciphertext := make([]byte, len(plaintext))
	cbc := make([]byte, block_size)
	cbc = iv //our XOR vector starts out as IV

	for i := 0; i < num_blocks; i++ {
		block_offset := i*block_size
		for j:= 0; j < block_size; j++ {
			cbc[j] ^= plaintext[block_offset+j] 
		}
		

		aes_cipher.Encrypt(cbc, cbc)
		// our XOR vector cbc is now ready for the next block, too

		copy(ciphertext[block_offset:], cbc)
	}

	return ciphertext
}

func generate_iv(length int) [] byte {
	arr := make([]byte, length)
	rand.Read(arr)
	return arr

}

func apply_pkcs_5_padding(m_prime [] byte) [] byte {

	n := len(m_prime) % 16
	padding_length := 16 - n

	if (padding_length == 0){
		padding_length = 16
	}

	ps := make([]byte, padding_length)
	
	for i := range ps {
		ps[i] = byte(padding_length)
	}

	m_double_prime := append(m_prime, ps...)
	return m_double_prime

}

func validate_padding(m_double_prime [] byte) ([] byte, error) {

	last_index := len(m_double_prime)-1
	last_byte := m_double_prime[last_index]

	pad_num := int(last_byte)
	var m_prime []byte
	
	for i := last_index; i > last_index - pad_num; i-- {
		if (int(m_double_prime[i]) != pad_num) {
			fmt.Println("INVALID PADDING")
			return m_prime, errors.New("INVALID PADDING")
		}
	}

	m_prime = m_double_prime[:last_index-pad_num+1]
	//fmt.Println("padding number: ", last_byte)

	return m_prime, nil

}

func validate_mac(k_mac, m, tag [] byte) ([] byte, error) {

	calculated_tag := hmac_sha256(k_mac, m)
	//var p []byte

	for i := range tag {
		if (calculated_tag[i] != tag[i]) {
			return m, errors.New("INVALID MAC")
		}
	}
	
	return m, nil
}


/*
	The HMAC-SHA256 algorithm. Based on the crypto/hmac package
	implementation.	
*/
func hmac_sha256(k_mac [] byte, m [] byte) [] byte {
	inner_hasher := sha256.New()
	outer_hasher := sha256.New()
	block_size := inner_hasher.BlockSize()
	ipad := make([]byte, block_size)
	opad := make([]byte, block_size)

	if (len(k_mac) > block_size) {
		//fmt.Println("k_mac > block_size (64 bytes)")
		inner_hasher.Write(k_mac)
		k_mac = inner_hasher.Sum(nil)
		inner_hasher.Reset()
	}

	copy(ipad, k_mac)
	copy(opad, k_mac)

	// Basically, XORing the k_mac with each repeating padding value
	for i := range ipad {
		ipad[i] ^= 0x36
	}
	for i := range opad {
		opad[i] ^= 0x5c
	}

	inner_hasher.Write(ipad)
	//hash_of_ipad_and_message := inner_hasher.Sum(m)
	inner_hasher.Write(m)
	hash_of_ipad_and_message := inner_hasher.Sum(nil)
	outer_hasher.Reset()
	outer_hasher.Write(opad)
	outer_hasher.Write(hash_of_ipad_and_message)
	tag := outer_hasher.Sum(nil)

	return tag

}

//function taken from gobyexample.com/reading-files
func check(e error) {
    if e != nil {
        panic(e)
    }
}