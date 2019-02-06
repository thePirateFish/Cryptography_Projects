#!/usr/local/bin/python3

# Matt Gigliotti
# Fall 2018

# models Alice's receipt of Bob's secret.

# usage: dh-alice2 <filename of message from Bob> <filename to read secret key>.
# Reads in Bob's message and Alice's stored secret, prints the shared secret g^ab.

import sys
from random import getrandbits, randrange


def main():

	key_length = 40
	k = 2 # primality tests to check

	message_file = sys.argv[1]
	secret_file = sys.argv[2]


	message = open(message_file, "r")
	secret = open(secret_file, "r")
	from_bob = message.readline()
	from_bob = from_bob.replace('(', '').replace(')','')
	g_b = int(from_bob)
	
	secrets = secret.readline()
	secrets = secrets.replace('(', '').replace(')','')
	values = secrets.split(',')
	p = int(values[0])
	g = int(values[1])
	a = int(values[2])
	g_a_b = exponentiate(g_b, a, p)


	print(g_a_b)

def exponentiate(x, y, n):
	result = 1
	while y > 0:
		if y % 2 == 1:
			result = (result * x) % n
		y = y // 2
		x = (x * x) % n

	return result


if __name__ == '__main__': main()
