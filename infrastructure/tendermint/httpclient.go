package tendermint

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	tendermintadapter "github.com/crypto-com/chainindex/adapter/tendermint"
	"github.com/crypto-com/chainindex/adapter/tendermint/types"
	"github.com/crypto-com/chainindex/internal/primptr"
)

type HTTPClient struct {
	httpClient *http.Client
	serverUrl  string
}

func NewHTTPClient(serverUrl string) *HTTPClient {
	httpClient := &http.Client{
		// TODO: configurable timeout
		Timeout: 10 * time.Second,
	}

	return &HTTPClient{
		httpClient,
		serverUrl,
	}
}

func (client *HTTPClient) Genesis() (*types.Genesis, error) {
	var err error

	rawRespBody, err := client.request("genesis")
	if err != nil {
		return nil, err
	}
	defer rawRespBody.Close()

	genesis, err := client.parseGenesisResp(rawRespBody)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}

func (client *HTTPClient) parseGenesisResp(rawRespReader io.Reader) (*types.Genesis, error) {
	var resp types.GenesisResp
	if err := jsoniter.NewDecoder(rawRespReader).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling Tendermint genesis response: %v", err)
	}

	return &types.Genesis{
		GenesisTime: resp.Result.Genesis.GenesisTime,
		ChainID:     resp.Result.Genesis.ChainID,
		AppHash:     resp.Result.Genesis.AppHash,
		AppState: types.GenesisAppState{
			CouncilNodes: client.parseGenesisCouncilNodes(resp.Result.Genesis.AppState.CouncilNodes),
			Distribution: client.parseGenesisDistribution(resp.Result.Genesis.AppState.Distribution),
		},
	}, nil
}

func (client *HTTPClient) parseGenesisCouncilNodes(rawNodes map[string][]interface{}) []types.GenesisCouncilNode {
	nodes := make([]types.GenesisCouncilNode, 0, len(rawNodes))
	for address, rawNode := range rawNodes {
		pubkey, _ := rawNode[2].(map[string]interface{})
		nodes = append(nodes, types.GenesisCouncilNode{
			StakingAccountAddress: address,
			Address:               tendermintadapter.AddressFromPubKey(pubkey["value"].(string)),
			Name:                  rawNode[0].(string),
			SecurityContact:       rawNode[1].(string),
			PubKeyType:            pubkey["type"].(string),
			PubKey:                pubkey["value"].(string),
		})
	}

	sort.SliceStable(nodes, func(i, j int) bool {
		return strings.Compare(nodes[i].Address, nodes[j].Address) < 0
	})
	return nodes
}

func (client *HTTPClient) parseGenesisDistribution(rawDistribution map[string][]string) []types.GenesisDistribution {
	distribution := make([]types.GenesisDistribution, 0, len(rawDistribution))
	for stakingAddress, rawEntry := range rawDistribution {
		distType := rawEntry[0]
		unit := primptr.String(rawEntry[1])

		entry := types.GenesisDistribution{
			StakingAccountAddress: stakingAddress,
			Bonded:                nil,
			Unbonded:              nil,
		}

		switch distType {
		case types.RAW_GENESIS_DISTRIBUTION_TYPE_BONDED:
			entry.Bonded = unit
		case types.RAW_GENESIS_DISTRIBUTION_TYPE_UNBONDED:
			entry.Unbonded = unit
		default:
			panic(fmt.Sprintf("error parsing genesis distribution type: unknown type %s", distType))
		}

		distribution = append(distribution, entry)
	}

	sort.SliceStable(distribution, func(i, j int) bool {
		return strings.Compare(distribution[i].StakingAccountAddress, distribution[j].StakingAccountAddress) < 0
	})
	return distribution
}

func (client *HTTPClient) LatestBlockHeight() (uint64, error) {
	var err error

	rawRespBody, err := client.request("block_results")
	if err != nil {
		return uint64(0), err
	}
	defer rawRespBody.Close()

	blockResults, err := client.parseBlockResultsResp(rawRespBody)
	if err != nil {
		return uint64(0), err
	}

	return blockResults.Height, nil
}

func (client *HTTPClient) BlockResults(height uint64) (*types.BlockResults, error) {
	var err error

	rawRespBody, err := client.request("block_results", "height="+strconv.FormatUint(height, 10))
	if err != nil {
		return nil, err
	}
	defer rawRespBody.Close()

	blockResults, err := client.parseBlockResultsResp(rawRespBody)
	if err != nil {
		return nil, err
	}

	return blockResults, nil
}

