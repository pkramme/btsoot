#!/usr/bin/env python3.6


import sys
import os
import time
import shutil


class color:
	HEADER = '\033[95m'
	OKBLUE = '\033[94m'
	OKGREEN = '\033[92m'
	WARNING = '\033[93m'
	FAIL = '\033[91m'
	ENDC = '\033[0m'
	BOLD = '\033[1m'
	UNDERLINE = '\033[4m'


try:
	from compare import compare
except ImportError:
	print(color.FAIL + "Failed to import compare library." + color.ENDC)
	print("BTSOOT can download the missing library.")
	print("This requires Git and an Internet connection.")
	if input("Should i try? ") == "y":
		os.system("git clone https://git.paukra.com/open-source/compare.git")
		from compare import compare
	else:
		print("Aborting. You have to manualy install it then and/or restart the program.")
		exit()


usage = f"""USAGE: {sys.argv[0]} <commands>

add <name> <path> <server/local>\tadd block
rm <name>\t\t\t\tremove added block
scan <name>\t\t\t\tscan added block
backup <name>\t\t\t\tbackup scanned block
update_dependencies\t\t\tupdate the needed libraries
"""


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
				f.write(root + "\n")
				for filename in files:
					file_path = os.path.join(root, filename)
					checksum = compare.crc(file_path)
					if verbose == True:
						print(file_path, checksum, end="\n")
					else:
						pass
					f.write(file_path + "," + checksum + "\n")
		print("Done.")
	except FileNotFoundError:
		print(color.FAIL + "There was a reading error... Probably os protected." + color.ENDC)


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
			#LIST FILES TO FIND SCANFILES
			#SORT OUT ANY UNINTERESTING FILES
			scanfilelist = []
			dirs = os.listdir("./")
			number_of_files = 0
			for file in dirs:
				blockname = split(file, ["_", "."])
				try:
					if blockname[4] == "btsscan" and blockname[2] == sys.argv[2]:
						number_of_files = number_of_files + 1
						scanfilelist.append(file)
					else:
						pass
				except IndexError:
					pass
			
			if number_of_files == 0:
				print("There aren't any scan files.")
				print(f"Create one by using\n{sys.argv[0]} scan <block name>.")
			
			elif number_of_files == 1:
				print("There is only one file. That means a complete backup must be created.")
				#TODO: TRANSFER FILE TO SERVER, RESOLVE SERVER ADDR
			
			else:
				print("Sufficient number of scan files were found.")
				splitted_timestamp = []

				#FIND LATEST TWO FILES
				#SPLIT EVERY FILE NAME TO GAIN TIMESTAMP
				for scanfile in scanfilelist:
					temp = split(scanfile, "_")
					splitted_timestamp.append(int(temp[0]))

				#GETS LATEST SCANFILE'S TIMESTAMP
				latest_timestamp = max(splitted_timestamp)

				#SETS MAX VALUE TO -1 TO FIND SECOND HIGHEST VALUE
				listcounter = 0
				for timestamp in splitted_timestamp:
					if timestamp == latest_timestamp:
						splitted_timestamp[listcounter] = -1
					listcounter = listcounter + 1

				#GET PREVIOUS FILE'S TIMESTAMP
				previous_timestamp = max(splitted_timestamp)

				dircounter = 0
				latest_scan_array_index = -1
				previous_scan_array_index = -1
				for file in scanfilelist:
					temp = split(file, "_")
					#print(f"Check {temp[0]} against {latest_timestamp} and {previous_timestamp}")
					if int(temp[0]) == latest_timestamp:
						latest_scan_array_index = dircounter
					elif int(temp[0]) == previous_timestamp:
						previous_scan_array_index = dircounter
					else:
						pass
					dircounter = dircounter + 1

				print("Latest scan: " + scanfilelist[latest_scan_array_index])
				print("Previous scan: " + scanfilelist[previous_scan_array_index] + "\n")

				#COMPARE THE TWO FILES AGAINST EACH OTHER
				latest_scan_fd = open(scanfilelist[latest_scan_array_index], "r")
				previous_scan_fd = open(scanfilelist[previous_scan_array_index], "r")
				transmit_list_fd = open("transmit.list", "w")

				latest_scan = latest_scan_fd.readlines()
				previous_scan = previous_scan_fd.readlines()

				file_same = 0
				file_new = 0
				file_total = 0


				for line in latest_scan:
					if line in previous_scan:
						file_same = file_same + 1
					else:
						print(color.OKGREEN + line + color.ENDC)
						transmit_list_fd.write(line)
						file_new = file_new + 1
					file_total = file_total + 1


				block_change_percentage = file_new / file_total * 100
				block_change_percentage = int(block_change_percentage)
				print(f"Total files: {file_total}")
				print(f"Unchanged files: {file_same}")
				print(f"New/Changed files: {file_new}")
				print(color.OKBLUE + f"Block changed by {block_change_percentage}%" + color.ENDC)


				previous_scan_fd.close()
				latest_scan_fd.close()


				#FIND SERVER ADDRESS
				searched_row = None
				f = open("btsoot.conf", "r")
				row = 0
				beginning_path = -1
				indentifier = "name=" + sys.argv[2] + '\n'
				lines = f.readlines()
				f.close()

				for line in lines:
					row = row + 1
					if line == indentifier:
						beginning_row = row
					if row == beginning_row + 2:
						searched_path = line
					else:
						pass

				serverstring = split(searched_path, "=")
				addr = serverstring[2]
				"""
				print(f"Sending to {addr}")
				transmitlist = transmit_list_fd.readlines()
				print(transmitlist)
				for line in transmitlist:
					print(line)
"""
				transmit_list_fd.close()
				
				
				
				


		elif sys.argv[1] == "update_dependencies":
			print("Updating dependecies.")
			print("This requires an internet connection. ")
			if input("Should i continue?") == "y":
				shutil.rmtree("compare")
				os.system("git clone https://git.paukra.com/open-source/compare.git")
				shutil.rmtree("datatransfer")
				os.system("git clone https://git.paukra.com/open-source/datatransfer.git")
			else:
				print(color.FAIL + "Abort." + color.ENDC)

		else:
			print(usage)

	except IndexError:
		print(usage)
		sys.exit()


if __name__ == "__main__":
	try:
		main()
	except KeyboardInterrupt:
		print("\nInterrupt by keyboard.\n")
		sys.exit()
