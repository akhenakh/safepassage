package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

func TestMainProgram(t *testing.T) {
	pubData, err := ioutil.ReadFile("./testdata/pubkey.asc")
	if err != nil {
		t.Error("can't read pub key", err)
	}
	pubKey := string(pubData)

	for _, tt := range []struct {
		name              string
		env               []string
		expectedDecrypted string
	}{
		{
			name: "default",
			env: append(os.Environ(),
				"PLUGIN_SECRETS=ONE_SECRET,SECOND_SECRET",
				"ONE_SECRET=one",
				"PLUGIN_SECOND_SECRET=two",
				fmt.Sprintf("PLUGIN_PUBKEY=%s", pubKey),
			),
			expectedDecrypted: `{secret_name: ONE_SECRET, secret: one}
{secret_name: SECOND_SECRET, secret: two}
`,
		},
		{
			name: "env",
			env: append(os.Environ(),
				"PLUGIN_SECRETS=ONE_SECRET,SECOND_SECRET",
				"ONE_SECRET=one",
				"PLUGIN_SECOND_SECRET=two",
				fmt.Sprintf("PLUGIN_PUBKEY=%s", pubKey),
				"PLUGIN_FORMAT=env",
			),
			expectedDecrypted: `ONE_SECRET=b25l
SECOND_SECRET=dHdv
`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "main.go")
			cmd.Env = tt.env

			out, err := cmd.Output()
			if err != nil {
				t.Fatal(err)
			}

			enc := string(out)

			privData, err := ioutil.ReadFile("./testdata/privkey.asc")
			if err != nil {
				t.Error("can't read priv key", err)
			}

			// decrypt armored encrypted message using the private key
			decrypted, err := helper.DecryptMessageArmored(string(privData), nil, enc)
			if err != nil {
				t.Error(err)
			}
			if decrypted != tt.expectedDecrypted {
				t.Fatal("invalid decrypted values got:", decrypted, "expected:", tt.expectedDecrypted)
			}
		})
	}
}
