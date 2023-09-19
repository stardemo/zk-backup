/*
Copyright © 2023 starliu <starliu1995@hotmail.com>
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stardemo/zk-backup/pkg/zookeeper"
	bolt "go.etcd.io/bbolt"
	"log"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import zookeeper data from a db file that contains data",
	Long: `import zookeeper data from a db file that contains data with key value
file could be generate with zk-backup cli`,
	Run: importCmdFunc,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func importCmdFunc(cmd *cobra.Command, args []string) {
	if !cmd.Flags().Lookup("src").Changed {
		log.Fatal("missing src flags.")
	}
	log.Printf("start import zookeeper data from %s to conn %s,start with %s,exclude %s", dbFile, srcAddr, rootPath, excludePath)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatalf("db init error: %s", err.Error())
	}
	defer db.Close()
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(zookeeper.BoltBackupBucket))
		return err
	}); err != nil {
		log.Fatalf("db bucket init error: %s", err.Error())
	}
	dstConn, _, err := zookeeper.DialZk(dstAddr)
	if err != nil {
		log.Fatalf("dst zk connect failed %s", err.Error())
	}
	if err := zookeeper.RestoreFromDB(db, dstConn); err != nil {
		log.Fatal(err)
	}
	log.Println("finished")
}
