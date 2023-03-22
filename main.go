package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type VerifyRequest struct {
	Account string `json:"account,omitempty"`
	Sign    string `json:"sign,omitempty"`
	Message string `json:"message,omitempty"`
}

func main() {
	http.HandleFunc("/api/verify", func(w http.ResponseWriter, r *http.Request) {
		var req VerifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body"+err.Error(), http.StatusBadRequest)
			return
		}
		msgStr := strings.TrimPrefix(req.Message, "0x")
		msg, err := hex.DecodeString(msgStr)
		if err != nil {
			http.Error(w, "invalid message: "+err.Error(), http.StatusBadRequest)
			return
		}

		if !verifySig(req.Account, req.Sign, msg) {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func verifySig(from, sigHex string, msg []byte) bool {
	sig := hexutil.MustDecode(sigHex)

	msg = accounts.TextHash(msg)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	return strings.EqualFold(from, recoveredAddr.Hex())
}
