package log

type Level int

const (
	// Level
	LOG_EMERG Level = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

var LevelNames = map[Level]string {
	LOG_EMERG: "EMERGENCY",
	LOG_ALERT: "ALERT",
	LOG_CRIT: "CRIT",
	LOG_ERR: "ERROR",
	LOG_WARNING: "Warning",
	LOG_NOTICE: "Notice",
	LOG_INFO: "Info",
	LOG_DEBUG: "Debug",
}