func (client *HTTPClient) parseBlockResultsResp(rawRespReader io.Reader) (*types.BlockResults, error) {
	var err error

	var resp types.BlockResultsResp
	if err = jsoniter.NewDecoder(rawRespReader).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling Tendermint block_results response: %v", err)
	}

	var txsResults [][]types.BlockResultsEvent
	if resp.Result.TxsEvents != nil {
		txsResults = client.parseBlockResultsTxsEvents(resp.Result.TxsEvents)
	}

	var beginBlockEvents []types.BlockResultsEvent
	if resp.Result.BeginBlockEvents != nil {
		beginBlockEvents = client.parseBlockResultsEvent(resp.Result.BeginBlockEvents)
	}

	height, err := strconv.ParseUint(resp.Result.Height, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting block height to unsigned integer: %v", err)
	}
	return &types.BlockResults{
		Height:           uint64(height),
		TxsEvents:        txsResults,
		BeginBlockEvents: beginBlockEvents,
		ValidatorUpdates: client.parseBlockResultsValidatorUpdates(resp.Result.ValidatorUpdates),
	}, nil
}

func (client *HTTPClient) parseBlockResultsTxsEvents(rawResults []types.RawBlockResultsTxResult) [][]types.BlockResultsEvent {
	results := make([][]types.BlockResultsEvent, 0, len(rawResults))
	for _, rawResult := range rawResults {
		results = append(results, client.parseBlockResultsEvent(rawResult.Events))
	}

	return results
}

func (client *HTTPClient) parseBlockResultsEvent(rawEvents []types.RawBlockResultsEvent) []types.BlockResultsEvent {
	if rawEvents == nil {
		return nil
	}

	events := make([]types.BlockResultsEvent, 0, len(rawEvents))
	for _, rawEvent := range rawEvents {
		attributes := make([]types.BlockResultsEventAttribute, 0, len(rawEvent.Attributes))
		for _, rawAttribute := range rawEvent.Attributes {
			attributes = append(attributes, types.BlockResultsEventAttribute{
				Key:   rawAttribute.Key,
				Value: rawAttribute.Value,
			})
		}
		events = append(events, types.BlockResultsEvent{
			Type:       rawEvent.Type,
			Attributes: attributes,
		})
	}

	return events
}

func (client *HTTPClient) parseBlockResultsValidatorUpdates(rawUpdates []types.RawBlockResultsValidator) []types.BlockResultsValidator {
	if rawUpdates == nil {
		return nil
	}

	updates := make([]types.BlockResultsValidator, 0, len(rawUpdates))
	for _, rawUpdate := range rawUpdates {
		updates = append(updates, types.BlockResultsValidator{
			PubKey: types.BlockResultsValidatorPubKey{
				Type:    rawUpdate.PubKey.Type,
				PubKey:  rawUpdate.PubKey.Data,
				Address: tendermintadapter.AddressFromPubKey(rawUpdate.PubKey.Data),
			},
			Power: rawUpdate.Power,
		})
	}

	return updates
}

func (client *HTTPClient) Block(height uint64) (*types.Block, error) {
	var err error

	rawRespBody, err := client.request("block", "height="+strconv.FormatUint(height, 10))
	if err != nil {
		return nil, err
	}
	defer rawRespBody.Close()

	block, err := client.parseBlockResp(rawRespBody)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (client *HTTPClient) parseBlockResp(rawRespReader io.Reader) (*types.Block, error) {
	var err error

	var resp types.BlockResp
	if err = jsoniter.NewDecoder(rawRespReader).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling Tendermint block_results response: %v", err)
	}

	height, err := strconv.ParseUint(resp.Result.Block.Header.Height, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting block height to unsigned integer: %v", err)
	}
	return &types.Block{
		Height:         height,
		Hash:           resp.Result.BlockID.Hash,
		Time:           resp.Result.Block.Header.Time,
		AppHash:        resp.Result.Block.Header.AppHash,
		PropserAddress: resp.Result.Block.Header.ProposerAddress,
		Txs:            resp.Result.Block.Data.Txs,
		Signatures:     client.parseBlockSignatures(resp.Result.Block.LastCommit.Signatures),
	}, nil
}

func (client *HTTPClient) parseBlockSignatures(rawSignatures []types.RawBlockSignature) []types.BlockSignature {
	if rawSignatures == nil {
		return nil
	}

	signatures := make([]types.BlockSignature, 0, len(rawSignatures))
	for _, rawSignature := range rawSignatures {
		if rawSignature.Signature == nil {
			continue
		}
		signatures = append(signatures, types.BlockSignature{
			ValidatorAddress: rawSignature.ValidatorAddress,
			Signature:        *rawSignature.Signature,
		})
	}

	return signatures
}

func (client *HTTPClient) request(method string, queryString ...string) (io.ReadCloser, error) {
	var err error

	url := client.serverUrl + "/" + method
	if len(queryString) > 0 {
		url += "?" + queryString[0]
	}
	rawResp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error requesting Tendermint %s endpoint: %v", method, err)
	}

	if rawResp.StatusCode != 200 {
		rawResp.Body.Close()
		return nil, fmt.Errorf("error requesting Tendermint %s endpoint: %s", method, rawResp.Status)
	}

	return rawResp.Body, nil
}
