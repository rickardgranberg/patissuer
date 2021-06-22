package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rickardgranberg/patissuer/pkg/auth"
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

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List users PATs",
		RunE:  list,
	}
)

const (
	configFileName      = "." + clientName
	flagTenantId        = "aad-tenant-id"
	flagClientId        = "aad-client-id"
	flagClientSecret    = "aad-client-secret"
	flagLoginMethod     = "login-method"
	flagLoginToken      = "login-token"
	flagOrganizationUrl = "org-url"
	flagTokenScope      = "token-scope"
	flagTokenTTL        = "token-ttl"
	flagOutput          = "output"
)

// Execute executes the root command.
func Execute(version, commit, buildTime string) error {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, time: %s)", version, commit, buildTime)
	return rootCmd.Execute()
}

func issue(cmd *cobra.Command, args []string) error {
	authClient, err := auth.NewAuthClient(viper.GetString(flagTenantId), viper.GetString(flagClientId), viper.GetString(flagClientSecret))

	if err != nil {
		return fmt.Errorf("failed to initialize auth client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	t, err := authClient.Login(ctx, viper.GetString(flagLoginMethod), viper.GetString(flagLoginToken))

	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	cl, err := devops.NewClient(viper.GetString(flagOrganizationUrl), t)

	if err != nil {
		log.Printf("Error creating DevOps client: %v", err)
		return err
	}

	validTo := time.Now().Add(viper.GetDuration(flagTokenTTL))

	pat, err := cl.IssuePat(ctx, args[0], viper.GetStringSlice(flagTokenScope), validTo)

	if err != nil {
		return err
	}

	format := viper.GetString(flagOutput)
	switch format {
	case "raw":
		fmt.Print(pat.Token)
	case "json":
		b, err := json.Marshal(pat)
		if err != nil {
			return err
		}
		fmt.Print(string(b))
	}
	return nil
}

func list(cmd *cobra.Command, args []string) error {
	authClient, err := auth.NewAuthClient(viper.GetString(flagTenantId), viper.GetString(flagClientId), viper.GetString(flagClientSecret))

	if err != nil {
		return fmt.Errorf("failed to initialize auth client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	t, err := authClient.Login(ctx, viper.GetString(flagLoginMethod), viper.GetString(flagLoginToken))

	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	cl, err := devops.NewClient(viper.GetString(flagOrganizationUrl), t)

	if err != nil {
		log.Printf("Error creating DevOps client: %v", err)
		return err
	}

	pats, err := cl.ListPats(ctx)

	if err != nil {
		return err
	}

	format := viper.GetString(flagOutput)
	switch format {
	case "raw":
		for _, t := range pats {
			fmt.Printf("%s %s %s\n", t.AuthorizationId, t.DisplayName, t.Scope)
		}
	case "json":
		b, err := json.Marshal(pats)
		if err != nil {
			return err
		}
		fmt.Print(string(b))
	}

	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ./%s.yaml)", configFileName))
	rootCmd.PersistentFlags().String(flagTenantId, "", "AAD Tenant Id")
	rootCmd.PersistentFlags().String(flagClientId, "", "AAD Client Id")
	rootCmd.PersistentFlags().String(flagClientSecret, "", fmt.Sprintf("AAD Client Secret. Only required for %s login", auth.LoginMethodDeviceCode))
	rootCmd.PersistentFlags().String(flagOrganizationUrl, "", "Azure DevOps Organization URL")
	rootCmd.PersistentFlags().String(flagOutput, "raw", "Output format, 'raw' or 'json'")
	rootCmd.PersistentFlags().String(flagLoginMethod, auth.LoginMethodInteractive, fmt.Sprintf("Login method, valid options are '%s', '%s' and '%s'", auth.LoginMethodInteractive, auth.LoginMethodDeviceCode, auth.LoginMethodBearerToken))
	rootCmd.PersistentFlags().String(flagLoginToken, "", "The bearer token when using 'token' login method")

	issueCmd.Flags().StringSlice(flagTokenScope, nil, "Azure DevOps PAT Token Scope")
	issueCmd.Flags().Duration(flagTokenTTL, time.Hour*24*30, "Azure DevOps PAT Token TTL")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatalf("Error: %v", err)
	}
	if err := viper.BindPFlags(issueCmd.Flags()); err != nil {
		log.Fatalf("Error: %v", err)
	}

	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(listCmd)
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
