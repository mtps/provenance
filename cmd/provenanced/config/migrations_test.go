package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ConfigMigrationsTestSuite struct {
	suite.Suite

	// Home is a temp directory that can be used to store files for a test.
	// It is different for each test function.
	Home string
}

func TestConfigMigrationsTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigMigrationsTestSuite))
}

func (s *ConfigMigrationsTestSuite) SetupTest() {
	s.Home = s.T().TempDir()
	s.T().Logf("%s Home: %s", s.T().Name(), s.Home)
}

func (s *ConfigMigrationsTestSuite) makeDummyCmd() *cobra.Command {
	dummyCmd, err := makeDummyCmd(s.Home)
	s.Require().NoError(err, "dummy command setup")
	return dummyCmd
}

func (s *ConfigMigrationsTestSuite) TestUniqueKeyEntries() {
	// This test makes sure that a key is only listed once among the
	// addedKeys, removedKeys, toDashes, and changedKeys variables.
	keySources := make(map[string][]string)
	addKey := func(key, source string) {
		keySources[key] = append(keySources[key], source)
	}
	for _, key := range addedKeys {
		addKey(key, "addedKeys")
	}
	for _, key := range removedKeys {
		addKey(key, "removedKeys")
	}
	for _, key := range toDashesKeys {
		addKey(key, "toDashesKeys")
	}
	for oldKey, newKey := range changedKeys {
		addKey(oldKey, "changedKeys-old")
		addKey(newKey, "changedKeys-new")
	}
	for key, sources := range keySources {
		s.Assert().Len(sources, 1, key)
	}
}

// TODO: Test MigrateUnpackedTMConfigTo35IfNeeded

func (s *ConfigMigrationsTestSuite) TestMigratePackedConfigToTM35IfNeeded() {
	dummyCmd := s.makeDummyCmd()
	s.Require().NoError(EnsureConfigDir(dummyCmd), "ensuring config dir")

	copyConf := func(conf map[string]string) map[string]string {
		rv := map[string]string{}
		for k, v := range conf {
			rv[k] = v
		}
		return rv
	}

	s.T().Run("empty", func(t *testing.T) {
		conf := map[string]string{}
		orig := copyConf(conf)
		MigratePackedConfigToTM35IfNeeded(dummyCmd, conf)
		packedFileExists := FileExists(GetFullPathToPackedConf(dummyCmd))
		assert.False(t, packedFileExists, "packedFileExists")
		assert.Equal(t, orig, conf, "pre and post configs")
	})

	s.T().Run("without changing entries", func(t *testing.T) {
		conf := map[string]string{
			// Two entries from the client config.
			"chain-id": "testing",
			"output": "json",
			// Two entries from the app config.
			"grpc-web.enable-unsafe-cors": "true",
			"rosetta.retries": "8",
			// Two unchanging entries from tendermint config.
			"rpc.unsafe": "true",
			"mempool.size": "6000",
		}
		orig := copyConf(conf)
		MigratePackedConfigToTM35IfNeeded(dummyCmd, conf)
		packedFileExists := FileExists(GetFullPathToPackedConf(dummyCmd))
		assert.False(t, packedFileExists, "packedFileExists")
		assert.Equal(t, orig, conf, "pre and post configs")
	})

	s.T().Run("only an added key", func(t *testing.T) {
		conf := map[string]string{
			"p2p.max-connections": "5",
		}
		orig := copyConf(conf)
		MigratePackedConfigToTM35IfNeeded(dummyCmd, conf)
		packedFileExists := FileExists(GetFullPathToPackedConf(dummyCmd))
		assert.False(t, packedFileExists, "packedFileExists")
		assert.Equal(t, orig, conf, "pre and post configs")
	})

	// Just to make sure it's not there.
	s.Require().NoError(deletePackedConfig(dummyCmd, false), "deleting packed config")

	tests := []struct{
		name string
		conf map[string]string
		expected map[string]string
	}{
		{
			name:     "removed key is removed",
			conf:     map[string]string{
				"mempool.wal_dir": "/a/b/c",
			},
			expected: map[string]string{},
		},
		{
			name:     "underscores to dashes",
			conf:     map[string]string{
				"log_format": "plain",
				"log_level": "debug",
			},
			expected: map[string]string{
				"log-format": "plain",
				"log-level": "debug",
			},
		},
		{
			name:     "all special cases",
			conf:     map[string]string{
				"fast_sync": "true",
				"fastsync.version": "v8",
				"priv_validator_key_file": "/some/priv-key",
				"priv_validator_laddr": "127.0.0.1:888",
				"priv_validator_state_file": "/some/state-file",
				"p2p.seed_mode": "true",
				"statesync.chunk_fetchers": "5",
				"tx_index.psql-conn": "127.0.0.1:5432",
				"tx_index.indexer": "yellow",
			},
			expected: map[string]string{
				"blocksync.enable": "true",
				"blocksync.version": "v8",
				"priv-validator.key-file": "/some/priv-key",
				"priv-validator.laddr": "127.0.0.1:888",
				"priv-validator.state-file": "/some/state-file",
				"mode": "seed",
				"statesync.fetchers": "5",
				"tx-index.psql-conn": "127.0.0.1:5432",
				"tx-index.indexer": `["yellow"]`,
			},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			MigratePackedConfigToTM35IfNeeded(dummyCmd, tc.conf)
			assert.Equal(t, tc.expected, tc.conf, "post-migrate config")
			confFile := GetFullPathToPackedConf(dummyCmd)
			jsonFromFile, err := os.ReadFile(confFile)
			require.NoError(t, err, "reading packed config file.")
			confFromFile := map[string]string{}
			err = json.Unmarshal(jsonFromFile, &confFromFile)
			require.NoError(t, err, "unmarshalling packed config json")
			assert.Equal(t, tc.expected, confFromFile, "config from file")
		})
		s.Require().NoError(deletePackedConfig(dummyCmd, false), "deleting packed config after %s", tc.name)
	}
}

