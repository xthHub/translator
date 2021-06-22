package parser

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type CSPData struct {
	Procs       map[int]string
	MsgChannels map[int]string
	Event       map[int]struct{}
	RootProc    map[int]struct{}
	CatchEvents map[int]struct{}

	Next   map[int]int
	End    map[int]int
	Init   map[int]map[int]struct{}
	Enable map[int][]int

	InitEvents map[int]map[int]struct{}
	InitProcs  []int
}

func Parse(filePath string) (*CSPData, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解析数据
	strData := string(data)
	dataArr := strings.Split(strData, "---")
	procs := dataArr[0]
	msgs := dataArr[1]
	events := dataArr[2]
	rootProcs := dataArr[3]
	next := dataArr[4]
	end := dataArr[5]
	init := dataArr[6]
	enable := dataArr[7]

	ret := CSPData{
		Event:       make(map[int]struct{}),
		RootProc:    make(map[int]struct{}),
		Procs:       make(map[int]string),
		MsgChannels: make(map[int]string),
		CatchEvents: make(map[int]struct{}),
		Next:        make(map[int]int),
		End:         make(map[int]int),
		Init:        make(map[int]map[int]struct{}),
		Enable:      make(map[int][]int),
		InitEvents:  make(map[int]map[int]struct{}),
		InitProcs:   make([]int, 0),
	}
	// 解析proc
	ret.parseIDValue(procs, false)
	// 解析消息通道
	ret.parseIDValue(msgs, true)
	// 标记event
	events = strings.Trim(events, "\n")
	eventArr := strings.Split(events, ",")
	for _, event := range eventArr {
		if event == "" {
			continue
		}
		id, err := strconv.Atoi(event)
		if err != nil {
			return nil, err
		}
		ret.Event[id] = struct{}{}
	}
	// 标记参与者进程
	rootProcs = strings.Trim(rootProcs, "\n")
	procArr := strings.Split(rootProcs, ",")
	for _, proc := range procArr {
		if proc == "" {
			continue
		}
		id, err := strconv.Atoi(proc)
		if err != nil {
			return nil, err
		}
		ret.RootProc[id] = struct{}{}
	}
	// 记录order
	parseOrder(next, ret.Next)
	parseOrder(end, ret.End)
	parseListOrder(init, ret.Init)
	// 记录event间的send-catch关系
	enableOrders := strings.Split(enable, "\n")
	for _, tmp := range enableOrders {
		enableOrder := strings.Split(tmp, ":")
		if len(enableOrder) != 3 {
			continue
		}
		sendEvents := strings.Split(enableOrder[1], ",")
		catchEvents := strings.Split(enableOrder[2], ",")
		sids := make([]int, 0)
		for _, ev := range sendEvents {
			id, err := strconv.Atoi(ev)
			if err != nil {
				continue
			}
			sids = append(sids, id)
		}
		cids := make([]int, 0)
		for _, cev := range catchEvents {
			cid, err := strconv.Atoi(cev)
			if err != nil {
				continue
			}
			cids = append(cids, cid)
			ret.CatchEvents[cid] = struct{}{}
		}
		if len(sids) == 0 || len(cids) == 0 {
			continue
		}
		for _, sid := range sids {
			ret.Enable[sid] = cids
		}
	}
	hasBefore := make(map[int]struct{})
	for _, next := range ret.Next {
		hasBefore[next] = struct{}{}
	}
	for _, enables := range ret.Enable {
		for _, enable := range enables {
			hasBefore[enable] = struct{}{}
		}
	}
	// 只有每个参与者的init才可以是初始进程
	initProcs := make(map[int]struct{})
	for rootProc := range ret.RootProc {
		inits := ret.GetInitEvents(rootProc)
		for init := range inits {
			initProcs[init] = struct{}{}
		}
	}
	for event := range ret.Event {
		if _, has := initProcs[event]; !has {
			continue
		}
		if _, has := hasBefore[event]; !has {
			ret.InitProcs = append(ret.InitProcs, event)
		}
	}
	return &ret, nil
}

func parseOrder(orders string, orderMap map[int]int) error {
	orderArr := strings.Split(orders, "\n")
	for _, order := range orderArr {
		pos := strings.Index(order, "=")
		if pos <= 0 {
			continue
		}
		strID := order[:pos]
		id, err := strconv.Atoi(strID)
		if err != nil {
			return err
		}

		strRightID := order[pos+1:]
		rightID, err := strconv.Atoi(strRightID)
		if err != nil {
			return err
		}
		orderMap[id] = rightID
	}
	return nil
}

func (c *CSPData) parseIDValue(data string, isMsgChannel bool) error {
	procArr := strings.Split(data, "\n")
	for _, proc := range procArr {
		pos := strings.Index(proc, "=")
		if pos <= 0 {
			continue
		}
		strID := proc[:pos]
		id, err := strconv.Atoi(strID)
		if err != nil {
			return err
		}
		name := proc[pos+1:]
		if !isMsgChannel {
			c.Procs[id] = name
		} else {
			c.MsgChannels[id] = name
		}
	}
	return nil
}

func parseListOrder(orders string, orderMap map[int]map[int]struct{}) error {
	orderArr := strings.Split(orders, "\n")
	for _, order := range orderArr {
		pos := strings.Index(order, "=")
		if pos <= 0 {
			continue
		}
		strID := order[:pos]
		id, err := strconv.Atoi(strID)
		if err != nil {
			return err
		}
		if _, ok := orderMap[id]; !ok {
			orderMap[id] = make(map[int]struct{})
		}

		rightArr := strings.Split(order[pos+1:], ",")
		for _, right := range rightArr {
			if right == "" {
				continue
			}
			rightID, err := strconv.Atoi(right)
			if err != nil {
				return err
			}
			orderMap[id][rightID] = struct{}{}
		}
	}
	return nil
}

func (c *CSPData) GetInitEvents(proc int) map[int]struct{} {
	if es, ok := c.InitEvents[proc]; ok {
		return es
	}
	if _, ok := c.Event[proc]; ok {
		return map[int]struct{}{proc: {}}
	}
	es := make(map[int]struct{})
	for k := range c.Init[proc] {
		data := c.GetInitEvents(k)
		for nk := range data {
			es[nk] = struct{}{}
		}
	}

	c.InitEvents[proc] = es
	return es
}
