package messages

type CommandSubmittedMsg struct{ Input string }
type LevelSelectedMsg struct{ Index int }
type LevelPassedMsg struct{}
type NextLevelMsg struct{}
type RestartLevelMsg struct{}
type ShowHintMsg struct{}
