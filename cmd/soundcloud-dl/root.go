package soundclouddl

import "os"

var Search bool
var DownloadPath string

// define flags and handle configuration
func InitConfigVars() {
	tmpDLdir, _ := os.Getwd()
	rootCmd.PersistentFlags().BoolVarP(&Search, "search", "s", false, "Check if the track exists or not.")
	rootCmd.PersistentFlags().StringVarP(&DownloadPath, "download-path", "p", tmpDLdir, "The download path where tracks are stored.")
}
