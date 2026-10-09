package dashboard

import "testing"

func TestR10_AttributionByCwd(t *testing.T) {
	TestNestedProjectAttribution(t)
	TestSiblingPrefixNotAttributed(t)
}

func TestR11_GlobalProjectTokenAndConfig(t *testing.T) {
	TestGlobalProjectSessionsAndProgress(t)
	TestTokenLivesInHomeDashboard(t)
	TestStopCommand(t)
}

func TestR17_ParityCapRereadOmissions(t *testing.T) {
	TestProjectsParityWithRegistryAndOmissions(t)
	TestProjectsCap(t)
	TestConfigRereadWithServerUp(t)
}
