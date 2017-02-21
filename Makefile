all:
	gcc -Wall -Wextra -O3 -march=native btsoot.c sqlite3.c crc.c backup.c copy.c -o btsoot -pthread -ldl
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
