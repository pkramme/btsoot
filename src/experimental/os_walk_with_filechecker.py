#!/usr/bin/env python3
# THIS FILE IS AN EXPERIMENTAL PROGRAM TO LEARN ABOUT OS_WALK

import os, sys, datetime
import compare

#dt = datetime.datetime(1970,1,1).total_seconds()
#	print(dt)

walk_dir = sys.argv[1]
scanfile = sys.argv[2]

with open(scanfile, "w") as f:
	f.write("path,checksum\n")

	print("SCANFROM " + walk_dir)

	for root, subdirs, files in os.walk(walk_dir):
		f.write(root + "\n")

		for filename in files:
			file_path = os.path.join(root, filename)
			checksum = compare.sha1sum(file_path)
			print(checksum)
			f.write(file_path + "," + checksum + "\n")