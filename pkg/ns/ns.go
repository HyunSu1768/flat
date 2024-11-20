package ns

import (
	"github.com/vishvananda/netns"
	"runtime"
	"testing"
)

func SetUpNetLinkTest(t *testing.T) func() {
	runtime.LockOSThread()
	var err error
	ns, err := netns.New()
	if err != nil {
		t.Fatalf("새로운 ns 를 생성하는데 실패하였습니다 : %v", err)
	}

	return func() {
		ns.Close()
		runtime.UnlockOSThread()
	}
}
