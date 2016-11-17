#!/usr/bin/python3
# THIS FILE IS AN EXPERIMENTAL PROGRAM TO LEARN ABOUT OS_WALK

import os, sys

walk_dir = sys.argv[1]

print("walk directory: " + walk_dir)

print("Walk directory (absolute) = " + os.path.abspath(walk_dir))

for root, subdirs, files in os.walk(walk_dir):
	print("root = " + root)
	list_file_path = os.path.join(root, "dirlist.txt")
	print("dirlist = ", list_file_path)

	#with open(list_file_path, 'wb') as list_file:
	for subdir in subdirs:
		print("\tsubdirectorys = ", subdir)

	for filename in files:
		file_path = os.path.join(root, filename)
		print("\t\tfile %s in full path = %s" % (filename, file_path))