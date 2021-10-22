package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	v34config "github.com/provenance-io/provenance/cmd/provenanced/config/legacy/tendermint_0_34/config"
	// TODO: Once Tendermint v0.35 is fully pulled in, replace this with: tmconfig "github.com/tendermint/tendermint/config"
	v35config "github.com/provenance-io/provenance/cmd/provenanced/config/legacy/tendermint_0_35/config"
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

func addTxIndexPsqlConnLineToConfig(t *testing.T, path, value string) {
	oldBz, err := os.ReadFile(path)
	require.NoError(t, err, "reading config file for update")
	lineAdded := false
	var newFile strings.Builder
	for _, line := range strings.Split(string(oldBz), "\n") {
		newFile.WriteString(line)
		newFile.WriteByte('\n')
		if !lineAdded && strings.HasPrefix(line, "[tx_index]") {
			newFile.WriteString(fmt.Sprintf("\npsql-conn = \"%s\"\n", value))
			lineAdded = true
		}
	}
	require.True(t, lineAdded, "tx_index.psql-conn line added")
	require.NoError(t, os.WriteFile(path, []byte(newFile.String()), 0644))
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

func (s *ConfigMigrationsTestSuite) TestMigrateUnpackedTMConfigTo35IfNeeded() {
	dummyCmd := s.makeDummyCmd()
	confFile := GetFullPathToTmConf(dummyCmd)

	touchConf := func(t *testing.T) {
		_, err := os.Stat(confFile)
		switch {
		case os.IsNotExist(err):
			require.NoError(t, os.WriteFile(confFile, []byte{}, 0644), "creating empty config file")
		case err != nil:
			require.NoError(t, err, "conf file stat")
		default:
			require.NoError(t, os.Chtimes(confFile, time.Now(), time.Now()), "chtimes config file")
		}
	}
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "config file does not exist",
			test: func(t *testing.T) {
				vpr := viper.New()
				require.NoError(t, MigrateUnpackedTMConfigTo35IfNeeded(dummyCmd, vpr), "calling MigrateUnpackedTMConfigTo35IfNeeded")
				confFileExists := FileExists(confFile)
				assert.False(t, confFileExists, "confFileExists")
			},
		},
		{
			name: "config is new version",
			test: func(t *testing.T) {
				conf := v35config.DefaultConfig()
				conf.LogFormat = "json"
				conf.LogLevel = "debug"
				conf.Moniker = t.Name()

				// Write an initial config and load it into viper.
				v35config.WriteConfigFile(s.Home, conf)
				vpr := viper.New()
				vpr.SetConfigFile(confFile)
				require.NoError(t, vpr.ReadInConfig(), "reading config into viper")

				// Now, delete the config file.
				// This makes it easy to tell if MigrateUnpackedTMConfigTo35IfNeeded thought it needed migrating.
				require.NoError(t, deleteConfigFile(dummyCmd, confFile, false), "deleting config file")

				require.NoError(t, MigrateUnpackedTMConfigTo35IfNeeded(dummyCmd, vpr), "calling MigrateUnpackedTMConfigTo35IfNeeded")

				confExists := FileExists(confFile)
				assert.False(t, confExists, "confExists")
			},
		},
		{
			name: "viper has just fast_sync - file is created",
			test: func(t *testing.T) {
				touchConf(t)
				vpr := viper.New()
				vpr.Set("fast_sync", false)

				require.NoError(t, MigrateUnpackedTMConfigTo35IfNeeded(dummyCmd, vpr), "calling MigrateUnpackedTMConfigTo35IfNeeded")

				confExists := FileExists(confFile)
				assert.True(t, confExists, "confExists")
			},
		},
		{
			name: "viper has just fast_sync - ends up with other stuff",
			test: func(t *testing.T) {
				shouldHaveKeys := []string{
					"blocksync.enable", "mode", "consensus.create-empty-blocks", "mempool.ttl-duration", "moniker",
				}

				shouldNotHaveKeys := []string{
					"p2p.seed_mode", "consensus.create_empty_blocks", "mempool.ttl_duration",
				}

				touchConf(t)
				vpr := viper.New()
				vpr.Set("fast_sync", true)

				require.NoError(t, MigrateUnpackedTMConfigTo35IfNeeded(dummyCmd, vpr), "calling MigrateUnpackedTMConfigTo35IfNeeded")

				for _, key := range shouldHaveKeys {
					actual := vpr.Get(key)
					assert.NotNil(t, actual, "vpr.Get(\"%s\")", key)
				}

				for _, key := range shouldNotHaveKeys {
					actual := vpr.Get(key)
					assert.Nil(t, actual, "vpr.Get(\"%s\")", key)
				}
			},
		},
		{
			name: "special cases migrate expectedly",
			test: func(t *testing.T) {
				expected := map[string]interface{}{
					"blocksync.enable": true,
					"blocksync.version": "v8",
					"priv-validator.key-file": "/some/priv-key",
					"priv-validator.laddr": "127.0.0.1:888",
					"priv-validator.state-file": "/some/state-file",
					"mode": "seed",
					"statesync.fetchers": int32(5),
					"tx-index.psql-conn": "127.0.0.1:5432",
					"tx-index.indexer": []string{"yellow"},
				}
				expectedKeys := make([]string, 0, len(expected))
				for key := range expected {
					expectedKeys = append(expectedKeys, key)
				}
				sortKeys(expectedKeys)

				checkConf := func(label string, conf *v35config.Config) {
					assert.Equal(t, expected["blocksync.enable"], conf.BlockSync.Enable, "%s blocksync.enable", label)
					assert.Equal(t, expected["blocksync.version"], conf.BlockSync.Version, "%s blocksync.version", label)
					assert.Equal(t, expected["priv-validator.key-file"], conf.PrivValidator.Key, "%s priv-validator.key-file", label)
					assert.Equal(t, expected["priv-validator.laddr"], conf.PrivValidator.ListenAddr, "%s priv-validator.laddr", label)
					assert.Equal(t, expected["priv-validator.state-file"], conf.PrivValidator.State, "%s priv-validator.state-file", label)
					assert.Equal(t, expected["mode"], conf.Mode, "%s mode", label)
					assert.Equal(t, expected["statesync.fetchers"], conf.StateSync.Fetchers, "%s statesync.fetchers", label)
					assert.Equal(t, expected["tx-index.psql-conn"], conf.TxIndex.PsqlConn, "%s tx-index.psql-conn", label)
					assert.Equal(t, expected["tx-index.indexer"], conf.TxIndex.Indexer, "%s tx-index.indexer", label)
				}

				oldConf := v34config.DefaultConfig()
				oldConf.FastSyncMode = expected["blocksync.enable"].(bool)
				oldConf.FastSync.Version = expected["blocksync.version"].(string)
				oldConf.PrivValidatorKey = expected["priv-validator.key-file"].(string)
				oldConf.PrivValidatorListenAddr = expected["priv-validator.laddr"].(string)
				oldConf.PrivValidatorState = expected["priv-validator.state-file"].(string)
				oldConf.P2P.SeedMode = true
				oldConf.StateSync.ChunkFetchers = expected["statesync.fetchers"].(int32)
				oldConf.TxIndex.PsqlConn = expected["tx-index.psql-conn"].(string)
				oldConf.TxIndex.Indexer = expected["tx-index.indexer"].([]string)[0]
				v34config.WriteConfigFile(confFile, oldConf)
				addTxIndexPsqlConnLineToConfig(t, confFile, oldConf.TxIndex.PsqlConn)

				vpr := viper.New()
				vpr.SetConfigFile(confFile)
				require.NoError(t, vpr.ReadInConfig(), "reading config into viper")
				require.NotNil(t, vpr.Get("fast_sync"), "setup error: config loaded into viper does ot have a fast_sync entry")

				require.NoError(t, MigrateUnpackedTMConfigTo35IfNeeded(dummyCmd, vpr), "calling MigrateUnpackedTMConfigTo35IfNeeded")

				freshConf := v35config.DefaultConfig()
				require.NoError(t, vpr.Unmarshal(freshConf), "unmarshalling fresh conf")
				checkConf("immediately after migrate", freshConf)

				// Clear out viper, and reload the file to make sure the file is the new version.
				vpr = viper.New()
				vpr.SetConfigFile(confFile)
				require.NoError(t, vpr.ReadInConfig(), "reading config into viper")
				require.Nil(t, vpr.Get("fast_sync"), "vpr.Get(\"fast_sync\")")

				newConf := v35config.DefaultConfig()
				require.NoError(t, vpr.Unmarshal(newConf), "unmarshalling new conf")
				checkConf("after file load", newConf)
			},
		},
	}

	s.Require().NoError(EnsureConfigDir(dummyCmd), "ensuring config dir")
	s.Require().NoError(deleteConfigFile(dummyCmd, confFile, false), "deleting config file at start")
	for _, tc := range tests {
		s.T().Run(tc.name, tc.test)
		s.Require().NoError(deleteConfigFile(dummyCmd, confFile, false), "deleting config file after %s", tc.name)
	}
}

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
				"log_format": "json",
				"log_level": "debug",
			},
			expected: map[string]string{
				"log-format": "json",
				"log-level": "debug",
			},
		},
		{
			name:     "all special cases",
			conf:     map[string]string{
				"fast_sync": "false",
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
				"blocksync.enable": "false",
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
		{
			name:     "default value removed",
			conf:     map[string]string{
				"log_format": v35config.DefaultConfig().LogFormat,
				"log_level": "debug",
			},
			expected: map[string]string{
				"log-level": "debug",
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
