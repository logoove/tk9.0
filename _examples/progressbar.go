package main

import "time"
import . "modernc.org/tk9.0"
import _ "modernc.org/tk9.0/themes/azure"

func main() {
	ActivateTheme("azure light")
	pb := TProgressbar()
	var start *TButtonWidget
	ch := make(chan float64)
	start = TButton(Txt("Start"), Command(func() {
		start.Configure(State("disabled"))
		go func() {
			for i := 0; i <= 100; i += 10 {
				ch <- float64(i)
				time.Sleep(300 * time.Millisecond)
			}
			ch <- -1
		}()
	}))
	NewTicker(100*time.Millisecond, func() {
		select {
		case v := <-ch:
			if v < 0 {
				start.Configure(State("enabled"))
				break
			}

			pb.Configure(Value(v))
		default:
		}
	})
	Pack(pb,
		start,
		TExit(),
		Padx("1m"), Pady("2m"), Ipadx("1m"), Ipady("1m"))
	App.SetResizable(false, false)
	App.Wait()
}
