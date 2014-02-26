package main

import (
	zaa "./zabbix_aggregate_agent"
	"flag"
	"log"
)

const (
	DefaultAddress = "127.0.0.1:10050"
	DefaultTimeout = 5
)

func main() {
	var (
		listen      string
		listFile    string
		listArg     string
		listCommand string
		timeout     int
		expires     int
	)

	flag.StringVar(&listen, "listen", DefaultAddress, "listen address e.g. 0.0.0.0:10050")
	flag.StringVar(&listFile, "list-file", "", "zabbix-agent list file")
	flag.StringVar(&listCommand, "list-command", "", "command which prints zabbix-agent list to stdout")
	flag.StringVar(&listArg, "list", "", "zabbix-agent list , separated. e.g. 'web.example.com:10050,192.168.1.1:10050'")
	flag.IntVar(&timeout, "timeout", DefaultTimeout, "network timeout with zabbix-agent (seconds)")
	flag.IntVar(&expires, "expires", 0, "list cache expires (seconds)")
	flag.Parse()

	agent := zaa.NewAgent()
	agent.Timeout = timeout

	if listFile != "" {
		agent.ListGenerator = zaa.ListFromFile
		agent.ListSource = listFile
	} else if listArg != "" {
		agent.ListGenerator = zaa.ListFromArg
		agent.ListSource = listArg
	} else if listCommand != "" {
		agent.ListGenerator = zaa.CachedListGenerator(zaa.ListFromCommand, expires)
		agent.ListSource = listCommand
	} else {
		log.Fatalln("option either --list, --list-file or --list-command is required.")
	}

	err := agent.Run(listen)
	if err != nil {
		log.Fatalln("Error", err)
	}
}
