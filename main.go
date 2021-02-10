package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/namsral/flag"
)

const (
	prefix = "PLUGIN"

	formatStd = "std"
	formatEnv = "env"
)

func main() {
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], prefix, 0)

	var (
		ssecrets = fs.String(
			"secrets",
			"",
			"comma separated secrets to recover",
		)

		pubKey = fs.String(
			"pubKey",
			"",
			"the OpenPGP public key to encode secrets to",
		)
		format = fs.String(
			"format",
			formatStd,
			fmt.Sprintf("The output format -- must be either %q or %q", formatStd, formatEnv),
		)
	)

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal("can't parse args")
	}

	if *format != formatStd && *format != formatEnv {
		log.Fatalf("Format must be either %q or %q but got %q", formatStd, formatEnv, *format)
	}

	var secrets []string

	for _, u := range strings.Split(*ssecrets, ",") {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			secrets = append(secrets, u)
		}
	}

	if len(secrets) == 0 {
		log.Println("no secrets passed in settings to recover")
		os.Exit(2)
	}

	if pubKey == nil || *pubKey == "" {
		log.Println("no public key passed")
		os.Exit(2)
	}

	var msg string

	for _, secret := range secrets {
		senv := os.Getenv(secret)
		if senv == "" {
			senv = os.Getenv("PLUGIN_" + secret)
			if senv == "" {
				log.Printf("can't find the secret named variable: %s\n", secret)
				os.Exit(2)
			}
		}
		switch *format {
		case formatEnv:
			msg += fmt.Sprintf("%s=%s\n", secret, base64.StdEncoding.EncodeToString([]byte(senv)))
		default:
			msg += fmt.Sprintf("{secret_name: %s, secret: %s}\n", secret, senv)
		}
	}

	// encrypt message using public key
	armor, err := helper.EncryptMessageArmored(*pubKey, msg)
	if err != nil {
		log.Fatal(err, *pubKey)
	}

	fmt.Println(armor)
}
