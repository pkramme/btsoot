#include"backup.h"

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	FILE *fp = fopen(fpath, "rb");
	crc_t checksum;
	char buffer[BUFSIZ];
	int total_read = 0;

	if(tflag == FTW_F)
	{
		checksum = crc_init();
		while((total_read = fread(buffer, BUFSIZ, 1, fp)) > 0)
		{
			checksum = crc_update(checksum, buffer, sizeof(buffer));
		}
		checksum = crc_finalize(checksum);
	}
	else
	{
		checksum = 0;
	}

	FILE *scanfile = fopen("test.scan", "a");
	fprintf(scanfile, "%-3s %2d %7jd %-40s %llx\n",
		(tflag == FTW_D) ?   "d"   : (tflag == FTW_DNR) ? "dnr" :
		(tflag == FTW_DP) ?  "dp"  : (tflag == FTW_F) ?   "f" :
		(tflag == FTW_NS) ?  "ns"  : (tflag == FTW_SL) ?  "sl" :
		(tflag == FTW_SLN) ? "sln" : "???",
		ftwbuf->level, (intmax_t) sb->st_size,
		fpath, (unsigned long long int) checksum
	);
	
	fclose(fp);
	fclose(scanfile);
	return 0;
}


int backup(job *job_import)
{
	/*DATABASE INIT*/
	char *error;
	error = db_init(job_import->block_name);
	printf("%s", error);

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
	return 0;
}

