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

#define PIP_PURP_ID_BACKUP 1
#define PIP_PURP_ID_RESTORE 2

#define CONFIG_PATH "btsoot.conf"

int test_last_char(const char *string)
{
	return (string && *string && string[strlen(string) - 1] == '/') ? 0 : 1;
}

int main(int argc, char *argv[])
{
	/* Config filestreams */
	FILE *config;
	config = fopen(CONFIG_PATH, "a+");

	/* Argument resolving code */
	if(argc < 2)
	{
		puts("USAGE");
	}
	
	struct job {
		char block_name[256];
		int pip_purp_id;
		char src_path[256];
		char dest_path[256];
	} *job;

	if(argc >= 5)
	{
		if(strcmp(argv[1], "add") == 0)
		{
			if(	test_last_char((const char *)argv[2]) == 0 || 
				test_last_char((const char *)argv[3]) == 0)
			{
				puts("Please remove suffixed slash from paths!");
				return 1;
			}
			fprintf(config, "%s,%s,%s\n", argv[2], argv[3], argv[4]);
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

		/*
		 * HERE COME THE REAL PIPELINES...
		 */

		else if(strcmp(argv[1], "backup") == 0)
		{
			job->pip_purp_id = PIP_PURP_ID_BACKUP;
			
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
	fclose(config);
	return 0;
}

