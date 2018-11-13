package cmd

import (
	"fmt"
	"github.com/lanfang/gcurl/config"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	version    string = "0.0.1"
	quit       <-chan os.Signal
	symbolList []string
	addr       string
	method     string
)

func parseArgs(cmd *cobra.Command, args []string) error {
	var err error
	n := len(args)
	if n == 1 {
		config.G_Conf.Addr = args[0]
	} else if n > 1 {
		config.G_Conf.Addr, config.G_Conf.SymbolList = args[0], args[1:]
	} else {
		fmt.Println(`gcurl: try 'gcurl --help' for more information`)
		os.Exit(1)
	}
	return err
}

var cmdDesc = &cobra.Command{
	Use:   "desc",
	Short: "desc symbol, show the detail info of the symbol",
	Long:  `show the detail info of the symbol. support multiple symbol. show  message definition, field type definition`,
	RunE:  commandDesc,
	Args:  parseArgs,
}

var rootCmd = &cobra.Command{
	Use:   "gcurl",
	Short: `interacting with the rpc server`,
	Long: `a command line tool for gRPC, like curl for HTTP. 
you can interact with the rpc server like this: 
gcurl host:port method -d '{"username":"gcurl", "password":"gcurl"}' or exec the subcommand`,
	Version:      "0.0.1",
	RunE:         commandDefault,
	Args:         parseArgs,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(cmdDesc)
}

func listenSignal(signals ...os.Signal) <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signals = append(signals, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)
	}
	signal.Notify(sig, signals...)
	return sig
}
func stop() <-chan os.Signal {
	return quit
}
func Execute() {
	errCh := make(chan error, 1)
	go func(errCh chan error) {
		err := rootCmd.Execute()
		errCh <- err
	}(errCh)
	waitQuit(errCh)
}
func waitQuit(errCh chan error) {
	ch := listenSignal()
	for {
		select {
		case <-ch:
			fmt.Println("gcurl exit...")
			os.Exit(0)
		case <-errCh:
			os.Exit(1)
		}
	}
}

type empty struct {
}
type Task struct {
	T   func(interface{})
	Arg interface{}
}

type Runner struct {
	worker  int
	pending []Task
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) AddTask(t Task) {
	r.pending = append(r.pending, t)
}

func (r *Runner) Start(n int) {
	r.worker = n
	wg := sync.WaitGroup{}
	concurrent := make(chan empty, r.worker)
	for _, task := range r.pending {
		concurrent <- empty{}
		wg.Add(1)
		go func(t Task) {
			defer func() {
				wg.Done()
				<-concurrent
			}()
			t.T(t.Arg)
		}(task)
	}
	wg.Wait()
}
