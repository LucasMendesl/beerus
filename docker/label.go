package docker

const beerusServiceLabel = "com.github.lucasmendesl.beerus.service"

type labeler interface {
	GetLabels() map[string]string
}

func (i Image) GetLabels() map[string]string {
	return i.Labels
}

// removeIgnored filters out items from the given slice that have labels matching any of the given labels or the built-in "com.github.lucasmendesl.beerus.service" label.
//
// It takes a slice of items that satisfy the labeler interface and a variable number of label strings.
// The labeler interface requires a GetLabels() method that returns a map of labels.
// The filtered slice is returned as the result.
func removeIgnored[T labeler](items []T, labels ...string) []T {
	ignoredList := []string{beerusServiceLabel}
	ignoredList = append(ignoredList, labels...)

	filteredItems := make([]T, 0, len(items))
	for _, item := range items {
		if !hasIgnoredLabels(item.GetLabels(), ignoredList) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func hasIgnoredLabels(labels map[string]string, ignoredLabels []string) bool {
	for _, label := range ignoredLabels {
		if _, ok := labels[label]; ok {
			return true
		}
	}
	return false
}
