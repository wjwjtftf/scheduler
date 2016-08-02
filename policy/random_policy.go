package policy

import "math/rand"

type RandomPolicy struct {
	retryMaxTimes int
	urls          []string
	retryTimes    int
}

func NewRandomPolicy(retryMaxTimes int, urls []string) *RandomPolicy {

	return &RandomPolicy{retryMaxTimes:retryMaxTimes, urls:urls}
}

func (this *RandomPolicy)GetNextUrl() string {

	if this.retryTimes > this.retryMaxTimes {
		return ""
	}

	i := rand.Intn(len(this.urls))
	this.retryTimes ++
	return this.urls[i]
}

