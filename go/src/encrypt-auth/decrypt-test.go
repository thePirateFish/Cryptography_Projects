// Author: Matt Gigliotti
// Fall 2018

package main 

import (
	"flag"
	"fmt"
	"io/ioutil"
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)


func main() {

	var key_string string = "38885e04d782b4b6ecb6b2b2eeef55954c2606b738867fbbf4801525cbb876cd"
	key, err := hex.DecodeString(key_string)
	if err != nil {
    	panic(err)
	}

	var ciphertext_file string
	flag.StringVar(&ciphertext_file, "i", "cipher.txt", "ciphertext file")
	flag.Parse()

	k_enc, k_mac := key[:16], key[16:]

	ct_byte_array, err := ioutil.ReadFile(ciphertext_file)
	check(err)

	p, err := decrypt(k_enc, k_mac, ct_byte_array)
	_=p

	if (err==nil) {
		fmt.Print("SUCCESS")
	} //else {
	//	//fmt.Println("NOT")
	//}


}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func decrypt(k_enc, k_mac, c [] byte) ([] byte, error) {

	iv, c_prime := c[:16], c[16:]
	m_double_prime := aes_cbc_dec(k_enc, iv, c_prime)
	m_prime, err := validate_padding(m_double_prime)
	//check(err)
	if (err!=nil) {
		fmt.Print("INVALID PADDING")
		return m_prime, errors.New("INVALID PADDING")
	}
	l := len(m_prime)
	if (l < 32) {
		fmt.Print("INVALID MAC")
		return m_prime, errors.New("INVALID MAC")
	}
	m, tag := m_prime[:l-32], m_prime[l-32:]
	plaintext, err2 := validate_mac(k_mac, m, tag)
	if (err2!=nil) {
		fmt.Print("INVALID MAC")
		return plaintext, errors.New("INVALID MAC")
		//os.Exit(1)
	}


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

func validate_padding(m_double_prime [] byte) ([] byte, error) {

	last_index := len(m_double_prime)-1
	last_byte := m_double_prime[last_index]

	pad_num := int(last_byte)
	var m_prime []byte
	
	for i := last_index; i > last_index - pad_num; i-- {
		if (int(m_double_prime[i]) != pad_num) {
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
	//fmt.Println("ipad xor m length: ", len(hash_of_ipad_and_message))
	outer_hasher.Reset()
	outer_hasher.Write(opad)
	outer_hasher.Write(hash_of_ipad_and_message)
	tag := outer_hasher.Sum(nil)

	return tag

}