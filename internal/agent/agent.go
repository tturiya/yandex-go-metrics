package agent

import (
	"time"

	"github.com/tturiya/iter5/internal/config"
	"github.com/tturiya/iter5/internal/mertic"
	"github.com/tturiya/iter5/internal/storage/memstorage"
)

func StartAgent() error {
	var (
		agentCfg      = config.NewAgentConfig()
		metrics       = memstorage.NewMemStorage()
		timerInterval = time.NewTicker(
			time.Duration(agentCfg.ReportInterval) * time.Second)
		timerPoll = time.NewTicker(
			time.Duration(agentCfg.PollInterval) * time.Second)
	)

	for {
		select {
		case <-timerPoll.C:
			mertic.WriteMetric(metrics)
		case <-timerInterval.C:
			err := mertic.SendMetric(agentCfg.Addr, metrics)
			if err != nil {
				return err
			}
		}
	}
}
