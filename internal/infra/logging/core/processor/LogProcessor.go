package processor

type LogProcessor struct {
	active bool
}

func NewLogProcessor() *LogProcessor {
	return &LogProcessor{
		active: true,
	}
}

func start() {

}

func waitClose() {

}

func drainLogsForever() {

}
