package gateway

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/bnb-chain/greenfield-go-sdk/client/sp"
	"github.com/bnb-chain/greenfield-go-sdk/keys"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bnb-chain/greenfield-storage-provider/model"
)

func TestVerifySignatureV1(t *testing.T) {
	// mock request
	urlmap := url.Values{}
	urlmap.Add("greenfield", "storage-provider")
	parms := io.NopCloser(strings.NewReader(urlmap.Encode()))
	req, err := http.NewRequest("POST", "gnfd.nodereal.com", parms)
	require.NoError(t, err)
	req.Header.Set(model.ContentTypeHeader, "application/x-www-form-urlencoded")
	req.Host = "testBucket.gnfd.nodereal.com"
	// mock
	privKey, _, addrInput := testdata.KeyEthSecp256k1TestPubAddr()
	keyManager, err := keys.NewPrivateKeyManager(hex.EncodeToString(privKey.Bytes()))
	require.NoError(t, err)
	// sign
	spClient, err := sp.NewSpClient("example.com", sp.WithKeyManager(keyManager), sp.WithSecure(false))
	require.NoError(t, err)
	err = spClient.SignRequest(req, sp.NewAuthInfo(false, ""))
	require.NoError(t, err)
	// check sign
	rc := &requestContext{
		request: req,
	}
	addrOutput, err := rc.verifySignature()
	assert.Equal(t, nil, err)
	assert.Equal(t, addrInput.String(), addrOutput.String())
}

func TestVerifySignatureV2(t *testing.T) {
	// mock request
	urlmap := url.Values{}
	urlmap.Add("greenfield", "storage-provider")
	parms := io.NopCloser(strings.NewReader(urlmap.Encode()))
	req, err := http.NewRequest("POST", "example.com", parms)
	require.NoError(t, err)
	req.Header.Set(model.ContentTypeHeader, "application/x-www-form-urlencoded")
	req.Host = "testBucket.gnfd.nodereal.com"
	// mock pk
	privKey, _, _ := testdata.KeyEthSecp256k1TestPubAddr()
	keyManager, err := keys.NewPrivateKeyManager(hex.EncodeToString(privKey.Bytes()))
	require.NoError(t, err)

	// sign
	spClient, err := sp.NewSpClient("example.com", sp.WithKeyManager(keyManager), sp.WithSecure(false))
	require.NoError(t, err)
	err = spClient.SignRequest(req, sp.NewAuthInfo(true, hex.EncodeToString([]byte("123"))))
	require.NoError(t, err)
	// check sign
	rc := &requestContext{
		request: req,
	}
	_, err = rc.verifySignature()
	assert.Equal(t, nil, err)
}

func TestParseRangeHeader(t *testing.T) {
	isRange, _, _ := parseRange("bytes=1")
	assert.Equal(t, false, isRange)

	isRange, start, end := parseRange("bytes=1-")
	assert.Equal(t, true, isRange)
	assert.Equal(t, 1, int(start))
	assert.Equal(t, -1, int(end))

	isRange, start, end = parseRange("bytes=1-100")
	assert.Equal(t, true, isRange)
	assert.Equal(t, 1, int(start))
	assert.Equal(t, 100, int(end))
}
