package main

import (
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

// IPCounter manages counts of ip hits and maintains a list of top IPs with
// their counts.
type IPCounter struct {
	Counts map[uint32]uint64
	Top    *list.List

	// How many top entries to maintain
	topCount int
}

// TopEntry is used for the top list
type TopEntry struct {
	IP    uint32
	Count uint64
}

func (t *TopEntry) String() string {
	return fmt.Sprintf("IP: %s, Count: %d", intToIP(t.IP), t.Count)
}

func NewIPCounter(topCount int) *IPCounter {
	return &IPCounter{
		Counts:   map[uint32]uint64{},
		Top:      list.New(),
		topCount: topCount,
	}
}

func (i *IPCounter) RequestHandled(ip string) error {
	intIP, err := ipToInt(ip)
	if err != nil {
		return err
	}

	// Increase the hit count
	i.Counts[intIP]++

	// Update the top 100 list
	count := i.Counts[intIP]
	endCount := uint64(0)
	end := i.Top.Back()
	if end != nil {
		endCount = end.Value.(*TopEntry).Count
	}
	// We will update the top list if we haven't reached the wanted length of
	// the top list, or the count for this IP is greater than or equal to the
	// current lowest count, which is at the end of the list.
	if i.Top.Len() < i.topCount || count >= endCount {
		i.updateTopList(intIP, count)
	}

	return nil
}

func (i *IPCounter) updateTopList(ip uint32, count uint64) {
	entry := &TopEntry{
		IP:    ip,
		Count: count,
	}

	topList := i.Top
	topLen := topList.Len()
	// Handle the easy case first
	if topLen == 0 {
		topList.PushFront(entry)
		return
	}

	// Find where this should go, starting from the back of the list as likely
	// this will be inserted near there
	var possibleLocation *list.Element
	for e := topList.Back(); e != nil; e = e.Prev() {
		entry := e.Value.(*TopEntry)
		if entry.IP == ip {
			// We found a previous entry for this IP, just update it
			entry.Count = count

			// We may need to adjust position to maintain count order.
			//
			// First, see if this just needs to go to the front
			if topList.Front().Value.(*TopEntry).Count <= count {
				topList.MoveBefore(e, topList.Front())
				return
			}
			// Otherwise backtrack until we find a value higher
			prev := e.Prev()
			if prev != nil && count >= prev.Value.(*TopEntry).Count {
				for prev != nil && count >= prev.Value.(*TopEntry).Count {
					prev = prev.Prev()
				}
				if prev != nil {
					// Now put this after the value we know is higher
					topList.MoveAfter(e, prev)
				}
			}

			return
		}
		if entry.Count >= count {
			possibleLocation = e
			break
		}
	}
	if possibleLocation != nil {
		if topLen < i.topCount {
			// If the list is not full yet, insert this after the possible location
			topList.InsertAfter(entry, possibleLocation)
		} else {
			// Otherwise replace the IP at this position. More recent IPs with
			// the same count will replace those already in the list.
			possibleLocation.Value.(*TopEntry).IP = ip
			possibleLocation.Value.(*TopEntry).Count = count
		}
	} else {
		// We couldn't find a place to put this IP, this may be because the list
		// is not yet full and it should be added to the end.
		if topLen < i.topCount {
			topList.InsertAfter(entry, topList.Back())
		}
	}
}

// This is called Top100 but in theory this struct can maintain a top list of
// any length, like in the tests
func (i *IPCounter) Top100() []*TopEntry {
	// This is now trivial and takes almost no time
	top := make([]*TopEntry, 0, i.Top.Len())

	for e := i.Top.Front(); e != nil; e = e.Next() {
		top = append(top, e.Value.(*TopEntry))
	}

	return top
}

func (i *IPCounter) Clear() {
	i.Counts = map[uint32]uint64{}
	// This will clear the list
	i.Top.Init()
}

func ipToInt(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}

	ip = ip.To4()

	return binary.BigEndian.Uint32(ip), nil
}

func intToIP(ipInt uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipInt)
	ip := net.IP(ipByte)
	return ip.String()
}

func makeIP(i int) string {
	return intToIP(uint32(i))
}

// This main can be used to insert some random IPs and make sure that the list
// is in the right order and see performance.
func main() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	ipc := NewIPCounter(100)

	start := time.Now()
	// Adjust these to try different variables to test performance
	counts := 500_000
	numberOfIPs := 500
	for i := 0; i < counts; i++ {
		ip := makeIP(r1.Intn(numberOfIPs))
		ipc.RequestHandled(ip)
	}
	totalTime := time.Since(start)
	fmt.Printf("Took %s to generate %d counts for a total of %d IPs\n", totalTime, counts, len(ipc.Counts))
	fmt.Printf("Average time to handle each IP is %s\n", totalTime/time.Duration(counts))

	start = time.Now()
	top := ipc.Top100()
	fmt.Printf("Took %s to get top 100\n", time.Since(start))
	for i, c := range top {
		fmt.Printf("%d. %s\n", i+1, c)
	}
}
