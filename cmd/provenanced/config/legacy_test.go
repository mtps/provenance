package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	v34config "github.com/provenance-io/provenance/cmd/provenanced/config/legacy/tendermint_0_34/config"
	// TODO: Once Tendermint v0.35 is fully pulled in, replace this with: tmconfig "github.com/tendermint/tendermint/config"
	v35config "github.com/provenance-io/provenance/cmd/provenanced/config/legacy/tendermint_0_35/config"
)

type LegacyTestSuite struct {
	suite.Suite

	// Home is a temp directory that can be used to store files for a test.
	Home string
}

func TestLegacyTestSuite(t *testing.T) {
	suite.Run(t, new(LegacyTestSuite))
}

func (s *LegacyTestSuite) SetupTest() {
	s.Home = s.T().TempDir()
	s.T().Logf("%s Home: %s", s.T().Name(), s.Home)
}

func (s *LegacyTestSuite) convertViperValToString(key string, val interface{}) string {
	switch key {
	case "statesync.rpc-servers", "statesync.rpc_servers":
		// This one is in the config file as a string, but the config object as a []string.
		// The entries are comma delimited in the string.
		if val == nil {
			return ""
		}
		valStr, ok := val.(string)
		if !ok {
			s.Require().NoError(fmt.Errorf("field [%s]: interface conversion: interface {} is %T, not string", key, val))
		}
		stringVals := []string{}
		if len(valStr) > 0 {
			for _, str := range strings.Split(valStr, ",") {
				stringVals = append(stringVals, strings.TrimSpace(str))
			}
		}
		val = stringVals
	case "rpc.cors-allowed-headers", "rpc.cors-allowed-methods", "rpc.cors-allowed-origins", "tx-index.indexer",
		"rpc.cors_allowed_headers", "rpc.cors_allowed_methods", "rpc.cors_allowed_origins":
		// These entries are all []string in both the config object and file.
		// However, viper reads them in as []interface{}, and we need to help tell it
		// that they are []string values.
		if val == nil {
			return "[]"
		}
		stringVals, ok := val.([]string)
		if !ok {
			valSlice, ok := val.([]interface{})
			if !ok {
				s.Require().NoError(fmt.Errorf("field [%s]: interface conversion: interface {} is %T, not []string or []interface {}", key, val))
			}
			stringVals = make([]string, len(valSlice))
			for i, v := range valSlice {
				stringVals[i] = v.(string)
			}
		}
		val = stringVals
	}
	return unquote(GetStringFromValue(reflect.ValueOf(val)))
}

