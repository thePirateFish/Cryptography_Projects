#!/usr/local/bin/python3

# Matt Gigliotti
# Fall 2018

# dl-efficient

# an implementation of Pollard's rho algorithm for solving discrete logarithms

# resources: https://en.wikipedia.org/wiki/Pollard%27s_rho_algorithm_for_logarithms
# HAC http://cacr.uwaterloo.ca/hac/about/chap3.pdf
# https://www.geeksforgeeks.org/gcd-in-python/
# https://en.wikipedia.org/wiki/Pohlig%E2%80%93Hellman_algorithm

# NOTE: This Pollard's rho implementation does NOT work with my dh-alice1.py
# because the 'g' supplied by that program is not a generator, but a high order element
# since alice1 generates safe primes to use

import sys
from random import getrandbits, randrange


beta = 0
alpha = 0
n = 0

def main():

	# prime = 1019
	# a = 10
	# generator = 2
	# sample_g_a = exponentiate(generator, a, prime)
	# example = open("example.txt", "w")
	# example.write("(%d, %d, %d)" % (prime, generator, sample_g_a))
	# example.close()

	message_file = sys.argv[1]


	message = open(message_file, "r")
	#message = open("example.txt", "r")
	from_alice = message.readline()
	from_alice = from_alice.replace('(', '').replace(')','')
	values = from_alice.split(',')
	p = int(values[0])
	g = int(values[1])
	g_a = int(values[2])

	global beta
	beta = g_a
	global alpha
	alpha = g
	global n
	n = p-1


	x = pollards_rho()
	print("%d" % x)


# alpha = generator g, beta = g^a, n = p-1 (order)
def pollards_rho():
	x = 1
	a = 0
	b = 0
	X = x
	A = a
	B = b

	i = 0
	for i in range(n):
		x, a, b = next(x, a, b)
		X, A, B = next(X, A, B)
		X, A, B = next(X, A, B)

		if (x == X):
			r = (B - b)
			if (r == 0):
				return -1
			x = (((a - A) % n) // r) % n
			return x
	return 0


def next(x, a, b):
	j = x % 3
	if (j == 0):
		x = (x * x) % (n+1)
		a = (2 * a) % n 
		b = (2 * b) % n
		return x, a, b
	if (j == 1):
		x = (x * alpha) % (n+1)
		a = (a + 1) % n 
		return x, a, b
	if (j == 2):
		x = (x * beta) % (n+1)
		b = (b + 1) % n 
		return x, a, b
	return 0, 0, 0

def exponentiate(x, y, n):
	result = 1
	while y > 0:
		if y % 2 == 1:
			result = (result * x) % n
		y = y // 2
		x = (x * x) % n

	return result

# NOTE - not actually used, was looking at the HAC simplified version
# The Euclidean Algorithm as shown in
# https://www.geeksforgeeks.org/gcd-in-python/
def gcd(x, y):
	while(y):
		x, y = y, x%y
	return x



if __name__ == '__main__': main()