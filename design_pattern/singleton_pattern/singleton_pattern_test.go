/**
 * @Author: Sean
 * @Date: 2022/5/19 16:02
 */

package singleton_pattern

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	s["name"] = "lee"

	// 验证唯一性
	s1 := New()
	if s1["name"] != "lee" {
		t.Error("singleton pattern error")
	}

	// change name
	s1["name"] = "anne"
	if s["name"] != "anne" {
		t.Error("singleton pattern error")
	}
}
