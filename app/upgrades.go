package app

import (
	"fmt"
	"strings"

	icqtypes "github.com/strangelove-ventures/async-icq/v6/types"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	attributekeeper "github.com/provenance-io/provenance/x/attribute/keeper"
	attributetypes "github.com/provenance-io/provenance/x/attribute/types"
	"github.com/provenance-io/provenance/x/exchange"
	"github.com/provenance-io/provenance/x/hold"
	ibchookstypes "github.com/provenance-io/provenance/x/ibchooks/types"
	"github.com/provenance-io/provenance/x/marker/types"
	msgfeetypes "github.com/provenance-io/provenance/x/msgfees/types"
	oracletypes "github.com/provenance-io/provenance/x/oracle/types"
	triggertypes "github.com/provenance-io/provenance/x/trigger/types"
)

// appUpgrade is an internal structure for defining all things for an upgrade.
type appUpgrade struct {
	// Added contains names of modules being added during an upgrade.
	Added []string
	// Deleted contains names of modules being removed during an upgrade.
	Deleted []string
	// Renamed contains info on modules being renamed during an upgrade.
	Renamed []storetypes.StoreRename
	// Handler is a function to execute during an upgrade.
	Handler func(sdk.Context, *App, module.VersionMap) (module.VersionMap, error)
}

// upgrades is where we define things that need to happen during an upgrade.
// If no Handler is defined for an entry, a no-op upgrade handler is still registered.
// If there's nothing that needs to be done for an upgrade, there still needs to be an
// entry in this map, but it can just be {}.
//
// On the same line as the key, there should be a comment indicating the software version.
// Entries currently in use (e.g. on mainnet or testnet) cannot be deleted.
// Entries should be in chronological order, earliest first. E.g. quicksilver-rc1 went to
// testnet first, then quicksilver-rc2 went to testnet, then quicksilver went to mainnet.
//
// If something is happening in the rc upgrade(s) that isn't being applied in the non-rc,
// or vice versa, please add comments explaining why in both entries.
var upgrades = map[string]appUpgrade{
	"rust-rc1": { // upgrade for v1.16.0-rc1
		Handler: func(ctx sdk.Context, app *App, vm module.VersionMap) (module.VersionMap, error) {
			var err error
			vm, err = runModuleMigrations(ctx, app, vm)
			if err != nil {
				return nil, err
			}

			removeInactiveValidatorDelegations(ctx, app)

			err = setAccountDataNameRecord(ctx, app.AccountKeeper, &app.NameKeeper)
			if err != nil {
				return nil, err
			}

			// We only need to call addGovV1SubmitFee on testnet.
			addGovV1SubmitFee(ctx, app)

			removeP8eMemorializeContractFee(ctx, app)

			fixNameIndexEntries(ctx, app)

			return vm, nil
		},
		Added: []string{triggertypes.ModuleName},
	},
	"rust": { // upgrade for v1.16.0
		Handler: func(ctx sdk.Context, app *App, vm module.VersionMap) (module.VersionMap, error) {
			var err error
			vm, err = runModuleMigrations(ctx, app, vm)
			if err != nil {
				return nil, err
			}

			removeInactiveValidatorDelegations(ctx, app)

			err = setAccountDataNameRecord(ctx, app.AccountKeeper, &app.NameKeeper)
			if err != nil {
				return nil, err
			}

			// No need to call addGovV1SubmitFee in here as mainnet already has it defined.

			removeP8eMemorializeContractFee(ctx, app)

			fixNameIndexEntries(ctx, app)

			return vm, nil
		},
		Added: []string{triggertypes.ModuleName},
	},
	"saffron-rc1": { // upgrade for v1.17.0-rc1
		Handler: func(ctx sdk.Context, app *App, vm module.VersionMap) (module.VersionMap, error) {
			var err error
			vm, err = runModuleMigrations(ctx, app, vm)
			if err != nil {
				return nil, err
			}
			// set ibchoooks defaults (no allowed async contracts)
			app.IBCHooksKeeper.SetParams(ctx, ibchookstypes.DefaultParams())

			removeInactiveValidatorDelegations(ctx, app)
			setupICQ(ctx, app)
			updateMaxSupply(ctx, app)
			setExchangeParams(ctx, app)

			return vm, nil
		},
		Added: []string{icqtypes.ModuleName, oracletypes.ModuleName, ibchookstypes.StoreKey, hold.ModuleName, exchange.ModuleName},
	},
	"saffron-rc2": { // upgrade for v1.17.0-rc2
		Handler: func(ctx sdk.Context, app *App, vm module.VersionMap) (module.VersionMap, error) {
			var err error
			vm, err = runModuleMigrations(ctx, app, vm)
			if err != nil {
				return nil, err
			}

			updateIbcMarkerDenomMetadata(ctx, app)

			return vm, nil
		},
	},
	"saffron-rc3": { // upgrade for v1.17.0-rc3
		Handler: func(ctx sdk.Context, app *App, vm module.VersionMap) (module.VersionMap, error) {
			var err error
			vm, err = runModuleMigrations(ctx, app, vm)
			if err != nil {
				return nil, err
			}

			updateIbcMarkerDenomMetadata(ctx, app)

			return vm, nil
		},
	},
	"saffron": { // upgrade for v1.17.0,
		Handler: func(ctx sdk.Context, app *App, vm module.VersionMap) (module.VersionMap, error) {
			var err error
			vm, err = runModuleMigrations(ctx, app, vm)
			if err != nil {
				return nil, err
			}

			// set ibchoooks defaults (no allowed async contracts)
			app.IBCHooksKeeper.SetParams(ctx, ibchookstypes.DefaultParams())

			removeInactiveValidatorDelegations(ctx, app)
			setupICQ(ctx, app)
			updateMaxSupply(ctx, app)
			setExchangeParams(ctx, app)
			updateIbcMarkerDenomMetadata(ctx, app)

			return vm, nil
		},
		Added: []string{icqtypes.ModuleName, oracletypes.ModuleName, ibchookstypes.StoreKey, hold.ModuleName, exchange.ModuleName},
	},
	// TODO - Add new upgrade definitions here.
}

