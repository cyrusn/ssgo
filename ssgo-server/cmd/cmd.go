package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/subosito/gotenv"
	auth "github.com/cyrusn/goJWTAuthHelper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	CONFIG_PATH       = "./config.yaml"
	PRIVATE_KEY       = "skill-vein-planet-neigh-envoi"
	DEFAULT_DSN       = "./database/master.sqlite"
	DEFAULT_OVERWRITE = true
	DEFAULT_PORT      = ":5000"
	STATIC_FOLDER_LOCATION = "./public"
	DEFAULT_LIFE_TIME = 30

	CONTEXT_KEY_NAME = "authClaim"
	JWT_KEY_NAME     = "jwt"
	ROLE_KEY_NAME    = "Role"
)

var (
	cfgFile string
	port    string
	staticFolderLocation string
	dsn             string
	isOverwrite     bool
	privateKey      string
	lifeTime        int64
	secret          auth.Secret
)

func initConfig() {
	// Try loading .env from current and parent directories
	gotenv.Load(".env")
	gotenv.Load("../.env")
	
	if os.Getenv("DEFAULT_SCHOOL_NAME") != "" {
		fmt.Println("Successfully loaded environment variables from .env")
	} else {
		fmt.Println("Warning: Could not find DEFAULT_SCHOOL_NAME in environment. .env may not be loaded.")
	}
	
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}

// initiate secret token for serveCmd
func initSecret() {
	secret = auth.New(
		CONTEXT_KEY_NAME, JWT_KEY_NAME, ROLE_KEY_NAME, []byte(privateKey),
	)
}

func initDefault() {
	viper.SetDefault("key", PRIVATE_KEY)
	viper.SetDefault("dsn", DEFAULT_DSN)
	viper.SetDefault("overwrite", DEFAULT_OVERWRITE)
	viper.SetDefault("port", DEFAULT_PORT)
	viper.SetDefault("static", STATIC_FOLDER_LOCATION)
	viper.SetDefault("time", DEFAULT_LIFE_TIME)
}

func initVariables() {
	privateKey = os.Getenv("JWT_PRIVATE_KEY")
	if privateKey == "" {
		privateKey = viper.GetString("key")
	}
	dsn = viper.GetString("dsn")
	isOverwrite = viper.GetBool("overwrite")
	port = viper.GetString("port")
	staticFolderLocation = viper.GetString("static")
	lifeTime = viper.GetInt64("time")
}

func init() {
	cobra.OnInitialize(initConfig, initDefault, initVariables, initSecret)
}

// Execute excute all commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
