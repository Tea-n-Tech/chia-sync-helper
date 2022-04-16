# Chia-Sync-Helper 🌱🩹

A small cli helping to sync chia full nodes.

## Download

Download prebuild binaries from the [releases page][releases-link] or build it yourself with [go][go-download].

```bash
git clone https://github.com/Tea-n-Tech/chia-sync-helper.git

go mod download
go build .
```

The  you can simply run the command line tool with:

```bash
./chia-sync-helper
```

[releases-link]: https://github.com/Tea-n-Tech/chia-sync-helper/releases
[go-download]: https://go.dev/dl/

## Why did I write this command line tool?

My full node was stuck syncing and I saw that every full node I'm connected to is behind me in height.
This happens quite frequently to me since the original software does not balance connections as it seems.

## What does it do?

This program removes full nodes which are behind in height and thus improves syncing.

Options:

- specify a height tolerance to allow nodes being behind (default 5000)
- run indefinitely every X minutes

## Important Note

This program is an expression of desperation, since the chia full node software is not helping me sufficiently.
By design this software disconnects full nodes which are far behind.
This is bad for the community since new full nodes will have a harder time syncing if everyone uses this software.
Depending on how many people will use this software I will address this issue. 

## Improvement Ideas

- [ ] Installation script for cron job use-case
- [ ] Use chia API instead of cli
- [ ] Allow a healthy balance of nodes being behind and nodes being further than us
- [ ] Disable the node removal if entirely synced
