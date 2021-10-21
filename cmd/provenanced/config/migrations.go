package config

import (
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tmconfig "github.com/tendermint/tendermint/config"
)

var (
	// The addedKeys, removedKeys, and toDashesKeys lists come from the output of the TestCompare34and35 test in legacy_test.go.
	// Manual analysis was then done to identify the changedKeys.
	// Each entry put into changedKeys was then removed from the addedKeys and removedKeys lists.

	// addedKeys are keys in v0.35 that were added since v0.34.
	// Keys that changed names are in the changedKeys map, but not here.
	addedKeys = []string{
		"other", "mempool.ttl-duration", "mempool.ttl-num-blocks", "mempool.version",
		"p2p.bootstrap-peers", "p2p.max-connections", "p2p.max-incoming-connection-attempts",
		"p2p.queue-type", "p2p.use-legacy",
		"priv-validator.client-certificate-file", "priv-validator.client-key-file",
		"priv-validator.home", "priv-validator.root-ca-file",
		"statesync.use-p2p",
	}
	// removedKeys are keys in v0.34 that were removed in v0.35
	// Keys that changed names are in the changedKeys map, but not here.
	removedKeys = []string{
		"mempool.wal_dir",
		"p2p.test_fuzz", "p2p.test_fuzz_config.maxdelay", "p2p.test_fuzz_config.mode",
		"p2p.test_fuzz_config.probdropconn", "p2p.test_fuzz_config.probdroprw", "p2p.test_fuzz_config.probsleep",
	}
	// toDashesKeys are keys that use underscores in v0.34 and are changing to dashes in v0.35.
	toDashesKeys = []string{
		"db_backend", "db_dir", "filter_peers", "genesis_file", "log_format", "log_level", "node_key_file", "proxy_app",
		"consensus.create_empty_blocks", "consensus.create_empty_blocks_interval", "consensus.double_sign_check_height",
		"consensus.peer_gossip_sleep_duration", "consensus.peer_query_maj23_sleep_duration", "consensus.skip_timeout_commit",
		"consensus.timeout_commit", "consensus.timeout_precommit", "consensus.timeout_precommit_delta",
		"consensus.timeout_prevote", "consensus.timeout_prevote_delta", "consensus.timeout_propose",
		"consensus.timeout_propose_delta", "consensus.wal_file",
		"instrumentation.max_open_connections", "instrumentation.prometheus_listen_addr",
		"mempool.cache_size", "mempool.max_batch_bytes", "mempool.max_tx_bytes", "mempool.max_txs_bytes",
		"p2p.addr_book_file", "p2p.addr_book_strict", "p2p.allow_duplicate_ip", "p2p.dial_timeout", "p2p.external_address",
		"p2p.flush_throttle_timeout", "p2p.handshake_timeout", "p2p.max_num_inbound_peers", "p2p.max_num_outbound_peers",
		"p2p.max_packet_msg_payload_size", "p2p.persistent_peers", "p2p.persistent_peers_max_dial_period",
		"p2p.private_peer_ids", "p2p.recv_rate", "p2p.send_rate", "p2p.test_dial_fail", "p2p.unconditional_peer_ids",
		"rpc.cors_allowed_headers", "rpc.cors_allowed_methods", "rpc.cors_allowed_origins", "rpc.grpc_laddr",
		"rpc.grpc_max_open_connections", "rpc.max_body_bytes", "rpc.max_header_bytes", "rpc.max_open_connections",
		"rpc.max_subscription_clients", "rpc.max_subscriptions_per_client", "rpc.pprof_laddr",
		"rpc.timeout_broadcast_tx_commit", "rpc.tls_cert_file", "rpc.tls_key_file",
		"statesync.chunk_request_timeout", "statesync.discovery_time", "statesync.rpc_servers", "statesync.temp_dir",
		"statesync.trust_hash", "statesync.trust_height", "statesync.trust_period",
		"tx_index.indexer",
	}
	// changedKeys is a map of old key name to new key name for non-trivial key name changes.
	changedKeys = map[string]string{
		"fast_sync": "blocksync.enable",
		"fastsync.version": "blocksync.version",
		"priv_validator_key_file": "priv-validator.key-file",
		"priv_validator_laddr": "priv-validator.laddr",
		"priv_validator_state_file": "priv-validator.state-file",
		"p2p.seed_mode": "mode",
		"statesync.chunk_fetchers": "statesync.fetchers",
		"tx_index.psql-conn": "tx-index.psql-conn",
	}
	// allChangedKeysRev is a map of new key name to old key name for all keys that have changed.
	// This is lazily loaded when needed and contains key/value swapped entries from changedKeys
	// as well as entries for each of the toDashesKeys.
	allChangedKeysRev = map[string]string{}
)

