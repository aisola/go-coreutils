package main


import "bufio"
import "bytes"
import "fmt"
import "io/ioutil"
import "strconv"
import "strings"
import "syscall"
import "time"


type Load struct {
	L1, L5, L15 float64
}

type Uptime struct {
	Time float64
}

func (self *Load) Get() error {
	line, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil
	}

	f := strings.Fields(string(line))

	self.L1, _ = strconv.ParseFloat(f[0], 64)
	self.L5, _ = strconv.ParseFloat(f[1], 64)
	self.L15, _ = strconv.ParseFloat(f[2], 64)

	return nil
}

func (self *Uptime) Get() error {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return err
	}

	self.Time = float64(sysinfo.Uptime)

	return nil
}

func (self *Uptime) Format() string {
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	uptime := uint64(self.Time)

	days := uptime / (60 * 60 * 24)

	if days != 0 {
		s := ""
		if days > 1 {
			s = "s"
		}
		fmt.Fprintf(w, "%d day%s, ", days, s)
	}

	minutes := uptime / 60
	hours := minutes / 60
	hours %= 24
	minutes %= 60

	fmt.Fprintf(w, "%2d:%02d", hours, minutes)

	w.Flush()
	return buf.String()
}

func Users() int { return 0 }

func main() {
	up := Uptime{}
	up.Get()
	load := Load{}
	load.Get()

	fmt.Printf(" %s up %s load average: %.2f, %.2f, %.2f\n",
		time.Now().Format("15:04:05"),
		up.Format(),
		load.L1, load.L5, load.L15)
}