# ruready

`ruready` is an application that exposes an HTTP endpoint to allow you to check if a machine is ready. It can be used in conjunction with `docker-compose` to make one container wait for another one to be ready. The ready state is determined by a command passed as argument.

[![Build Status](https://travis-ci.org/victorhurdugaci/ruready.svg?branch=master)](https://travis-ci.org/victorhurdugaci/ruready)

## Downloads

- The latest compiled version is available on the [releases page](https://github.com/victorhurdugaci/ruready/releases/latest)
- To compile from sources, build the repository using [go](https://golang.org/)

## Full working example

*The `client` container waits for `mysql` on the `db` container to become available.*

Create a file called `docker-compose.yml` with the content below:

```
version: '2'
services:
  db:
    image: "mysql"
    environment:
        - MYSQL_ROOT_PASSWORD=1q2w3e
    ports:
     - "8099:8099"
    command: bash -c "apt-get update && apt-get -y install curl && curl -L https://github.com/victorhurdugaci/ruready/releases/download/1.0.0-alpha.1/ruready > ruready && chmod +x ruready && (./ruready -c mysql -- --host=localhost --port=3306 --user=root --password=1q2w3e --execute=quit & ./entrypoint.sh mysqld)"
    
  client:
    image: "busybox"
    command: sh -c "while [ true ]; do wget -q -O - http://db:8099/ready && break; echo 'Not ready...'; sleep 1s; done; echo 'Ready!';"
```

Run `docker-compose up` in the folder where you created the file.

## Command line arguments

| Argument           | Required | Default | Description                                                             |
| ------------------ | -------- | ------- | ----------------------------------------------------------------------- |
| `-c`/`--command`   | Yes      |         | The command that checks if the machine is ready                         |
| `-t`/`--cachetime` | No       | 3       | Number of seconds to cache the result of the command before  reinvoking |
| `-p`/`--port`      | No       | 8099    | Server port                                                             |
| `-v`/`--version`   | No       |         | Shows version information                                               |
| `--`               | No       |         | Anything that follows `--` is passed as argument to `command`

**Examples**

- `ruready -c ls -- ./opt/app/started.txt` 
    Waits for the file `./opt/app/started.txt` to exist. Endpoint `<hostname>:8099/ready`

- `ruready -p 3000 -c ls -- ./opt/app/started.txt` 
    Waits for the file `./opt/app/started.txt` to exist. Endpoint `<hostname>:3000/ready`


## Response

- If the machine is ready, `/ready` endpoint returns `200 (OK)` and the message `ruready: Ready`
- If the machine is not ready, `/ready` endpoint returns `503 (Service Unavailable)` and the message `ruready: Not Ready`

## Checking if the machine is available

- **curl**: `while [ true ]; do curl -I http://<host>:8099/ready && break; echo 'Not ready...'; sleep 1s; done;`
- **wget**: `while [ true ]; do wget -q -O - http://<host>:8099/ready && break; echo 'Not ready...'; sleep 1s; done;`
