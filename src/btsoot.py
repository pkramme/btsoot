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
import database
import sys
import datatransfer-lib


def usage():
	print("Usage: btsoot command")
	print("\tcreate config		| Creates configfile")
	print("\tcreate database 	| Creates database")
	print("\tdebug=true/false	| Sets debug mode")
	print("MORE IS COMMING. Please report any bugs to")
	print("https://github.com/paulkramme/btsoot/")


def backup(name):
	database = open("btsootdb", "r")
	if name in database:
		print("Positive match. Proceding.")
		
	else:
		print("Negative Match. Aborting.")


def main():	
	#DEBUG MODE LOADER
	try:
		loadconfig = open("btsoot.conf", "r")
		if "debug=true" in loadconfig.readline():
			debug = True
		elif "debug=false" in loadconfig.readline():
			debug = False
		else:
			pass
		loadconfig.close()
	except FileNotFoundError:
		print("Configfile not found. You should create one with 'create config'.")
	
	#SYSTEM ARGS
	if(len(sys.argv) > 1):
		if "create" in sys.argv:
			if "config" in sys.argv:
				config.create()
			if "database" in sys.argv:
				database.create()
		elif "debug=true" in sys.argv:
			config.configfile = open("btsoot.conf", "w")
			config.configfile.write("debug=true\n")
			config.configfile.close()
		elif "add" in sys.argv:
			if "block" in sys.argv:
				newblock = input("New blocks name: ")
				databaseblock = open("btsootdb", "r+")
				if newblock in databaseblock.readline():
					print("The block '" + newblock + "' already exists")
				else:
					path = input("Path: ")
					databaseblock.write(newblock + " " + path + "\n")
					databaseblock.close()
			else:
				print("block?")
		elif "backup" in sys.argv:
			print("Backup in Progress...")
			backup()
		elif "version" in sys.argv:
			print("BTSOOT 0.1.0")
		elif "help" or "usage" in sys.argv:
			usage()
		else:
			usage()


if __name__ == __name__:
	try:
		main()
	except KeyboardInterrupt:
		print("Stopping program.")
		sys.exit()
