package social

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// GetConfig returns the twitter configuration
func GetConfig() (TwitterConfig, error) {
	configPath, _ := filepath.Abs("../candybar/config.json")

	file, _ := os.Open(configPath)

	decoder := json.NewDecoder(file)

	twitterConfig := TwitterConfig{}

	if err := decoder.Decode(&twitterConfig); err == nil {
		return twitterConfig, nil

	} else {
		return TwitterConfig{}, err
	}
}

// TwitterInstance get an instance of the twitter client
func TwitterInstance() (*twitter.Client, error) {
	if config, err := GetConfig(); err != nil {
		return nil, err

	} else {
		flags := flag.NewFlagSet("user-auth", flag.ExitOnError)

		consumerKey := flags.String("consumer-key", config.Twitter.ConsumerKey, "Twitter Consumer Key")
		consumerSecret := flags.String("consumer-secret", config.Twitter.ConsumerSecret, "Twitter Consumer Secret")

		accessToken := flags.String("access-token", config.Twitter.AccessToken, "Twitter App Access Token")
		accessSecret := flags.String("access-secret", config.Twitter.AccessTokenSecret, "Twitter Access Secret")

		flags.Parse(os.Args[1:])

		flagutil.SetFlagsFromEnv(flags, "TWITTER")

		if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
			log.Fatal("Access Tokens/Secrets are missing")
		}

		oAuthConfig := oauth1.NewConfig(*consumerKey, *consumerSecret)

		token := oauth1.NewToken(*accessToken, *accessSecret)

		httpClient := oAuthConfig.Client(oauth1.NoContext, token)

		client := twitter.NewClient(httpClient)

		accountParams := &twitter.AccountVerifyParams{
			SkipStatus:   twitter.Bool(true),
			IncludeEmail: twitter.Bool(true),
		}

		user, _, _ := client.Accounts.VerifyCredentials(accountParams)

		fmt.Printf("User's account:\n%+v\n", user)

		return client, nil
	}
}
