// +build freebsd

package gopsutil

import (
	"os/exec"
	"strconv"
	"strings"
)

func NetIOCounters(pernic bool) ([]NetIOCountersStat, error) {
	out, err := exec.Command("/usr/bin/netstat", "-ibdn").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	ret := make([]NetIOCountersStat, 0, len(lines)-1)

	for _, line := range lines {
		values := strings.Fields(line)
		if len(values) < 1 || values[0] == "Name" {
			continue
		}
		base := 1
		// sometimes Address is ommitted
		if len(values) < 13 {
			base = 0
		}

		parsed := make([]uint64, 0, 8)
		vv := []string{
			values[base+3],  // PacketsRecv
			values[base+4],  // Errin
			values[base+5],  // Dropin
			values[base+6],  // BytesRecvn
			values[base+7],  // PacketSent
			values[base+8],  // Errout
			values[base+9],  // BytesSent
			values[base+11], // Dropout
		}
		for _, target := range vv {
			if target == "-" {
				parsed = append(parsed, 0)
				continue
			}

			t, err := strconv.ParseUint(target, 10, 64)
			if err != nil {
				return nil, err
			}
			parsed = append(parsed, t)
		}

		n := NetIOCountersStat{
			Name:        values[0],
			PacketsRecv: parsed[0],
			Errin:       parsed[1],
			Dropin:      parsed[2],
			BytesRecv:   parsed[3],
			PacketsSent: parsed[4],
			Errout:      parsed[5],
			BytesSent:   parsed[6],
			Dropout:     parsed[7],
		}
		ret = append(ret, n)
	}

	return ret, nil
}
