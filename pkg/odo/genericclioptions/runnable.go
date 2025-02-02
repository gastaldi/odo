package genericclioptions

import (
	"flag"
	"os"

	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/odo/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Runnable interface {
	Complete(name string, cmd *cobra.Command, args []string) error
	Validate() error
	Run() error
}

func GenericRun(o Runnable, cmd *cobra.Command, args []string) {

	// CheckMachineReadableOutput
	// fixes / checks all related machine readable output functions
	CheckMachineReadableOutputCommand()

	// Run completion, validation and run.
	util.LogErrorAndExit(o.Complete(cmd.Name(), cmd, args), "")
	util.LogErrorAndExit(o.Validate(), "")
	util.LogErrorAndExit(o.Run(), "")
}

// CheckMachineReadableOutputCommand performs machine-readable output functions required to
// have it work correctly
func CheckMachineReadableOutputCommand() {

	// Check that the -o flag has been correctly set as json.
	outputFlag := pflag.Lookup("o")
	if outputFlag != nil && outputFlag.Changed && outputFlag.Value.String() != "json" {
		log.Error("Please input a valid output format for -o, available format: json")
		os.Exit(1)
	}

	// Before running anything, we will make sure that no verbose output is made
	// This is a HACK to manually override `-v 4` to `-v 0` (in which we have no glog.V(0) in our code...
	// in order to have NO verbose output when combining both `-o json` and `-v 4` so json output
	// is not malformed / mixed in with normal logging
	if log.IsJSON() {
		flag.Set("v", "0")
	}
}
