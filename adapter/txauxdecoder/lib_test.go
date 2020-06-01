package txauxdecoder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/adapter/txauxdecoder"
	"github.com/crypto-com/chainindex/internal/primptr"
)

var _ = Describe("DecodeBase64", func() {
	It("should return decoded tx of Transfer transaction", func() {
		decodedTx, err := txauxdecoder.DecodeBase64(TRANFER_TX)
		Expect(err).To(BeNil())
		Expect(*decodedTx).To(Equal(txauxdecoder.DecodedTx{
			TxType: "Transfer",
			Inputs: []txauxdecoder.DecodedTxInput{
				{
					ID:    "06e297d37c2cfb2df425cd287994eea7568bc3498407d5bc27e2c301a49a3905",
					Index: uint32(1),
				},
			},
			OutputCount:           primptr.Uint32(uint32(2)),
			StakingAccountAddress: nil,
			CouncilNode:           nil,
		}))
	})

	It("should return decoded tx of NodeJoin transaction", func() {
		decodedTx, err := txauxdecoder.DecodeBase64(NODE_JOIN_TX)
		Expect(err).To(BeNil())
		Expect(*decodedTx).To(Equal(txauxdecoder.DecodedTx{
			TxType:                "NodeJoin",
			Inputs:                nil,
			OutputCount:           nil,
			StakingAccountAddress: primptr.String("0xb328a39002ede64c33bb60f1dc43f5df9eb47043"),
			CouncilNode: &txauxdecoder.DecodedTxCouncilNode{
				Name:            "canadacentral_validator_2",
				SecurityContact: "",
				ConsensusPubKey: txauxdecoder.DecodedTxConsensusPubKey{
					Type:  "tendermint/PubKeyEd25519",
					Value: "rXhu7xhqYBtJftVLKxvKN0XnpyOzxFnUEfAhD1dEF/8=",
				},
				ConfidentialInit: txauxdecoder.DecodedTxConfidentialInit{
					Cert: "RklYTUU=",
				},
			},
		}))
	})
})

const (
	TRANFER_TX   = "AAAEBuKX03ws+y30Jc0oeZTup1aLw0mEB9W8J+LDAaSaOQUBAAIAAAAAAAAAAACILvI8Gcck+IwmflRlBHvHQzIQdFoMSrOXl9a/ZgETkCyFWpwZl8tmJ0mhjGn2bPB7O/srnVyUKxV58j0C1BKFskIFJIr+lmr75vuM9m2DnyL1g2+TDu9nlZWvl+JHR49h0ToA36f0XshYVjyWG6r6O322w4HHP1ApOvuMpq458Hfy3HHnlgaKbcFrwUEbTDrgRslwpgR24jRvVY38IhkEIEHDgXVl8rsXsc04/UsaLQyYlaDi0zqks3qIuPguTQkQ5++/qDvpK6OromSAdCVHhnfdzf6UKsFFlKv9jMcUSeyg5Obx9u/76C2Gfww/iJ0dTZ/Ob0gVQGp+tWyyfa+2OMac++nhGckwru+mw5RoosoFiLMvmeMQrcHR3mcPIaBBx9p4dndcp0JIUEhzFbo29cm1Z8vsbd5nQHbNOxSN5/daIoEJUcQ="
	NODE_JOIN_TX = "AQIAAAAAAAAAAACzKKOQAu3mTDO7YPHcQ/XfnrRwQwBCAAAAAAAAAAAAZGNhbmFkYWNlbnRyYWxfdmFsaWRhdG9yXzIAAK14bu8YamAbSX7VSysbyjdF56cjs8RZ1BHwIQ9XRBf/FEZJWE1FAAEbC33EXRoay1jVBgObgohxS0Q3NFDj/IprjkML6Vj/+nqbrwYBykRAsVxPFXKB6E+qa6II57Ngb3iStwu4Awfx"
)
