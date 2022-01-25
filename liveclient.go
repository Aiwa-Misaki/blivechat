package blivechat

import (
	"fmt"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/aynakeya/blivedm"
)

var Client *blivedm.BLiveWsClient
var ClientSet = make(chan int)

func SetupDanmuClient(g *gocui.Gui, cl *blivedm.BLiveWsClient) {
	Client = cl
	Client.RegHandler(blivedm.CmdDanmaku, func(context *blivedm.Context) {
		msg, _ := context.ToDanmakuMessage()
		danmuV, _ := g.View(ViewDanmu)
		if danmuV == nil {
			return
		}
		if danmuV.LinesHeight() > 2048 {
			tmp := danmuV.ViewBufferLines()[danmuV.LinesHeight()-1024 : danmuV.LinesHeight()]
			danmuV.Clear()
			for _, l := range tmp {
				if (len(l)) > 0 {
					viewPrintln(danmuV, l)
				}
			}
		}
		PrintDanmu(danmuV, msg)
		g.Update(func(gui *gocui.Gui) error {
			return nil
		})
	})

	err := g.SetKeybinding(ViewSend, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		debugV, err := g.View(ViewDebug)
		if err != nil {
			return err
		}
		msg := v.Buffer()
		if len(msg) == 0 {
			return nil
		}
		v.Clear()
		viewPrintWithTime(debugV, fmt.Sprintf("try send msg: %s", msg))
		if Client.Account.UID == 0 {
			viewPrintWithTime(debugV, "Send Msg fail, not login")
			return v.SetCursor(0, 0)
		}
		message, err := Client.SendMessage(blivedm.DanmakuSendForm{
			Bubble:   SendFormConfig.Bubble,
			Message:  msg,
			Mode:     SendFormConfig.Mode,
			Color:    SendFormConfig.Color,
			Fontsize: SendFormConfig.Fontsize,
			Rnd:      int(time.Now().Unix()),
		})
		if err != nil {
			viewPrintWithTime(debugV, fmt.Sprintf("Send Msg fail, %s", err))
			return v.SetCursor(0, 0)
		}
		viewPrintWithTime(debugV, fmt.Sprintf("send result - code:%d msg:%s", message.Code, message.Message))
		return v.SetCursor(0, 0)
	})
	if err != nil {
		return
	}

	g.Update(func(gui *gocui.Gui) error {
		debugV, err := g.View(ViewDebug)
		if err != nil {
			return err
		}
		go func() {
			viewPrintWithTime(debugV, fmt.Sprintf("try get room info"))
			viewPrintWithTime(debugV, fmt.Sprintf("GetRoomInfo: %t", Client.GetRoomInfo()))
			g.Update(func(gui *gocui.Gui) error { return nil })
			viewPrintWithTime(debugV, fmt.Sprintf("try danmu room info"))
			viewPrintWithTime(debugV, fmt.Sprintf("GetDanmuInfo: %t", Client.GetDanmuInfo()))
			g.Update(func(gui *gocui.Gui) error { return nil })
			viewPrintWithTime(debugV, fmt.Sprintf("try connect to danmu server"))
			viewPrintWithTime(debugV, fmt.Sprintf("ConnectToDanmuServer: %t", Client.ConnectDanmuServer()))
			g.Update(func(gui *gocui.Gui) error { return nil })
			roomV, err := g.View(ViewRoom)
			if err != nil {
				//viewPrintWithTime(debugV,err)
				return
			}
			var upname string
			if info, err := blivedm.ApiGetUpInfo(Client.RoomInfo.Uid); err != nil {
				upname = "Unknown"
			} else {
				upname = info.Data.Name
			}
			viewPrint(roomV,
				fmt.Sprintf("%s > %s | %s (%d) - Live: %t | %s (%d) | login as uid=%d",
					Client.RoomInfo.ParentAreaName, Client.RoomInfo.AreaName,
					Client.RoomInfo.Title, Client.RoomInfo.RoomId, Client.RoomInfo.LiveStatus == 1,
					upname, Client.RoomInfo.Uid,
					Client.Account.UID))
			width, height := g.Size()
			PrintToDebug(g, fmt.Sprintf("room info widget size: %d*%d", width, height))
			ClientSet <- 1
			g.Update(func(gui *gocui.Gui) error { return nil })
		}()
		return nil
	})
	return
}
