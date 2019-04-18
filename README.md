# tc
TeamCity command

## First time launching (no config file)
When you use `tc` for the first time you will be prompted with some informations to fill:
* **url**: http://teamcity.instance.tld
* **username**: your username on your instance
* **password**: your password on your instance
* **build id for dev**: copy from your teamcity instance url the argument buildTypeId (example: http://teamcity.instance.tld/viewType.html?buildTypeId=xxx)

## How to use it?
### Preamble
For example we have :
* 1 environment named dev
* 1 service named go-service

All commands are made at the root of your project go-service in bash.

### Deploy project
```
tc run dev
```

### Last Build deployed
```
tc last-build dev
```

### Open teamcity website page
```
tc op dev
```

## Molotov Team only
You can dowload our conf file in tc-config repository.