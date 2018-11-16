# Taipan

Go-configuration library for multiple environments, powered by [spf13/viper](https://github.com/spf13/viper)

## Overview

Taipan merges configuration from multiple files to enable shared configuration between different environments. 
It returns a [spf13/viper.Viper](https://github.com/spf13/viper/blob/master/viper.go) with easy access to the configuration.

Usage:
```
taipan.Get().GetString("db.username")
```
    

## Defaults
Base config should be named `config.yaml`. 

Files are loaded from from `resources/` 
(overridable with `TAIPAN_CONFIG_PATH`).

If `config-local.yaml` is present in the current working directory it will be applied last. 

## Multiple profiles
Multiple configs can be merged by specifying a comma-separated list in `TAIPAN_PROFILES`

Setting `TAIPAN_PROFILES=staging` will result in the following files being loaded in order (by default from `resources/`)
```
config.yaml
config-staging.yaml
config-local.yml (if present)
```

Given these config files

`resources/config.yaml`
```yaml
port: 8080

db:
  password: default-password
```
`resources/config-staging.yaml`
```yaml
db:
  password: staging-password
```

The result is

```go
config := taipan.Get()
config.GetInt("port") // 8080
config.GetString("db.password") // "staging-password"
``` 