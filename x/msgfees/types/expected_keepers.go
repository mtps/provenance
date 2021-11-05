package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosauthtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authtypes.AccountI) (stop bool))
	GetAccount(sdk.Context, sdk.AccAddress) authtypes.AccountI
	SetAccount(sdk.Context, authtypes.AccountI)
	NewAccount(sdk.Context, authtypes.AccountI) authtypes.AccountI
}

// BankKeeper defines the expected bank keeper (keeper, sendkeeper, viewkeeper) (noalias)
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

}

// MsgBasedFeeKeeper for additional msg fees.
type MsgBasedFeeKeeper interface {
	GetMsgBasedFee(ctx sdk.Context, msgType string) (*MsgBasedFee, error)
	GetFeeCollectorName() string
	DeductFees(bankKeeper cosmosauthtypes.BankKeeper, ctx sdk.Context, acc cosmosauthtypes.AccountI, fees sdk.Coins) error
}

// FeegrantKeeper defines the expected feegrant keeper.
type FeegrantKeeper interface {
	GetAllowance(ctx sdk.Context, granter sdk.AccAddress, grantee sdk.AccAddress) (feegrant.FeeAllowanceI, error)
	UseGrantedFees(ctx sdk.Context, granter, grantee sdk.AccAddress, fee sdk.Coins, msgs []sdk.Msg) error
}