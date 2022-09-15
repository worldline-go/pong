# PONG 🏓

[![Codecov](https://img.shields.io/codecov/c/github/worldline-go/pong?logo=codecov&style=flat-square)](https://app.codecov.io/gh/worldline-go/pong)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/worldline-go/pong/Test?logo=github&style=flat-square&label=ci)](https://github.com/worldline-go/pong/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/worldline-go/pong.svg)](https://pkg.go.dev/github.com/worldline-go/pong)

Pong status checker.

Currently support only REST.

## Usage

Give a json or yaml file with the following structure:

```yaml
# LogLevel is the log level, default info
log_level: "debug"

check:
  rest:
  - concurrent: 1
    check:
      # URL could be multiple URLs, separated by space
    - url: "https://api.punkapi.com/v2/beers/1 https://api.punkapi.com/v2/beers/2"
      # Method is the HTTP method to use, default is GET
      method: GET
      # Timeout is in millisecond, default is 0 (no timeout)
      timeout: 1000
      # Status to check, default 200
      status: 200
```

```sh
pong test.yaml
```

## With Ansible

Get pong binary in release page add the plugin modules area.

```sh
make build
```

```sh
docker run -it --rm -v ${PWD}:/workspace williamyeh/ansible:debian9 /bin/bash
```

Inside the container

```sh
echo localhost ansible_connection=local > /etc/ansible/hosts

mkdir -p /usr/share/ansible/plugins/modules/
cp /workspace/pong /usr/share/ansible/plugins/modules

ansible-playbook /workspace/testdata/ansible/deploy_check.yml
```

Example of playbook

```yaml
- name: Check image exists
  pong:
    check:
      rest:
        - concurrent: 1
          check:
            - url: "{% for k,item in yaml_return.ansible_facts.compose.services.items() %} https://hub.docker.com/v2/repositories/{{ item.image.split(':')[0] }}/tags/{{ item.image.split(':')[1] }} {% endfor %}"
              timeout: 1000
  register: http_result
  failed_when: http_result.failed
```

## Development

Run `make` command to see available commands.

Local generate binary and docker image

```sh
goreleaser release --snapshot --rm-dist
```
