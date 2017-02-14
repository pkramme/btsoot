#ifndef COPY_H_INCLUDED
#define COPY_H_INCLUDED

#include<stdio.h>
#include<fcntl.h>
#include<sys/stat.h>
#include<unistd.h>
#include<sys/sendfile.h>

int copy(char *source, char *destination);

#endif
