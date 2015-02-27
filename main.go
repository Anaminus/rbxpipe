package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	path "path/filepath"
	"strings"
	"time"
)

// C:\Program Files (x86)\Roblox\Versions\version-829e90842df44a35>RobloxStudioBeta -testMode -avatar -fileLocation "C:\Place1.rbxl" -script "loadfile(\"http://www.roblox.com/game/visit.ashx?universeId=0\")() wait() print('TEST') Instance.new('Part',Workspace) wait(3) game:GetService('TestService'):DoCommand('ShutdownClient')"

func findStudio() string {
	//	const allUsers = path.Clean("C:/Program Files (x86)/Roblox/Versions")
	//	const currentUser = path.Join(os.Getenv("LOCALAPPDATA"), "Roblox/Versions")
	return ""
}

// Input is received via the studio's -script option, which can run scripts of
// any length. The script runs very early in the process, so wait() is called
// to ensure that everything has been initialized before running the test.
//
// The output of the test is capture via LogService. Messages are stored in a
// single array, delimited by null characters. Null is safe to use because
// Roblox's API truncates null characters from strings anyway.
//
// When the test finishes, the log data is encoded in the following format: a
// byte indicating if the main thread ran to completion (0) or if an error
// occurred (1). If an error occurred, the remaining data is the error
// message. If there was no error, then the format is the following: a 32-bit
// unsigned integer N, indicating the number of messages; a byte array of
// length N, each byte indicating the type of each message; a run of messages,
// with each message terminated by a null character.
//
// A plugin with an ID of 0 is created via PluginManager.CreatePlugin. The
// plugin's SetSetting method saves the encoded data under a specified key,
// which can be unique per instance, allowing multiple instances to be run at
// once.
//
// Settings are saved in JSON format, though Roblox's JSON encoder is
// unstable. So, to ensure the data is saved correctly, it is encoded in
// base64.
const scriptTemplate = `--[[%-53.53s]]
local function main(key, test)
	local byte = string.byte
	local char = string.char
	local yield = coroutine.yield
	local concat = table.concat

	yield()

	local plugin = PluginManager():CreatePlugin()
	local testService = game:GetService('TestService')

	local base64 do
		local encodeByte = {
			[ 0]='A',[ 1]='B',[ 2]='C',[ 3]='D',[ 4]='E',[ 5]='F',[ 6]='G',[ 7]='H',
			[ 8]='I',[ 9]='J',[10]='K',[11]='L',[12]='M',[13]='N',[14]='O',[15]='P',
			[16]='Q',[17]='R',[18]='S',[19]='T',[20]='U',[21]='V',[22]='W',[23]='X',
			[24]='Y',[25]='Z',[26]='a',[27]='b',[28]='c',[29]='d',[30]='e',[31]='f',
			[32]='g',[33]='h',[34]='i',[35]='j',[36]='k',[37]='l',[38]='m',[39]='n',
			[40]='o',[41]='p',[42]='q',[43]='r',[44]='s',[45]='t',[46]='u',[47]='v',
			[48]='w',[49]='x',[50]='y',[51]='z',[52]='0',[53]='1',[54]='2',[55]='3',
			[56]='4',[57]='5',[58]='6',[59]='7',[60]='8',[61]='9',[62]='-',[63]='_',
		}
		function base64(data)
			local out = {}
			local nout = 1
			local ndata = #data
			for i = 0, ndata-1, 3 do
				local b1 = byte(data, i+1) or 0
				local b2 = byte(data, i+2) or 0
				local b3 = byte(data, i+3) or 0

				local b1r = b1/4
				local b2r = b2/16
				local b3r = b3/64

				local x = b1%%4*16
				local v = b2%%16*4
				local y = b2r%%256 - b2r%%1
				local w = b3r%%256 - b3r%%1

				local x0,x1,x2,x3,x4,x5,x6,x7,x8 = x%%1,x%%2,x%%4,x%%8,x%%16,x%%32,x%%64,x%%128,x%%256
				local y0,y1,y2,y3,y4,y5,y6,y7,y8 = y%%1,y%%2,y%%4,y%%8,y%%16,y%%32,y%%64,y%%128,y%%256
				local v0,v1,v2,v3,v4,v5,v6,v7,v8 = v%%1,v%%2,v%%4,v%%8,v%%16,v%%32,v%%64,v%%128,v%%256
				local w0,w1,w2,w3,w4,w5,w6,w7,w8 = w%%1,w%%2,w%%4,w%%8,w%%16,w%%32,w%%64,w%%128,w%%256

				local or0 = 0
				+ ((x1 - x0 ~= 0 or y1 - y0 ~= 0) and 1 or 0)
				+ ((x2 - x1 ~= 0 or y2 - y1 ~= 0) and 2 or 0)
				+ ((x3 - x2 ~= 0 or y3 - y2 ~= 0) and 4 or 0)
				+ ((x4 - x3 ~= 0 or y4 - y3 ~= 0) and 8 or 0)
				+ ((x5 - x4 ~= 0 or y5 - y4 ~= 0) and 16 or 0)
				+ ((x6 - x5 ~= 0 or y6 - y5 ~= 0) and 32 or 0)
				+ ((x7 - x6 ~= 0 or y7 - y6 ~= 0) and 64 or 0)
				+ ((x8 - x7 ~= 0 or y8 - y7 ~= 0) and 128 or 0)

				local or1 = 0
				+ ((v1 - v0 ~= 0 or w1 - w0 ~= 0) and 1 or 0)
				+ ((v2 - v1 ~= 0 or w2 - w1 ~= 0) and 2 or 0)
				+ ((v3 - v2 ~= 0 or w3 - w2 ~= 0) and 4 or 0)
				+ ((v4 - v3 ~= 0 or w4 - w3 ~= 0) and 8 or 0)
				+ ((v5 - v4 ~= 0 or w5 - w4 ~= 0) and 16 or 0)
				+ ((v6 - v5 ~= 0 or w6 - w5 ~= 0) and 32 or 0)
				+ ((v7 - v6 ~= 0 or w7 - w6 ~= 0) and 64 or 0)
				+ ((v8 - v7 ~= 0 or w8 - w7 ~= 0) and 128 or 0)

				out[nout] = encodeByte[b1r%%256 - b1r%%1]
				out[nout+1] = encodeByte[or0] or '='
				out[nout+2] = ndata-i > 1 and encodeByte[or1] or '='
				out[nout+3] = ndata-i > 2 and encodeByte[b3 %% 64] or '='
				nout = nout + 4
			end
			return concat(out)
		end
	end

	local function uint32(n)
		local s = ''
		for i = 1,4 do
			s = s .. char(n %% 256)
			n = n/256
			n = n - n%%1
		end
		return s
	end

	local messages = {}
	local messageTypes = {}
	local messagesLen = 0
	game:GetService('LogService').MessageOut:connect(function(message, messageType)
		messagesLen = messagesLen + 1
		messageTypes[messagesLen] = string.char(messageType.Value)
		messages[messagesLen] = message .. '\0'
	end)

	local success, err = pcall(test)

	yield()

	local out
	if not success then
		out = {'\1', err}
	else
		out = {'\0', uint32(messagesLen)}
		for i = 1,messagesLen do
			out[#out+1] = messageTypes[i]
		end
		for i = 1,messagesLen do
			out[#out+1] = messages[i]
		end
	end

	plugin:SetSetting(key, base64(concat(out)))

	testService:DoCommand('ShutdownClient')
end

main("%s", function()
%s
end)
`