// InstallCustomUpgradeHandlers sets upgrade handlers for all entries in the upgrades map.
func InstallCustomUpgradeHandlers(app *App) {
	// Register all explicit appUpgrades
	for name, upgrade := range upgrades {
		// If the handler has been defined, add it here, otherwise, use no-op.
		var handler upgradetypes.UpgradeHandler
		if upgrade.Handler == nil {
			handler = func(ctx sdk.Context, plan upgradetypes.Plan, versionMap module.VersionMap) (module.VersionMap, error) {
				ctx.Logger().Info(fmt.Sprintf("Applying no-op upgrade to %q", plan.Name))
				return versionMap, nil
			}
		} else {
			ref := upgrade
			handler = func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
				ctx.Logger().Info(fmt.Sprintf("Starting upgrade to %q", plan.Name), "version-map", vm)
				newVM, err := ref.Handler(ctx, app, vm)
				if err != nil {
					ctx.Logger().Error(fmt.Sprintf("Failed to upgrade to %q", plan.Name), "error", err)
				} else {
					ctx.Logger().Info(fmt.Sprintf("Successfully upgraded to %q", plan.Name), "version-map", newVM)
				}
				return newVM, err
			}
		}
		app.UpgradeKeeper.SetUpgradeHandler(name, handler)
	}
}

