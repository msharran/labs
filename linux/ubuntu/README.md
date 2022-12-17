# Ubuntu dev environment

These are helper scripts to spin up a persistent ubuntu dev environment

## Start the docker container in detached mode to keep it always running

```bash
$ ./start.sh
```

## Exec commands inside ubuntu

```bash
$ ./exec.sh cat /etc/os-release
PRETTY_NAME="Ubuntu 22.04.1 LTS"
NAME="Ubuntu"
VERSION_ID="22.04"
VERSION="22.04.1 LTS (Jammy Jellyfish)"
VERSION_CODENAME=jammy
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=jammy
```

## "Bash" into Ubuntu

```bash
$ ./exec.sh bash
```

## To stop the container without loosing any installed packages inside the container

```bash 
$ ./stop.sh
```
