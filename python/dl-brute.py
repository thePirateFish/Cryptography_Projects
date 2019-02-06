#!/usr/local/bin/python3

# Matt Gigliotti
# Fall 2018

# dl-brute

import sys
from random import getrandbits, randrange


def main():

	message_file = sys.argv[1]

	message = open(message_file, "r")
	from_alice = message.readline()
	from_alice = from_alice.replace('(', '').replace(')','')
	values = from_alice.split(',')
	p = int(values[0])
	g = int(values[1])
	g_a = int(values[2])

	x = find_x(p, g, g_a)
	print("%d" % x)

def find_x(p, g, g_a):
	while True:
		x = get_random_secret(p)
		if (exponentiate(g, x, p) == g_a):
			return x

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