// MigrateUnpackedTMConfigTo35IfNeeded migrates the config.toml file to v0.35 if needed.
// This assumes that the config.toml file (old or new) has already been read using viper.
// If the config.toml file is already the v0.35 version (or does not exist), this does nothing.
// If the config.toml is the old version, this will:
// 1) Migrate it to the v0.35 version:
// 2) Reloads the Tendermint config using the new v0.35 version.
func MigrateUnpackedTMConfigTo35IfNeeded(cmd *cobra.Command, vpr *viper.Viper) error {
	tmConfigPath := GetFullPathToTmConf(cmd)
	// If the file doesn't exist, there's nothing to worry about.
	if !FileExists(tmConfigPath) {
		return nil
	}
	// The new version of the config file does not have a "fast_sync" field.
	// So if that's not there, there's nothing to worry about.
	if vpr.Get("fast_sync") == nil {
		return nil
	}

	// Create a new config and field value map for it.
	tmConfig := tmconfig.DefaultConfig()
	tmConfigMap := MakeFieldValueMap(tmConfig, false)
	removeUndesirableTmConfigEntries(tmConfigMap)

	// Update the new config object using the old info that was loaded into viper.
	for newKey := range tmConfigMap {
		oldKey := getOldKey(newKey)
		if len(oldKey) == 0 {
			continue
		}
		oldValue := getValueStringFromViper(vpr, oldKey)
		newValue := getMigratedValue(oldKey, oldValue)
		err := tmConfigMap.SetFromString(newKey, newValue)
		if err != nil {
			return err
		}
	}
	writeUnpackedTMConfigFile(cmd, tmConfigPath, tmConfig, false)
	return loadUnpackedTMConfig(vpr, tmConfigPath)
}

// MigratePackedConfigToTM35IfNeeded migrates the packed-conf.json file to reflect changes in Tendermint v0.35.
// If the packed config doesn't contain any keys with changes, this does nothing.
// If the packed config does require a change, this will:
// 1) Update the packedConf map, adding/deleting keys appropriately, and updating values as needed.
// 2) Save the updated packedConf map over the existing packed-conf.json.
func MigratePackedConfigToTM35IfNeeded(cmd *cobra.Command, packedConf map[string]string) {
	needUpdate := false
	for key := range packedConf {
		newKey := getNewKey(key)
		if len(newKey) == 0 || key != newKey {
			needUpdate = true
			break
		}
	}
	if !needUpdate {
		return
	}

	tmDefaults := MakeFieldValueMap(tmconfig.DefaultConfig(), false)
	initialKeys := []string{}
	for key := range packedConf {
		initialKeys = append(initialKeys, key)
	}
	for _, oldKey := range initialKeys {
		newKey := getNewKey(oldKey)
		if len(newKey) == 0 {
			delete(packedConf, oldKey)
			continue
		}
		oldValue := packedConf[oldKey]
		newValue := getMigratedValue(oldKey, oldValue)
		defaultValue := unquote(tmDefaults.GetStringOf(newKey))
		if newValue == defaultValue {
			delete(packedConf, oldKey)
			continue
		}
		packedConf[newKey] = newValue
		if oldKey != newKey {
			delete(packedConf, oldKey)
		}
	}

	writePackedConfig(cmd, packedConf, false)
	return
}

// getNewKey converts a v0.34 key string to a v0.35 key string.
// A result of "" means the oldKey was removed in v0.35.
func getNewKey(oldKey string) string {
	if newKey, ok := changedKeys[oldKey]; ok {
		return newKey
	}
	for _, remKey := range removedKeys {
		if oldKey == remKey {
			return ""
		}
	}
	return strings.ReplaceAll(oldKey, "_", "-")
}

// getOldKey converts a v0.35 key string to a v0.34 key string.
// A result of "" means the newKey was added in v0.35 (and does not have a counterpart in v0.34).
func getOldKey(newKey string) string {
	if len(allChangedKeysRev) == 0 {
		for k, v := range changedKeys {
			allChangedKeysRev[v] = k
		}
		for _, k := range toDashesKeys {
			allChangedKeysRev[strings.ReplaceAll(k, "_", "-")] = k
		}
	}
	if oldKey, ok := allChangedKeysRev[newKey]; ok {
		return oldKey
	}
	for _, addedKey := range addedKeys {
		if newKey == addedKey {
			return ""
		}
	}
	return newKey
}

// getValueStringFromViper extracts the string representation of a key from Viper.
// This is needed for the []string entries because viper doesn't know enough just by
// reading a file, to know what type is expected for a value.
func getValueStringFromViper(vpr *viper.Viper, key string) string {
	val := vpr.Get(key)
	switch key {
	case "statesync.rpc-servers", "statesync.rpc_servers":
		// This one is in the config file as a string, but the config object as a []string.
		// The entries are comma delimited in the string.
		valStr := val.(string)
		stringVals := []string{}
		if len(valStr) > 0 {
			for _, str := range strings.Split(valStr, ",") {
				stringVals = append(stringVals, strings.TrimSpace(str))
			}
		}
		val = stringVals
	case "rpc.cors-allowed-headers", "rpc.cors-allowed-methods", "rpc.cors-allowed-origins", "tx-index.indexer",
	"rpc.cors_allowed_headers", "rpc.cors_allowed_methods", "rpc.cors_allowed_origins", "tx_index.indexer":
		// These entries are all []string in both the config object and file.
		// However, viper reads them in as []interface{}, and we need to help tell it
		// that they are []string values.
		valSlice := val.([]interface{})
		stringVals := make([]string, len(valSlice))
		for i, v := range valSlice {
			stringVals[i] = v.(string)
		}
		val = stringVals
	}
	return unquote(GetStringFromValue(reflect.ValueOf(val)))
}

// getMigratedValue converts the oldValue to a string representation of the v0.35 data type for the oldKey field.
func getMigratedValue(oldKey, oldValue string) string {
	switch oldKey {
	case "tx_index.indexer":
		// tx_index.indexer string -> tx-index.indexer []string
		if len(oldValue) == 0 {
			return "[]"
		}
		return "[\"" + oldValue + "\"]"
	case "p2p.seed_mode":
		// if seed_mode = true then mode = "seed" else mode = "full"
		if strings.EqualFold(oldValue, "true") {
			return "seed"
		}
		return "full"
	default:
		return oldValue
	}
}
