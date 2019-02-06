#!/usr/local/bin/python3

# Matt Gigliotti
# Fall 2018

# elg-decrypt

# Reads in the ciphertext produced by the previous program and a stored secret, prints
# the recovered message or error.

import sys
import hashlib
from cryptography.hazmat.primitives.ciphers.aead import AESGCM
import binascii

def main():

	ciphertext_file = sys.argv[1]
	secret_key_file = sys.argv[2]


	ciphertext_file = open(ciphertext_file, "r")
	line = ciphertext_file.readline()
	line = line.replace('(', '').replace(')','')
	values = line.split(',')
	g_b = int(values[0])
	ct = values[1].strip()

	secret_key_file = open(secret_key_file, "r")
	line = secret_key_file.readline()
	line = line.replace('(', '').replace(')','')
	values = line.split(',')
	p = int(values[0])
	g = int(values[1])
	a = int(values[2])

	g_a = exponentiate(g, a, p)
	g_a_b = exponentiate(g_b, a, p)
	
	b_concat = (str(g_a) + " " + str(g_b) + " " + str(g_a_b)).encode()

	hasher = hashlib.sha256()
	hasher.update(b_concat)
	dec_key = hasher.digest()

	ct_bytes = binascii.unhexlify(ct)
	nonce = ct_bytes[:16]
	ct = ct_bytes[16:]

	aesgcm = AESGCM(dec_key)
	plaintext = aesgcm.decrypt(nonce, ct, None).decode()
	print(plaintext)

def exponentiate(x, y, n):
	result = 1
	while y > 0:
		if y % 2 == 1:
			result = (result * x) % n
		y = y // 2
		x = (x * x) % n

	return result
if __name__ == '__main__': main()