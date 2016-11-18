#!/usr/bin/python3
# THIS FILE IS AN EXPERIMENTAL PROGRAM TO LEARN ABOUT OS_WALK

import os, sys, datetime

#dt = datetime.datetime(1970,1,1).total_seconds()
#	print(dt)

walk_dir = sys.argv[1]

with open("fsscan.scan", "w") as f:

	print("SCANFROM" + walk_dir)

	for root, subdirs, files in os.walk(walk_dir):
		f.write(root + "\n")

	for filename in files:
		file_path = os.path.join(root, filename)
		f.write(file_path + filename + "\n")