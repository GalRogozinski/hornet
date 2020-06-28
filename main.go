package main

import (
	"flag"
	"fmt"
	"github.com/gohornet/hornet/pkg/shutdown"
	daemon "github.com/iotaledger/hive.go/daemon"
	"github.com/iotaledger/hive.go/node"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/gohornet/hornet/pkg/config"
	"github.com/gohornet/hornet/pkg/toolset"
	"github.com/gohornet/hornet/plugins/autopeering"
	"github.com/gohornet/hornet/plugins/cli"
	"github.com/gohornet/hornet/plugins/coordinator"
	"github.com/gohornet/hornet/plugins/dashboard"
	"github.com/gohornet/hornet/plugins/database"
	"github.com/gohornet/hornet/plugins/gossip"
	"github.com/gohornet/hornet/plugins/gracefulshutdown"
	"github.com/gohornet/hornet/plugins/graph"
	"github.com/gohornet/hornet/plugins/metrics"
	"github.com/gohornet/hornet/plugins/monitor"
	"github.com/gohornet/hornet/plugins/mqtt"
	"github.com/gohornet/hornet/plugins/peering"
	"github.com/gohornet/hornet/plugins/profiling"
	"github.com/gohornet/hornet/plugins/prometheus"
	"github.com/gohornet/hornet/plugins/snapshot"
	"github.com/gohornet/hornet/plugins/spammer"
	"github.com/gohornet/hornet/plugins/tangle"
	"github.com/gohornet/hornet/plugins/tipselection"
	"github.com/gohornet/hornet/plugins/warpsync"
	"github.com/gohornet/hornet/plugins/webapi"
	"github.com/gohornet/hornet/plugins/zmq"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var f *os.File
var err error



func memProfiler() {
	if *memprofile != "" {
		f, err = os.Create(*memprofile)
		if err != nil {
			fmt.Errorf("could not create memory profile: %s", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Errorf("could not write memory profile: %s", err)
		}
	}
}

func closeProfiler() {
	pprof.StopCPUProfile()
	f.Close()
}


func main() {
	runtime.SetMutexProfileFraction(50)

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Errorf("could not create CPU profile: %s", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Errorf("could not start CPU profile: %s", err)
		}
		defer pprof.StopCPUProfile()
		defer f.Close()
	}
	
	daemon.BackgroundWorker("Profiler", func(shutdownSignal <-chan struct{}) {
		<-shutdownSignal
		fmt.Println("shutting down profiler")
		closeProfiler()
	}, shutdown.PriorityProfiler)

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6067", nil))
	}()

	cli.PrintVersion()
	cli.ParseConfig()
	toolset.HandleTools()
	cli.PrintConfig()

	plugins := []*node.Plugin{
		cli.PLUGIN,
		gracefulshutdown.PLUGIN,
		profiling.PLUGIN,
		database.PLUGIN,
		autopeering.PLUGIN,
		webapi.PLUGIN,
	}

	if !config.NodeConfig.GetBool(config.CfgNetAutopeeringRunAsEntryNode) {
		plugins = append(plugins, []*node.Plugin{
			gossip.PLUGIN,
			tangle.PLUGIN,
			peering.PLUGIN,
			warpsync.PLUGIN,
			tipselection.PLUGIN,
			metrics.PLUGIN,
			snapshot.PLUGIN,
			dashboard.PLUGIN,
			zmq.PLUGIN,
			mqtt.PLUGIN,
			graph.PLUGIN,
			monitor.PLUGIN,
			spammer.PLUGIN,
			coordinator.PLUGIN,
			prometheus.PLUGIN,
		}...)
	}

	node.Run(node.Plugins(plugins...))

}
