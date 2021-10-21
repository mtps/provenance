package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
// TODO: Test MigratePackedConfigToTM35IfNeeded

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
