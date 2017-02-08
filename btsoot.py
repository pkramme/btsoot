#!/usr/bin/env python3.6

#CONFIGURATION################################################

#STORAGE CONFIG
configpath = ""
scanstorage = ""

#SAFETY GUARD CONFIG
safetyguard = True 
#Input min value in percent for cloning file override safety guard. 
#Backup will be aborted if change counter passes this value.
minwarningvalue = 75

#COPY CONFIG
copypath = ""

##############################################################
#DO NOT EDIT BELOW HERE!


import os, sys, time, shutil, zlib


#STARTUP CODE

if configpath == "":
	configpath = "/etc/btsoot/btsoot.conf"
if scanstorage == "":
	scanstorage = "/etc/btsoot/scans/"
if copypath == "":
	copypath = "/usr/local/bin/btsoot-copy"
if os.path.exists("/etc/btsoot") == True:
	pass
else:
	try:
		os.makedirs("/etc/btsoot/scans")
	except PermissionError:
		print("BTSOOT needs root permissions")
		sys.exit()


class color:
	HEADER = '\033[95m'
	OKBLUE = '\033[94m'
	OKGREEN = '\033[92m'
	WARNING = '\033[93m'
	FAIL = '\033[91m'
	ENDC = '\033[0m'
	BOLD = '\033[1m'
	UNDERLINE = '\033[4m'

#DEBUG FUNCTION AND SETVAR

debug = False
if "--debug" in sys.argv:
	debug = True
def dprint(message):
	if debug == True:
		print(f"DEBUG: {message}") 

def shouldcontinue(quit = True):
	if input("Should i continue? (yes/No)") == "yes":
		return 0
	else:
		if quit == True:
			sys.exit()
		else:
			return 1

def crc(filepath):
	previous = 0
	try:
		for line in open(filepath,"rb"):
			previous = zlib.crc32(line, previous)
	except OSError:
		print("CRC ERROR: OSERROR")
	return "%X"%(previous & 0xFFFFFFFF)


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
				final.insert(i + 1, x)
				final.insert(i + 2, right)
	return final


def scandirectory(walk_dir, scanfile, verbose = False):
	try:
		current_scan = []
		for root, subdirs, files in os.walk(walk_dir):
			current_scan.extend([f"{root}\n"])
			for filename in files:
				file_path = os.path.join(root, filename)
				checksum = crc(file_path)
				current_scan.extend([f"{file_path},{checksum}\n"])
		with open(scanfile, "w") as current_scan_file:
			current_scan_file.writelines(current_scan)			
	except FileNotFoundError:
		if verbose == True:
			print(color.FAIL + "SCAN ERROR: FILE NOT FOUND" + color.ENDC)


