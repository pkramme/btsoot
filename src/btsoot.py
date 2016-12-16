#!/usr/bin/env python3

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
	print("MORE IS COMMING. Please report any bugs to https://paukra.com/paulkramme/btsoot/")
	print("https://github.com/paulkramme/btsoot/")


def backup(name):
	database = open("btsootdb", "r")
	if name in database:
		print("Positive match. Proceeding.")
		
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


if __name__ == "__main__":
	try:
		main()
	except KeyboardInterrupt:
		print("Keyboard Interrupt. Exiting.")
		sys.exit()

