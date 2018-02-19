package interruptor

type Logger struct {
	//   stdout *log.Logger

	//   stderr *log.Logger
}

func DefaultLogger() *Logger {
	return &Logger{
	//     stdout: log.Logger(os.Stdout),
	//     stderr: log.Logger(os.Stderr),
	}
}

func (l *Logger) Error(msg interface{}) {
	//   l.stderr.Println(string)
}
