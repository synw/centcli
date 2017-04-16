# Centcli

Terminal client for the [Centrifugo](https://github.com/centrifugal/centrifugo/) websockets server.

## Usage

Edit `config.json` and add servers

## Available commands

### Initialization

Server to use (must be set in config): start the program with the flag `-s=server_name` will initialize a connection

**`USE`**: `use server_name`

Server actually in use:

**`USING`**: `using`

To set the client's Centrifugo username use the flag `-u=clientusername`.

### Statistics

Get stats about Centrifugo:

**`STATS`**: options:
- `stats all`
- `stats node`
- `stats client`
- `stats http`

Get a particular statistic:

**`STAT`**: `stat node_num_channels`

All the Centrifugo statistics are available: 
[check the complete list](https://fzambia.gitbooks.io/centrifugal/content/server/stats.html)

### Actions

**`PUBLISH`**: `publish channel_name {"foo":"bar"}

Note: do not use spaces in your payload

**`LISTEN`**: `listen channel_name`

**`STOP`**: `stop listen channel_name`

### Channels

**`COUNT`**: `count chans`

Get all channels on server:

**`CHANS`**: `chans`

List of users in a channel:

**`PRESENCE`**: `presence channel_name`

Channel history:

**`HISTORY`**: `history channel_name`