// GetUpgradeStoreLoader creates an StoreLoader for use in an upgrade.
// Returns nil if no upgrade info is found or the upgrade doesn't need a store loader.
func GetUpgradeStoreLoader(app *App, info upgradetypes.Plan) baseapp.StoreLoader {
	upgrade, found := upgrades[info.Name]
	if !found {
		return nil
	}

	if len(upgrade.Renamed) == 0 && len(upgrade.Deleted) == 0 && len(upgrade.Added) == 0 {
		app.Logger().Info("No store upgrades required",
			"plan", info.Name,
			"height", info.Height,
		)
		return nil
	}

	storeUpgrades := storetypes.StoreUpgrades{
		Added:   upgrade.Added,
		Renamed: upgrade.Renamed,
		Deleted: upgrade.Deleted,
	}
	app.Logger().Info("Store upgrades",
		"plan", info.Name,
		"height", info.Height,
		"upgrade.added", storeUpgrades.Added,
		"upgrade.deleted", storeUpgrades.Deleted,
		"upgrade.renamed", storeUpgrades.Renamed,
	)
	return upgradetypes.UpgradeStoreLoader(info.Height, &storeUpgrades)
}

// runModuleMigrations wraps standard logging around the call to app.mm.RunMigrations.
// In most cases, it should be the first thing done during a migration.
//
// If state is updated prior to this migration, you run the risk of writing state using
// a new format when the migration is expecting all state to be in the old format.
func runModuleMigrations(ctx sdk.Context, app *App, vm module.VersionMap) (module.VersionMap, error) {
	// Even if this function is no longer called, do not delete it. Keep it around for the next time it's needed.
	ctx.Logger().Info("Starting module migrations. This may take a significant amount of time to complete. Do not restart node.")
	newVM, err := app.mm.RunMigrations(ctx, app.configurator, vm)
	if err != nil {
		ctx.Logger().Error("Module migrations encountered an error.", "error", err)
		return nil, err
	}
	ctx.Logger().Info("Module migrations completed.")
	return newVM, nil
}

// Create a use of runModuleMigrations so that the linter neither complains about it not being used,
// nor complains about a nolint:unused directive that isn't needed because the function is used.
var _ = runModuleMigrations

// addGovV1SubmitFee adds a msg-fee for the gov v1 MsgSubmitProposal if there isn't one yet.
// TODO: Remove with the rust handlers.
func addGovV1SubmitFee(ctx sdk.Context, app *App) {
	typeURL := sdk.MsgTypeURL(&govtypesv1.MsgSubmitProposal{})

	ctx.Logger().Info(fmt.Sprintf("Creating message fee for %q if it doesn't already exist.", typeURL))
	// At the time of writing this, the only way GetMsgFee returns an error is if it can't unmarshall state.
	// If that's the case for the v1 entry, we want to fix it anyway, so we just ignore any error here.
	fee, _ := app.MsgFeesKeeper.GetMsgFee(ctx, typeURL)
	// If there's already a fee for it, do nothing.
	if fee != nil {
		ctx.Logger().Info(fmt.Sprintf("Message fee for %q already exists with amount %q. Nothing to do.", fee.MsgTypeUrl, fee.AdditionalFee.String()))
		return
	}

	// Copy the fee from the beta entry if it exists, otherwise, just make it fresh.
	betaTypeURL := sdk.MsgTypeURL(&govtypesv1beta1.MsgSubmitProposal{})
	// Here too, if there's an error getting the beta fee, just ignore it.
	betaFee, _ := app.MsgFeesKeeper.GetMsgFee(ctx, betaTypeURL)
	if betaFee != nil {
		fee = betaFee
		fee.MsgTypeUrl = typeURL
		ctx.Logger().Info(fmt.Sprintf("Copying %q fee to %q.", betaTypeURL, fee.MsgTypeUrl))
	} else {
		fee = &msgfeetypes.MsgFee{
			MsgTypeUrl:           typeURL,
			AdditionalFee:        sdk.NewInt64Coin("nhash", 100_000_000_000), // 100 hash
			Recipient:            "",
			RecipientBasisPoints: 0,
		}
		ctx.Logger().Info(fmt.Sprintf("Creating %q fee.", fee.MsgTypeUrl))
	}

	// At the time of writing this, SetMsgFee always returns nil.
	_ = app.MsgFeesKeeper.SetMsgFee(ctx, *fee)
	ctx.Logger().Info(fmt.Sprintf("Successfully set fee for %q with amount %q.", fee.MsgTypeUrl, fee.AdditionalFee.String()))
}

