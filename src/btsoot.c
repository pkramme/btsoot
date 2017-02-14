/*
 * BTSOOT by Paul Kramme
 * Single folder redundancy offsite-backup utility
 * 
 * Licensed under MIT License
 */

#include<stdio.h>
#include<string.h>

#include"copy.h"

int main(int argc, char *argv[])
{
	/* 
	 * Argument resolving code
	 */

	if(argc < 2)
	{
		puts("USAGE");
	}

	if(strcmp(argv[1], "add") == 0 && argc >= 5)
	{
		printf("Adding %s with src=%s and dest=%s\n", argv[2], argv[3], argv[4]);
		/*TODO: Add code for config creation, adding*/
	}
	else if(strcmp(argv[1], "rm") == 0 && argc >= 3)
	{
		printf("Removing %s\n", argv[2]);
		/*TODO: Add code for config deletion*/
	}
	else if(strcmp(argv[1], "list") == 0 && argc >= 3)
	{
		printf("Listing %s\n", argv[2]);
		/*TODO: Add code for config listing*/
	}
	else if(strcmp(argv[1], "backup") == 0 && argc >= 3)
	{
		printf("Backing up %s\n", argv[2]);
		/*TODO: Add code for starting backup pipeline routine*/
	}
	else if(strcmp(argv[1], "restore") == 0 && argc >= 3)
	{
		printf("Restoring %s\n", argv[2]);
		/*TODO: Add code for starting restore pipeline routine*/
	}
	else
	{
		puts("USAGE");
	}
	return 0;
}

