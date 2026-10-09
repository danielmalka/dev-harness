package doctor

import "testing"

func TestR9_SettingsHintWithAndWithoutEntry(t *testing.T) { TestHarnessSettingsHint(t) }

func TestR14_DoctorModeIgnoreAndBothFolders(t *testing.T) {
	TestHarnessRepoRegistered(t)
	TestHarnessUnregisteredAndBoth(t)
	TestHarnessGlobalAndNone(t)
	TestHarnessLocalYamlIgnore(t)
}
