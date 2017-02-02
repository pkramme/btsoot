#include<stdio.h>

int fcopy(FILE *src, FILE *dest);

int main(int argc, char *argv[])
{
	if(argc != 3)
	{
		return 1;
	}
	fcopy(argv[1], argv[2]);
	return 0;
}

int fcopy(FILE *src, FILE *dest)
{
	char buffer[BUFSIZ];
	size_t read_check;

	while((read_check = fread(buffer, sizeof(char), sizeof(buffer), src)) > 0)
	{
		if(fwrite(buffer, sizeof(char), read_check, dest) != read_check)
		{
			return -1;
		}
	}
	return 0;
}

