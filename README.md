# BTSOOT 
## Build Status
Master:     [![build status](https://git.paukra.com/open-source/btsoot/badges/master/build.svg)](https://git.paukra.com/open-source/btsoot/commits/master)  
Production: [![build status](https://git.paukra.com/open-source/btsoot/badges/production/build.svg)](https://git.paukra.com/open-source/btsoot/commits/production)
## How to get it
i386/amd64: [btsoot_0.3.0-i386.deb](https://drive.google.com/open?id=0B2sQy9J1YjgJbHBEcm5DbjFKWjA)  
**Until i spinn up my own aptly server, the distribution will unfortunately happen over google drive.  
I apologize.**

---

## What is this
`tl;dr: A backup/cloning tool`  
First of: BTSOOT should only be used under very special circumstaces. If these requirements are
met, BTSOOT could be **your** backup/cloning solution. A list of this circumstaces:
- You want to create offsite backups
- You have much faster internal drive speed (in my case 400 - 600 MB/s) than your LAN/WAN
- You have a slow connection to the remote device (far below your diskspeed) and/or
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
However, if you want an offsite backup, your software doesn't care. It will still copy any file,
every Movie, every single MP3, everything, even if only your "Files" folder changed, aside from 
maybe 3 added Music files.  
This software fixes that, as it will only copy the real changed data.

---

## Dependencies
- [Python 3.6](https://python.org) or above
- Any Linux OS

---

## Usage
### Installation
1. Clone the repository to the folder where you want it
2. Start btsoot.py
3. If this is your first start, it'll ask you to let btsoot download dependencies
4. Done

### Create a block
`./btsoot add <block-name> <path> <path-to-remote-dir`  
This is written to a file named `btsoot.conf` which is created inside the directory where BTSOOT lies.

### Scan a block
`./btsoot scan <block-name>`  
This creates a scanfile at the folder where BTSOOT lies. The filename identifies the time the scan 
where initiated, and the block name. The file ending is `.btscan`.

### Backup a block
`./btsoot backup <block-name>`  
The program will search for the latest two scanfiles, and compare them for changed files, which it then copyies to their
paths on the remote location.  
**This also means that you MUST NOT change the remote files per hand. BTSOOT will not know about any changed file that it didn't changed itself.**

### Restore a block
This is not implemented yet. Incase of a dataloss, you have to copy them manually back to the source folder.

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
As BTSOOT is currently, as of 52a445fa, single threaded, the performance is not as good as it could be. Limiting factors include, but are not limited to:
- Slow disk speed
- Many little files, which slow down the CRC algorithm

|  COMMIT  | Data   |  Time             |  
| -------- | :----- | :---------------- |  
| 52a445fa | 1,9TB  |  187 Min, 22 Secs |  

---

## Roadmap and Known Problems  
- Going to add application file format with sqlite
- Going to add Multithreading
- ~~Going to add installer~~
- Going to add safety guard that aborts file copying if suddenly no file is found (Disk failure or unmount)
- File with "," in name corrupts the transmit list, solution is to stop using csv files
- Copying files is slow, solution is to stop using shutil's copy2
- Going to add verbose mode, and silence current mode for performance reasons
- Refactor all code to use common practices (such as functions, what have i done), for maintainabilitie's sake
- Going to make program a service with front and backend for usability reasons
- Going to add better outputs during scans, such as percent of files completed
- ~~Going to distribute it with usable format (format currently unclear, snaps are insecure, debs hard to make, etc)~~ Using .deb now.
- ~~Requires the compare lib, which may have a higher overhead than directly included functions~~ Introduced CRC function into BTSOOT