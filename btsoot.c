/*
 * BTSOOT by Paul Kramme
 * Single folder redundancy offsite-backup utility
 * 
 * Licensed under MIT License
 */

#include"btsoot.h"

int test_last_char(const char *string);

int test_last_char(const char *string)
{
	return (string && *string && string[strlen(string) - 1] == '/') ? 0 : 1;
}

int main(int argc, char *argv[])
{
	job j;
	job *jptr = &j;

	/* Argument resolving code */
	if(argc < 2)
	{
		puts("USAGE");
	}

	if(argc >= 5)
	{
		if(strcmp(argv[1], "add") == 0)
		{
			FILE *config = fopen(CONFIG_PATH, "a+");
			if(	test_last_char((const char *)argv[2]) == 0 || 
				test_last_char((const char *)argv[3]) == 0)
			{
				fprintf(stderr, "Please remove suffixed slash from paths!");
				fclose(config);
				return 1;
			}
			fprintf(config, "%s,%s,%s\n", argv[2], argv[3], argv[4]);
			fclose(config);
		}
	}
	else if(argc >= 3)
	{
		if(strcmp(argv[1], "rm") == 0)
		{
			FILE *config = fopen(CONFIG_PATH, "r");
			FILE *copyconfig = fopen(COPY_CONFIG_PATH, "w");
			char buffer[8448];
			while(fgets(buffer, sizeof(buffer), config) != NULL)
			{
				if(strstr(buffer, argv[2]))
				{
					/*DO NOTHING*/
				}
				else
				{
					fprintf(copyconfig, "%s", buffer);
				}
			}
			fclose(config);
			fclose(copyconfig);
			rename(COPY_CONFIG_PATH, CONFIG_PATH);
		}
		else if(strcmp(argv[1], "list") == 0)
		{
			printf("Listing %s\n", argv[2]);
			/*TODO: Add code for config listing blocks*/
		}
		else if(strcmp(argv[1], "backup") == 0)
		{
			j.pip_purp_id = PIP_PURP_ID_BACKUP;
			strcpy(j.block_name, argv[2]);

			FILE *config = fopen(CONFIG_PATH, "r");
			char buffer[8448];
			while(fgets(buffer, sizeof(buffer), config) != NULL)
			{
				if(strstr(buffer, argv[2]))
				{
					strcpy(j.block_name, strtok(buffer, ","));
					strcpy(j.src_path, strtok(NULL, ","));
					strcpy(j.dest_path, strtok(NULL, ","));

					strcpy(j.dest_path, strtok(j.dest_path, (char *) "\n"));

					backup(jptr);
				}
			}
			fclose(config);
		}
		else if(strcmp(argv[1], "restore") == 0)
		{
			j.pip_purp_id = PIP_PURP_ID_RESTORE;
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

