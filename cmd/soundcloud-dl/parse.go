package soundclouddl

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AYehia0/soundcloud-dl/internal"
	"github.com/AYehia0/soundcloud-dl/pkg/theme"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "sc <url>",
	Short: "Sc is a simple CLI application to download soundcloud tracks",
	Long: `A blazingly fast go program to download tracks from soundcloud 
		using just the URL, with some cool features and beautiful UI.
	`,
	Args:    cobra.ArbitraryArgs,
	Version: "v1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		// get the URL
		if len(args) < 1 && !Search {
			if err := cmd.Usage(); err != nil {
				log.Fatal(err)
			}
			return
		}
		// run the core app
		// FIXME: Probably not the best thing to do lol, it's better to just pass it to the function, who cares.
		internal.Sc(args, DownloadPath, BestQuality, Search)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if len(args) != 0 {
				// check if <url> is passed with --search-and-download flag
				if strings.HasPrefix(args[0], "https") && Search {
					fmt.Printf("Can't use/pass a %s with --%s flag\n\n", theme.Green("<url>"), theme.Red(f.Name))
					cmd.Usage()
					os.Exit(0)
				}
				// check if the best quality is passed with --quality
				if BestQuality && TrackQuality != "" {
					qFlag := cmd.Flag("quality").Name
					fmt.Printf("Can't use/pass a %s with --%s flag\n\n", theme.Green(qFlag), theme.Red(f.Name))
					cmd.Usage()
					os.Exit(0)
				}
			}
		})
	},
}

func Execute() {
	// initialize the arg parser variables
	InitConfigVars()

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Something went wrong : %s\n", err)
	}
}
