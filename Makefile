btsoot:	sqlite3.o btsoot.o xxhash.o db.o backup.o copy.o
	gcc -Wall -Wextra -O3 btsoot.o backup.o db.o copy.o xxhash.o sqlite3.o -o btsoot -pthread -ldl
sql:
	gcc -c sqlite3.c
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
