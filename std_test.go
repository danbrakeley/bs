package bs

import (
	"os"
	"testing"

	"github.com/magefile/mage/mg"
)

const altEnvVarName = "BS_VERBOSE_TEST"

func Test_IsVerbose(t *testing.T) {
	os.Unsetenv(altEnvVarName)
	os.Unsetenv(initialVerboseEnvVar)

	if mg.Verbose() {
		t.Errorf(`expected mg.Verbose() to be false when "%s" is unset`, initialVerboseEnvVar)
	}
	if IsVerbose() {
		t.Errorf(`expected IsVerbose() to be false when "%s" is unset`, initialVerboseEnvVar)
	}

	os.Setenv(initialVerboseEnvVar, "true")

	if !mg.Verbose() {
		t.Errorf(`unable to make mg.Verbose() return true by setting "%s" to "true"`, initialVerboseEnvVar)
	}
	if !IsVerbose() {
		t.Errorf(`unable to make IsVerbose() return true by setting "%s" to "true"`, initialVerboseEnvVar)
	}

	os.Unsetenv(initialVerboseEnvVar)

	if mg.Verbose() {
		t.Errorf(`unable to clear mg.Verbose() when "%s" is unset`, initialVerboseEnvVar)
	}
	if IsVerbose() {
		t.Errorf(`unable to clear IsVerbose() when "%s" is unset`, initialVerboseEnvVar)
	}

	SetVerbose(true)

	if !mg.Verbose() {
		t.Errorf(`unable to make mg.Verbose() return true by calling SetVerbose(true)`)
	}
	if !IsVerbose() {
		t.Errorf(`unable to make IsVerbose() return true by calling SetVerbose(true)`)
	}

	SetVerbose(false)

	if mg.Verbose() {
		t.Errorf(`unable to make mg.Verbose() return false by calling SetVerbose(false)`)
	}
	if IsVerbose() {
		t.Errorf(`unable to make IsVerbose() return false by calling SetVerbose(false)`)
	}

	SetVerbose(true)
	SetVerboseEnvVarName(altEnvVarName)

	if !mg.Verbose() {
		t.Errorf(`expected mg.Verbose() to still be true after calling SetVerboseEnvVarName("%s")`, altEnvVarName)
	}
	if !IsVerbose() {
		t.Errorf(`expected IsVerbose() to still return true after calling SetVerboseEnvVarName`)
	}

	os.Unsetenv(initialVerboseEnvVar)

	if mg.Verbose() {
		t.Errorf(`expected mg.Verbose() to be false after calling SetVerboseEnvVarName("%s") and then unsetting "%s"`, altEnvVarName, initialVerboseEnvVar)
	}
	if !IsVerbose() {
		t.Errorf(`expected IsVerbose() after SetVerboseEnvVarName("%s") to still return true after unsetting "%s"`, altEnvVarName, initialVerboseEnvVar)
	}

	SetVerbose(false)

	if mg.Verbose() {
		t.Errorf(`expected mg.Verbose() to still be false`)
	}
	if IsVerbose() {
		t.Errorf(`expected IsVerbose() to be false after SetVerbose(false)`)
	}
}
