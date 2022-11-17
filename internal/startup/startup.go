// Helper function for servers startup

package startup

import (
	"crypto/rsa"
	"io/ioutil"
	"log"

	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/dgrijalva/jwt-go"
)

// GetAuthenticator creates an authenticator based on the given configuration
func GetAuthenticator(cfg *config.Config) (auth.Authenticator, error) {
	if cfg.Auth.Debug {
		return auth.DebugAuthenticator{}, nil
	}

	keyContents, err := ioutil.ReadFile(cfg.Auth.PrivateKeyFile)
	if err != nil {
		log.Fatalf("main: Reading auth private key: %v", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		log.Fatalf("main: Parsing auth private key: %v", err)
	}

	publicKeyLookup := auth.NewSingleKeyFunc(cfg.Auth.KeyID, key.Public().(*rsa.PublicKey))

	authenticator, err := auth.NewDefaultAuthenticator(key, cfg.Auth.KeyID, cfg.Auth.Algorithm, publicKeyLookup)
	if err != nil {
		log.Fatalf("main: Constructing authenticator: %v", err)
	}

	return authenticator, err
}