var input = flag.String("i", "", "A Lua file that will be executed by the studio. If unspecified, then the standard input is read instead.")
var output = flag.String("o", "", "A file to write the results to. If unspecified, then the output is written to the standard output.")
var studio = flag.String("studio", "", "A path to the studio executable. If unspecified, then the studio will be located automatically, assuming it has been installed.")
var file = flag.String("file", "", "A Roblox place file to open with the studio.")
var play = flag.Bool("play", false, "If given, the studio's `Play Solo` state will be mimicked by starting the RunService and inserting a character.")
var timeout = flag.Duration("timeout", time.Duration(-1), "Terminates the studio process after the given duration (e.g. '30s' for 30 seconds). If less than 0, then the timeout is disabled.")
var filter = flag.String("filter", "oiwe", "Filters the output by message type. Each character includes messages of a certain type: 'o' for regular output, 'i' for info, 'w' for warnings, and 'e' for errors.")

var filterMap = map[byte]byte{
	0: 'o',
	1: 'i',
	2: 'w',
	3: 'e',
}

func main() {
	flag.Parse()

	*filter = strings.ToLower(*filter)

	if *studio == "" {
		*studio = findStudio()
		if *studio == "" {
			fmt.Fprintln(os.Stderr, "studio location must be provided (-studio option)")
			return
		}
	}

	args := []string{"-testMode"}

	if *play {
		args = append(args, "-avatar")
	}

	if *file != "" {
		var err error
		*file, err = path.Abs(*file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "could not resolve absolute path of file:", err)
			return
		}

		args = append(args, "-fileLocation", *file)
	}

	var script string
	if *input != "" {
		s, err := ioutil.ReadFile(*input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "could not read input file:", err)
			return
		}

		script = string(s)
	} else {
		s, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "could not read standard input:", err)
			return
		}

		script = string(s)
	}

	if script != "" {
		key := "TestScript"
		args = append(args, "-script", fmt.Sprintf(strings.Replace(scriptTemplate, "\n", " ", -1), key, key, script))
	}

	cmd := exec.Command(*studio, args...)

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "could not run studio:", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "error waiting process:", err)
		return
	}

	var data struct {
		TestScript []byte
	}

	f, err := os.Open(`C:\Users\admin\AppData\Local\Roblox\InstalledPlugins\0\settings.json`)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading output:", err)
		return
	}
	jd := json.NewDecoder(f)
	if err := jd.Decode(&data); err != nil {
		f.Close()
		fmt.Fprintln(os.Stderr, "error decoding output:", err)
		return
	}
	f.Close()

	o := data.TestScript
	if len(o) == 0 {
		fmt.Fprintln(os.Stderr, "format error: data length is 0")
		return
	}

	switch o[0] {
	case 0:
		if len(o) < 5 {
			fmt.Fprintln(os.Stderr, "format error: data length does not accommodate array length")
			return
		}

		o = o[1:]
		messagesLen := int(binary.LittleEndian.Uint32(o[:4]))
		o = o[4:]
		if len(o) < messagesLen {
			fmt.Fprintln(os.Stderr, "format error: data length does not accommodate message type array")
			return
		}

		messageTypes := o[:messagesLen]
		messagesRaw := o[messagesLen:]

		messages := make([]string, 0, messagesLen)
		for i := 0; len(messagesRaw) > 0; i++ {
			n := bytes.IndexByte(messagesRaw, '\000')
			if n < 0 {
				break
			}

			if bytes.IndexByte([]byte(*filter), filterMap[messageTypes[i]]) > -1 {
				messages = append(messages, string(messagesRaw[:n]))
			}

			messagesRaw = messagesRaw[n+1:]
		}

		data := []byte(strings.Join(messages, "\n"))
		if *output != "" {
			if err := ioutil.WriteFile(*output, data, 0666); err != nil {
				fmt.Fprintln(os.Stderr, "could not write to output file:", err)
				return
			}
		} else {
			if _, err := os.Stdout.Write(data); err != nil {
				fmt.Fprintln(os.Stderr, "could not write to stdout:", err)
				return
			}
		}

	case 1:
		fmt.Fprintln(os.Stderr, "an error occurred in the main thread:\n", string(o[1:]))
		return
	}
}
