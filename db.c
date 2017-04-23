/*
 * Copyright (C) Paul Kramme 2017
 * 
 * Part of BTSOOT
 * Single folder redundancy offsite-backup utility
 * 
 * Licensed under MIT License
 */

#include"db.h"

int db_init(char path[4096])
{
	/*DATABASE INIT*/
	sqlite3 *database;
	char *errormessage = 0;

	sqlite3_open(path, &database);
	sqlite3_exec(database, 
		"CREATE TABLE IF NOT EXISTS files(filename TEXT, path TEXT, type INT, hash BIGINT, size UNSIGNED BIGINT, level INT, scantime NUMERIC)", 
		0, 
		0, 
		&errormessage
	);
	if(errormessage != NULL)
	{
		fprintf(stderr, "%s\n", errormessage);
	}

	return 0;
}
