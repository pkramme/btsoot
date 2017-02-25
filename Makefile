all:
	gcc -Wall -Wextra -O3 btsoot.c crc.c backup.c db.c copy.c sqlite3.c -o btsoot -pthread -ldl
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
