#!/usr/bin/python3
# THIS FILE IS AN EXPERIMENTAL PROGRAM TO LEARN ABOUT OS_WALK

import os, sys

walk_dir = sys.argv[1]

print("walk directory: " + walk_dir)

print("Walk directory (absolute) = " + os.path.abspath(walk_dir))
print("\n\n\n\n\n\n\n\n\n")

for root, subdirs, files in os.walk(walk_dir):
	print(root)
	#list_file_path = os.path.join(root, "dirlist.txt")
	#print("dirlist = ", list_file_path)

	#for subdir in subdirs:
	#	print(subdir)

	for filename in files:
		file_path = os.path.join(root, filename)
		print(file_path + filename)