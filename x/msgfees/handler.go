package msgfees

import (
	"github.com/provenance-io/provenance/x/msgfees/keeper"
	"github.com/provenance-io/provenance/x/msgfees/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewHandler returns a handler for msg based fee messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.CreateMsgBasedFeeRequest:
			res, err := msgServer.CreateMsgBasedFee(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddMsgBasedFeeProposal:
			return keeper.HandleAddMsgBasedFeeProposal(ctx, k, c)
		case *types.UpdateMsgBasedFeeProposal:
			return keeper.HandleUpdateMsgBasedFeeProposal(ctx, k, c)
		case *types.RemoveMsgBasedFeeProposal:
			return keeper.HandleRemoveMsgBasedFeeProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized marker proposal content type: %T", c)
		}
	}
}
