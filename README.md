# BTSOOT [![build status](https://git.paukra.com/open-source/btsoot/badges/master/build.svg)](https://git.paukra.com/open-source/btsoot/commits/master)  

## What is BTSOOT
`tl;dr: A data selective backup/cloning tool that only manipulates changed data`  
First of: BTSOOT should only be used under very special circumstaces. If these requirements are
met, BTSOOT could be **your** backup/cloning solution. A list of this circumstaces:
- You want to create offsite backups
- You have much faster internal drive speed (in my case 400 - 600 MB/s) than your external connection (LAN/WAN/remote HDD etc...) and/or
- You have much unchanged data

### Practical Example
Lets say you have
- 2 TBs Movies
- 200 GBs Music
- 100 GBs Files, like PDFs, DOCXs, Mails, whatever...
- 50 GB Photos

I bet you do not change 90% of these files. You put movies one time on your server, 
and then don't touch them anymore. Same goes for Music and Photos. Your files however,
are daily used. You move a few gigabyte here, add a few megabytes there, delete something.
However, if you want an offsite backup, most backup software doesn't care. It will still copy any file,
every Movie, every single MP3, everything, even if only your "Files" folder changed, aside from 
maybe 3 added Music files.  
This software fixes that, as it will only copy the real changed data.
BTSOOT now scans all files in that directory and gives them a checksum. I now knows, when it scans them again,
if they have been changes, or if files have been added or removed. After that it copies the changed and new files to
the remote directory, and deletes the non existing ones. This happens every scan + backup phase. If you want to do
on a regular basis, and don't want to do this manualy, it might be a good idea to add the commands 
`btsoot scan <blockname>` and `btsoot backup <blockname>` to your cron list.
If, for some reason, no file is found at a new scan, the backup gets aborted to prevent corrupting your functioning
backup. Same goes for Specific deletion levels, which can be configured.

---

## Current project status
BTSOOT is still under heavy development, thus it is only in 0.x release. The project has left its prototyping implementation
in Python, and will be rewritten in C with the currently anticipated features, that are (partly) available in 0.6-stable.
The 0.6-stable release will be maintained with bugfixes, but no new features will be added. 0.7-stable release will be the
first C implementation of BTSOOT and will hopefully be faster and better executed than the prototype.  
Altough it's a prototype, 0.6-stable is totally usable for the public, and I am encouraging you to use it, as it is a good
product.

---

## Dependencies
- 64bit Linux OS (no specific distribution)  

If you want to build from source:  
- clang
- make

---

## Usage
### Installation
1. Clone the repository to the folder where you want it
2. Check out your desired release branch e.g. `0.4-stable`, so `git checkout 0.4-stable`
2. Execute `make` and `sudo make install`
4. Done

BTSOOT is now in `/usr/local/bin/` and can be used. It needs root permissions to run. The config or all scans are in `/etc/btsoot/`. If you want to 
uninstall it, run `sudo make uninstall`. This will delete all scans and your config, too.

