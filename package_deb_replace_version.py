import sys

def split(string, splitters):
	final = [string]
	for x in splitters:
	for i,s in enumerate(final):
			if x in s and x != s:
				left, right = s.split(x, 1)
				final[i] = left
				final.insert(i + 1, x)
				final.insert(i + 2, right)
	return final

path = "btsoot/DEBIAN/control"
fullversion = sys.argv[1]
version = split(fullversion, "v")
version = version[1]
control_content = f"""Package: btsoot
Version: {version}
Section: base
Priority: optional
Architecture: i386
Depends: build-essential
Maintainer: Paul Kramme <pjkramme@gmail.com>
Description: BTSOOT
 Folder redundancy offsite-backup utility.
"""
print("DEB PACKAGE VERSION REPLACER")
# yes, i wrote a tool for this...
with open(path, "a") as f:
	f.write(control_content)
print("Done.")
