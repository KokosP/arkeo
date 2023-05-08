package indexer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/types"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// type attributeProvider interface {
// 	attributes() map[string]string
// }

type attributes func() map[string]string

func parseEventToEventOpenContract(event interface{}) (atypes.EventOpenContract, error) {
	eventData := make(map[string]string)

	prefix := "arkeo.arkeo.EventOpenContract."
	switch evt := event.(type) {
	case ctypes.ResultEvent:
		for key, attribute := range evt.Events {
			k := strings.TrimPrefix(key, prefix)
			v := strings.Trim(attribute[0], `"`)
			eventData[k] = v
		}
	case abcitypes.Event:
		for _, attribute := range evt.Attributes {
			key := strings.TrimPrefix(string(attribute.GetKey()), prefix)
			value := strings.Trim(string(attribute.GetValue()), `"`)
			eventData[key] = value
		}
	default:
		return atypes.EventOpenContract{}, fmt.Errorf("unsupported event type: %T", evt)
	}

	type eventOpenContractAlias atypes.EventOpenContract
	eventOpenContract := struct {
		ContractId         string `json:"contract_id,omitempty"`
		Height             string `json:"height,omitempty"`
		Duration           string `json:"duration,omitempty"`
		Rate               string `json:"rate,omitempty"`
		Deposit            string `json:"deposit,omitempty"`
		Type               string `json:"type"`
		OpenCost           string `json:"open_cost"`
		SettlementDuration string `json:"settlement_duration"`
		Authorization      string `json:"authorization"`
		QueriesPerMinute   string `json:"queries_per_minute"`
		eventOpenContractAlias
	}{}

	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}

	if err := json.Unmarshal(jsonData, &eventOpenContract); err != nil {
		return atypes.EventOpenContract{}, err
	}

	result := atypes.EventOpenContract(eventOpenContract.eventOpenContractAlias)

	// make conversions
	result.QueriesPerMinute, err = strconv.ParseInt(eventOpenContract.QueriesPerMinute, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.OpenCost, err = strconv.ParseInt(eventOpenContract.OpenCost, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.Deposit, _ = cosmos.NewIntFromString(eventOpenContract.Deposit)
	result.ContractId, err = strconv.ParseUint(eventOpenContract.ContractId, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.Height, err = strconv.ParseInt(eventOpenContract.Height, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.SettlementDuration, err = strconv.ParseInt(eventOpenContract.SettlementDuration, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.Duration, err = strconv.ParseInt(eventOpenContract.Duration, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	err = json.Unmarshal([]byte(eventOpenContract.Rate), &result.Rate)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.Authorization = atypes.ContractAuthorization(atypes.ContractAuthorization_value[eventOpenContract.Authorization])
	result.Type = atypes.ContractType(atypes.ContractType_value[eventOpenContract.Type])
	return result, nil
}

func parseEventToEventModProvider(event interface{}) (atypes.EventModProvider, error) {
	eventData := make(map[string]string)

	prefix := "arkeo.arkeo.EventModProvider."
	switch evt := event.(type) {
	case ctypes.ResultEvent:
		for key, attribute := range evt.Events {
			k := strings.TrimPrefix(key, prefix)
			v := strings.Trim(attribute[0], `"`)
			eventData[k] = v
		}
	case abcitypes.Event:
		for _, attribute := range evt.Attributes {
			key := strings.TrimPrefix(string(attribute.GetKey()), prefix)
			value := strings.Trim(string(attribute.GetValue()), `"`)
			eventData[key] = value
		}
	default:
		return atypes.EventModProvider{}, fmt.Errorf("unsupported event type: %T", evt)
	}

	type eventModProviderAlias atypes.EventModProvider
	eventModProvider := struct {
		MaxContractDuration string `json:"max_contract_duration,omitempty"`
		MinContractDuration string `json:"min_contract_duration,omitempty"`
		SettlementDuration  string `json:"settlement_duration,omitempty"`
		MetadataNonce       string `json:"metadata_nonce,omitempty"`
		SubscriptionRate    string `json:"subscription_rate"`
		PayAsYouGoRate      string `json:"pay_as_you_go_rate"`
		Status              string `json:"status"`
		eventModProviderAlias
	}{}

	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return atypes.EventModProvider{}, err
	}

	if err := json.Unmarshal(jsonData, &eventModProvider); err != nil {
		return atypes.EventModProvider{}, err
	}

	result := atypes.EventModProvider(eventModProvider.eventModProviderAlias)

	// make conversions
	result.MaxContractDuration, err = strconv.ParseInt(eventModProvider.MaxContractDuration, 10, 64)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	result.MinContractDuration, err = strconv.ParseInt(eventModProvider.MinContractDuration, 10, 64)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	result.SettlementDuration, err = strconv.ParseInt(eventModProvider.SettlementDuration, 10, 64)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	result.MetadataNonce, err = strconv.ParseUint(eventModProvider.MetadataNonce, 10, 64)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	err = json.Unmarshal([]byte(eventModProvider.SubscriptionRate), &result.SubscriptionRate)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	err = json.Unmarshal([]byte(eventModProvider.PayAsYouGoRate), &result.PayAsYouGoRate)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	result.Status = atypes.ProviderStatus(atypes.ProviderStatus_value[eventModProvider.Status])
	return result, nil
}

func wsAttributeSource(src ctypes.ResultEvent) func() map[string]string {
	attribs := make(map[string]string, len(src.Events))
	for k, v := range src.Events {
		if len(v) > 0 {
			key := k
			if sl := strings.Split(k, "."); len(sl) > 1 {
				key = sl[len(sl)-1]
			}
			if _, ok := attribs[key]; ok {
				log.Debugf("key %s already in results with value %s, overwriting with %s", key, attribs[key], v[0])
			}
			attribs[key] = strings.Trim(v[0], `"`)
		}
		if len(v) > 1 {
			log.Warnf("attrib %s has %d array values: %v", k, len(v), v)
		}
	}
	attribs["eventHeight"] = attribs["height"]
	return func() map[string]string { return attribs }
}

func tmAttributeSource(tx tmtypes.Tx, evt abcitypes.Event, height int64) func() map[string]string {
	attribs := make(map[string]string, 0)
	for _, attr := range evt.Attributes {
		attribs[string(attr.Key)] = string(attr.Value)
	}

	if tx != nil {
		if _, ok := attribs["hash"]; !ok {
			attribs["hash"] = strings.ToUpper(hex.EncodeToString(tx.Hash()))
		}
	}

	attribs["eventHeight"] = strconv.FormatInt(height, 10)
	if _, ok := attribs["height"]; !ok {
		attribs["height"] = attribs["eventHeight"]
	}

	return func() map[string]string { return attribs }
}

func (a *IndexerApp) handleValidatorPayoutEvent(evt types.ValidatorPayoutEvent) error {
	log.Infof("receieved validatorPayoutEvent %#v", evt)
	if evt.Paid < 0 {
		return fmt.Errorf("received negative paid amt: %d for tx %s", evt.Paid, evt.TxID)
	}
	if evt.Paid == 0 {
		return nil
	}
	log.Infof("upserting validator payout event for tx %s", evt.TxID)
	if _, err := a.db.UpsertValidatorPayoutEvent(evt); err != nil {
		return errors.Wrapf(err, "error upserting validator payout event")
	}
	return nil
}

func (a *IndexerApp) consumeEvents(clients []*tmclient.HTTP) error {
	// splitting across multiple tendermint clients as websocket allows max of 5 subscriptions per client
	blockEvents := subscribe(clients[0], "tm.event = 'NewBlock'")
	bondProviderEvents := subscribe(clients[0], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'")
	modProviderEvents := subscribe(clients[0], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'")
	openContractEvents := subscribe(clients[1], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")
	closeContractEvents := subscribe(clients[1], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'")
	claimContractIncomeEvents := subscribe(clients[1], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	log.Infof("beginning realtime event consumption")
	for {
		select {
		case evt := <-blockEvents:
			data, ok := evt.Data.(tmtypes.EventDataNewBlock)
			if !ok {
				log.Errorf("event not block: %T", evt.Data)
				continue
			}
			log := log.WithField("height", strconv.FormatInt(data.Block.Height, 10))
			log.Debugf("received block: %d", data.Block.Height)

			if err := a.handleBlockEvent(data.Block); err != nil {
				log.Errorf("error handling block event %d: %+v", data.Block.Height, err)
			}

			endBlockEvents := data.ResultEndBlock.Events
			log.Debugf("block %d with %d endBlock events", data.Block.Height, len(endBlockEvents))
			for _, evt := range endBlockEvents {
				switch evt.GetType() {
				case "validator_payout":
					validatorPayoutEvent := types.ValidatorPayoutEvent{}
					if err := convertEvent(tmAttributeSource(nil, evt, data.Block.Height), &validatorPayoutEvent); err != nil {
						log.Errorf("error converting validator_payout event: %+v", err)
						break
					}
					if err := a.handleValidatorPayoutEvent(validatorPayoutEvent); err != nil {
						log.Errorf("error handling validator_payout event: %+v", err)
					}
				case "contract_settlement":
					contractSettlementEvent := types.ContractSettlementEvent{}
					if err := convertEvent(tmAttributeSource(nil, evt, data.Block.Height), &contractSettlementEvent); err != nil {
						log.Errorf("error converting contract_settlement event: %+v", err)
						break
					}
					if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
						log.Errorf("error handling close_contract contract_settlement event: %+v", err)
					}
				}
			}
		case evt := <-openContractEvents:
			log.Debugf("received open contract event")
			openContractEvent, err := parseEventToEventOpenContract(evt)
			if err != nil {
				log.Errorf("error converting open_contract event: %+v", err)
				break
			}
			if err := a.handleOpenContractEvent(openContractEvent); err != nil {
				log.Errorf("error handling open_contract event: %+v", err)
			}
		case evt := <-bondProviderEvents:
			log.Debugf("received bond provider event")
			bondProviderEvent := types.BondProviderEvent{}
			if err := convertEvent(wsAttributeSource(evt), &bondProviderEvent); err != nil {
				log.Errorf("error converting bond_provider event: %+v", err)
				break
			}
			if err := a.handleBondProviderEvent(bondProviderEvent); err != nil {
				log.Errorf("error handling bond_provider event: %+v", err)
			}
		case evt := <-modProviderEvents:
			log.Debugf("received mod provider event")
			modProviderEvent, err := parseEventToEventModProvider(evt)
			if err != nil {
				log.Errorf("error converting mod_provider event: %+v", err)
				break
			}
			if err := a.handleModProviderEvent(modProviderEvent); err != nil {
				log.Errorf("error handling mod_provider event: %+v", err)
			}
		case evt := <-claimContractIncomeEvents:
			log.Debugf("received claim contract income event")
			claimContractIncomeEvent := types.ClaimContractIncomeEvent{}
			attribs := wsAttributeSource(evt)
			// hack contract_settlement height
			wrapped := func() map[string]string {
				tmp := attribs()
				if ev, ok := evt.Events["contract_settlement.height"]; ok && len(ev) > 0 {
					tmp["height"] = evt.Events["contract_settlement.height"][0]
				}
				return tmp
			}
			if err := convertEvent(wrapped, &claimContractIncomeEvent); err != nil {
				log.Errorf("error converting open_contract event: %+v", err)
				break
			}
			if err := a.handleContractSettlementEvent(claimContractIncomeEvent.ContractSettlementEvent); err != nil {
				log.Errorf("error handling claim contract income event: %+v", err)
			}
		case evt := <-closeContractEvents:
			log.Debugf("received close_contract event")
			closeContractEvent := types.CloseContractEvent{}
			attribs := wsAttributeSource(evt)
			wrapped := func() map[string]string {
				tmp := attribs()
				if ev, ok := evt.Events["contract_settlement.height"]; ok && len(ev) > 0 {
					tmp["height"] = evt.Events["contract_settlement.height"][0]
				}
				return tmp
			}

			if err := convertEvent(wrapped, &closeContractEvent); err != nil {
				log.Errorf("error converting close_contract event: %+v", err)
				break
			}

			if err := a.handleCloseContractEvent(closeContractEvent); err != nil {
				log.Errorf("error handling close_contract event: %+v", err)
			}
		case <-quit:
			log.Infof("received os quit signal")
			return nil
		}
	}
}

func (a *IndexerApp) consumeHistoricalBlock(client *tmclient.HTTP, bheight int64) (result *db.Block, err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var block *ctypes.ResultBlock
	var blockResults *ctypes.ResultBlockResults
	var blockErr, resultsErr error

	go func() {
		defer wg.Done()
		start := time.Now()
		block, blockErr = client.Block(context.Background(), &bheight)
		if time.Since(start) > 500*time.Millisecond {
			log.Warnf("%.3f elapsed reading block %d", time.Since(start).Seconds(), bheight)
		}
	}()

	go func() {
		defer wg.Done()
		start := time.Now()
		blockResults, resultsErr = client.BlockResults(context.Background(), &bheight)
		if time.Since(start) > 500*time.Millisecond {
			log.Warnf("%.3f elapsed reading block results %d", time.Since(start).Seconds(), bheight)
		}
	}()
	wg.Wait()

	if blockErr != nil {
		return nil, errors.Wrapf(blockErr, "error reading block")
	}
	if resultsErr != nil {
		return nil, errors.Wrapf(resultsErr, "error reading block results")
	}

	log := log.WithField("height", strconv.FormatInt(block.Block.Height, 10))
	for _, transaction := range block.Block.Txs {
		txInfo, err := client.Tx(context.Background(), transaction.Hash(), false)
		if err != nil {
			log.Warnf("failed to get transaction data for %s", transaction.Hash())
			continue
		}

		for _, event := range txInfo.TxResult.Events {
			log.Debugf("received %s txevent", event.Type)
			if err := a.handleAbciEvent(event, transaction, block.Block.Height); err != nil {
				log.Errorf("error handling abci event %#v\n%+v", event, err)
			}
		}
	}

	for _, event := range blockResults.EndBlockEvents {
		log.Debugf("received %s endblock event", event.Type)
		if err := a.handleAbciEvent(event, nil, block.Block.Height); err != nil {
			log.Errorf("error handling abci event %#v\n%+v", event, err)
		}
	}

	r := &db.Block{
		Height:    block.Block.Height,
		Hash:      block.Block.Hash().String(),
		BlockTime: block.Block.Time,
	}
	return r, nil
}

func (a *IndexerApp) handleAbciEvent(event abcitypes.Event, transaction tmtypes.Tx, height int64) error {
	var err error
	switch event.Type {
	case "provider_bond":
		bondProviderEvent := types.BondProviderEvent{}
		if err = convertEvent(tmAttributeSource(transaction, event, height), &bondProviderEvent); err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		if err = a.handleBondProviderEvent(bondProviderEvent); err != nil {
			log.Errorf("error handling %s event: %+v", event.Type, err)
		}
	case "provider_mod":
		modProviderEvent, err := parseEventToEventModProvider(event)
		if err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		if err = a.handleModProviderEvent(modProviderEvent); err != nil {
			log.Errorf("error handling %s event: %+v", event.Type, err)
		}
	case "open_contract":
		openContractEvent, err := parseEventToEventOpenContract(event)
		if err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		if err = a.handleOpenContractEvent(openContractEvent); err != nil {
			log.Errorf("error handling %s event: %+v", event.Type, err)
		}
	case "claim_contract_income":
		contractSettlementEvent := types.ContractSettlementEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &contractSettlementEvent); err != nil {
			log.Errorf("error converting claim_contract_income event: %+v", err)
			break
		}
		if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
			log.Errorf("error handling claim contract income event: %+v", err)
		}
	case "validator_payout":
		validatorPayoutEvent := types.ValidatorPayoutEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &validatorPayoutEvent); err != nil {
			log.Errorf("error converting validatorPayoutEvent event: %+v", err)
			break
		}
		if err := a.handleValidatorPayoutEvent(validatorPayoutEvent); err != nil {
			log.Errorf("error handling claim contract income event: %+v", err)
		}
	case "contract_settlement":
		contractSettlementEvent := types.ContractSettlementEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &contractSettlementEvent); err != nil {
			log.Errorf("error converting contractSettlementEvent: %+v", err)
			break
		}
		if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
			log.Errorf("error handling contractSettlementEvent: %+v", err)
		}
	case "close_contract":
		log.Debugf("received close_contract event")
		closeContractEvent := types.CloseContractEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &closeContractEvent); err != nil {
			log.Errorf("error converting close_contract event: %+v", err)
			break
		}
		if err := a.handleCloseContractEvent(closeContractEvent); err != nil {
			log.Errorf("error handling close contract event: %+v", err)
		}
	default:
		log.Debugf("ignored event %s", event.Type)
	}
	return nil
}

// copy attributes of map given by attributeFunc() to target which must be a pointer (map/slice implicitly ptr)
func convertEvent(attributeFunc attributes, target interface{}) error {
	return mapstructure.WeakDecode(attributeFunc(), target)
}

func subscribe(client *tmclient.HTTP, query string) <-chan ctypes.ResultEvent {
	out, err := client.Subscribe(context.Background(), "", query)
	if err != nil {
		log.Errorf("failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}
	return out
}
