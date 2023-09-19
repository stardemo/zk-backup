# zk-backup

simple tool to backup zk tree to a backup zk cluster

## Installation

```bash
go install github.com/stardemo/zk-backup@latest
```

# Usage

```
Usage:
  zk-backup [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  export      export zookeeper data to boltDB
  help        Help about any command
  import      import zookeeper data from a db file that contains data
  transfer    transfer a zookeeper to another with connection

Flags:
  -d, --dst string       dest zookeeper addr
      --exclude string   exclude path
  -f, --file string      file to export or import (default ./out.db) (default "./out.db")
  -h, --help             help for zk-backup
      --root string      transfer data start path,default / (default "/")
  -s, --src string       src zookeeper addr
  -t, --toggle           Help message for toggle

Use "zk-backup [command] --help" for more information about a command.
```

# Example
- export data 
```shell
zk-backup export -s 127.0.0.1:2181 -f ./test.db
```

- import data
```shell
zk-backup import -d 127.0.0.1:2181 -f ./test.db
```

- transfer data
```shell
zk-backup transfer -s 127.0.0.1:2181 -d 192.168.6.222:2181
```