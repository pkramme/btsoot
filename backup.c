/*
 * Copyright (C) Paul Kramme 2017
 * 
 * Part of BTSOOT
 * Single folder redundancy offsite-backup utility
 * 
 * Licensed under MIT License
 */

#include"backup.h"

typedef struct file_info {
	char name[256];
	const char *path;
	uint64_t checksum;
	int8_t type;
	int8_t level;
	size_t size;
	time_t scantime;
	int8_t thread_number;
} file_t;

typedef struct node {
	file_t link;
	struct node *next;
} node_t;

static time_t t0;
static time_t tsearched;

static node_t *files_head = NULL;
static node_t *current_node = NULL;

static node_t *old_files_head = NULL;
//static node_t *old_current_node = NULL;

static size_t total_size = 0;
static size_t max_allowed_thread_size = 0;

static void push(node_t *head, file_t filestruct)
{
	node_t *current = head;
	while(current->next != NULL)
	{
		current = current->next;
	}
	current->next = malloc(sizeof(node_t));
	current->next->link = filestruct;
	current->next->next = NULL;
}

static void delete(node_t *head)
{
	node_t *current;
	while((current = head) != NULL)
	{
		head = head->next;
		free((void*)current->link.path);
		free(current);
		current = NULL;
	}
}

static void print_list(node_t *head)
{
	node_t *current = head;
	while(current != NULL)
	{
		printf("path=%s\nname=%s\ntype=%i\nsize=%li\ntime=%li\nchecksum=%"PRIu64"\nthrnmb=%i\n\n", 
			current->link.path, 
			current->link.name, 
			current->link.type, 
			current->link.size, 
			current->link.scantime,
			current->link.checksum,
			current->link.thread_number
		);
		if(current->next != NULL)
		{
			current = current->next;
		}
		else
		{
			break;
		}
	}
}

static uint64_t hash(const char path[4096], size_t size)
{
	if(size > FILEBUFFER)
	{
		size = FILEBUFFER;
	}

	FILE *fp = fopen(path, "rb");
	if(fp == NULL)
	{
		return 0;
	}

	int8_t buffer[size];
	XXH64_state_t state64;	
	size_t total_read = 1;
		
	XXH64_reset(&state64, 0);
	while(total_read)
	{
		total_read = fread(buffer, 1, size, fp);	
		XXH64_update(&state64, buffer, size);
	}

	fclose(fp);	
	return XXH64_digest(&state64);
}

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	static int8_t thread_number = 0;

	file_t current_file = {0};
	current_file.path = strdup(fpath);
	strcpy(current_file.name, fpath + ftwbuf->base);
	current_file.size = sb->st_size;
	current_file.type = tflag;
	current_file.scantime = t0;
	current_file.thread_number = thread_number;

	static size_t current_thread_size;
	current_thread_size += sb->st_size;
	if(current_thread_size > max_allowed_thread_size)
	{
		++thread_number;
		current_thread_size = 0;
	}

	current_node->link = current_file;
	current_node->next = malloc(sizeof(node_t));
	current_node = current_node->next;
	current_node->next = NULL;

	return 0;
}

static int filewalk_size_callback(const char *fpath, const struct stat *sb, int typeflag)
{
	total_size += sb->st_size;
	return 0;
}

void *thread_hash(void* t)
{
	node_t *current = files_head;
	while(current != NULL)
	{
		if(current->link.thread_number == (uintptr_t)t)
		{
			current->link.checksum = hash(current->link.path, current->link.size);
		}
		if(current->next != NULL)
		{
			current = current->next;
		}
		else
		{
			break;
		}
	}

	pthread_exit((void*) t);
}

static int get_old_scantime_callback(void *notused, int argc, char **argv, char **azcolname)
{
	for(int i = 0; i < argc; i++)
	{
		if(strcmp(azcolname[i], "scantime") == 0)
		{
			tsearched = argv[i] ? atoi(argv[i]) : -1;
		}
	}
	return 1;
}

static int read_old_from_database(node_t *head, sqlite3 *database)
{
	node_t *current = head;
	sqlite3_exec(database, "SELECT scantime FROM files ORDER BY scantime DESC", get_old_scantime_callback, NULL, NULL);
	if(tsearched == -1)
	{
		puts("Make first scan...");
		return 1;
	}

	char *zsql = sqlite3_mprintf("SELECT * FROM files WHERE scantime = %li", tsearched);
	sqlite3_exec(database, zsql, NULL, NULL, NULL);
	
	//-> make current global so callback can easiely access it?
	//Read latest timestamp
	//Read all data with given timestamp into linked list starting at "head"
	sqlite3_free(zsql);
	return 0;
}

static int write_to_server(void)
{
	return 0;
}

static int write_to_db(node_t *head, sqlite3 *database)
{
	node_t *current = head;
	char *sqlerrormessage = 0;

	while(current != NULL)
	{
		char zsql[8192];
		sqlite3_snprintf(sizeof(zsql), zsql, "INSERT INTO files (filename, path, type, size, level, scantime, hash) VALUES('%q', '%q', %i, %lli, %i, %i, %lli)", 
			current->link.name, 
			current->link.path, 
			current->link.type, 
			current->link.size, 
			current->link.level, 
			current->link.scantime, 
			current->link.checksum);
		if(sqlite3_exec(database, zsql, NULL, NULL, &sqlerrormessage))
		{
			printf("ERROR: %s\n", sqlerrormessage);
			sqlite3_free(sqlerrormessage);
		}

		if(current->next != NULL)
		{
			current = current->next;
		}
		else
		{
			break;
		}
		memset(zsql, '\0', sizeof(zsql));
	}
	return 0;
}

int backup(job_t *job_import)
{
	t0 = time(0);

	if(ftw(job_import->src_path, &filewalk_size_callback, 1))
	{
		printf("SIZE CALC ERROR\n");
		exit(EXIT_FAILURE);
	}

	job_import->max_threads = 4;
	max_allowed_thread_size = total_size / job_import->max_threads;

	files_head = malloc(sizeof(node_t));
	current_node = files_head;

	if(nftw(job_import->src_path, filewalk_info_callback, 20, 0) == -1)
	{
		fprintf(stderr, "ERROR NFTW\n");
		exit(EXIT_FAILURE);
	}

	pthread_attr_t attr;
	pthread_attr_init(&attr);
	pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_JOINABLE);
	void *status;

	pthread_t threads[job_import->max_threads];
	int rc;
	for(long t = 0; t < job_import->max_threads; t++)
	{
		rc = pthread_create(&threads[t], NULL, thread_hash, (void*) t);
		if(rc)
		{
			printf("ERROR FROM PTHREAD\nrc is %i\n", rc);
			exit(EXIT_FAILURE);
		}
	}
	
	for(long t = 0; t < job_import->max_threads; t++)
	{
		rc = pthread_join(threads[t], &status);
		if(rc)
		{
			puts("PTHREAD ERROR");
			exit(EXIT_FAILURE);
		}
	}


	// Open and init database
	sqlite3 *database;
	db_init(job_import->db_path);
	sqlite3_open(job_import->db_path, &database);

	// Read old from database
	read_old_from_database(old_files_head, database);

	//Write to database
	//print_list(files_head);
	write_to_db(files_head, database);

	delete(files_head);
	sqlite3_close(database);
	pthread_exit(NULL);
	return 0;
}
