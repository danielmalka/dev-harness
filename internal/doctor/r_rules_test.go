package doctor

import "testing"

func TestR9_SettingsHintWithAndWithoutEntry(t *testing.T) { TestHarnessSettingsHint(t) }

func TestR14_DoctorModeIgnoreAndBothFolders(t *testing.T) {
	TestHarnessRepoRegistered(t)
	TestHarnessUnregisteredAndBoth(t)
	TestHarnessGlobalAndNone(t)
	TestHarnessLocalYamlIgnore(t)
}

// PRD-017 R2 (QA, T-1703-05)
func TestR2_DoctorPdocsOnlyNoBoth(t *testing.T) { TestHarnessPdocsOnlyNoBothWarning(t) }
