package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/provenance-io/provenance/x/metadata/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	uuid "github.com/google/uuid"
)

const (
	FlagSigners = "signers"
)

// NewTxCmd is the top-level command for attribute CLI transactions.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases:                    []string{"m"},
		Short:                      "Transaction commands for the metadata module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		AddMetadataScopeCmd(),
		RemoveMetadataScopeCmd(),

		AddOsLocatorCmd(),
		RemoveOsLocatorCmd(),
		ModifyOsLocatorCmd(),

		AddContractSpecificationCmd(),
		RemoveContractSpecificationCmd(),

		AddRecordSpecificationCmd(),
		RemoveRecordSpecificationCmd(),
	)

	return txCmd
}

// AddMetadataScopeCmd creates a command for adding a metadata scope.
func AddMetadataScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-scope [scope-uuid] [spec-id] [owner-addresses] [data-access] [value-owner-address]",
		Short: "Add/Update a metadata scope to the provenance blockchain",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			scopeUUID, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Printf("Invalid uuid for scope uuid: %s", args[0])
				return err
			}
			specUUID, err := uuid.Parse(args[1])
			if err != nil {
				fmt.Printf("Invalid uuid for specification uuid: %s", args[0])
				return err
			}

			specID := types.ScopeSpecMetadataAddress(specUUID)

			ownerAddresses := strings.Split(args[2], ",")
			owners := make([]types.Party, len(ownerAddresses))
			for i, ownerAddr := range ownerAddresses {
				owners[i] = types.Party{Address: ownerAddr, Role: types.PartyType_PARTY_TYPE_OWNER}
			}
			dataAccess := strings.Split(args[3], ",")
			valueOwnerAddress := args[4]

			signers, err := parseSigners(cmd, &clientCtx)
			if err != nil {
				return err
			}

			scope := *types.NewScope(
				types.ScopeMetadataAddress(scopeUUID),
				specID,
				owners,
				dataAccess,
				valueOwnerAddress)

			if err := scope.ValidateBasic(); err != nil {
				fmt.Printf("Failed to validate scope %s : %v", scope.String(), err)
				return err
			}

			msg := types.NewMsgAddScopeRequest(scope, signers)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addSignerFlagCmd(cmd)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// RemoveMetadataScopeCmd creates a command for removing a scope.
func RemoveMetadataScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-scope [scope-address] [signers]",
		Short: "Remove a metadata scope to the provenance blockchain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			scopeUUID, err := uuid.Parse(args[0])
			if err != nil {
				fmt.Printf("Invalid uuid for scope id: %s", args[0])
				return err
			}

			scopeMetaAddress := types.ScopeMetadataAddress(scopeUUID)

			signers, err := parseSigners(cmd, &clientCtx)
			if err != nil {
				return err
			}

			deleteScope := *types.NewMsgDeleteScopeRequest(scopeMetaAddress, signers)
			if err := deleteScope.ValidateBasic(); err != nil {
				fmt.Printf("Failed to validate remove scope %s : %v", deleteScope.String(), err)
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &deleteScope)
		},
	}

	addSignerFlagCmd(cmd)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// RemoveMetadataScopeCmd creates a command for removing a scope.
func AddOsLocatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-locator [owner] [uri]",
		Short: "Add a uri to an owner address on the provenance blockchain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if _, errAddr := sdk.AccAddressFromBech32(args[0]); errAddr != nil {
				fmt.Printf("failed to add locator for a given owner address, invalid address: %s\n", args[0])
				return fmt.Errorf("invalid address: %w", errAddr)
			}

			objectStoreLocator := types.ObjectStoreLocator{
				LocatorUri: args[1], Owner: args[0],
			}

			addOSLocator := *types.NewMsgBindOSLocatorRequest(objectStoreLocator)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &addOSLocator)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// RemoveOsLocatorCmd creates a command for removing a os locator
func RemoveOsLocatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-locator [owner] [uri]",
		Short: "Remove an os locator already associated owner address on the provenance blockchain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if _, errAddr := sdk.AccAddressFromBech32(args[0]); errAddr != nil {
				fmt.Printf("failed to remove locator for a given owner address, invalid address: %s\n", args[0])
				return fmt.Errorf("invalid address: %w", errAddr)
			}

			objectStoreLocator := types.ObjectStoreLocator{
				LocatorUri: args[1], Owner: args[0],
			}

			deleteOSLocator := *types.NewMsgDeleteOSLocatorRequest(objectStoreLocator)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &deleteOSLocator)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// ModifyOsLocatorCmd creates a command for modifying os locator
func ModifyOsLocatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-locator [owner] [uri]",
		Short: "Modify a uri already associated owner address on the provenance blockchain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if _, errAddr := sdk.AccAddressFromBech32(args[0]); errAddr != nil {
				fmt.Printf("failed to add locator for a given owner address, invalid address: %s\n", args[0])
				return fmt.Errorf("invalid address: %w", errAddr)
			}
			if err != nil {
				fmt.Printf("Invalid uuid for scope id: %s", args[0])
				return err
			}

			objectStoreLocator := types.ObjectStoreLocator{
				LocatorUri: args[1], Owner: args[0],
			}

			modifyOSLocator := *types.NewMsgModifyOSLocatorRequest(objectStoreLocator)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &modifyOSLocator)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// AddContractSpecificationCmd creates a command to add/update contract specifications
func AddContractSpecificationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-contract-specification [contractspec-id] [owners] [parties-involved] [source-type] [source-value] [classname] [description-name] [description] [website-url] [icon-url]",
		Short: "Add/Update metadata contract specification on the provenance blockchain",
		Long: `Add/Update metadata contract specification on the provenance blockchain
[contractspec-id] - contract specification metaaddress
[owners] - comma delimited list of bech32 owner addresses
[parties-involved] - comma delimited list of party types.  Accepted values: originator,servicer,investor,custodian,owner,affiliate,omnibus,provenance
[source-type] - accepted values: hash or resourceid
[source-value] - source identifier of type hash or resourceid
[classname] - Name of contract specification
[description-name] - optional- description name identifier 
[description] - optional - description text
[website-url] - optional - address of website
[icon-url] - optional - address to a image to be used as an icon
		`,
		Args: cobra.RangeArgs(6, 10),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var specificationID types.MetadataAddress
			specificationID, err = types.MetadataAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if !specificationID.IsContractSpecificationAddress() {
				return fmt.Errorf("invalid contract specification id: %s", args[0])
			}

			partiesInvolved := parsePartyTypes(args[2])
			description := parseDescription(args[6:])
			contractSpecification := types.ContractSpecification{SpecificationId: specificationID,
				Description:     description,
				OwnerAddresses:  strings.Split(args[1], ","),
				PartiesInvolved: partiesInvolved,
				ClassName:       args[5],
			}
			switch s := strings.ToUpper(args[3]); s {
			case "RESOURCEID":
				var recordID types.MetadataAddress
				recordID, err = types.MetadataAddressFromBech32(args[4])
				if err != nil {
					return err
				}
				contractSpecification.Source = &types.ContractSpecification_ResourceId{
					ResourceId: recordID,
				}
			case "HASH":
				contractSpecification.Source = &types.ContractSpecification_Hash{
					Hash: args[4],
				}
			default:
				return fmt.Errorf("incorrect source type for contract specification: %s", s)
			}

			signers, err := parseSigners(cmd, &clientCtx)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddContractSpecificationRequest(contractSpecification, signers)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	addSignerFlagCmd(cmd)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func AddRecordSpecificationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-record-specification [specification-id] [name] [input-specifications] [type-name] [result-types] [responsible-parties]",
		Short: "Add/Update metadata record specification to the provenance blockchain",
		Long: fmt.Sprintf(`Add/Update metadata record specification to the provenance blockchain.
[specification-id] - record specification metaaddress
[name] - record name
[input-specifications] - semi-colon delimited list of input specifications <name>,<type-name>,<source>,<source-value>
[type-name] - contract specification type name
[result-types] - result definition type.  Accepted values: proposed,record,record_list
[responsible-parties] - comma delimited list of party types.  Accepted values: originator,servicer,investor,custodian,owner,affiliate,omnibus,provenance

Example: 
$ %s tx metadata recspec1qh... recordname inputname1,typename1,hash,hashvalue;inputename2,typename2,recordid,<recordmetaaddress> record_list owner,originator

`, version.AppName),
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			specificationID, err := types.MetadataAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			if !specificationID.IsRecordSpecificationAddress() {
				return fmt.Errorf("invalid record specification id: %s", args[0])
			}

			recordName := args[1]
			contractSpecUUID, _ := specificationID.ContractSpecUUID()
			expectedRecordID := types.RecordSpecMetadataAddress(contractSpecUUID, recordName)
			if expectedRecordID.String() != specificationID.String() {
				return fmt.Errorf("invalid record name does not produce correct record id actual: %v expected %v", specificationID.String(), expectedRecordID.String())
			}

			inputs, err := inputSpecification(args[2])
			if err != nil {
				return err
			}

			resultType := definitionType(args[4])
			partyTypes := parsePartyTypes(args[5])
			signers, err := parseSigners(cmd, &clientCtx)
			if err != nil {
				return err
			}

			recordSpecification := types.RecordSpecification{
				SpecificationId:    specificationID,
				Name:               args[1],
				Inputs:             inputs,
				TypeName:           args[3],
				ResultType:         resultType,
				ResponsibleParties: partyTypes,
			}

			msg := *types.NewMsgAddRecordSpecificationRequest(recordSpecification, signers)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	addSignerFlagCmd(cmd)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// inputSpecification converts cli delimited argument and converts it to InputSpecifications
