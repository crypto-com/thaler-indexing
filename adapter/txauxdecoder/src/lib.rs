use base64;
use parity_scale_codec::Decode;
use std::os::raw::c_char;
use std::ffi::{CStr, CString};
use serde::Serialize;
use chain_core::state::account::{CouncilNode, StakedStateAddress};
use chain_core::tx::data::input::{TxoPointer, TxoSize};
use chain_core::tx::{TxAux, TxEnclaveAux, TxPublicAux};
use chain_tx_validation::witness::verify_tx_recover_address;

#[derive(Serialize)]
struct DecodedTx<'a> {
    tx_type: DecodedTxType,
    inputs: Option<&'a Vec<TxoPointer>>,
    output_count: Option<&'a TxoSize>,
    staked_state_address: Option<StakedStateAddress>,
    council_node_meta: Option<&'a CouncilNode>
}

#[derive(Serialize)]
enum DecodedTxType {
    Transfer,
    Deposit,
    Unbond,
    Withdraw,
    NodeJoin,
    Unjail,
}

#[no_mangle]
pub extern "C" fn decode_base64(encoded_tx_ptr: *const c_char) -> *mut c_char {
    let encoded_tx_cstr = unsafe {
        assert!(!encoded_tx_ptr.is_null(), "called passed a null pointer argument");

        CStr::from_ptr(encoded_tx_ptr)
    };

    let encoded_tx = encoded_tx_cstr.to_str().unwrap();

    let encoded_tx = base64::decode(encoded_tx).expect("error base64 decoding transaction");
    let tx_aux = TxAux::decode(&mut encoded_tx.as_slice()).expect("error decoding TxAux");

    let decoded_tx = decoded_tx_from_tx_aux(&tx_aux);

    let decoded_tx = serde_json::to_string(&decoded_tx).expect("error serializing decoded transaction into json");
    let decoded_tx_cstr = CString::new(decoded_tx).unwrap();

    decoded_tx_cstr.into_raw()
}

fn decoded_tx_from_tx_aux(tx_aux: &TxAux) -> DecodedTx {
    match tx_aux {
        TxAux::EnclaveTx(enclave_tx) => match enclave_tx {
            TxEnclaveAux::TransferTx {
                inputs,
                no_of_outputs,
                ..
            } => {
                DecodedTx {
                    tx_type: DecodedTxType::Transfer,
                    inputs: Some(inputs),
                    output_count: Some(no_of_outputs),
                    staked_state_address: None,
                    council_node_meta: None,
                } 
            },
            TxEnclaveAux::DepositStakeTx {tx, ..} => {
                DecodedTx {
                    tx_type: DecodedTxType::Deposit,
                    inputs: Some(&tx.inputs),
                    output_count: None,
                    staked_state_address: None,
                    council_node_meta: None,
                } 
            }
            TxEnclaveAux::WithdrawUnbondedStakeTx {no_of_outputs, witness, payload} => {
                let staked_state_address = verify_tx_recover_address(&witness, &payload.txid).expect("error recovering staked state address");
                DecodedTx {
                    tx_type: DecodedTxType::Withdraw,
                    inputs: None,
                    output_count: Some(no_of_outputs),
                    staked_state_address: Some(staked_state_address),
                    council_node_meta: None,
                } 
            },
        },
        TxAux::PublicTx(public_tx) => match public_tx {
            TxPublicAux::UnbondStakeTx( unbond_tx, .. ) => {
                DecodedTx {
                    tx_type: DecodedTxType::Unbond,
                    inputs: None,
                    output_count: None,
                    staked_state_address: Some(unbond_tx.from_staked_account),
                    council_node_meta: None,
                }
            },
            TxPublicAux::NodeJoinTx (node_join_tx, ..) => {
                DecodedTx {
                    tx_type: DecodedTxType::NodeJoin,
                    inputs: None,
                    output_count: None,
                    staked_state_address: Some(node_join_tx.address),
                    council_node_meta: Some(&node_join_tx.node_meta),
                }
            },
            TxPublicAux::UnjailTx(unjail_tx, ..) => {
                DecodedTx {
                    tx_type: DecodedTxType::Unjail,
                    inputs: None,
                    output_count: None,
                    staked_state_address: Some(unjail_tx.address),
                    council_node_meta: None,
                }
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn decode_free(s: *mut c_char) {
    unsafe {
        if s.is_null() {
            return;
        }
        CString::from_raw(s)
    };
}
