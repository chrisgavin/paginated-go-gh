package merge

func MergeResponses(target *interface{}, overlay *interface{}) {
	switch target := (*target).(type) {
	case []interface{}:
		overlay := (*overlay).([]interface{})
		target = append(target, overlay...)
	case map[string]interface{}:
		overlay := (*overlay).(map[string]interface{})
		for key, value := range overlay {
			if _, ok := target[key]; ok {
				targetValue := target[key]
				MergeResponses(&targetValue, &value)
			} else {
				target[key] = value
			}
		}
	default:
		target = overlay
	}
}
