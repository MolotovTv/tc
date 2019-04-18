# tc
TeamCity command

## How to use it?
### Preamble
For example we have :
* 1 environment named dev
* 1 service named go-service

All commands are made at the root of your project go-service in bash.

### Deploy project
tc run dev

### Last Build deployed
tc last-build dev

### Open teamcity website page
tc op dev

## Molotov Team only
You can dowload our conf file in tc-config repository.