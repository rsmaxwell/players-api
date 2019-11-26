package debug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/rsmaxwell/players-api/internal/basic/version"
	"github.com/rsmaxwell/players-api/internal/common"
)

// Package type
type Package struct {
	name  string
	level int
}

// Function type
type Function struct {
	pkg   *Package
	name  string
	level int
}

const (
	// ErrorLevel trace level
	ErrorLevel = 10

	// WarningLevel trace level
	WarningLevel = 20

	// InfoLevel trace level
	InfoLevel = 30

	// APILevel trace level
	APILevel = 40

	// VerboseLevel trace level
	VerboseLevel = 50

	minUint uint = 0 // binary: all zeroes

	maxUint = ^minUint // binary: all ones

	maxInt = int(maxUint >> 1) // binary: all ones except high bit

	minInt = ^maxInt // binary: all zeroes except high bit

)

var (
	level                int
	defaultPackageLevel  int
	defaultFunctionLevel int
	dumpRoot             string
)

func init() {
	level, _ = getEnvInteger("DEBUG_LEVEL", InfoLevel)
	defaultPackageLevel, _ = getEnvInteger("DEBUG_DEFAULT_PACKAGE_LEVEL", InfoLevel)
	defaultFunctionLevel, _ = getEnvInteger("DEBUG_DEFAULT_FUNCTION_LEVEL", InfoLevel)

	dir, ok := os.LookupEnv("DEBUG_DUMP_DIR")
	if ok {
		dumpRoot = dir
	} else {
		dumpRoot = common.HomeDir() + "/players-api-dump"
	}
}

func getEnvInteger(name string, def int) (int, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return def, nil
	}
	return strconv.Atoi(value)
}

func getEnvString(name string, def string) (string, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return def, nil
	}
	return value, nil
}

// NewPackage function
func NewPackage(name string) *Package {
	m := &Package{name: name, level: defaultPackageLevel}

	value, ok := os.LookupEnv("DEBUG_PACKAGE_LEVEL_" + name)
	if ok {
		number, err := strconv.Atoi(value)
		if err == nil {
			m.level = number
		}
	}

	return m
}

// NewFunction function
func NewFunction(pkg *Package, name string) *Function {

	d := &Function{pkg: pkg, name: name, level: defaultFunctionLevel}

	value, ok := os.LookupEnv("DEBUG_FUNCTION_LEVEL_" + pkg.name + "_" + name)
	if ok {
		number, err := strconv.Atoi(value)
		if err == nil {
			d.level = number
		}
	}

	return d
}

// --------------------------------------------------------

// DebugError prints an 'error' message
func (f *Function) DebugError(format string, a ...interface{}) {
	f.Debug(ErrorLevel, format, a...)
}

// DebugWarn prints an 'warning' message
func (f *Function) DebugWarn(format string, a ...interface{}) {
	f.Debug(WarningLevel, format, a...)
}

// DebugInfo prints an 'info' message
func (f *Function) DebugInfo(format string, a ...interface{}) {
	f.Debug(InfoLevel, format, a...)
}

// DebugAPI prints an 'error' message
func (f *Function) DebugAPI(format string, a ...interface{}) {
	f.Debug(APILevel, format, a...)
}

// DebugVerbose prints an 'error' message
func (f *Function) DebugVerbose(format string, a ...interface{}) {
	f.Debug(VerboseLevel, format, a...)
}

// --------------------------------------------------------

// Errorf prints an 'error' message
func (f *Function) Errorf(format string, a ...interface{}) {
	f.Println(ErrorLevel, format, a...)
}

// Warnf prints an 'warning' message
func (f *Function) Warnf(format string, a ...interface{}) {
	f.Println(WarningLevel, format, a...)
}

// Infof prints an 'info' message
func (f *Function) Infof(format string, a ...interface{}) {
	f.Println(InfoLevel, format, a...)
}

// APIf prints an 'error' message
func (f *Function) APIf(format string, a ...interface{}) {
	f.Println(APILevel, format, a...)
}

// Verbosef prints an 'error' message
func (f *Function) Verbosef(format string, a ...interface{}) {
	f.Println(VerboseLevel, format, a...)
}

// --------------------------------------------------------

// Fatalf prints a 'fatal' message
func (f *Function) Fatalf(format string, a ...interface{}) {
	f.Debug(ErrorLevel, format, a...)
	os.Exit(1)
}

// Debug prints the function name
func (f *Function) Debug(l int, format string, a ...interface{}) {
	if l <= level {
		if l <= f.pkg.level {
			if l <= f.level {
				line1 := fmt.Sprintf(format, a...)
				line2 := fmt.Sprintf("%s.%s %s", f.pkg.name, f.name, line1)
				fmt.Fprintln(os.Stderr, line2)
			}
		}
	}
}

// Printf prints a debug message
func (f *Function) Printf(l int, format string, a ...interface{}) {
	if l <= level {
		if l <= f.pkg.level {
			if l <= f.level {
				fmt.Printf(format, a...)
			}
		}
	}
}

