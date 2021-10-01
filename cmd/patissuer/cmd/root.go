package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	configFileName      = clientName
	flagTenantId        = "aad-tenant-id"
	flagClientId        = "aad-client-id"
	flagLoginMethod     = "login-method"
	flagLoginToken      = "login-token"
	flagLoginRetry      = "login-retry"
	flagOrganizationUrl = "org-url"
	flagTokenScope      = "token-scope"
	flagTokenTTL        = "token-ttl"
	flagOutput          = "output"
	flagOutputFile      = "output-file"
)

// Execute executes the root command.
func Execute(version, commit, buildTime string) error {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, time: %s)", version, commit, buildTime)
	return rootCmd.Execute()
}

func issue(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt(flagLoginRetry))*time.Minute)
	defer cancel()

	cl, err := loginAndCreateClient(ctx)

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
		outputContent(pat.Token)
	case "json":
		b, err := json.Marshal(pat)
		if err != nil {
			return err
		}
		outputContent(string(b))
	}
	return nil
}

func list(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt(flagLoginRetry))*time.Minute)
	defer cancel()

	cl, err := loginAndCreateClient(ctx)

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
		b := strings.Builder{}

		for _, t := range pats {
			b.WriteString(fmt.Sprintf("%s %s %s\n", t.AuthorizationId, t.DisplayName, t.Scope))
		}
		outputContent(b.String())
	case "json":
		b, err := json.Marshal(pats)
		if err != nil {
			return err
		}
		outputContent(string(b))
	}

	return nil
}

func loginAndCreateClient(ctx context.Context) (*devops.Client, error) {
	authClient, err := auth.NewAuthClient(viper.GetString(flagTenantId), viper.GetString(flagClientId))

	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth client: %w", err)
	}

	var t string

	for i := 0; i < viper.GetInt(flagLoginRetry); i++ {
		sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		t, err = authClient.Login(sctx, viper.GetString(flagLoginMethod), viper.GetString(flagLoginToken))
		if err == context.Canceled {
			return nil, fmt.Errorf("login canceled %w", err)
		}
		if err != nil {
			log.Printf("failed to login with error: %v\n", err)
			log.Printf("Retrying login (%d of %d)...", i+1, viper.GetInt(flagLoginRetry))
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	cl, err := devops.NewClient(viper.GetString(flagOrganizationUrl), t)

	if err != nil {
		log.Printf("Error creating DevOps client: %v", err)
		return nil, err
	}

	return cl, nil
}

func outputContent(format string, a ...interface{}) error {
	fn := viper.GetString(flagOutputFile)

	if fn == "" {
		fmt.Printf(format, a...)
	} else {
		f, err := os.Create(fn)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		f.WriteString(fmt.Sprintf(format, a...))
	}

	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ./%s.yaml)", configFileName))
	rootCmd.PersistentFlags().String(flagTenantId, "", "AAD Tenant Id")
	rootCmd.PersistentFlags().String(flagClientId, "", "AAD Client Id")
	rootCmd.PersistentFlags().String(flagOrganizationUrl, "", "Azure DevOps Organization URL")
	rootCmd.PersistentFlags().String(flagOutput, "raw", "Output format, 'raw' or 'json'")
	rootCmd.PersistentFlags().String(flagOutputFile, "", "File name to save output in, instead of printing to stdout")
	rootCmd.PersistentFlags().String(flagLoginMethod, auth.LoginMethodInteractive, fmt.Sprintf("Login method, valid options are '%s', '%s' and '%s'", auth.LoginMethodInteractive, auth.LoginMethodDeviceCode, auth.LoginMethodBearerToken))
	rootCmd.PersistentFlags().String(flagLoginToken, "", "The bearer token when using 'token' login method")
	rootCmd.PersistentFlags().Int(flagLoginRetry, 3, "The number of times to retry the login phase")

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
		viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", configFileName))
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
