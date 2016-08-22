
### FSM 入门示例代码：

```

package main

import (
	"fmt"
	"src/github.com/looplab/fsm"
)

type Door struct {
	To  string
	FSM *fsm.FSM
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
			{Name: "abandon", Src: []string{"closed"}, Dst: "abandoned"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.enterState(e) },
			"before_abandon": func(e *fsm.Event) { d.beforeAbandon(e) },
			"before_close": func(e *fsm.Event) { d.beforeClose(e) },
		},
	)

	return d
}

func (d *Door) enterState(e *fsm.Event) {
	fmt.Println("I am enter state %s")
}

func (d *Door) beforeAbandon(e *fsm.Event) {
	fmt.Println("I am before abandon")
}

func (d *Door) beforeClose(e *fsm.Event) {
	fmt.Println("I am before close")
}



func main() {
	door := NewDoor("heaven")

	err := door.FSM.Event("open")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(door.FSM.Current());

	err = door.FSM.Event("close")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(door.FSM.Current());

	err = door.FSM.Event("abandon")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(door.FSM.Current());
}
```

### 代码运行结果
```
I am enter state %s closed
open
I am before close
I am enter state %s open
closed
I am before abandon
I am enter state %s closed
abandoned
```