// Println prints a debug message
func (f *Function) Println(l int, format string, a ...interface{}) {
	if l <= level {
		if l <= f.pkg.level {
			if l <= f.level {
				fmt.Println(fmt.Sprintf(format, a...))
			}
		}
	}
}

// Level returns the effective trace level
func (f *Function) Level() int {

	effectiveLevel := maxInt

	if level < effectiveLevel {
		effectiveLevel = level
	}

	if f.pkg.level < effectiveLevel {
		effectiveLevel = f.pkg.level
	}

	if f.level < effectiveLevel {
		effectiveLevel = f.level
	}

	return effectiveLevel
}

// DebugRequest traces the http request
func (f *Function) DebugRequest(req *http.Request) {

	if f.Level() >= APILevel {
		f.DebugAPI("%s %s %s %s", req.Method, req.Proto, req.Host, req.URL)

		for name, headers := range req.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				f.DebugAPI("%v: %v", name, h)
			}
		}
	}
}

// DebugRequestBody traces the http request body
func (f *Function) DebugRequestBody(data []byte) {

	if f.Level() >= APILevel {
		text1 := string(data) // multi-line json

		space := regexp.MustCompile(`\s+`)
		text2 := space.ReplaceAllString(text1, " ") // may contain a 'password' field

		text3 := text2
		var m map[string]interface{}
		err := json.Unmarshal([]byte(text2), &m)
		if err == nil {
			text3 = "{ "
			sep := ""
			for k, v := range m {
				v2 := v
				if strings.ToLower(k) == "password" {
					v2 = interface{}("********")
				}
				text3 = fmt.Sprintf("%s%s\"%s\": \"%s\"", text3, sep, k, v2)
				sep = ", "
			}
			text3 = text3 + " }"
		}
		f.DebugAPI("request: %s", text3) // sanitised!
	}
}

// Dump type
type Dump struct {
	GroupID       string `json:"groupidid"`
	Artifact      string `json:"artifact"`
	Classifier    string `json:"classifier"`
	RepositoryURL string `json:"repositoryurl"`
	Timestamp     string `json:"timestamp"`
	TimeUnix      int64  `json:"timeunix"`
	TimeUnixNano  int64  `json:"timeunixnano"`
	Package       string `json:"package"`
	Function      string `json:"function"`
	FuncForPC     string `json:"funcforpc"`
	Filename      string `json:"filename"`
	Line          int    `json:"line"`
	Version       string `json:"version"`
	BuildDate     string `json:"builddate"`
	GitCommit     string `json:"gitcommit"`
	GitBranch     string `json:"gitbranch"`
	GitURL        string `json:"giturl"`
	Message       string `json:"message"`
}

// Dump function
func (f *Function) Dump(format string, a ...interface{}) (string, error) {

	t := time.Now()
	now := fmt.Sprintf(t.Format("20060102-150405"))
	dumpDir := dumpRoot + "/" + now

	f.DebugError("DUMP: writing dump:[%s]", dumpDir)
	err := os.MkdirAll(dumpDir, 0755)
	if err != nil {
		f.DebugError("DUMP: %v", err)
		return "", err
	}

	pc, fn, line, ok := runtime.Caller(1)
	if ok {
		fmt.Println(fmt.Sprintf("package.function: %s.%s", f.pkg.name, f.name))
		fmt.Println(fmt.Sprintf("package.function: %s", runtime.FuncForPC(pc).Name()))
		fmt.Println(fmt.Sprintf("filename: %s[%d]", fn, line))
	}

	// *****************************************************************
	// * Main dump info
	// *****************************************************************
	dump := new(Dump)
	dump.GroupID = "com.rsmaxwell.players"
	dump.Artifact = "players-api"
	dump.Classifier = "test"
	dump.RepositoryURL = "https://server.rsmaxwell.co.uk/archiva"
	dump.Timestamp = now
	dump.TimeUnix = t.Unix()
	dump.TimeUnixNano = t.UnixNano()
	dump.Package = f.pkg.name
	dump.Function = f.name
	dump.FuncForPC = runtime.FuncForPC(pc).Name()
	dump.Filename = fn
	dump.Line = line
	dump.Version = version.Version()
	dump.BuildDate = version.BuildDate()
	dump.GitCommit = version.GitCommit()
	dump.GitBranch = version.GitBranch()
	dump.GitURL = version.GitURL()
	dump.Message = fmt.Sprintf(format, a...)

	json, err := json.Marshal(dump)
	if err != nil {
		f.DebugError("DUMP: %v", err)
		return "", err
	}

	filename := dumpDir + "/dump.json"

	err = ioutil.WriteFile(filename, json, 0644)
	if err != nil {
		f.DebugError("DUMP: %v", err)
		return "", err
	}

	// *****************************************************************
	// * Call stack
	// *****************************************************************
	stacktrace := debug.Stack()
	filename = dumpDir + "/callstack.txt"

	err = ioutil.WriteFile(filename, stacktrace, 0644)
	if err != nil {
		f.DebugError("DUMP: %v", err)
		return "", err
	}

	return dumpDir, nil
}
