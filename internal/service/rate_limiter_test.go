package service

import "testing"

func TestRateLimiterCheck(t *testing.T) {
	limiter := &RateLimiter{attempts: make(map[string]*attemptRecord)}

	allowed, _ := limiter.Check("ip1", 2, 1)
	if !allowed {
		t.Fatal("first attempt should be allowed")
	}
	allowed, _ = limiter.Check("ip1", 2, 1)
	if !allowed {
		t.Fatal("second attempt should be allowed")
	}
	allowed, wait := limiter.Check("ip1", 2, 1)
	if allowed || wait == 0 {
		t.Fatal("third attempt should be blocked with wait time")
	}

	limiter.Reset("ip1")
	allowed, _ = limiter.Check("ip1", 2, 1)
	if !allowed {
		t.Fatal("attempt should be allowed after reset")
	}
}
