package sandbox

type Settings struct {
	LogSandbox    bool
	IgnoreWarning bool
}

var defaultSettings = Settings{LogSandbox: true, IgnoreWarning: false}

// Option represents a sandbox option.
type Option func(Settings) Settings

// IgnoreWarnings sets whether sandbox warnings should be silenced.
func IgnoreWarnings(ignore bool) Option {
	return func(o Settings) Settings {
		o.IgnoreWarning = ignore
		return o
	}
}

// EnableSandboxLogs sets whether sandbox logs should be printed.
func EnableSandboxLogs(enable bool) Option {
	return func(o Settings) Settings {
		o.LogSandbox = enable
		return o
	}
}

// MakeSettings creates a Settings object for sandboxes using options passed
func MakeSettings(options ...Option) Settings {
	setting := defaultSettings
	for _, option := range options {
		setting = option(setting)
	}
	return setting
}
