/*
Copyright © 2023 starliu <starliu1995@hotmail.com>
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stardemo/zk-backup/pkg/zookeeper"
	bolt "go.etcd.io/bbolt"
	"log"
	"strings"
)

// exportCmd represents the export command
var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "export zookeeper data to boltDB",
		Long: `Export data to backup & transfer & archived ,
data will generate db file with bolt db format.`,
		Run: exportCmdFunc,
	}
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

func exportCmdFunc(cmd *cobra.Command, args []string) {
	if !cmd.Flags().Lookup("src").Changed {
		log.Fatal("missing src flags.")
	}
	log.Printf("start export zookeeper data from %s to file %s,start with %s,exclude %s", srcAddr, dbFile, rootPath, excludePath)
	sourceConn, _, err := zookeeper.DialZk(srcAddr)
	if err != nil {
		log.Fatalf("source zk connect failed %s", err.Error())
	}
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
	var exPaths []string
	if excludePath != "" {
		exPaths = strings.Split(excludePath, ",")
	}
	zookeeper.WalkIntoDB(rootPath, sourceConn, db, exPaths)
}
