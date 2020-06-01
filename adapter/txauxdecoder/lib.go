package txauxdecoder

/*
#cgo LDFLAGS: -L${SRCDIR}/./target/release -ltxauxdecoder -ldl
#include "bindings.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

func MustDecodeBase64(rawTx string) *DecodedTx {
	decodedTx, err := DecodeBase64(rawTx)
	if err != nil {
		panic(err)
	}

	return decodedTx
}

func DecodeBase64(rawTx string) (*DecodedTx, error) {
	var err error

	arg := C.CString(rawTx)
	defer C.free(unsafe.Pointer(arg))

	decodedTxPtr := C.decode_base64(arg)
	defer C.decode_free(decodedTxPtr)

	var decodedTx DecodedTx
	if err = jsoniter.UnmarshalFromString(C.GoString(decodedTxPtr), &decodedTx); err != nil {
		return nil, fmt.Errorf("error deserializing decoded tx: %v", err)
	}
	return &decodedTx, nil
}

type DecodedTx struct {
	// TODO: Parse TxType to pre-defined constants
	TxType                string                `json:"tx_type"`
	Inputs                []DecodedTxInput      `json:"inputs"`
	OutputCount           *uint32               `json:"output_count"`
	StakingAccountAddress *string               `json:"staked_state_address"`
	CouncilNode           *DecodedTxCouncilNode `json:"council_node_meta"`
}

type DecodedTxInput struct {
	ID    string `json:"id"`
	Index uint32 `json:"index"`
}

type DecodedTxCouncilNode struct {
	Name             string                    `json:"name"`
	SecurityContact  string                    `json:"security_contact"`
	ConsensusPubKey  DecodedTxConsensusPubKey  `json:"consensus_pubkey"`
	ConfidentialInit DecodedTxConfidentialInit `json:"confidential_init"`
}

type DecodedTxConsensusPubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type DecodedTxConfidentialInit struct {
	Cert string `json:"cert"`
}
