package cbatch

import (
	"CraneFrontEnd/internal/util"
	"github.com/spf13/cobra"
	"os"
)

var (
	FlagNodes         uint32
	FlagCpuPerTask    float64
	FlagGpus          uint64
	FlagNtasksPerNode uint32
	FlagTime          string
	FlagMem           string
	FlagPartition     string
	FlagJob           string
	FlagOutput        string
	FlagAccount       string
	FlagQos           string
	FlagCwd           string
	FlagRepeat        uint32
	FlagNodelist      string
	FlagExclude       string

	FlagConfigFilePath string
)

func ParseCmdArgs() {
	rootCmd := &cobra.Command{
		Use:   "cbatch",
		Short: "submit batch jobs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Cbatch(args[0])
		},
	}

	rootCmd.PersistentFlags().StringVarP(&FlagConfigFilePath, "config", "C",
		util.DefaultConfigPath, "Path to configuration file")
	rootCmd.Flags().Uint32VarP(&FlagNodes, "nodes", "N", 0, " number of nodes on which to run (N = min[-max])")
	rootCmd.Flags().Float64VarP(&FlagCpuPerTask, "cpus-per-task", "c", 0, "number of cpus required per task")
	rootCmd.Flags().Uint64Var(&FlagGpus, "gpus", 0, "num of gpus required per task")
	rootCmd.Flags().Uint32Var(&FlagNtasksPerNode, "ntasks-per-node", 0, "number of tasks to invoke on each node")
	rootCmd.Flags().StringVarP(&FlagTime, "time", "t", "", "time limit")
	rootCmd.Flags().StringVar(&FlagMem, "mem", "", "minimum amount of real memory")
	rootCmd.Flags().StringVarP(&FlagPartition, "partition", "p", "", "partition requested")
	rootCmd.Flags().StringVarP(&FlagOutput, "output", "o", "", "file for batch script's standard output")
	rootCmd.Flags().StringVarP(&FlagJob, "job-name", "J", "", "name of job")
	rootCmd.Flags().StringVarP(&FlagAccount, "account", "A", "", "account used by the task")
	rootCmd.Flags().StringVar(&FlagCwd, "chdir", "", "working directory of the task")
	rootCmd.Flags().StringVarP(&FlagQos, "qos", "q", "", "quality of service")
	rootCmd.Flags().Uint32Var(&FlagRepeat, "repeat", 1, "submit the task multiple times")
	rootCmd.Flags().StringVarP(&FlagNodelist, "nodelist", "w", "", "List of specific nodes to be allocated to the job")
	rootCmd.Flags().StringVarP(&FlagExclude, "exclude", "x", "", "exclude a specific list of hosts")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
