package zookeeper

import "github.com/go-zookeeper/zk"

const (
	PERM_FILE        = zk.PermAdmin | zk.PermRead | zk.PermWrite
	PermDirectory    = zk.PermAdmin | zk.PermCreate | zk.PermDelete | zk.PermRead | zk.PermWrite
	BoltBackupBucket = "zk-backup"
)
