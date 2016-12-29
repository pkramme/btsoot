#!/usr/bin/env python3

import sys
import os
import time

try:
	from compare import compare
except ImportError:
	print("Failed to import compare library.")
	print("BTSOOT can download the missing library. This requires Git and an Internet connection.")
	if input("Should i try?") == "y" or input("Should i try?") == "":
		system("git clone https://paukra.com/open-source/compare.git")
	else:
		print("Aborting. You have to manualy install it then and/or restart the program.")
	exit()

try:
	from datatransfer import datalib
except ImportError:
	print("Failed to import datatransfer library.")
	print("BTSOOT can download the missing library. This requires Git and an Internet connection.")
	print("Should i try? (Y/n)")
    if input("Should i try?") == "y" or input("Should i try?") == "":
		system("git clone https://paukra.com/open-source/datatransfer.git")
	else:
		print("Aborting. You have to manualy install it then and/or restart the program.")
	exit()

usage = f"USAGE: {sys.argv[0]} add/rm/scan "

class color:
	HEADER = '\033[95m'
	OKBLUE = '\033[94m'
	OKGREEN = '\033[92m'
	WARNING = '\033[93m'
	FAIL = '\033[91m'
	ENDC = '\033[0m'
	BOLD = '\033[1m'
	UNDERLINE = '\033[4m'


def split(string, splitters): #MAY RESOLVE ALL PROBLEMS WITH CSV
	final = [string]
	for x in splitters:
		for i,s in enumerate(final):
			if x in s and x != s:
				left, right = s.split(x, 1)
				final[i] = left
				final.insert(i+1, x)
				final.insert(i+2, right)
	return final


def scandirectory(walk_dir, scanfile, verbose = False):
	try:
		print("Initializing scan...")
		with open(scanfile, "w") as f:
			f.write("path,checksum\n")
			for root, subdirs, files in os.walk(walk_dir):
				#f.write(root + "\n")
				for filename in files:
					file_path = os.path.join(root, filename)
					checksum = compare.md5sum(file_path)
					if verbose == True:
						print(file_path, checksum, end="\n")
					else:
						pass
					f.write(file_path + "," + checksum + "\n")
		print("Done.")
	except FileNotFoundError:
		print("There was a reading error... Probably os protected.")


def main():
	try:
		if sys.argv[1] == "add":
			try:
				name = sys.argv[2]
				path = sys.argv[3]
				server = sys.argv[4]
			except IndexError:
				print("Usage: " + sys.argv[0] + " add name path server")
				print("Type local if you want to use local filesystem.")
				exit()
			with open("btsoot.conf", "a") as f:
				f.write("name=" + name + '\n')
				f.write("path=" + path + '\n')
				f.write("server=" + server + '\n')

		elif sys.argv[1] == "rm":
			try:
				name = sys.argv[2]
			except IndexError:
				print("Usage: " + sys.argv[0] + "rm name")
				exit()
			try:
				f = open("btsoot.conf", "r")
				row = 0
				beginning_row = -10
				indentifier = "name=" + name + '\n'
				lines = f.readlines()
				f.close()
				f = open("btsoot.conf", "w")
				for line in lines:
					row = row + 1
					if line == indentifier:
						beginning_row = row
					elif row == beginning_row + 1 or row == beginning_row + 2:
						pass
					else:
						f.write(line)
				f.close()
			except FileNotFoundError:
				print("Configfile not found. Create one with 'add'.")

		elif sys.argv[1] == "scan":
			print("Execute scan...")
			#NAME NEEDS TO BE RESOLVED TO CORRECT DIRECTORY!
			#USE PPARTIAL FUNCTION FROM 'rm'
			searched_path = None			
			try:
				name = sys.argv[2]
				scanfilename = "{}_{}.btsscan".format(int(time.time()), name)
			except IndexError:
				print("Usage: " + sys.argv[0] + "scan name")
			try:
				f = open("btsoot.conf", "r")
				row = 0
				beginning_row = -1 #set counter to a negative state so it's not finding any rows
				identifier = "name=" + name + '\n'
				lines = f.readlines()
				f.close()
				for line in lines:
					row = row + 1
					if line == identifier:
						beginning_row = row
					elif row == beginning_row + 1:
						searched_path = line
						break
					else:
						pass
				path_with_newline = split(searched_path, "=")
				tempstring = path_with_newline[2]
				path = tempstring.rstrip()#GETS RID OF NEWLINE
				print(path)

			except FileNotFoundError:
				print("Configfile not found. Create one with 'add'.")
			scandirectory(path, scanfilename, True)

		elif sys.argv[1] == "backup":
			print("Filetransfer initialized...")
			#TODO
			# FIGURE OUT A WAY TO CREATE ALL NECESSARY DIRECTORIES
			# FIGURE OUT A WAY TO LEAVE THE CONNECTION OPEN AND SEND THE 
			#  FILENAME

			# Open two scanfiles and compare them
			#scandirectory("./", "localscanlist")
			#scanlistfile = open("localscanlist", "r")
			dirs = os.listdir("./")
			number_of_files = 0
			for file in dirs:
				fileending = split(file, ".")
				#print(fileending)
				try:
					if fileending[2] == "btsscan":
						print(file)
						number_of_files = number_of_files + 1
					else:
						pass
				except IndexError:
					pass
			if number_of_files == 0:
				print("There arent any files. Create one by using 'scan'.")
			elif number_of_files == 1:
				print("There is only one file. I'll need two.")
			else:
				print("Sufficient number of files were found.")

		else:
			print(usage)

	except IndexError:
		print(usage)
		exit()


if __name__ == "__main__":
	try:
		main()
	except KeyboardInterrupt or IndexError:
		print("\nStopping program.")
		sys.exit()