NOTE:  
I am currently engaged in a project in which a [cURL package manager](https://github.com/eddyx9/install.paukra.com/) is developed. I'm hopeing it can be used to 
distribute this project in the future. There are `.deb` packages available, but it is recommended to build the project yourself, as it will give you the best 
possible performance.

### Create a block
`btsoot add <block-name> <source-path> <path-to-remote-dir>`  
This is written to a file named `btsoot.conf` which is created inside the directory where BTSOOT lies.

### Backup a block
`btsoot backup <block-name>`  
The program will search for the latest changes, and compare them with the previous, which it then copyies to their 
paths on the remote location.  
Prior to 0.5 there were two commands `scan` and `backup`. This wasn't practical, 
as both were used together while one on its own was unecessary, so they got removed.
**This also means that you MUST NOT change the remote files per hand. BTSOOT will not know about any changed file that 
it didn't changed itself.** BTSOOT does not check the remote files for integrity. It is your responsibility that they 
are not broken

### Restore a block
`btsoot restore <block-name>`  
The program will delete the entire source folder, and then reload the files from the latest scan. If no files are
found, everything is lost. Fortunately there is no scanning involved, so its just copying, and this is fast.  
`--override` will override the safety time, so the restoring will begin immediately. Use with caution.

---

## Backstory
### What i have
I do not have
- Much time
- Much money

What i do have
- nearly 3 TB's of data
- 2 Routers with NAS features
- 2 Raspberry Pi's
- Fast RAIDz1 on my primary NAS (450 - 600 MB/s)
- Gigabits LAN
- Very important files, like backups and other mission critical files, that are irreplaceable
- Sysadmin and programming skills

Now i have a problem. I don't have much money, so buying new hardware for an offsite 
(which is defined here as not in the same room or the same server) NAS is not an option.
So i have to work with what i have. Let's see:
#### Raspberry Pi
Raspberry Pi's are not suited for NAS's. You can read this anywhere on the internet.
They are slow, because they lack of a dedicated Ethernet Interface, everything is runned over
one single USB 2.0 Host. This is slow. On the plus side, they are small, silent, have low energy
usage, and are cheap - 35 Euros, here in Germany.
#### FritzBox Router
I have two of them, one is my primary gateway, one my Wifi AP, both have rudimentary NAS features like SMB and FTP.
Also, one of them has USB3.0 which theoreticly should be fast, when it is combined with my 3TB HDD from WD, which also has 
a USB 3.0 interface. It even has dedicated Gigabits LAN (not so surprising, as it is a router) but for a NAS, it's great.  
### Testing
Well. The Raspberry Pi was slow. The very definition of slow. 7 MB/s at best with ext4 for the HDD, NFSv4 and NFSv3
and all speed fixes i know. My other Raspberry Pi, a second generation device, brought it up to 15MB/s max. So i turned to the FritzBox.
And was disappointed again: 10MB/s without even the possibility of speed improvements, because hacking a Router OS is beyond what one should do
with absolute mission critical hardware. If the FritzBox get broken, i have to spend money to fix it, which i do not have, and my family members and me
are out of internet until a replacement device arrives. So, no. I have to work with this. 
### Planning and Development
At this point i decided that there was not much i could do about the performance.
I decided to use the FritzBox, because i only could use the slower of them, the other one were redisigned to a PiHole DNS server, and the speed 
where more unstable that the router's. The initial datatransfer wouldn't be a problem, as i could do it on every device and then just 
plug it into the FritzBox after that. So i had to find a way to speed up the transfer, without increasing connection speed, or tuning protocol.
I tought about the thing every sysadmin would have though about: rsync. I would just mount the SMB on my NAS, use rsync to transfer any changed parts of the file
and would therefore only have to deal with a minimum filesize, as most of the files are media files, which doesn't change that often. There were a problem though:
What happens when a file gets renamed, or a folder moved, or deleted? Rsync doesn't cover that. I a file is deleted, rsync doesn't remove it, it will just remain.
Same with renamed directories.  
In this moment i changed from sysadmin to developer:  

**I had to create a program which identifies a file, and mirrors a tree over the network, without having to send already existing files.**  

I thought about a project which i already planned before: A program where you can create certain blocks, and monitor them.
Early drafts of this are found in [this](https://git.paukra.com/open-source/backup) repository, where i was eager to write it in system dependent languages, like C
for Linux and C# for Windows. I wasn't experienced enough to write it, and altough i could have learned it, it would need time, which i don't had.
So i created a new project, this project, called BTSOOT, and began to write it in Python.
I wasn't good in Python, but it is way faster to learn than C with all Linux system calls. And here i am. The program as it is runs on a Linux host under Python3.6 (formatted srings
were to nice to ignore them) and copies changed files to a mounted network folder on the host.

---

## Performance
As BTSOOT is currently, as of 52a445fa, single threaded, the performance is not as good as it could be. I am currently 
trying to eliminate bottlenecks, like the random writing to disk, to improve performance with small files. This, 
however, drasticly raises RAM usage. The worst case is <number-of-files> * (256 + 8 + 2) * 2 (256 path lengh; 8 CRC  
Hash; 2 Escape and Comma; 2 scanfiles). Limiting factors include, but are not limited to:
- Slow disk speed
- Many little files, which slow down the CRC algorithm and arrays which hold the data

The table below might get you an impression of speed.

|  COMMIT  | Data   |  Time             |  
| -------- | :----- | :---------------- |  
| 52a445fa | 1,9TB  |  187 Min, 22 Secs |  

I will rewrite most of the python performance parts in C, as i did with the copy code. This alone gave promising
results. For better or worse, there will be a C backend, python as a frontend, which only launches the backend pipeline.
The key is to reduce as much disk writing as possible.

---

## Project Information
I am currently enforcing GitLab Flow, with [semantic versioning](http://semver.com) in release channels.

## Dependencies
For the sake of portability, i try to statically link anything possible. [A list of additional dependencies and their licenses](https://git.paukra.com/open-source/btsoot/wikis/additional-dependencies).