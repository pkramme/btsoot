#include"backup.h"

static int display_info(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
    printf("%2d %7jd %-40s %d %s\n", ftwbuf->level, (intmax_t) sb->st_size, fpath, ftwbuf->base, fpath + ftwbuf->base);
    return 0;           /* To tell nftw() to continue */
}

int backup(job *job_import)
{
	if(nftw(job_import->src_path, display_info, 20, 0) == -1)
	{
		printf("ERROR NFTW\n");
		exit(EXIT_FAILURE);
	}
	/* BACKUP PIPELINE
	 * 
	 * Needed functions:
	 *  - scan for files and directories, record size of files
	 *  - scan for crc values
	 *  - diff this scan with the last
	 *  - execute all necessary changes
	 */
	return 0;
}
