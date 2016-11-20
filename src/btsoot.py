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


import sys


class color:
	HEADER = '\033[95m'
	OKBLUE = '\033[94m'
	OKGREEN = '\033[92m'
	WARNING = '\033[93m'
	FAIL = '\033[91m'
	ENDC = '\033[0m'
	BOLD = '\033[1m'
	UNDERLINE = '\033[4m'


def scandirectory():
	with open(scanfile, "w") as f:
		f.write("path,checksum\n")

		#print("SCANFROM " + walk_dir)

		for root, subdirs, files in os.walk(walk_dir):
			f.write(root + "\n")
			for filename in files:
				file_path = os.path.join(root, filename)
				checksum = compare.sha1sum(file_path)
				print(checksum)
				f.write(file_path + "," + checksum + "\n")


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
		if sys.argv[1] == "rm":
			try:
				name = sys.argv[2]
			except IndexError:
				print("Usage: " + sys.argv[0] + "rm name")
				exit()
			with open("btsoot.conf", "rw") as f:
				row = 0
				beginning_row = 0
				lines = f.readlines()
				for line in lines:
					row = row + 1
					if line == "name=" + name:
						beginning_row = row
					elif row == beginning_row + 1 or row == beginning_row + 2:
						pass
					else:
						f.write(line)


	except IndexError:
		print("Usage.")
		exit()


if __name__ == "__main__":
	try:
		main()
	except KeyboardInterrupt:
		print("Stopping program.")
		sys.exit()
