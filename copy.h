#pragma once

#include<stdio.h>
#include<fcntl.h>
#include<sys/stat.h>
#include<unistd.h>
#include<sys/sendfile.h>

int copy(char *source, char *destination);

int copy_fallback(char *source, char *destination);

