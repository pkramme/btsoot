#!/usr/bin/env python3.5
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

debug = False
import config
import sys

def main():
	print("BTSOOT 0.1.0")
		
	if(len(sys.argv) > 1):
		if "create" in sys.argv:
			if config in sys.argv:
				config.create()
		else:
			print("Usage: BTSOOT command")
	else:
		print("Console Mode.") #Should be used primarily for Win
		while 1:
			print(">", end="")
			consoleinput = input()
			if consoleinput == "debug=true":
				debug = True
				print("Debug enabled.")
			elif consoleinput == "debug=false":
				debug = False
				print("Debug disabled.")
			elif consoleinput == "create config":
				config.create()
			elif consoleinput == "exit":
				sys.exit()
			else:
				print("Command not found")
				
if __name__ == __name__:
	try:
		main()
	except KeyboardInterrupt:
		print("Stopping program.")
		sys.exit()
	except Exception:
		print("Unknown critical exception")
		sys.exit()