package soundclouddl

import (
	"fmt"
	"log"

	"github.com/AYehia0/soundcloud-dl/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sc <url>",
	Short: "Sc is a simple CLI application to download soundcloud tracks",
	Long: `A blazingly fast go program to download tracks from soundcloud 
		using just the URL, with some cool features and beautiful UI.
	`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// get the URL
		// TODO: check cobra docs for a cleaner way to do this
		if len(args) < 1 {
			if err := cmd.Usage(); err != nil {
				log.Fatal(err)
			}
			return
		}
		// run the core app
		internal.Sc(args)
	},
}

func Execute() {
	// initialize the arg parser variables
	InitConfigVars()

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Something went wrong : %s\n", err)
	}
}
