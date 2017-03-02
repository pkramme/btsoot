all:
	clang-3.9 -Weverything -O3 -march=native btsoot.c backup.c copy.c -o btsoot
install:
	cp btsoot /usr/local/bin/btsoot
	mkdir -p /etc/btsoot/scans
update:
	cp btsoot /usr/local/bin/btsoot
remove:
	rm /usr/local/bin/btsoot
	rm -rf /etc/btsoot
clean:
	rm btsoot