func (s *LegacyTestSuite) TestPrintChangesBetween34and35() {
	v34 := v34config.DefaultConfig()
	v35 := v35config.DefaultConfig()

	knownChanges34To35 := map[string]string{
		"fast_sync": "blocksync.enable",
		"fastsync.version": "blocksync.version",
		"priv_validator_key_file": "priv-validator.key-file",
		"priv_validator_laddr": "priv-validator.laddr",
		"priv_validator_state_file": "priv-validator.state-file",
		"p2p.seed_mode": "mode",
		"statesync.chunk_fetchers": "statesync.fetchers",
		"tx_index.psql-conn": "tx-index.psql-conn",
	}
	knownChanges34 := []string{}
	knownChanges35To34 := map[string]string{}
	for k34, k35 := range knownChanges34To35 {
		knownChanges34 = append(knownChanges34, k34)
		knownChanges35To34[k35] = k34
	}
	sortKeys(knownChanges34)

	v34Map := MakeFieldValueMap(v34, true)
	v35Map := MakeFieldValueMap(v35, true)

	for _, k34 := range knownChanges34 {
		k35 := knownChanges34To35[k34]
		_, ok34 := v34Map[k34]
		_, ok35 := v35Map[k35]
		s.Assert().True(ok34, "known change v0.34 key [%s] not found", k34)
		s.Assert().True(ok35, "known change v0.35 key [%s] not found", k35)
	}

	v34Types := map[string]string{}
	v35Types := map[string]string{}

	unchanged := []string{}
	added := []string{}
	removed := []string{}
	toDashes := []string{}
	asDashes := []string{}

	stringsContains := func(vals []string, lookFor string) bool {
		for _, val := range vals {
			if val == lookFor {
				return true
			}
		}
		return false
	}

	for key34 := range v34Map {
		v34Types[key34] = v34Map[key34].Type().String()
		if _, ok := knownChanges34To35[key34]; ok {
			continue
		}
		if _, ok := v35Map[key34]; ok {
			unchanged = append(unchanged, key34)
			continue
		}
		key35 := strings.ReplaceAll(key34, "_", "-")
		if _, ok := v35Map[key35]; ok {
			toDashes = append(toDashes, key34)
			asDashes = append(asDashes, key35)
		} else {
			removed = append(removed, key34)
		}
	}

	for key35 := range v35Map {
		v35Types[key35] = v35Map[key35].Type().String()
		if _, ok := knownChanges35To34[key35]; ok {
			continue
		}
		if _, ok := v34Map[key35]; ok {
			continue
		}
		if stringsContains(asDashes, key35) {
			continue
		}
		added = append(added, key35)
	}

	sortKeys(unchanged)
	sortKeys(added)
	sortKeys(removed)
	sortKeys(toDashes)

	toV35Key := func(key34 string) string {
		if key35, ok := knownChanges34To35[key34]; ok {
			return key35
		}
		if stringsContains(removed, key34) {
			return ""
		}
		return strings.ReplaceAll(key34, "_", "-")
	}

	toCompareTypes := []string{}
	toCompareTypes = append(toCompareTypes, knownChanges34...)
	toCompareTypes = append(toCompareTypes, unchanged...)
	toCompareTypes = append(toCompareTypes, toDashes...)
	sortKeys(toCompareTypes)
	typeChanges := []string{}
	for _, key34 := range toCompareTypes {
		key35 := toV35Key(key34)
		if len(key35) == 0 {
			continue
		}
		type34 := v34Types[key34]
		type35 := v35Types[key35]
		if type34 != type35 {
			typeChanges = append(typeChanges, fmt.Sprintf("%s %s -> %s %s", key34, type34, key35, type35))
		}
	}

	knownChanges := make([]string, len(knownChanges34))
	for i, key34 := range knownChanges34 {
		knownChanges[i] = fmt.Sprintf("%s -> %s", key34, knownChanges34To35[key34])
	}
	dashChanges := make([]string, len(toDashes))
	for i, key34 := range toDashes {
		dashChanges[i] = fmt.Sprintf("%s -> %s", key34, strings.ReplaceAll(key34, "_", "-"))
	}

	printStrings := func(header string, vals []string) {
		fmt.Printf("%s (%d):\n", header, len(vals))
		for _, val := range vals {
			fmt.Printf("  %s\n", val)
		}
		fmt.Printf("\n")
	}

	printStrings("unchanged", unchanged)
	printStrings("added", added)
	printStrings("removed", removed)
	printStrings("dash changes", dashChanges)
	printStrings("non-trivial changes", knownChanges)
	printStrings("type changes", typeChanges)

	printStringsAsVar := func(varName string, vals []string) {
		fmt.Printf("var %s = []string{\n", varName)
		fmt.Printf("\t\"%s\"\n", strings.Join(vals, `", "`))
		fmt.Printf("}\n")
	}
	printStringsAsVar("addedKeys", added)
	printStringsAsVar("removedKeys", removed)
	printStringsAsVar("toDashesKeys", toDashes)
	fmt.Printf("var changedKeys = map[string]string{\n")
	for k, v := range knownChanges34To35 {
		fmt.Printf("\t\"%s\": \"%s\"\n", k, v)
	}
	fmt.Printf("}\n")
}

func (s *LegacyTestSuite) TestPrintDefaultConfigAndTypes34() {
	conf := v34config.DefaultConfig()
	confMap := MakeFieldValueMap(conf, false)
	removeUndesirableTmConfigEntries(confMap)
	confKeys := confMap.GetSortedKeys()

	byType := map[string][]string{}

	for _, key := range confKeys {
		val := confMap.GetStringOf(key)
		valType := fmt.Sprintf("%T", confMap[key].Interface())
		fmt.Printf("%s %s = %s\n", key, valType, val)
		byType[valType] = append(byType[valType], fmt.Sprintf("%s = %s", key, val))
	}
	fmt.Printf("\n")

	valTypes := []string{}
	for valType := range byType {
		valTypes = append(valTypes, valType)
	}
	sortKeys(valTypes)

	for _, valType := range valTypes {
		fmt.Printf("%s entries (%d):\n", valType, len(byType[valType]))
		for _, entry := range byType[valType] {
			fmt.Printf("\t%s\n", entry)
		}
	}
	fmt.Printf("\n")
}

func (s *LegacyTestSuite) TestPrintDefaultConfigAndTypes35() {
	conf := v35config.DefaultConfig()
	confMap := MakeFieldValueMap(conf, false)
	removeUndesirableTmConfigEntries(confMap)
	confKeys := confMap.GetSortedKeys()

	byType := map[string][]string{}

	for _, key := range confKeys {
		val := confMap.GetStringOf(key)
		valType := fmt.Sprintf("%T", confMap[key].Interface())
		fmt.Printf("%s %s = %s\n", key, valType, val)
		byType[valType] = append(byType[valType], fmt.Sprintf("%s = %s", key, val))
	}
	fmt.Printf("\n")

	valTypes := []string{}
	for valType := range byType {
		valTypes = append(valTypes, valType)
	}
	sortKeys(valTypes)

	for _, valType := range valTypes {
		fmt.Printf("%s entries (%d):\n", valType, len(byType[valType]))
		for _, entry := range byType[valType] {
			fmt.Printf("\t%s\n", entry)
		}
	}
	fmt.Printf("\n")
}