func (s *ConfigMigrationsTestSuite) TestGetNewKey() {
	// Just some spot checking here.
	tests := []struct {
		oldKey string
		expected string
	}{
		// A removed key
		{
			oldKey:   "mempool.wal_dir",
			expected: "",
		},
		// A to-dashed key
		{
			oldKey:   "log_format",
			expected: "log-format",
		},
		// A couple non-trivial change keys
		{
			oldKey:   "priv_validator_key_file",
			expected: "priv-validator.key-file",
		},
		{
			oldKey:   "fastsync.version",
			expected: "blocksync.version",
		},
		// An unchanged key
		{
			oldKey:   "instrumentation.prometheus",
			expected: "instrumentation.prometheus",
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.oldKey, func(t *testing.T) {
			actual := getNewKey(tc.oldKey)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func (s *ConfigMigrationsTestSuite) TestGetOldKey() {
	// Just some spot checking here.
	tests := []struct {
		newKey string
		expected string
	} {
		// An added key
		{
			newKey:   "mempool.ttl-duration",
			expected: "",
		},
		// A to-dashed key
		{
			newKey:   "instrumentation.max-open-connections",
			expected: "instrumentation.max_open_connections",
		},
		// A couple non-trivial change keys
		{
			newKey:   "tx-index.psql-conn",
			expected: "tx_index.psql-conn",
		},
		{
			newKey:   "blocksync.enable",
			expected: "fast_sync",
		},
		// An unchanged key
		{
			newKey:   "rpc.unsafe",
			expected: "rpc.unsafe",
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.newKey, func(t *testing.T) {
			actual := getOldKey(tc.newKey)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func (s *ConfigMigrationsTestSuite) TestGetValueStringFromViper() {

	type testCase struct {
		name string
		value interface{}
		expected string
	}
	tests := []struct {
		keys []string
		cases []testCase
	}{
		{
			keys:  []string{"statesync.rpc-servers", "statesync.rpc_servers"},
			cases: []testCase{
				{
					name:     "empty string",
					value:    "",
					expected: `[]`,
				},
				{
					name:     "one entry",
					value:    "banana",
					expected: `["banana"]`,
				},
				{
					name:     "four entries",
					value:    "apple, orange, peach, grape",
					expected: `["apple", "orange", "peach", "grape"]`,
				},
			},
		},
		{
			keys:  []string{
				"rpc.cors-allowed-headers", "rpc.cors-allowed-methods", "rpc.cors-allowed-origins", "tx-index.indexer",
				"rpc.cors_allowed_headers", "rpc.cors_allowed_methods", "rpc.cors_allowed_origins", "tx_index.indexer",
			},
			cases: []testCase{
				{
					name:     "empty",
					value:    []interface{}{},
					expected: `[]`,
				},
				{
					name:     "one entry",
					value:    []interface{}{"penny"},
					expected: `["penny"]`,
				},
				{
					name:     "four entries",
					value:    []interface{}{"nickle", "dime", "quarter", "looney"},
					expected: `["nickle", "dime", "quarter", "looney"]`,
				},
			},
		},
		{
			keys:  []string{"non-special", "normal"},
			cases: []testCase{
				{
					name:     "empty string",
					value:    "",
					expected: "",
				},
				{
					name:     "filled string",
					value:    "filled",
					expected: "filled",
				},
				{
					name:     "number",
					value:    88,
					expected: "88",
				},
				{
					name:     "bool true",
					value:    true,
					expected: "true",
				},
				{
					name:     "bool false",
					value:    false,
					expected: "false",
				},
				{
					name:     "time.Duration",
					value:    time.Duration(800) * time.Millisecond,
					expected: "800ms",
				},
			},
		},
	}

	for _, tg := range tests {
		for _, key := range tg.keys {
			for _, tc := range tg.cases {
				s.T().Run(fmt.Sprintf("%s %s", key, tc.name), func(t *testing.T) {
					vpr := viper.New()
					vpr.Set(key, tc.value)
					actual := getValueStringFromViper(vpr, key)
					assert.Equal(t, tc.expected, actual)
				})
			}
		}
	}
}

func (s *ConfigMigrationsTestSuite) TestGetMigratedValue() {
	tests := []struct {
		name string
		oldKey string
		oldValue string
		expected string
	}{
		{
			name:     "tx_index.indexer empty",
			oldKey:   "tx_index.indexer",
			oldValue: "",
			expected: "[]",
		},
		{
			name:     "tx_index.indexer with value",
			oldKey:   "tx_index.indexer",
			oldValue: "value1",
			expected: `["value1"]`,
		},
		{
			name:     "p2p.seed_mode true",
			oldKey:   "p2p.seed_mode",
			oldValue: "true",
			expected: "seed",
		},
		{
			name:     "p2p.seed_mode false",
			oldKey:   "p2p.seed_mode",
			oldValue: "false",
			expected: "full",
		},
		{
			name:     "p2p.seed_mode TRUE",
			oldKey:   "p2p.seed_mode",
			oldValue: "TRUE",
			expected: "seed",
		},
		{
			name:     "p2p.seed_mode FALSE",
			oldKey:   "p2p.seed_mode",
			oldValue: "FALSE",
			expected: "full",
		},
		{
			name:     "non-special key empty value",
			oldKey:   "non-special",
			oldValue: "",
			expected: "",
		},
		{
			name:     "non-special key with value",
			oldKey:   "non-special-key",
			oldValue: "non-special value",
			expected: "non-special value",
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			actual := getMigratedValue(tc.oldKey, tc.oldValue)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