// removeP8eMemorializeContractFee removes the message fee for the now-non-existent MsgP8eMemorializeContractRequest.
// TODO: Remove with the rust handlers.
func removeP8eMemorializeContractFee(ctx sdk.Context, app *App) {
	typeURL := "/provenance.metadata.v1.MsgP8eMemorializeContractRequest"

	ctx.Logger().Info(fmt.Sprintf("Removing message fee for %q if one exists.", typeURL))
	// Get the existing fee for log output, but ignore any errors so we try to delete the entry either way.
	fee, _ := app.MsgFeesKeeper.GetMsgFee(ctx, typeURL)
	// At the time of writing this, the only error that RemoveMsgFee can return is ErrMsgFeeDoesNotExist.
	// So ignore any error here and just use fee != nil for the different log messages.
	_ = app.MsgFeesKeeper.RemoveMsgFee(ctx, typeURL)
	if fee == nil {
		ctx.Logger().Info(fmt.Sprintf("Message fee for %q already does not exist. Nothing to do.", typeURL))
	} else {
		ctx.Logger().Info(fmt.Sprintf("Successfully removed message fee for %q with amount %q.", fee.MsgTypeUrl, fee.AdditionalFee.String()))
	}
}

// removeInactiveValidatorDelegations unbonds all delegations from inactive validators, triggering their removal from the validator set.
// This should be applied in most upgrades.
func removeInactiveValidatorDelegations(ctx sdk.Context, app *App) {
	unbondingTimeParam := app.StakingKeeper.GetParams(ctx).UnbondingTime
	ctx.Logger().Info(fmt.Sprintf("removing all delegations from validators that have been inactive (unbonded) for %d days", int64(unbondingTimeParam.Hours()/24)))
	removalCount := 0
	validators := app.StakingKeeper.GetAllValidators(ctx)
	for _, validator := range validators {
		if validator.IsUnbonded() {
			inactiveDuration := ctx.BlockTime().Sub(validator.UnbondingTime)
			if inactiveDuration >= unbondingTimeParam {
				ctx.Logger().Info(fmt.Sprintf("validator %v has been inactive (unbonded) for %d days and will be removed", validator.OperatorAddress, int64(inactiveDuration.Hours()/24)))
				valAddress, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
				if err != nil {
					ctx.Logger().Error(fmt.Sprintf("invalid operator address: %s: %v", validator.OperatorAddress, err))
					continue
				}
				delegations := app.StakingKeeper.GetValidatorDelegations(ctx, valAddress)
				for _, delegation := range delegations {
					ctx.Logger().Info(fmt.Sprintf("undelegate delegator %v from validator %v of all shares (%v)", delegation.DelegatorAddress, validator.OperatorAddress, delegation.GetShares()))
					_, err = app.StakingKeeper.Undelegate(ctx, delegation.GetDelegatorAddr(), valAddress, delegation.GetShares())
					if err != nil {
						ctx.Logger().Error(fmt.Sprintf("failed to undelegate delegator %s from validator %s: %v", delegation.GetDelegatorAddr().String(), valAddress.String(), err))
						continue
					}
				}
				removalCount++
			}
		}
	}
	ctx.Logger().Info(fmt.Sprintf("a total of %d inactive (unbonded) validators have had all their delegators removed", removalCount))
}

// fixNameIndexEntries fixes the name module's address to name index entries.
// TODO: Remove with the rust handlers.
func fixNameIndexEntries(ctx sdk.Context, app *App) {
	ctx.Logger().Info("Fixing name module store index entries.")
	app.NameKeeper.DeleteInvalidAddressIndexEntries(ctx)
	ctx.Logger().Info("Done fixing name module store index entries.")
}

