/*
COPY RETURN CODES
Return codes are incase of error displayed and handled by BTSOOT.

return -1 = incorrect number of args
return 0 = ok
return 1 = couldnt open source fd
return 2 = couldnt open dest fd
return 3 = couldnt write to dest fd - NOT EXISTENT IN COPY()
return 4 = couldnt read from source fd - NOT EXISTENT IN COPY()
return 5 = couldnt close source fd
return 6 = couldnt close dest fd
*/

#include"copy.h"

int copy(char *source, char *destination)
{
	int fd_source;
	int fd_destination;
	int dest_flags;
	mode_t permissions;
	
	fd_source = open(source, O_RDONLY);
	if(fd_source == -1)
	{
		return 1;
	}
	
	dest_flags = O_CREAT | O_WRONLY | O_TRUNC;
	
	permissions = S_IRUSR | S_IWUSR | S_IRGRP | S_IWGRP | S_IROTH | S_IWOTH;

	fd_destination = open(destination, dest_flags, permissions);
	if(fd_destination == -1)
	{
		return 2;
	}

	struct stat stat_source;
	fstat(fd_source, &stat_source);

	sendfile(fd_destination, fd_source, 0, (size_t)stat_source.st_size);
	
	/*close fds*/
	
	if(close(fd_source) == -1)
	{
		return 5;
	}
	if(close(fd_destination) == -1)
	{
		return 6;
	}
	
	return 0;
}

int copy_fallback(char *source, char *destination)
{
	int fd_source;
	int fd_destination;
	char buffer[BUFSIZ];
	int dest_flags;
	mode_t permissions;
	ssize_t read_check;

	/*compiler complains...*/
	read_check = 0;

	fd_source = open(source, O_RDONLY);
	if(fd_source == -1)
	{
		return 1;
	}
	
	dest_flags = O_CREAT | O_WRONLY | O_TRUNC;
	
	permissions = S_IRUSR | S_IWUSR | S_IRGRP | S_IWGRP | S_IROTH | S_IWOTH;

	fd_destination = open(destination, dest_flags, permissions);
	if(fd_destination == -1)
	{
		return 2;
	}
	
	while((read_check = read(fd_source, buffer, BUFSIZ)) > 0)
	{
		if(write(fd_destination, buffer, (size_t)read_check) != read_check)
		{
			return 3;
		}
	}
	if(read_check == -1)
	{
		return 4;
	}
	
	/*close fds*/
	if(close(fd_source) == -1)
	{
		return 5;
	}
	if(close(fd_destination) == -1)
	{
		return 6;
	}
	
	return 0;
}

