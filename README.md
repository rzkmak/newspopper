# Newspopper

Subscribe to any feed or web, deliver to your telegram bot/channel

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

- Copy environment file from `sites.yaml.example` to be `sites.yaml`, add your favourite fansub there
- Verify and download dependencies `make dep`
- Run the app `make run`

## How to Simulate

- Setup dev environment properly
- Update sites.yaml to be desired value, put `stdout` value in your target.

For example:
```$xslt
listener:
- type: "rss"
  url: "http://your-site.com/feed"
  interval: 1m
  target: "stdout"

```
