package codec_fixtures

import (
	"fmt"
	"os"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
)

func TestCodecs(t *testing.T) {
	dirs, err := os.ReadDir("../fixtures/")
	if err != nil {
		t.Fatalf("failed to open fixtures dir: %v", err)
	}

	for _, dir := range dirs {
		fixtureName := dir.Name()
		if !dir.IsDir() {
			continue
		}
		if reason, blacklisted := FixtureBlacklist[fixtureName]; blacklisted {
			fmt.Printf("Skipping fixture '%v': %v\n", fixtureName, reason)
			continue
		}
		t.Run(fixtureName, func(t *testing.T) {
			data, err := loadFixture(fixtureName)
			if err != nil {
				t.Fatalf("failed to load fixture: %v", err)
			}
			for fromCodec := range data {
				for toCodec := range data {
					msg := fmt.Sprintf("decode(%v)->encode(%v)", fromCodec, toCodec)
					verifyCid(t, msg, data[fromCodec].value, codecs[toCodec], data[toCodec].cid)
				}
			}
		})
	}
}

func verifyCid(t *testing.T, desc string, node ipld.Node, toEnc ipld.LinkPrototype, expected cid.Cid) {
	actual, err := nodeToCid(toEnc, node)
	if err != nil {
		t.Fatalf("failed to convert node to CID: %v", err)
	}
	if !expected.Equals(actual) {
		t.Fatalf("[%v] generated CID (%v) does not match expected (%v)", desc, expected.String(), actual.String())
	}
}
