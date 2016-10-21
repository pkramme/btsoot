#!/usr/bin/env python3
#MIT License
#
#Copyright (c) 2016 Paul Kramme
#
#Permission is hereby granted, free of charge, to any person obtaining a copy
#of this software and associated documentation files (the "Software"), to deal
#in the Software without restriction, including without limitation the rights
#to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
#copies of the Software, and to permit persons to whom the Software is
#furnished to do so, subject to the following conditions:
#
#The above copyright notice and this permission notice shall be included in all
#copies or substantial portions of the Software.
#
#THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
#SOFTWARE.

import socket
import sys
import os


def transmit(path, TCP_IP, TCP_PORT = 8000):
	print(path)
	file = open(path, "r")
	blocksize = os.path.getsize(path)
	sock = socket.socket()
	sock.connect((TCP_IP, TCP_PORT))
	offset = 0
	while 1:
		sent = os.sendfile(sock.fileno(), file.fileno(), offset, blocksize)
		if sent == 0:
			print("Done.")
			break
		offset += sent


def receive(filename, TCP_PORT = 8000, TCP_IP = ''):
	#TCP_IP = '' #FILL WITH IP WHICH REPRESENTS THE SERVER
	file = open(filename, "wb")

	s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	s.bind((TCP_IP, TCP_PORT))
	s.listen(1)

	conn, addr = s.accept()
	print('Connection address:', addr)
	while 1:
		data = conn.recv(4096)
		if not data:
			print("Done")
			break
		file.write(data)
	file.close()
	conn.close()


def main():
	if sys.argv[1] == "receive":
		receive(sys.argv[2])#, sys.argv[3], sys.argv[4])
	elif sys.argv[1] == "transmit":
		transmit(sys.argv[2], sys.argv[3])#, sys.argv[4])
	else:
		print("USAGE: " + sys.argv[0] + " transmit path ip port\nreceive path port ip")

if __name__ == __name__:
	main()
elif:
	print("Not a library.")
