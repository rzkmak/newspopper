# Anipoke

`Anipoke` is configurable notification bot for web based anime fansub

## Overview

## Prerequisite

To run this program, you will need

### List System & App Dependencies

```$xslt
- Golang 1.10+
- Go Mod Enabled
- Redis
```

## How to Run

- Copy environment file from `env.example` to be `.env` or use `ENVIRONMENT VARIABLE` directly
- Copy environment file from `fansubs.yaml.example` to be `fansubs.yml`, add your favourite fansub there
- Verify and download dependencies `make dep`
- Run the app `make run`
