package config

import (
	"flag"
	"os"
	"strconv"
)

// Client configuration
type AgentConfig struct {
	Addr           string
	ReportInterval int
	PollInterval   int
}

// Constructs AgentConfig with input data from flags and env
func NewAgentConfig() *AgentConfig {
	ac := &AgentConfig{}
	ac.parseAgent()

	return ac
}

func (ac *AgentConfig) parseAgent() {
	flag.StringVar(&ac.Addr, "a", ":8080", "server address")
	flag.IntVar(&ac.ReportInterval, "r", 10, "report interval time in sec")
	flag.IntVar(&ac.PollInterval, "p", 2, "poll interval time in sec")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		ac.Addr = envAddress
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportIntervalVal, err := strconv.Atoi(envReportInterval)
		if err != nil {
			panic(err)
		}
		ac.ReportInterval = envReportIntervalVal
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		envPoolIntervalVal, err := strconv.Atoi(envPoolInterval)
		if err != nil {
			panic(err)
		}
		ac.PollInterval = envPoolIntervalVal
	}
}

// Server configuration
type ServerConfig struct {
	Addr string
}

// Constructs ServerConfig with input data from flags and env
func NewServerConfig() *ServerConfig {
	sc := &ServerConfig{}
	sc.parseServer()

	return sc
}

func (sc *ServerConfig) parseServer() {
	flag.StringVar(&sc.Addr, "a", ":8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		sc.Addr = envRunAddr
	}
}
