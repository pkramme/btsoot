/*
 * BTSOOT by Paul Kramme
 * Single folder redundancy offsite-backup utility
 * 
 * Licensed under MIT License
 */

#include<stdio.h>
#include<stdlib.h>
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
	
	if(argc >= 5)
	{
		if(strcmp(argv[1], "add") == 0)
		{
			printf("Adding %s with src=%s and dest=%s\n", argv[2], argv[3], argv[4]);
			/*TODO: Add code for config creation, adding*/
		}
	}
	else if(argc >= 3)
	{
		if(strcmp(argv[1], "rm") == 0)
		{
			printf("Removing %s\n", argv[2]);
			/*TODO: Add code for config deletion*/
		}
		else if(strcmp(argv[1], "list") == 0)
		{
			printf("Listing %s\n", argv[2]);
			/*TODO: Add code for config listing*/
		}
		else if(strcmp(argv[1], "backup") == 0)
		{
			printf("Backing up %s\n", argv[2]);
			/*TODO: Add code for starting backup pipeline routine*/
		}
		else if(strcmp(argv[1], "restore") == 0)
		{
			printf("Restoring %s\n", argv[2]);
			/*TODO: Add code for starting restore pipeline routine*/
		}
		else
		{
			puts("USAGE");
		}
	}
	else
	{
		puts("Not enough args given");
	}
	return 0;
}

