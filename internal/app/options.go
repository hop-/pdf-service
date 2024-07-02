package app

type HttpOptions struct {
	Enabled bool
	Port    uint16
	Secure  bool
	Cert    string
	Key     string
}

type KafkaOptions struct {
	Enabled         bool
	Host            string
	ConsumerGroupId string
	RequestsTopic   string
	ResponsesTopic  string
}

type Options struct {
	EngineType  string
	Concurrency uint8
	Http        HttpOptions
	Kafka       KafkaOptions
}

// Default options
func defaultOptions() Options {
	return Options{
		EngineType:  "chromedp",
		Concurrency: 4,
		Http: HttpOptions{
			Enabled: false,
			Port:    0,
			Secure:  false,
			Cert:    "",
			Key:     "",
		},
		Kafka: KafkaOptions{
			Enabled:         false,
			Host:            "",
			ConsumerGroupId: "",
			RequestsTopic:   "",
			ResponsesTopic:  "",
		},
	}
}

type OptionModifier func(*Options)

// Option modifiers

func WithHttp(port uint16, secure bool, certFile string, keyFile string) OptionModifier {
	return func(o *Options) {
		o.Http = HttpOptions{
			Enabled: true,
			Port:    port,
			Secure:  secure,
			Cert:    certFile,
			Key:     keyFile,
		}
	}
}

func WithKafka(host string, consumerGroup string, requestsTopic string, responsesTopic string) OptionModifier {
	return func(o *Options) {
		o.Kafka = KafkaOptions{
			Enabled:         true,
			Host:            host,
			ConsumerGroupId: consumerGroup,
			RequestsTopic:   requestsTopic,
			ResponsesTopic:  responsesTopic,
		}
	}
}

func WithConcurrency(concurrent uint8) OptionModifier {
	return func(o *Options) {
		o.Concurrency = concurrent
	}
}

func WithEngine(engine string) OptionModifier {
	return func(o *Options) {
		o.EngineType = engine
	}
}
