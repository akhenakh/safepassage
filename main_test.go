package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

const validMsg = `{secret_name: ONE_SECRET, secret: one}
{secret_name: SECOND_SECRET, secret: two}
`

func TestMainProgram(t *testing.T) {
	pubData, err := ioutil.ReadFile("./testdata/pubkey.asc")
	if err != nil {
		t.Error("can't read pub key", err)
	}
	pubKey := string(pubData)

	cmd := exec.Command("go", "run", "main.go")
	cmd.Env = append(os.Environ(),
		"PLUGIN_SECRETS=ONE_SECRET,SECOND_SECRET",
		"ONE_SECRET=one",
		"PLUGIN_SECOND_SECRET=two",
		fmt.Sprintf("PLUGIN_PUBKEY=%s", pubKey),
	)

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
	if decrypted != validMsg {
		t.Fatal("invalid decrypted values got:", decrypted, "expected:", validMsg)
	}
}
