package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPCounter(t *testing.T) {
	ipc := NewIPCounter(10)

	ipc.RequestHandled("1.2.3.4")
	ipc.RequestHandled("1.2.3.4")
	ipc.RequestHandled("1.2.3.4")
	ipc.RequestHandled("1.2.3.5")
	ipc.RequestHandled("1.2.3.6")
	ipc.RequestHandled("1.2.3.4")
	ipc.RequestHandled("1.2.3.5")
	ipc.RequestHandled("1.2.3.8")
	ipc.RequestHandled("1.2.3.9")
	ipc.RequestHandled("1.2.3.8")
	ipc.RequestHandled("1.2.3.10")
	ipc.RequestHandled("1.2.3.11")
	ipc.RequestHandled("1.2.3.12")
	ipc.RequestHandled("1.2.3.13")
	ipc.RequestHandled("1.2.3.14")
	ipc.RequestHandled("1.2.3.15")
	ipc.RequestHandled("1.2.3.16")
	ipc.RequestHandled("1.2.3.17")
	ipc.RequestHandled("1.2.3.18")
	ipc.RequestHandled("1.2.3.9")

	top := ipc.Top100()
	require.Equal(t, top, []*TopEntry{
		{IP: mustIPToInt("1.2.3.4"), Count: 4},
		{IP: mustIPToInt("1.2.3.9"), Count: 2},
		{IP: mustIPToInt("1.2.3.8"), Count: 2},
		{IP: mustIPToInt("1.2.3.5"), Count: 2},
		{IP: mustIPToInt("1.2.3.6"), Count: 1},
		{IP: mustIPToInt("1.2.3.10"), Count: 1},
		{IP: mustIPToInt("1.2.3.11"), Count: 1},
		{IP: mustIPToInt("1.2.3.12"), Count: 1},
		{IP: mustIPToInt("1.2.3.13"), Count: 1},
		{IP: mustIPToInt("1.2.3.18"), Count: 1},
	})
}

func mustIPToInt(ip string) uint32 {
	n, err := ipToInt(ip)
	if err != nil {
		panic(err)
	}
	return n
}
