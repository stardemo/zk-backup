/*
Copyright © 2023 starliu <starliu1995@hotmail.com>
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "zk-backup",
		Short: "A tool for zookeeper to transfer or backup",
		Long:  `A tool for zookeeper to transfer or backup`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
	srcAddr     string
	dstAddr     string
	rootPath    string
	excludePath string
	dbFile      string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&srcAddr, "src", "s", "", "src zookeeper addr")
	rootCmd.PersistentFlags().StringVarP(&dstAddr, "dst", "d", "", "dest zookeeper addr")
	rootCmd.PersistentFlags().StringVar(&rootPath, "root", "/", "transfer data start path,default /")
	rootCmd.PersistentFlags().StringVar(&excludePath, "exclude", "", "exclude path")
	rootCmd.PersistentFlags().StringVarP(&dbFile, "file", "f", "./out.db", "file to export or import (default ./out.db)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