def main():
	try:
		if sys.argv[1] == "add":
			name = sys.argv[2]
			path = sys.argv[3]
			server = sys.argv[4]
			with open(configpath, "a") as conf:
				conf.write(f"{name},{path},{server}\n")


		elif sys.argv[1] == "rm":
			name = sys.argv[2]
			try:
				lines = []
				with open(configpath, "r") as conf:
					lines = conf.readlines()
				with open(configpath, "w") as conf:
					for line in lines:
						split_line = split(line, ",")
						if split_line[0] != name:
							conf.write(line)

			except FileNotFoundError:
				print(color.FAIL + "Configfile not found." + color.ENDC)
				print("Create one with 'add'.")


		elif sys.argv[1] == "list":
			try:
				with open(configpath, "r") as conf:
					for line in conf:
						split_line = split(line, ",")
						print(f"BLOCKNAME: {split_line[0]}")
						print(f"\tSRC:  {split_line[2]}")
						print(f"\tDEST: {split_line[4].rsplit()}")
			except FileNotFoundError:
				print(color.FAIL + "Configfile not found." + color.ENDC)
				print("Create one with 'add'.")


		elif sys.argv[1] == "backup":
			#REMOVE ENTREE FROM BTSOOT CONFIG
			searched_path = None
			name = sys.argv[2]
			scanfilename = "{}_{}.btsscan".format(int(time.time()), name)
			try:
				path = ""
				with open(configpath, "r") as conf:
					for line in conf:
						split_line = split(line, ",")
						path = split_line[2].rstrip()

				print(color.OKBLUE + f"Executing scan for block {sys.argv[2]}" + color.ENDC)

			except FileNotFoundError:
				print(color.FAIL + "Configfile not found." + color.ENDC)
				print("Create one with 'add'.")
			
			#SCAN
			scandirectory(path, f"{scanstorage}{scanfilename}", False)

			#LIST FILES TO FIND SCANFILES
			#SORT OUT ANY UNINTERESTING FILES
			scanfilelist = []
			#LIST DIRS
			dirs = os.listdir(scanstorage)
			number_of_files = 0
			
			#SEARCH FOR SCANFILES
			for singlefile in dirs:
				blockname = split(singlefile, ["_", "."])
				try:
					if blockname[4] == "btsscan" and blockname[2] == sys.argv[2]:
						number_of_files = number_of_files + 1
						scanfilelist.append(singlefile)
				except IndexError:
					pass
			
			#LIST CONFIG ENTREES
			serverlocation = ""
			sourcelocation = ""
				
			with open(configpath, "r") as conf:
				for line in conf:
					split_line = split(line, ",")
					if split_line[0] == sys.argv[2]:
						sourcelocation = split_line[2]
						serverlocation = split_line[4].rstrip() #Last entree has nline
					else:
						print(color.FAIL + f"No block {sys.argv[2]} found." + color.ENDC)

			if number_of_files == 1:
				print("One scan found. Complete backup of ALL data will be created.")
				print(color.OKBLUE + "Executing datatransfer." + color.ENDC)

				with open(f"{scanstorage}{scanfilename}", "r") as scan:
					lines = scan.readlines()
					for line in lines:
						path = split(line, ",")
						if len(path) == 1:
							os.makedirs(f"{serverlocation}{line.rstrip()}", exist_ok=True)
						elif len(path) == 3:
							path = path[0]
							path = path.replace(" ", "\ ")
							path = path.replace("(", "\(")
							path = path.replace(")", "\)")
							status = os.system(f"{copypath} {path} {serverlocation}{path}")
							exit_status = os.WEXITSTATUS(status)
							if exit_status != 0:
								print(color.FAIL + f"COPY ERROR: {exit_status}" + color.ENDC)
						else:
							print(color.FAIL + "Corrupt: " + line + color.ENDC)
				sys.exit()


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
			for singlefile in scanfilelist:
				temp = split(singlefile, "_")
				if int(temp[0]) == latest_timestamp:
					latest_scan_array_index = dircounter
				elif int(temp[0]) == previous_timestamp:
					previous_scan_array_index = dircounter
				dircounter = dircounter + 1

			print("Latest scan: " + scanfilelist[latest_scan_array_index])
			print("Previous scan: " + scanfilelist[previous_scan_array_index] + "\n")

			#COMPARE THE TWO FILES AGAINST EACH OTHER
			latest_scan_fd = open(f"{scanstorage}{scanfilelist[latest_scan_array_index]}", "r")
			previous_scan_fd = open(f"{scanstorage}{scanfilelist[previous_scan_array_index]}", "r")
			transmit_list = []

			latest_scan = latest_scan_fd.readlines()
			previous_scan = previous_scan_fd.readlines()

			file_same = 0
			file_new = 0
			file_total_old = 0
			file_total_latest = 0
			file_deleted = 0 #DELETED LINES COUNTER

			#REMOVE DELETED OR CHANGED FILES
			for oldline in previous_scan:
				if oldline not in latest_scan:
					checkifdir = split(oldline, ",")
					if len(checkifdir) == 1:
						#IF DIRECTORY, HASH WILL BE "directory".
						#THAT IS NEEDED DURING DIRECTORY REMOVAL
						transmit_list.extend([f"{oldline.rstrip()},directory,-\n"])
						print(color.FAIL + f"- {oldline}" + color.ENDC, end='')
					else:
						transmit_list.extend([f"{oldline.rstrip()},-\n"])
						print(color.FAIL + f"- {oldline}" + color.ENDC, end='')
						file_deleted = file_deleted + 1
				file_total_old = file_total_old + 1


			#FIND OUT CHANGED OR NEW FILES
			for line in latest_scan:
				if line in previous_scan:
					file_same = file_same + 1
				else:
					checkifdir = split(line, ",")
					if len(checkifdir) == 1:
						#IF DIRECTORY, HASH WILL BE "directory".
						#THAT IS NEEDED DURING DIRECTORY CREATION
						transmit_list.extend([f"{line.rstrip()},directory,+\n"])
					else:
						transmit_list.extend([f"{line.rstrip()},+\n"])
						file_new = file_new + 1
				file_total_latest = file_total_latest + 1


			#FILE STATS
			print(f"\nUnchanged files: {file_same}")
			print(f"New/Changed files: {file_new}")
			print(f"Deleted files: {file_deleted}")
			print(f"Total files in latest scan: {file_total_latest}")
			print(f"Total files in previous scan: {file_total_old}")

			#SAFETY GUARD: SEE ISSUE #8
			if safetyguard == True:
				if file_deleted >= file_total_old / 100 * minwarningvalue:
					print(f"SAFETY GUARD: MORE THAN {minwarningvalue}% DELETED")
					shouldcontinue()
				elif file_total_latest == 0:
					print("SAFETY GUARD: NO FILES FOUND.")
					shouldcontinue()
			else:
				pass

			#TRANSMITTER
			print(color.OKBLUE + "Executing datatransfer." + color.ENDC)
			for line in transmit_list:
				line = split(line.rstrip(), ",")
				if len(line) > 5:
					print(color.FAIL + f"Cannot backup file {line}." + color.ENDC)
					print("Path would brick BTSOOT.")
				else:
					if line[4] == "-":
						if line[2] == "directory":
							try:
								shutil.rmtree(f"{serverlocation}{line[0]}")
							except FileNotFoundError:
								pass
						else:
							try:
								os.remove(f"{serverlocation}{line[0]}")
							except FileNotFoundError:
								pass
					elif line[4] == "+":
						if line[2] == "directory":
							os.makedirs(f"{serverlocation}{line[0]}", exist_ok=True)
						else:
							path = line[0]
							path = path.replace(" ", "\ ")
							path = path.replace("(", "\(")
							path = path.replace(")", "\)")
							status = os.system(f"{copypath} {path} {serverlocation}{path}")
							exit_status = os.WEXITSTATUS(status)
							if exit_status != 0:
								print(color.FAIL + f"COPY ERROR: {exit_status}"+ color.ENDC)
					else:
						print(color.WARNING + "TRANSMIT CORRUPTION:" + color.ENDC)
						print(color.WARNING + line + color.ENDC)

			previous_scan_fd.close() 
			latest_scan_fd.close()
			print(color.OKGREEN + "Done." + color.ENDC)

		elif sys.argv[1] == "restore":
			print(color.FAIL + "WARNING! This will remove all files from source.")
			print("IF NO FILES ARE FOUND INSIDE THE BACKUP FOLDER, EVERYTHING IS LOST.")
			print("Abort using CTRL+C within 15 seconds." + color.ENDC)
			if not "--override" in sys.argv:
				time.sleep(15)
	
			serverlocation = ""
			sourcelocation = ""

			with open(configpath, "r") as conf:
				for line in conf:
					split_line = split(line, ",")
					if split_line[0] == sys.argv[2]:
						sourcelocation = split_line[2]
						serverlocation = split_line[4].rstrip()
			print(color.OKBLUE + "Deleting source." + color.ENDC)
			shutil.rmtree(sourcelocation)
			os.makedirs(sourcelocation)
			print(color.OKBLUE + "Executing datatransfer." + color.ENDC)
			print("This may take a long time.")

			#LIST FILES TO FIND SCANFILES
			#SORT OUT ANY UNINTERESTING FILES
			scanfilelist = []
			#LIST DIRS
			dirs = os.listdir(scanstorage)
			number_of_files = 0
			
			#SEARCH FOR SCANFILES
			for singlefile in dirs:
				blockname = split(singlefile, ["_", "."])
				try:
					if blockname[4] == "btsscan" and blockname[2] == sys.argv[2]:
						number_of_files = number_of_files + 1
						scanfilelist.append(singlefile)
				except IndexError:
					pass
			
			splitted_timestamp = []

			#FIND LATEST TWO FILES
			#SPLIT EVERY FILE NAME TO GAIN TIMESTAMP
			for scanfile in scanfilelist:
				temp = split(scanfile, "_")
				splitted_timestamp.append(int(temp[0]))

			#GETS LATEST SCANFILE'S TIMESTAMP
			latest_timestamp = max(splitted_timestamp)


			dircounter = 0
			latest_scan_array_index = -1
			previous_scan_array_index = -1
			for singlefile in scanfilelist:
				temp = split(singlefile, "_")
				if int(temp[0]) == latest_timestamp:
					latest_scan_array_index = dircounter
				dircounter = dircounter + 1

			print("Latest scan: " + scanfilelist[latest_scan_array_index])
			latest_scan_fd = open(f"{scanstorage}{scanfilelist[latest_scan_array_index]}", "r")

			for line in latest_scan_fd:
				split_line = split(line, ",")
				if len(split_line) == 1:
					path = split_line[0]
					path = path.rstrip()
					os.makedirs(path, exist_ok=True)
				elif len(split_line) == 3:
					path = split_line[0]
					path = path.replace(" ", "\ ")
					path = path.replace("(", "\(")
					path = path.replace(")", "\)")

					status = os.system(f"/etc/btsoot/copy {serverlocation}{path} {path}")
					exit_status = os.WEXITSTATUS(status)
					if exit_status != 0:
						print(color.FAIL + f"COPY ERROR: {exit_status}"+ color.ENDC)
			else:
				pass
			latest_scan_fd.close()
			print(color.OKGREEN + "Done." + color.ENDC)
			

		else:
			print(usage)


	except IndexError:
		print("INDEX ERROR")
		print(usage)
		sys.exit()


if __name__ == "__main__":
	try:
		main()
	except KeyboardInterrupt:
		print("\nQuitting.\n")
		sys.exit()
