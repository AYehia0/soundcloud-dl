package soundclouddl

var Search bool

// define flags and handle configuration
func InitConfigVars() {
	rootCmd.PersistentFlags().BoolVarP(&Search, "search", "s", false, "Does the track exist ?")
}
