btsoot:	sqlite3.o btsoot.o xxhash.o db.o backup.o copy.o
	clang-5.0 -Wall -Wextra -O2 -march=native -frename-registers -std=c11 btsoot.o backup.o db.o copy.o xxhash.o sqlite3.o -o btsoot -pthread -ldl
install:
	cp btsoot /usr/local/bin/btsoot
	mkdir -p /etc/btsoot/scans
update:
	cp btsoot /usr/local/bin/btsoot
remove:
	rm /usr/local/bin/btsoot
	rm -rf /etc/btsoot
clean:
	rm btsoot *.o
