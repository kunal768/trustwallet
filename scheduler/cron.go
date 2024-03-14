package scheduler

import (
	"context"

	"github.com/kunal768/trustwallet/parser"
)

type CronJob interface {
	StoreSubscriberTxns()
}

type cronSvc struct {
	parserSvc parser.Parser
}

func NewSchedulerService(svc parser.Parser) CronJob {
	return &cronSvc{
		parserSvc: svc,
	}
}

func (svc cronSvc) StoreSubscriberTxns() {
	ctx := context.Background()
	svc.parserSvc.UpdateSubscriberTxns(ctx)
}
