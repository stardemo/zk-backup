package zookeeper

import (
	"github.com/go-zookeeper/zk"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"path"
)

func WalkIntoDB(root string, conn *zk.Conn, tDB *bolt.DB, excludePaths []string) {
	children, _, err := conn.Children(root)
	if err != nil {
		log.Printf("error, when get children of %s, %s\n", root, err)
		os.Exit(1)
	}

	for _, node := range children {
		fullPath := path.Join(root, node)
		if IsPathExcluded(excludePaths, fullPath) {
			return
		}
		data, stat, _ := conn.Get(fullPath)
		if stat.EphemeralOwner == 0 {
			// ignore ephemeral node
			if err := tDB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(BoltBackupBucket))
				err := b.Put([]byte(fullPath), data)
				return err
			}); err != nil {
				log.Printf("error, when insert node in target db, %v\n", err)
				os.Exit(1)
			}
			log.Println(fullPath, " backup into db success")
		}
		WalkIntoDB(fullPath, conn, tDB, excludePaths)
	}
}

func RestoreFromDB(db *bolt.DB, conn *zk.Conn) error {
	return db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(BoltBackupBucket))
		return b.ForEach(func(k, v []byte) error {
			if _, err := CreateRecursive(conn, string(k), v, 0, zk.WorldACL(zk.PermAll)); err != nil {
				log.Printf("error, when create node in target zk, %v\n", err)
				return err
			}
			log.Printf("path: %s restore success", string(k))
			return nil
		})
	})
}
