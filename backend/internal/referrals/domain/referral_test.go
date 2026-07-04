package domain

import "testing"

func TestReferralIsOpen(t *testing.T) {
	if !(&Referral{}).IsOpen() {
		t.Error("a referral with no referrer should be open")
	}
	if (&Referral{ReferrerID: "alice"}).IsOpen() {
		t.Error("a directed referral should not be open")
	}
}

func TestValidOutcomes(t *testing.T) {
	for _, o := range []string{
		OutcomeApplicationSubmitted, OutcomeInterviewing, OutcomeOffer,
		OutcomeHired, OutcomeRejected, OutcomeWithdrawn,
	} {
		if !ValidOutcomes[o] {
			t.Errorf("expected %q to be a valid outcome", o)
		}
	}
	if ValidOutcomes["bogus"] {
		t.Error("did not expect an unknown outcome to be valid")
	}
}
