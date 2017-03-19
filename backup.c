#include"backup.h"

typedef struct file_info {
	char name[256];
	char path[4096];
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

node_t *files_head = NULL;
node_t *current_node = NULL;

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
		free(current);
		current = NULL;
	}
}

static void print_list(node_t *head)
{
	node_t *current = head;
	while(current != NULL)
	{
		printf("%s\n%s\n%i\n%li\n%li\n%"PRIu64"\n\n", 
			current->link.path, 
			current->link.name, 
			current->link.type, 
			current->link.size, 
			current->link.scantime,
			current->link.checksum
		);
		current = current->next;
	}
}

static uint64_t hash(char path[4096], size_t size)
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
	uint64_t h64;

	int8_t buffer[size];
	XXH64_state_t state64;	
	size_t total_read = 1;
		
	XXH64_reset(&state64, 0);
	while(total_read)
	{
		total_read = fread(buffer, 1, size, fp);	
		XXH64_update(&state64, buffer, size);
	}
	h64 = XXH64_digest(&state64);

	fclose(fp);	
	return h64;
}

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	file_t current_file = {0};
	strcpy(current_file.path, fpath);
	strcpy(current_file.name, fpath + ftwbuf->base);
	current_file.size = sb->st_size;
	current_file.type = tflag;
	current_file.scantime = t0;

	push(files_head, current_file);
	return 0;
}

int backup(job_t *job_import)
{
	t0 = time(0);

	files_head = malloc(sizeof(node_t));
	current_node = files_head;

	//Execute filewalker
	if(nftw(job_import->src_path, filewalk_info_callback, 20, 0) == -1)
	{
		fprintf(stderr, "ERROR NFTW\n");
		exit(EXIT_FAILURE);
	}

	//print_list(files_head);
	delete(files_head);
	return 0;
}

