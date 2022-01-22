package blivechat

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/aynakeya/blivedm"
	"github.com/spf13/cast"
	"time"
)

func viewPrintln(v *gocui.View, a ...interface{}) {
	fmt.Fprintln(v, a...)
}

func viewPrint(v *gocui.View, a interface{}) {
	fmt.Fprint(v, a)
}

func viewPrintWithTime(v *gocui.View, a interface{}) {
	debugTime := time.Now().Format("2006/01/02 15:04:05")
	debugMsg := SetForegroundColor(RGB{90, 90, 90}, cast.ToString(a))
	fmt.Fprintf(v, "%s >\n%s\n", SetForegroundColor(RGB{70, 70, 70}, debugTime), debugMsg)
}

func PrintToDebug(g *gocui.Gui, a interface{}) {
	view, err := g.View(ViewDebug)
	if err != nil {
		return
	}
	viewPrintWithTime(view, a)
}

func printDanmuColor(v *gocui.View, msg blivedm.DanmakuMessage, paintMsgColor bool) {
	nameColor := HexToRGB(msg.UnameColor)
	if !paintMsgColor {
		nameColor = RGB{220, 150, 180}
	}
	name := SetForegroundColor(nameColor, msg.Uname)
	medal := "[Unknown](0)"
	if len(msg.MedalName) > 0 {
		medal = SetForegroundColor(IntToRGB(int(msg.MedalColor)),
			fmt.Sprintf("[%s](%d)", msg.MedalName, msg.MedalLevel))
	}
	if paintMsgColor {
		viewPrintln(v,
			fmt.Sprintf("%s %s: %s",
				medal, name, msg.Msg))
	} else {
		viewPrintln(v,
			fmt.Sprintf("%s %s: %s",
				medal, name, SetForegroundColor(RGB{178, 223, 238}, msg.Msg)))
	}

}

func PrintDanmu(v *gocui.View, msg blivedm.DanmakuMessage) {
	printDanmuColor(v, msg, Config.VisualColorMode)
}
