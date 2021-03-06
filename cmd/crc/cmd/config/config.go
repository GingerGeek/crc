package config

import (
	"sort"
	"strings"

	cfg "github.com/code-ready/crc/pkg/crc/config"
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/spf13/cobra"
)

var (
	// Start command settings in config

	Bundle         = cfg.AddSetting("bundle", nil, []cfg.ValidationFnType{cfg.ValidateBundle}, []cfg.SetFn{cfg.SuccessfullyApplied})
	CPUs           = cfg.AddSetting("cpus", constants.DefaultCPUs, []cfg.ValidationFnType{cfg.ValidateCPUs}, []cfg.SetFn{cfg.RequiresRestartMsg})
	Memory         = cfg.AddSetting("memory", constants.DefaultMemory, []cfg.ValidationFnType{cfg.ValidateMemory}, []cfg.SetFn{cfg.RequiresRestartMsg})
	NameServer     = cfg.AddSetting("nameserver", nil, []cfg.ValidationFnType{cfg.ValidateIPAddress}, []cfg.SetFn{cfg.SuccessfullyApplied})
	PullSecretFile = cfg.AddSetting("pull-secret-file", nil, []cfg.ValidationFnType{cfg.ValidatePath}, []cfg.SetFn{cfg.SuccessfullyApplied})

	DisableUpdateCheck   = cfg.AddSetting("disable-update-check", nil, []cfg.ValidationFnType{cfg.ValidateBool}, []cfg.SetFn{cfg.SuccessfullyApplied})
	ExperimentalFeatures = cfg.AddSetting("enable-experimental-features", nil, []cfg.ValidationFnType{cfg.ValidateBool}, []cfg.SetFn{cfg.SuccessfullyApplied})

	// Proxy Configuration

	HTTPProxy   = cfg.AddSetting("http-proxy", nil, []cfg.ValidationFnType{cfg.ValidateURI}, []cfg.SetFn{cfg.SuccessfullyApplied})
	HTTPSProxy  = cfg.AddSetting("https-proxy", nil, []cfg.ValidationFnType{cfg.ValidateURI}, []cfg.SetFn{cfg.SuccessfullyApplied})
	NoProxy     = cfg.AddSetting("no-proxy", nil, []cfg.ValidationFnType{cfg.ValidateNoProxy}, []cfg.SetFn{cfg.SuccessfullyApplied})
	ProxyCAFile = cfg.AddSetting("proxy-ca-file", nil, []cfg.ValidationFnType{cfg.ValidatePath}, []cfg.SetFn{cfg.SuccessfullyApplied})
)

var (
	configCmd = &cobra.Command{
		Use:   "config SUBCOMMAND [flags]",
		Short: "Modify crc configuration",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func init() {
	configCmd.Long = `Modifies crc configuration properties.
Configurable properties (enter after SUBCOMMAND): ` + "\n\n" + configurableFields()
}

func isPreflightKey(key string) bool {
	return strings.HasPrefix(key, "skip-") || strings.HasPrefix(key, "warn-")
}

// less is used to sort the config keys. We want to sort first the regular keys, and
// then the keys related to preflight starting with a skip- or warn- prefix. We want
// these preflight keys to be grouped by pair: 'skip-bar', 'warn-bar', 'skip-foo', 'warn-foo'
// would be sorted in that order.
func less(lhsKey, rhsKey string) bool {
	if isPreflightKey(lhsKey) {
		if isPreflightKey(rhsKey) {
			// lhs is preflight, rhs is preflight
			if lhsKey[4:] == rhsKey[4:] {
				// we want skip-foo before warn-foo
				return lhsKey < rhsKey
			}
			// ignore skip-/warn- prefix
			return lhsKey[4:] < rhsKey[4:]
		}
		// lhs is preflight, rhs is not preflight
		return false
	}

	if isPreflightKey(rhsKey) {
		// lhs is not preflight, rhs is preflight
		return true
	}

	// lhs is not preflight, rhs is not preflight
	return lhsKey < rhsKey
}

func configurableFields() string {
	var fields []string
	var keys = make([]string, len(cfg.AllConfigs()))

	for key := range cfg.AllConfigs() {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return less(keys[i], keys[j])
	})
	for _, key := range keys {
		fields = append(fields, " * "+key)
	}
	return strings.Join(fields, "\n")
}

func GetConfigCmd() *cobra.Command {
	/* Delay generation of configCmd.Long as much as possible as some parts of crc may have registered more
	 * fields after init() time but before the command is registered
	 */
	configCmd.Long = `Modifies crc configuration properties.
Configurable properties (enter as SUBCOMMAND): ` + "\n\n" + configurableFields()

	return configCmd
}
