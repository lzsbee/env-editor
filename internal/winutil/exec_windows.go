//go:build windows

package winutil

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func ShellExecuteRunAs(exe, args string) error {
	shell32 := windows.NewLazySystemDLL("shell32.dll")
	shellExecuteW := shell32.NewProc("ShellExecuteW")

	verbPtr, err := windows.UTF16PtrFromString("runas")
	if err != nil {
		return err
	}
	exePtr, err := windows.UTF16PtrFromString(exe)
	if err != nil {
		return err
	}
	argsPtr, err := windows.UTF16PtrFromString(args)
	if err != nil {
		return err
	}

	const swHide = 0
	ret, _, callErr := shellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argsPtr)),
		0,
		swHide,
	)

	if ret <= 32 {
		if ret == 0 {
			return fmt.Errorf("已取消 UAC 授权或未授予管理员权限")
		}
		return fmt.Errorf("无法启动 UAC 提权（错误码 %d）: %v", ret, callErr)
	}
	if ret == 1223 {
		return fmt.Errorf("已取消 UAC 授权")
	}
	return nil
}

func WaitForResultFile(resultFile string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(resultFile)
		if err == nil {
			_ = os.Remove(resultFile)
			result := strings.TrimSpace(string(data))
			if result == "OK" {
				return nil
			}
			if result == "" {
				return fmt.Errorf("提权操作失败")
			}
			return fmt.Errorf("%s", result)
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("操作超时，请查看是否被 UAC 拦截或杀毒软件阻止")
}

func BroadcastEnvChange() error {
	user32 := windows.NewLazySystemDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	const (
		hwndBroadcast   = 0xFFFF
		wmSettingChange = 0x001A
		smtoAbortifHung = 0x0002
	)

	envName, err := windows.UTF16PtrFromString("Environment")
	if err != nil {
		return err
	}

	ret, _, callErr := sendMessageTimeout.Call(
		uintptr(hwndBroadcast),
		uintptr(wmSettingChange),
		0,
		uintptr(unsafe.Pointer(envName)),
		uintptr(smtoAbortifHung),
		5000,
		0,
	)

	if ret == 0 {
		if callErr != nil && callErr != syscall.Errno(0) {
			return fmt.Errorf("broadcast env change: %w", callErr)
		}
	}

	return nil
}

func IsProcessElevated() bool {
	var token windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token); err != nil {
		return false
	}
	defer token.Close()
	return token.IsElevated()
}
