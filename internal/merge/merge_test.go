package merge_test

import (
	"testing"

	"github.com/chrisgavin/paginated-go-gh/internal/merge"
)

func TestMergeResponsesBasic(t *testing.T) {
	target := []int{1, 2, 3}
	overlay := []int{4, 5, 6}
	merge.MergeResponses(&target, &overlay)
	if len(target) != 6 {
		t.Errorf("Expected target to have length 6, got %d", len(target))
	}

}
