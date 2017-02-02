#include<stdio.h>
#include<fcntl.h>
#include<sys/stat.h>
#include<unistd.h>

#define BUFFERSIZE 4096

int copy(const char *source, const char *destination);

int main(int argc, char *argv[])
{
	int error;
	
	/*
	return -1 = incorrect number of args
	return 0 = ok
	return 1 = couldnt open source fd
	return 2 = couldnt open dest fd
	return 3 = couldnt write to dest fd
	return 4 = couldnt read from source fd
	return 5 = couldnt close source fd
	return 6 = couldnt close dest fd
	*/

	if(argc != 3)
	{
		return -1;
	}

	error = copy(argv[1], argv[2]);
	return error;
}

int copy(const char *source, const char *destination)
{
	int fd_source;
	int fd_destination;
	char buffer[BUFFERSIZE];
	int dest_flags;
	mode_t permissions;
	ssize_t read_check;

	read_check = 0;

	fd_source = open(source, O_RDONLY);
	if(fd_source == -1)
	{
		return 1;
	}
	
	dest_flags = O_CREAT | O_WRONLY | O_TRUNC;
	permissions = S_IRUSR | S_IRGRP | S_IROTH | S_IWUSR | S_IWGRP | S_IWOTH;
	
	fd_destination = open(destination, dest_flags, permissions);
	if(fd_destination == -1)
	{
		return 2;
	}
	
	while((read_check == read(fd_source, buffer, BUFFERSIZE)) > 0)
	{
		if(write(fd_destination, buffer, read_check) != read_check)
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

