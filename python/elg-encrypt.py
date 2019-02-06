#!/usr/local/bin/python3

# Matt Gigliotti
# Fall 2018

# elg-encrypt 

# Reads in the public key ( p, g, g^a ) produced by elg-keygen. Generates b and com-
# putes k = SHA256(g^a||g^b||g^ab). Outputs ( gb; AESGCM_k(M) ) to a ciphertextfle, where
# the latter value is encoded as a hexadecimal string.


import sys
from random import getrandbits, randrange
import hashlib
import os
from cryptography.hazmat.primitives.ciphers.aead import AESGCM
import binascii


def main():

	message_text = sys.argv[1]
	public_key_file = sys.argv[2]
	ciphertext_file = sys.argv[3]


	pk_output = open(public_key_file, "r")
	pk = pk_output.readline()
	pk = pk.replace('(', '').replace(')','')
	values = pk.split(',')
	p = int(values[0])
	g = int(values[1])
	g_a = int(values[2])

	b = get_random_secret(p)
	g_b = exponentiate(g, b, p)
	g_a_b = exponentiate(g_a, b, p)

	b_concat = (str(g_a) + " " + str(g_b) + " " + str(g_a_b)).encode()

	hasher = hashlib.sha256()
	hasher.update(b_concat)
	enc_key = hasher.digest()


	aesgcm = AESGCM(enc_key)
	nonce = os.urandom(16)
	ct = aesgcm.encrypt(nonce, message_text.encode(), None)

	ciphertext_file = open(ciphertext_file, "w")
	ciphertext_file.write("(%d, %s)" % (g_b, binascii.hexlify(nonce+ct).decode()))


def exponentiate(x, y, n):
	result = 1
	while y > 0:
		if y % 2 == 1:
			result = (result * x) % n
		y = y // 2
		x = (x * x) % n

	return result

def get_random_secret(p):
	x = randrange(1, p-2)
	return x

if __name__ == '__main__': main()