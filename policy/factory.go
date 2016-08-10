package policy

import (
	"github.com/wjwjtftf/scheduler/entity"
	"strings"
)

type Factory struct {
}

func (this *Factory) FindPolicy(jobinfo *entity.JobInfo) Policy {

	switch jobinfo.Type {
	case "PRIORITY":
		return &PriorityPolicy{urls: strings.Split(jobinfo.Urls, ","), retryMaxTimes: 1}
	default:
		return &RandomPolicy{urls: strings.Split(jobinfo.Urls, ","), retryMaxTimes: 1}

	}
}
