package cmd

import (
	"time"

	defaultlog "log"

	. "github.com/app-sre/user-validator/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	cfgFile  string
	logLevel string

	rootCmd = &cobra.Command{
		Use:   "user-validator",
		Short: "user-validator integration",
		Long:  "Integration, that verifies the content of users in app-interface",
	}

	userValidatorCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate users",
		Long:  "Run validations for pgp keys, usernames and github logins",
		Run: func(cmd *cobra.Command, args []string) {
			userValidator()
		},
	}

	accountNotifierCmd = &cobra.Command{
		Use:   "account-notifier",
		Short: "Sent out notifications on new passwords",
		Long:  "Send pgp encrypted passwords to users",
		Run: func(cmd *cobra.Command, args []string) {
			accountNotifier()
		},
	}
)

// Execute executes the rootCmd
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(userValidatorCmd)
	rootCmd.AddCommand(accountNotifierCmd)
	rootCmd.PersistentFlags().StringVarP(&logLevel, "logLevel", "l", "info", "Log level")
	userValidatorCmd.Flags().StringVarP(&cfgFile, "cfgFile", "c", "", "Configuration File")
	accountNotifierCmd.Flags().StringVarP(&cfgFile, "cfgFile", "c", "", "Configuration File")

	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(configureLogging)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		Log().Debugw("Using configuration", "config", cfgFile)
	}
}

func configureLogging() {
	loggerConfig := zap.NewDevelopmentConfig()

	switch logLevel {
	case "info":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "debug":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "error":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "warn":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "fatal":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	zap.ReplaceGlobals(logger)

	if err != nil {
		defaultlog.Fatal(err)
	}
}
