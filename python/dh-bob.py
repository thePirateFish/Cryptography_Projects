#!/usr/local/bin/python3

# Matt Gigliotti
# Fall 2018

#models Bob's receipt of the message from Alice, and outputs a response message back to Alice.

# usage: dh-bob <filename of message from Alice> <filename of message back to Alice>.
# Reads in Alice's message, outputs ( g^b ) to Alice, prints the shared secret g^ab.

import sys
from random import getrandbits, randrange


def main():

	message_file = sys.argv[1]
	response_file = sys.argv[2]

	message = open(message_file, "r")
	from_alice = message.readline()
	from_alice = from_alice.replace('(', '').replace(')','')
	values = from_alice.split(',')
	p = int(values[0])
	g = int(values[1])
	g_a = int(values[2])
	b = get_random_secret(p)
	g_b = exponentiate(g, b, p)
	g_a_b = exponentiate(g_a, b, p)

	response = open(response_file, "w")
	response.write("(%d)" % g_b)
	print(g_a_b)


def exponentiate(x, y, n):
	result = 1
	while y > 0:
		if y % 2 == 1:
			result = (result * x) % n
		y = y // 2
		x = (x * x) % n

	return result

def get_random_secret(p):
	x = randrange(1, p-1)
	return x

if __name__ == '__main__': main()