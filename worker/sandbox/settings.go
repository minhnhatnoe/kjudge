package sandbox

type Settings struct {
	LogSandbox bool
	IgnoreWarning bool
}

var defaultSettings = Settings{LogSandbox: true, IgnoreWarning: false}

type Option func(Settings) Settings

func IgnoreWarnings(ignore bool) Option {
	return func(o Settings) Settings {
		o.IgnoreWarning = ignore
		return o
	}
}

func EnableSandboxLogs(enable bool) Option {
	return func(o Settings) Settings {
		o.LogSandbox = enable
		return o
	}
}

func MakeSettings(options ...Option) Settings {
	setting := defaultSettings
	for _, option := range options {
		setting = option(setting)
	}
	return setting
}