// setAccountDataNameRecord makes sure the account data name record exists, is restricted,
// and is owned by the attribute module. An error is returned if it fails to make it so.
// TODO: Remove with the rust handlers.
func setAccountDataNameRecord(ctx sdk.Context, accountK attributetypes.AccountKeeper, nameK attributetypes.NameKeeper) (err error) {
	return attributekeeper.EnsureModuleAccountAndAccountDataNameRecord(ctx, accountK, nameK)
}

// setupICQ sets the correct default values for ICQKeeper.
// TODO: Remove with the saffron handlers.
func setupICQ(ctx sdk.Context, app *App) {
	ctx.Logger().Info("Updating ICQ params")
	app.ICQKeeper.SetParams(ctx, icqtypes.NewParams(true, []string{"/provenance.oracle.v1.Query/Oracle"}))
	ctx.Logger().Info("Done updating ICQ params")
}

// updateMaxSupply sets the value of max supply to the current value of MaxTotalSupply.
// TODO: Remove with the saffron handlers.
func updateMaxSupply(ctx sdk.Context, app *App) {
	ctx.Logger().Info("Updating MaxSupply marker param")
	params := app.MarkerKeeper.GetParams(ctx)
	//nolint:staticcheck // Populate new param with deprecated param
	params.MaxSupply = math.NewIntFromUint64(params.MaxTotalSupply)
	app.MarkerKeeper.SetParams(ctx, params)
	ctx.Logger().Info("Done updating MaxSupply marker param")
}

// setExchangeParams sets exchange module's params to the defaults.
// TODO: Remove with the saffron handlers.
func setExchangeParams(ctx sdk.Context, app *App) {
	ctx.Logger().Info("Ensuring exchange module params are set.")
	params := app.ExchangeKeeper.GetParams(ctx)
	if params != nil {
		ctx.Logger().Info("Exchange module params are already defined.")
	} else {
		params = exchange.DefaultParams()
		ctx.Logger().Info("Setting exchange module params to defaults.")
		app.ExchangeKeeper.SetParams(ctx, params)
	}
	ctx.Logger().Info("Done ensuring exchange module params are set.")
}

// updateIbcMarkerDenomMetadata iterates markers and creates denom metadata for ibc markers
// TODO: Remove with the saffron handlers.
func updateIbcMarkerDenomMetadata(ctx sdk.Context, app *App) {
	ctx.Logger().Info("Updating ibc marker denom metadata")
	app.MarkerKeeper.IterateMarkers(ctx, func(record types.MarkerAccountI) bool {
		if !strings.HasPrefix(record.GetDenom(), "ibc/") {
			return false
		}

		hash, err := transfertypes.ParseHexHash(strings.TrimPrefix(record.GetDenom(), "ibc/"))
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("invalid denom trace hash: %s, error: %s", hash.String(), err))
			return false
		}
		denomTrace, found := app.TransferKeeper.GetDenomTrace(ctx, hash)
		if !found {
			ctx.Logger().Error(fmt.Sprintf("trace not found: %s, error: %s", hash.String(), err))
			return false
		}

		parts := strings.Split(denomTrace.Path, "/")
		if len(parts) == 2 && parts[0] == "transfer" {
			ctx.Logger().Info(fmt.Sprintf("Adding metadata to %s", record.GetDenom()))
			chainID := app.Ics20MarkerHooks.GetChainID(ctx, parts[0], parts[1], app.IBCKeeper)
			markerMetadata := banktypes.Metadata{
				Base:        record.GetDenom(),
				Name:        chainID + "/" + denomTrace.BaseDenom,
				Display:     chainID + "/" + denomTrace.BaseDenom,
				Description: denomTrace.BaseDenom + " from " + chainID,
			}
			app.BankKeeper.SetDenomMetaData(ctx, markerMetadata)
		}

		return false
	})
	ctx.Logger().Info("Done updating ibc marker denom metadata")
}
