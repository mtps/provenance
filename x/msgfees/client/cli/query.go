package cli

import (
	"context"
	"encoding/base64"
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/provenance-io/provenance/x/msgfees/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the top-level command for msgfees CLI queries.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the msgfees module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(AllMsgBasedFeesCmd())
	return queryCmd
}

// AllMsgBasedFeesCmd is the CLI command for listing all msg based fees.
func AllMsgBasedFeesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List all the msg based fees on the Provenance Blockchain",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(withPageKeyDecoded(cmd.Flags()))
			if err != nil {
				return err
			}

			var response *types.QueryAllMsgBasedFeesResponse
			if response, err = queryClient.QueryAllMsgBasedFees(
				context.Background(),
				&types.QueryAllMsgBasedFeesRequest{Pagination: pageReq},
			); err != nil {
				fmt.Printf("failed to query msg based fees: %s\n", err.Error())
				return nil
			}
			return clientCtx.PrintProto(response)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "markers")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// sdk ReadPageRequest expects binary but we encoded to base64 in our marshaller
func withPageKeyDecoded(flagSet *flag.FlagSet) *flag.FlagSet {
	encoded, err := flagSet.GetString(flags.FlagPageKey)
	if err != nil {
		panic(err.Error())
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err.Error())
	}
	_ = flagSet.Set(flags.FlagPageKey, string(raw))
	return flagSet
}