func inputSpecification(cliDelimitedValue string) ([]*types.InputSpecification, error) {
	delimitedInputs := strings.Split(cliDelimitedValue, ";")
	inputs := make([]*types.InputSpecification, len(delimitedInputs))
	for i, delimitedInput := range delimitedInputs {
		values := strings.Split(delimitedInput, ",")
		if len(values) != 4 {
			return nil, fmt.Errorf("invalid number of values for input specification: %v", len(values))
		}
		name := values[0]
		typeName := values[1]
		switch s := strings.ToUpper(values[2]); s {
		case "RECORDID":
			recordID, err := types.MetadataAddressFromBech32(values[3])
			if err != nil {
				return nil, err
			}
			inputs[i] = &types.InputSpecification{
				Name: name,
				Source: &types.InputSpecification_RecordId{
					RecordId: recordID,
				},
				TypeName: typeName,
			}
		case "HASH":
			inputs[i] = &types.InputSpecification{
				Name:     name,
				Source:   &types.InputSpecification_Hash{Hash: values[3]},
				TypeName: typeName,
			}
		default:
			return nil, fmt.Errorf("incorrect source type for input specification: %s", s)
		}
	}
	return inputs, nil
}

func addSignerFlagCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagSigners, "", "comma delimited list of bech32 addresses")
}

func parseSigners(cmd *cobra.Command, client *client.Context) ([]string, error) {
	flagSet := cmd.Flags()
	if flagSet.Changed(FlagSigners) {
		signerList, _ := flagSet.GetString(FlagSigners)
		signers := strings.Split(signerList, ",")
		for _, signer := range signers {
			_, err := sdk.AccAddressFromBech32(signer)
			if err != nil {
				fmt.Printf("signer address must be a Bech32 string: %v", err)
				return nil, err
			}
		}
		return signers, nil
	}
	return []string{client.GetFromAddress().String()}, nil
}

func parsePartyTypes(delimitedPartyTypes string) []types.PartyType {
	parties := strings.Split(delimitedPartyTypes, ",")
	partyTypes := make([]types.PartyType, len(parties))
	for i, party := range parties {
		partyValue := types.PartyType_value[fmt.Sprintf("PARTY_TYPE_%s", strings.ToUpper(party))]
		partyTypes[i] = types.PartyType(partyValue)
	}
	return partyTypes
}

func definitionType(cliValue string) types.DefinitionType {
	typeValue := types.DefinitionType_value[fmt.Sprintf("DEFINITION_TYPE_%s", strings.ToUpper(cliValue))]
	return types.DefinitionType(typeValue)
}

func parseDescription(cliArgs []string) *types.Description {
	if len(cliArgs) == 0 {
		return nil
	}

	description := types.Description{}
	if len(cliArgs) >= 1 {
		description.Name = cliArgs[0]
	}
	if len(cliArgs) >= 2 {
		description.Description = cliArgs[1]
	}
	if len(cliArgs) >= 3 {
		description.WebsiteUrl = cliArgs[2]
	}
	if len(cliArgs) >= 4 {
		description.IconUrl = cliArgs[3]
	}
	return &description
}

// RemoveContractSpecificationCmd creates a command to remove a contract specification
func RemoveContractSpecificationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-contract-specification [specification-id]",
		Short: "Removes a contract specification on the provenance blockchain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var specificationID types.MetadataAddress
			specificationID, err = types.MetadataAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if !specificationID.IsContractSpecificationAddress() {
				return fmt.Errorf("invalid contract specification id: %s", args[0])
			}

			signers, err := parseSigners(cmd, &clientCtx)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteContractSpecificationRequest(specificationID, signers)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addSignerFlagCmd(cmd)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func RemoveRecordSpecificationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-record-specification [specification-id]",
		Short: "Remove record specification from the provenance blockchain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var specificationID types.MetadataAddress
			specificationID, err = types.MetadataAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if !specificationID.IsRecordSpecificationAddress() {
				return fmt.Errorf("invalid contract specification id: %s", args[0])
			}
			signers, err := parseSigners(cmd, &clientCtx)
			if err != nil {
				return err
			}
			msg := *types.NewMsgDeleteRecordSpecificationRequest(specificationID, signers)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	addSignerFlagCmd(cmd)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