func (s *LegacyTestSuite) TestCompareConfigToFileEntries34() {
	confDir := filepath.Join(s.Home, "config")
	s.Require().NoError(os.MkdirAll(confDir, os.ModePerm), "creating config directory")

	expectedNotInFile := []string{"tx_index.psql-conn"}

	v34Config := v34config.DefaultConfig()
	confFile := filepath.Join(confDir, "config.toml")
	v34config.WriteConfigFile(confFile, v34Config)

	vpr := viper.New()
	vpr.SetConfigFile(confFile)
	s.Require().NoError(vpr.ReadInConfig(), "reading config into viper")

	v34ConfigObjMap := MakeFieldValueMap(v34Config, true)
	removeUndesirableTmConfigEntries(v34ConfigObjMap)
	objSettings := map[string]string{}
	objKeys := []string{}
	for key := range v34ConfigObjMap {
		objKeys = append(objKeys, key)
		objSettings[key] = unquote(v34ConfigObjMap.GetStringOf(key))
	}

	fileKeys := vpr.AllKeys()
	sortKeys(fileKeys)
	fileSettings := map[string]string{}
	for _, key := range fileKeys {
		fileSettings[key] = s.convertViperValToString(key, vpr.Get(key))
	}

	inObjNotFile := []string{}
	inFileNotObj := []string{}
	different := []string{}

	for _, key := range objKeys {
		fileValue, ok := fileSettings[key]
		if !ok {
			inObjNotFile = append(inObjNotFile, key)
			continue
		}
		objValue := objSettings[key]
		if fileValue != objValue {
			different = append(different, fmt.Sprintf("%s: (%s) != (%s)", key, objValue, fileValue))
		}
	}

	for _, key := range fileKeys {
		if _, ok := objSettings[key]; !ok {
			inFileNotObj = append(inFileNotObj, key)
		}
	}

	s.Assert().Equal(inObjNotFile, expectedNotInFile, "In object but not file")
	s.Assert().Len(inFileNotObj, 0, "In file but not object")
	s.Assert().Len(different, 0, "Different")
}

func (s *LegacyTestSuite) TestCompareConfigToFileEntries35() {
	confDir := filepath.Join(s.Home, "config")
	s.Require().NoError(os.MkdirAll(confDir, os.ModePerm), "creating config directory")

	v35Config := v35config.DefaultConfig()
	v35config.WriteConfigFile(s.Home, v35Config)
	confFile := filepath.Join(confDir, "config.toml")

	vpr := viper.New()
	vpr.SetConfigFile(confFile)
	s.Require().NoError(vpr.ReadInConfig(), "reading config into viper")

	v35ConfigObjMap := MakeFieldValueMap(v35Config, true)
	removeUndesirableTmConfigEntries(v35ConfigObjMap)
	objSettings := map[string]string{}
	objKeys := []string{}
	for key := range v35ConfigObjMap {
		objKeys = append(objKeys, key)
		objSettings[key] = unquote(v35ConfigObjMap.GetStringOf(key))
	}

	fileKeys := vpr.AllKeys()
	sortKeys(fileKeys)
	fileSettings := map[string]string{}
	for _, key := range fileKeys {
		fileSettings[key] = s.convertViperValToString(key, vpr.Get(key))
	}

	inObjNotFile := []string{}
	inFileNotObj := []string{}
	different := []string{}

	for _, key := range objKeys {
		fileValue, ok := fileSettings[key]
		if !ok {
			inObjNotFile = append(inObjNotFile, key)
			continue
		}
		objValue := objSettings[key]
		if fileValue != objValue {
			different = append(different, fmt.Sprintf("%s: (%s) != (%s)", key, objValue, fileValue))
		}
	}

	for _, key := range fileKeys {
		if _, ok := objSettings[key]; !ok {
			inFileNotObj = append(inFileNotObj, key)
		}
	}

	s.Assert().Len(inObjNotFile, 0, "In object but not file")
	s.Assert().Len(inFileNotObj, 0, "In file but not object")
	s.Assert().Len(different, 0, "Different")
}

func (s *LegacyTestSuite) TestRead34FileWith35Struct() {
	v34 := v34config.DefaultConfig()
	confFile := filepath.Join(s.Home, "config.toml")
	v34config.WriteConfigFile(confFile, v34)

	vpr := viper.New()
	vpr.SetConfigFile(confFile)
	err := vpr.ReadInConfig()
	s.Require().NoError(err, "reading config into viper")

	v35 := v35config.DefaultConfig()
	err = vpr.Unmarshal(v35)
	s.Require().NoError(err, "unmarshaling conf from viper")

	otherKeys := make([]string, 0, len(v35.Other))
	for key := range v35.Other {
		otherKeys = append(otherKeys, key)
	}
	sortKeys(otherKeys)
	for _, key := range otherKeys {
		val := v35.Other[key]
		fmt.Printf("%s: %#v\n", key, val)
	}
	s.Assert().Len(otherKeys, 14, "other keys")
}
