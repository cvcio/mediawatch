package commands

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	verbose  bool
	logLevel string
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Version: "2.0.0",
		Use:     "mediawatch",
		Short:   "MediaWatch CLI",
		Long: `MediaWatch CLI

...

  Find more information at: https://github.com/cvcio/mediawatch
`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// override loglevel if verbose logging
			if verbose {
				logLevel = "debug"
			}

			return nil
		},
	}
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
	// Initialize CLI config for the first time
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose logging")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", logrus.ErrorLevel.String(), "log level, allowed values {debug|info|warn|error|fatal|panic}")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mediawatch/config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// get home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cfgFile = home + "/.mediawatch/config.yaml"
		// set the config file type
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		// create home sub-directory if not exists
		if _, err := os.Stat(home + "/.mediawatch"); os.IsNotExist(err) {
			err := os.MkdirAll(home+"/.mediawatch", os.ModePerm)
			cobra.CheckErr(err)
		}
		viper.AddConfigPath(home + "/.mediawatch")
	}

	// check if config file exitst and otherwise create it
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		if _, err := os.Create(cfgFile); err != nil { // perm 0666
			cobra.CheckErr(err)
		}
	}

	// read config and create on error
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			cobra.CheckErr(err)
		} else {
			// Config file was found but another error was produced
			cobra.CheckErr(err)
		}
	}

	// don't use viper's build in read/write methods
	// as it uses yaml v1 and mallforms the yaml objects
	// we can bypass this problem by providing the appropriate
	// build tags, but is not suggested

	// read the configuration file
	// config, err := models.NewConfigFromFile(viper.ConfigFileUsed())
	// if err != nil {
	// 	cobra.CheckErr(err)
	// }

	// do some validation of the configuration here
	//
	// save the configuration file
	// config.WriteConfig(viper.ConfigFileUsed())
}
