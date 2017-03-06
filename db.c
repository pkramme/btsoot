#include"db.h"

static int sqlite_callback(void *notused, int argc, char **argv, char **azcolumnname)
{
	int i;
	for(i = 0; i < argc; i++)
	{
		printf("%s = %s\n", azcolumnname[i], argv[i] ? argv[i] : "NULL");
	}
	printf("\n");
	return 0;
}


char db_init(char blockname[256])
{
	/*DATABASE INIT*/
	sqlite3 *database;
	char *errormessage = 0;
	int recall;

	strcat(blockname, ".dat");
	recall = sqlite3_open(blockname, &database);

	/*	TABLE CREATION*/
	sqlite3_exec(database, 
		"CREATE TABLE IF NOT EXISTS files(filename TEXT, path TEXT, type TEXT, crc TEXT, size INT, level INT)", 
		sqlite_callback, 
		0, 
		&errormessage
	);

	return &errormessage;
}

