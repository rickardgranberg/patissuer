package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/rickardgranberg/patissuer/pkg/devops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const clientName = "patissuer"

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
	flagTenantId   = "aad-tenant-id"
	flagClientId   = "aad-client-id"
	flagTokenScope = "token-scope"
)

// Execute executes the root command.
func Execute(version, commit, buildTime string) error {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, time: %s)", version, commit, buildTime)
	return rootCmd.Execute()
}

func issue(cmd *cobra.Command, args []string) error {
	cl, err := devops.NewClient(viper.GetString(flagTenantId), viper.GetString(flagClientId))

	if err != nil {
		log.Printf("Error creating DevOps client: %v", err)
		return err
	}

	pat, err := cl.IssuePat(viper.GetStringSlice(flagTokenScope))

	if err != nil {
		return err
	}

	fmt.Print(pat)
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ./%s.yaml)", configFileName))
	rootCmd.PersistentFlags().String(flagTenantId, "", "AAD Tenant Id")
	rootCmd.PersistentFlags().String(flagClientId, "", "AAD Client Id")
	rootCmd.PersistentFlags().StringSlice(flagTokenScope, nil, "Azure DevOps Token Scope")

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
