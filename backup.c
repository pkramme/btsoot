#include"backup.h"

static sqlite3 *database = NULL;

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	FILE *fp = fopen(fpath, "rb");
	FILE *scanfile = fopen("test.scan", "a");

	crc_t checksum = crc_init();
	char buffer[BUFSIZ];
	int total_read = 0;

	if(tflag == FTW_F)
	{
		while((total_read = fread(buffer, sizeof(buffer), 1, fp)) > 0)
		{
			checksum = crc_update(checksum, (void *)buffer, strlen(buffer));
		}
		checksum = crc_finalize(checksum);

		printf("0x%llx %s\n", (unsigned long long int) checksum, fpath);
		fprintf(scanfile, "%-3s %2d %7jd %-40s 0x%llx\n",
			(tflag == FTW_D) ?   "d"   : (tflag == FTW_DNR) ? "dnr" :
			(tflag == FTW_DP) ?  "dp"  : (tflag == FTW_F) ?   "f" :
			(tflag == FTW_NS) ?  "ns"  : (tflag == FTW_SL) ?  "sl" :
			(tflag == FTW_SLN) ? "sln" : "???",
			ftwbuf->level, (intmax_t) sb->st_size,
			fpath, (unsigned long long int) checksum
		);
	}
	else
	{
		checksum = 0;
	}


	fclose(fp);
	fclose(scanfile);
	return 0;
}


int backup(job *job_import)
{
	/*DATABASE CREATE*/
	int *error;
	error = (int) db_init(job_import->block_name);
	if(error != NULL)
	{
		sqlite3_free(error);
		printf("%s\n", error);
	}

	/*CURRENT DATABASE INIT*/
		/*USE CLEAR FROM CREATE*/
	sqlite3_open(job_import->block_name, &database);

	/*FILEWALKER*/
	printf("%s\n", job_import->src_path);
	if(nftw(job_import->src_path, filewalk_info_callback, 20, 0) == -1)
	{
		fprintf(stderr, "ERROR NFTW\n");
		exit(EXIT_FAILURE);
	}

	/*CRC CHECK*/

	/*EXECUTOR*/




	/**
	 * BACKUP PIPELINE
	 * 
	 * Needed functions:
	 *  - scan for files and directories, record size of files
	 *  - scan for crc values
	 *  - diff this scan with the last
	 *  - execute all necessary changes
	 */
	sqlite3_close(database);
	return 0;
}

