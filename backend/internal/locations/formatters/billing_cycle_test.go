package formatters

import (
	"testing"
)

func TestFormatBillingCycleToBackend(t *testing.T) {
    /*
        GIVEN a monthly frontend format
        WHEN FormatBillingCycleToBackend() is called
        THEN the method returns the monthly backend format
    */
    t.Run(`FormatBillingCycleToBackend() converts monthly frontend to backend format`, func(t *testing.T) {
        got := FormatBillingCycleToBackend(FrontendMonthly)
        if got != BackendMonthly {
            t.Errorf("FormatBillingCycleToBackend() = %v, want %v", got, BackendMonthly)
        }
    })

    /*
        GIVEN a quarterly frontend format
        WHEN FormatBillingCycleToBackend() is called
        THEN the method returns the quarterly backend format
    */
    t.Run(`FormatBillingCycleToBackend() converts quarterly frontend to backend format`, func(t *testing.T) {
        got := FormatBillingCycleToBackend(FrontendQuarterly)
        if got != BackendQuarterly {
            t.Errorf("FormatBillingCycleToBackend() = %v, want %v", got, BackendQuarterly)
        }
    })

    /*
        GIVEN an annual frontend format
        WHEN FormatBillingCycleToBackend() is called
        THEN the method returns the annual backend format
    */
    t.Run(`FormatBillingCycleToBackend() converts annual frontend to backend format`, func(t *testing.T) {
        got := FormatBillingCycleToBackend(FrontendAnnually)
        if got != BackendAnnually {
            t.Errorf("FormatBillingCycleToBackend() = %v, want %v", got, BackendAnnually)
        }
    })

    /*
        GIVEN an unknown frontend format
        WHEN FormatBillingCycleToBackend() is called
        THEN the method returns the monthly backend format as default
    */
    t.Run(`FormatBillingCycleToBackend() returns monthly for unknown format`, func(t *testing.T) {
        got := FormatBillingCycleToBackend("unknown")
        if got != BackendMonthly {
            t.Errorf("FormatBillingCycleToBackend() = %v, want %v", got, BackendMonthly)
        }
    })
}

func TestFormatBillingCycleToFrontend(t *testing.T) {
    /*
        GIVEN a monthly backend format
        WHEN FormatBillingCycleToFrontend() is called
        THEN the method returns the monthly frontend format
    */
    t.Run(`FormatBillingCycleToFrontend() converts monthly backend to frontend format`, func(t *testing.T) {
        got := FormatBillingCycleToFrontend(BackendMonthly)
        if got != FrontendMonthly {
            t.Errorf("FormatBillingCycleToFrontend() = %v, want %v", got, FrontendMonthly)
        }
    })

    /*
        GIVEN a quarterly backend format
        WHEN FormatBillingCycleToFrontend() is called
        THEN the method returns the quarterly frontend format
    */
    t.Run(`FormatBillingCycleToFrontend() converts quarterly backend to frontend format`, func(t *testing.T) {
        got := FormatBillingCycleToFrontend(BackendQuarterly)
        if got != FrontendQuarterly {
            t.Errorf("FormatBillingCycleToFrontend() = %v, want %v", got, FrontendQuarterly)
        }
    })

    /*
        GIVEN an annual backend format
        WHEN FormatBillingCycleToFrontend() is called
        THEN the method returns the annual frontend format
    */
    t.Run(`FormatBillingCycleToFrontend() converts annual backend to frontend format`, func(t *testing.T) {
        got := FormatBillingCycleToFrontend(BackendAnnually)
        if got != FrontendAnnually {
            t.Errorf("FormatBillingCycleToFrontend() = %v, want %v", got, FrontendAnnually)
        }
    })

    /*
        GIVEN an unknown backend format
        WHEN FormatBillingCycleToFrontend() is called
        THEN the method returns the monthly frontend format as default
    */
    t.Run(`FormatBillingCycleToFrontend() returns monthly for unknown format`, func(t *testing.T) {
        got := FormatBillingCycleToFrontend("unknown")
        if got != FrontendMonthly {
            t.Errorf("FormatBillingCycleToFrontend() = %v, want %v", got, FrontendMonthly)
        }
    })
}
