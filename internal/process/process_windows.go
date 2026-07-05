//go:build windows

package process

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"env-editor/internal/model"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

func List() ([]model.ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("list processes: %w", err)
	}

	portMap := localPortsByPID()

	items := make([]model.ProcessInfo, 0, len(procs))
	for _, p := range procs {
		name, err := p.Name()
		if err != nil || name == "" {
			name = "—"
		}

		exePath, _ := p.Exe()
		items = append(items, model.ProcessInfo{
			PID:     p.Pid,
			Name:    name,
			ExePath: exePath,
			Ports:   formatPorts(portMap[p.Pid]),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Name == items[j].Name {
			return items[i].PID < items[j].PID
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items, nil
}

func EnvVars(pid int32) ([]model.EnvVar, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("open process: %w", err)
	}

	exists, err := p.IsRunning()
	if err != nil {
		return nil, fmt.Errorf("check process: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("进程 %d 已退出", pid)
	}

	entries, err := p.Environ()
	if err != nil {
		return nil, fmt.Errorf("无法读取进程环境（权限不足或受保护进程）: %w", err)
	}

	vars := make([]model.EnvVar, 0, len(entries))
	for _, entry := range entries {
		name, value, ok := strings.Cut(entry, "=")
		if !ok || name == "" {
			continue
		}
		vars = append(vars, model.EnvVar{Name: name, Value: value})
	}

	sort.Slice(vars, func(i, j int) bool {
		return strings.ToLower(vars[i].Name) < strings.ToLower(vars[j].Name)
	})

	return vars, nil
}

func Kill(pid int32) error {
	if pid <= 0 {
		return fmt.Errorf("无效的进程 ID")
	}
	if pid == 4 {
		return fmt.Errorf("无法结束 System 进程")
	}
	if pid == int32(os.Getpid()) {
		return fmt.Errorf("无法结束当前应用进程")
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("打开进程失败: %w", err)
	}

	exists, err := p.IsRunning()
	if err != nil {
		return fmt.Errorf("检查进程失败: %w", err)
	}
	if !exists {
		return fmt.Errorf("进程 %d 已退出", pid)
	}

	if err := p.Kill(); err != nil {
		return fmt.Errorf("结束进程失败（权限不足或受保护进程）: %w", err)
	}

	return nil
}

func localPortsByPID() map[int32][]uint32 {
	conns := collectConnections()
	if len(conns) == 0 {
		return nil
	}

	listenSets := make(map[int32]map[uint32]struct{})
	otherSets := make(map[int32]map[uint32]struct{})
	for _, c := range conns {
		if c.Pid <= 0 || c.Laddr.Port == 0 {
			continue
		}

		target := otherSets
		if strings.EqualFold(c.Status, "LISTEN") {
			target = listenSets
		}
		if target[c.Pid] == nil {
			target[c.Pid] = make(map[uint32]struct{})
		}
		target[c.Pid][c.Laddr.Port] = struct{}{}
	}

	pids := make(map[int32]struct{})
	for pid := range listenSets {
		pids[pid] = struct{}{}
	}
	for pid := range otherSets {
		pids[pid] = struct{}{}
	}

	result := make(map[int32][]uint32, len(pids))
	for pid := range pids {
		seen := make(map[uint32]struct{})
		ports := make([]uint32, 0)

		appendPorts := func(set map[uint32]struct{}) {
			for port := range set {
				if _, ok := seen[port]; ok {
					continue
				}
				seen[port] = struct{}{}
				ports = append(ports, port)
			}
		}

		appendPorts(listenSets[pid])
		appendPorts(otherSets[pid])

		sort.Slice(ports, func(i, j int) bool { return ports[i] < ports[j] })
		result[pid] = ports
	}

	return result
}

func collectConnections() []net.ConnectionStat {
	conns, err := net.Connections("all")
	if err == nil && len(conns) > 0 {
		return conns
	}

	merged := make([]net.ConnectionStat, 0)
	for _, kind := range []string{"tcp", "udp", "tcp4", "udp4", "tcp6", "udp6"} {
		part, err := net.Connections(kind)
		if err != nil {
			continue
		}
		merged = append(merged, part...)
	}
	return merged
}

func formatPorts(ports []uint32) string {
	if len(ports) == 0 {
		return ""
	}

	parts := make([]string, len(ports))
	for i, port := range ports {
		parts[i] = strconv.FormatUint(uint64(port), 10)
	}
	return strings.Join(parts, ", ")
}
