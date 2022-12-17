# Ubuntu dev environment

These are helper scripts to spin up a persistent ubuntu dev environment

## Install Ubuntu CLI

```bash
./install_ubuntu_cli.sh
```

***NOTE:** Source the shell or Open new terminal*

## Usage

### Start the docker container in detached mode to keep it always running

```bash
$ ubuntu start
```

### Exec commands inside ubuntu

```bash
$ ubuntu start exec cat /etc/os-release
```

### "Bash" into Ubuntu

```bash
$ ubuntu exec bash
```

### To stop the container without loosing any installed packages inside the container

```bash 
$ ubuntu stop
```
