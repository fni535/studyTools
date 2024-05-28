package id

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
)

type id int64

var (
	node *Node
	once sync.Once
	err  error
)

func SnowflakeStr() string {
	return strconv.FormatInt(Snowflake(), 10)
}

func Snowflake() int64 {
	once.Do(func() {
		node, err = newNode(nodeID())
		if err != nil {
			panic(err)
		}
	})
	baseNumber := node.generate().int64()
	return baseNumber
}

func IsContainSnowflake(id string) bool {
	l := len(SnowflakeStr())

	return regexp.MustCompile(fmt.Sprintf(".*[0-9]{%d,}.*", l)).MatchString(id)
}

func (f id) int64() int64 {
	return int64(f)
}

func nodeID() int64 {
	getContainerId := func() (int64, error) {
		content, err := os.ReadFile("/proc/self/cgroup")
		if err != nil {
			return 0, err
		}

		line := bytes.Split(content, []byte("\n"))
		if len(line) == 0 {
			return 0, errors.New("LinesEmpty")
		}

		words := bytes.Split(line[0], []byte("/"))
		if len(words) == 0 {
			return 0, errors.New("WordsEmpty")
		}

		idStr := words[len(words)-1]
		if len(idStr) > 3 {
			idStr = idStr[:3]
		}

		id, err := strconv.ParseInt(string(idStr), 16, 64)
		if err != nil {
			return 0, err
		}

		return id, nil
	}

	id, err := getContainerId()
	if err != nil {
		s, _ := os.Hostname()
		id = int64(xxhash.Sum64String(s)) & nodeMax
		fmt.Printf("SnowflakeNodeID.ByHostnameSupport:%d\n", id)
	} else {
		id = id & nodeMax
		fmt.Printf("SnowflakeNodeID.ByContainerIDSupport:%d\n", id)
	}
	return id
}

var (
	epoch     int64 = time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local).UnixMilli()
	nodeBits  uint8 = 12
	stepBits  uint8 = 10
	mu        sync.Mutex
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	nodeMask        = nodeMax << stepBits
	stepMask  int64 = -1 ^ (-1 << stepBits)
	timeShift       = nodeBits + stepBits
	nodeShift       = stepBits
)

func newNode(node int64) (*Node, error) {
	if nodeBits+stepBits > 22 {
		return nil, fmt.Errorf("remember, you have a total 22 bits to share between Node/Step")
	}

	mu.Lock()
	nodeMax = -1 ^ (-1 << nodeBits)
	nodeMask = nodeMax << stepBits
	stepMask = -1 ^ (-1 << stepBits)
	timeShift = nodeBits + stepBits
	nodeShift = stepBits
	mu.Unlock()

	n := Node{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << nodeBits)
	n.nodeMask = n.nodeMax << stepBits
	n.stepMask = -1 ^ (-1 << stepBits)
	n.timeShift = nodeBits + stepBits
	n.nodeShift = stepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	n.epoch = curTime.Add(time.Unix(epoch/1000, (epoch%1000)*1e6).Sub(curTime))

	return &n, nil
}

type Node struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

func (n *Node) generate() id {
	n.mu.Lock()

	now := time.Since(n.epoch).Milliseconds()

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := id((now)<<n.timeShift |
		(n.node << n.nodeShift) |
		(n.step),
	)

	n.mu.Unlock()
	return r
}
