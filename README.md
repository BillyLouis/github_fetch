# github_fetch
Golang code to fetch all repositories that Organization|User contains. This will clone the repos and rearrange them into alike languages such as all C++ program will be cloned into one directory, the python repos will be placed into one python directory and so on.  The code is only intends for public repositories.

## Directories:

## Pre-requisites:
[A] - Have a linux machine (Ubuntu 22.04 is prefered)

[B] - Have an external hardrive of more than 4TB:
```shell
* The reason is some organisations may have more tan 6,000 repos and that will take at least 3TB of storage space.
* Constantly monitor the cloning process since once in a while, 1 or 2 hours in the cloning process the sudo token session
  may expire and it would require to re-enter the sudo password.

```

# Workflow:
Connect the external hardrive.
Move to the path [/home/{YOUR_PATH}/github_fetch/cmd/github_fetch/fetch_with_token/](/home/{YOUR_PATH}/github_fetch/cmd/github_fetch/fetch_with_token/).
Add you github JWT Token to the variable "const githubToken" in the main.go file and issue the commands:
```shell
# Check your harddrive name on the system
    lsblk
# Which should output something similar:
    sda                8:0    0 953.9G  0 disk 
    ├─sda1             8:1    0   512M  0 part /boot/efi
    └─sda2             8:2    0 953.4G  0 part 
    ├─vgkubuntu-root
    │              252:0    0 930.4G  0 lvm  /var/snap/firefox/common/host-hunspell
    │                                        /
    └─vgkubuntu-swap_1
                    252:1    0   976M  0 lvm  [SWAP]
    sdb                8:16   0   3.6T  0 disk 
    ├─sdb1             8:17   0    16M  0 part 
    └─sdb2             8:18   0   3.6T  0 part /home/...

#Then
sudo mount /dev/sdb2 ~/home/<YOUR PATH>github_fetch/cmd/github_fetch/fetch_with_token/drive_mount/
run main.go

#in the prompt, enter the organisation name like XMSS, microsoft etc ...
# When done, unmount the hardrive with the command:
sudo umount /dev/sdb2

``` 

## Authors
- [Billy Louis](): Golang code to fetch multiple github repos at once from a public organisation.

## Badges
Hardware Team: [NSAL.com](https://NSAL.com/)

[![NSA License](https://img.shields.io/badge/License-NSAL-green.svg)](https://choosealicense.com/licenses/nsal/)
