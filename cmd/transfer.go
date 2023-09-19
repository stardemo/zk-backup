/*
Copyright © 2023 starliu <starliu1995@hotmail.com>
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stardemo/zk-backup/pkg/zookeeper"
	"log"
	"strings"
)

// transferCmd represents the transfer command
var (
	transferCmd = &cobra.Command{
		Use:   "transfer",
		Short: "transfer a zookeeper to another with connection",
		Long:  `with two zookeeper connection,to transfer data`,
		Run:   transferCmdFunc,
	}
)

func init() {
	rootCmd.AddCommand(transferCmd)
}

func transferCmdFunc(cmd *cobra.Command, args []string) {
	if !cmd.Flags().Lookup("src").Changed {
		log.Fatal("missing src flags.")
	}
	if !cmd.Flags().Lookup("dst").Changed {
		log.Fatal("missing dst flags.")
	}
	log.Printf("start transfer zookeeper data from %s to %s,start with %s,exclude %s", srcAddr, dstAddr, rootPath, excludePath)
	var exPaths []string
	if excludePath != "" {
		exPaths = strings.Split(excludePath, ",")
	}
	sourceConn, _, err := zookeeper.DialZk(srcAddr)
	if err != nil {
		log.Fatalf("source zk connect failed %s", err.Error())
	}
	dstConn, _, err := zookeeper.DialZk(dstAddr)
	if err != nil {
		log.Fatalf("dst zk connect failed %s", err.Error())
	}
	zookeeper.Walk(rootPath, sourceConn, dstConn, exPaths)
}
