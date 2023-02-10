package logger

type PrefixHook struct{}

func (h *PrefixHook) Levels() []Level {
	return AllLevels
}

func (h *PrefixHook) Fire(entry *Entry) error {
	for k, v := range entry.Data {
		if len(k) == 0 || v == nil {
			delete(entry.Data, k)
		}
	}
	return nil
}
