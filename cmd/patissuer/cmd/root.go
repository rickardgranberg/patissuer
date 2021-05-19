package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const clientName = "patsissuer"

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   clientName,
		Short: "Azure DevOps PATS Issuer",
	}

	issueCmd = &cobra.Command{
		Use:   "issue [name]",
		Short: "Issue a PAT",
		Args:  cobra.ExactArgs(1),
		RunE:  issue,
	}
)

const (
	configFileName = "." + clientName
	tenantIdFlag   = "aad-tenant-id"
	clientIdFlag   = "aad-client-id"
	tokenScopeFlag = "token-scope"
)

// Execute executes the root command.
func Execute(version, commit, buildTime string) error {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, time: %s)", version, commit, buildTime)
	log.Printf("%s version %s", rootCmd.Use, rootCmd.Version)
	return rootCmd.Execute()
}

func issue(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ./%s.yaml)", configFileName))
	rootCmd.PersistentFlags().String(tenantIdFlag, "", "AAD Tenant Id")
	rootCmd.PersistentFlags().String(clientIdFlag, "", "AAD Client Id")
	rootCmd.PersistentFlags().StringSlice(tokenScopeFlag, nil, "Azure DevOps Token Scope")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatalf("Error: %v", err)
	}
	if err := viper.BindPFlags(issueCmd.Flags()); err != nil {
		log.Fatalf("Error: %v", err)
	}

	rootCmd.AddCommand(issueCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in working directory with name ".apigw" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName(configFileName)
	}

	viper.SetEnvPrefix(rootCmd.Use)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
