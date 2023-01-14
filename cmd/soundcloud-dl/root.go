package soundclouddl

import "os"

var Search bool
var DownloadPath string
var BestQuality bool

// define flags and handle configuration
func InitConfigVars() {
	tmpDLdir, _ := os.Getwd()
	rootCmd.PersistentFlags().BoolVarP(&Search, "search-and-download", "s", false, "Search for tracks by title and prompt one for download ")
	rootCmd.PersistentFlags().StringVarP(&DownloadPath, "download-path", "p", tmpDLdir, "The download path where tracks are stored.")
	rootCmd.PersistentFlags().BoolVarP(&BestQuality, "best", "b", false, "Download with the best available quality.")
}
