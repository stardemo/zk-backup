# zk-backup
simple tool to backup zk tree to a backup zk cluster

## Installation

```bash
go install github.com/stardemo/zk-backup@latest
```
# Usage
```
Usage of ./zk-backup:
  -excludepath string
        exclude path (default "/r3/failover/history,/r3/failover/doing")
  -sourceaddr string
        source zk cluster address
  -targetaddr string
        target zk cluster address
  -rootpath string
        transfer data start path,default /
Example:
    zk-backup -sourceaddr 127.0.0.1:2181 -targetaddr 127.0.0.1:2182 -rootpath /
```
