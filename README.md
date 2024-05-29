# smtp-cli

Simple SMTP client

## Usage

```bash
# See ./smtp-cli --help for all supported options
echo My email body | ./smtp-cli --subject ...
```

## Local testing

```bash
podman compose -f test/docker-compose.yml up -d
go test ./...
podman compose -f test/docker-compose.yml down
```
