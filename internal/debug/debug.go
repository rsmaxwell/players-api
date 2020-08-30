package debug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/rsmaxwell/players-api/internal/basic"
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
	rootDir              string
	rootDumpDir          string
	dumpFields           map[string]string
	dumpGroupID          string
	dumpArtifact         string
	dumpRepositoryURL    string
)

func init() {
	level, _ = getEnvInteger("DEBUG_LEVEL", InfoLevel)
	defaultPackageLevel, _ = getEnvInteger("DEBUG_DEFAULT_PACKAGE_LEVEL", InfoLevel)
	defaultFunctionLevel, _ = getEnvInteger("DEBUG_DEFAULT_FUNCTION_LEVEL", InfoLevel)

	path, ok := os.LookupEnv("DEBUG_DUMP_DIR")
	if !ok {
		callinfo, ok := basic.GetCallInfo(1)
		if !ok {
			panic("common.GetCallInfo failed")
		}
		path = filepath.Join(basic.HomeDir(), callinfo.ProjectName)
	}

	rootDir = path
	rootDumpDir = filepath.Join(rootDir, "dump")
	os.MkdirAll(rootDumpDir, 0755)

	dumpFields = make(map[string]string)
}

// RootDir returns the application root dir
func RootDir() string {
	return rootDir
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

// InitDump initialise the static dump fields
func InitDump(groupID, artifact, repositoryURL string) {
	dumpGroupID = groupID
	dumpArtifact = artifact
	dumpRepositoryURL = repositoryURL
}

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
		f.DebugAPI("request body: %s", text3) // sanitised!
	}
}

// Dump type
type Dump struct {
	Directory string
	Err       error
}

// DumpInfo type
type DumpInfo struct {
	GroupID       string `json:"groupidid"`
	Artifact      string `json:"artifact"`
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
func (f *Function) Dump(format string, a ...interface{}) *Dump {

	dump := new(Dump)

	t := time.Now()
	now := fmt.Sprintf(t.Format("2006-01-02_15-04-05.999999999"))
	dump.Directory = filepath.Join(rootDumpDir, now)

	f.DebugError("DUMP: writing dump:[%s]", dump.Directory)
	err := os.MkdirAll(dump.Directory, 0755)
	if err != nil {
		dump.Err = err
		return dump
	}

	// *****************************************************************
	// * Main dump info
	// *****************************************************************
	info := new(DumpInfo)
	info.GroupID = dumpGroupID
	info.Artifact = dumpArtifact
	info.RepositoryURL = dumpRepositoryURL
	info.Timestamp = now
	info.TimeUnix = t.Unix()
	info.TimeUnixNano = t.UnixNano()
	info.Version = basic.Version()
	info.BuildDate = basic.BuildDate()
	info.GitCommit = basic.GitCommit()
	info.GitBranch = basic.GitBranch()
	info.GitURL = basic.GitURL()
	info.Message = fmt.Sprintf(format, a...)
	info.Package = f.pkg.name
	info.Function = f.name

	pc, fn, line, ok := runtime.Caller(1)
	if ok {
		info.FuncForPC = runtime.FuncForPC(pc).Name()
		info.Filename = fn
		info.Line = line
	}

	json, err := json.Marshal(info)
	if err != nil {
		dump.Err = err
		return dump
	}

	filename := dump.Directory + "/dump.json"

	err = ioutil.WriteFile(filename, json, 0644)
	if err != nil {
		dump.Err = err
		return dump
	}

	// *****************************************************************
	// * Call stack
	// *****************************************************************
	stacktrace := debug.Stack()
	filename = dump.Directory + "/callstack.txt"

	err = ioutil.WriteFile(filename, stacktrace, 0644)
	if err != nil {
		dump.Err = err
		return dump
	}

	return dump
}

// AddString method
func (d *Dump) AddString(title string, data string) {
	d.AddByteArray(title, []byte(data))
}

// AddByteArray method
func (d *Dump) AddByteArray(title string, data []byte) {

	if d.Err != nil {
		return
	}

	filename := d.Directory + "/" + title

	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return
	}
}

// MarkDumps type
type MarkDumps struct {
	dumps map[string]bool
	err   error
}

// Mark method
func Mark() *MarkDumps {

	mark := new(MarkDumps)

	files, err := ioutil.ReadDir(rootDumpDir)
	if err != nil {
		mark.err = err
		return mark
	}

	mark.dumps = map[string]bool{}

	for _, file := range files {
		if file.IsDir() {
			mark.dumps[file.Name()] = true
		}
	}

	return mark
}

// ListNewDumps method
func (mark *MarkDumps) ListNewDumps() ([]*Dump, error) {

	if mark.err != nil {
		return nil, mark.err
	}

	files, err := ioutil.ReadDir(rootDumpDir)
	if err != nil {
		mark.err = err
		return nil, err
	}

	newdumps := []*Dump{}

	for _, file := range files {
		if file.IsDir() {
			if !mark.dumps[file.Name()] {

				dump := new(Dump)
				dump.Directory = rootDumpDir + "/" + file.Name()

				newdumps = append(newdumps, dump)
			}
		}
	}

	return newdumps, nil
}

// ListDumps method
func ListDumps() ([]*Dump, error) {

	files, err := ioutil.ReadDir(rootDumpDir)
	if err != nil {
		return nil, err
	}

	newdumps := []*Dump{}

	for _, file := range files {
		if file.IsDir() {
			dump := new(Dump)
			dump.Directory = rootDumpDir + "/" + file.Name()

			newdumps = append(newdumps, dump)
		}
	}

	return newdumps, nil
}

// Remove function
func (d *Dump) Remove() error {

	err := os.RemoveAll(d.Directory)
	if err != nil {
		return err
	}

	return nil
}

// GetInfo function
func (d *Dump) GetInfo() (*DumpInfo, error) {

	infofile := d.Directory + "/dump.json"

	data, err := ioutil.ReadFile(infofile)
	if err != nil {
		return nil, err
	}

	var info DumpInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// ClearDumps function
func ClearDumps() error {

	dumps, err := ListDumps()
	if err != nil {
		return err
	}

	for _, dump := range dumps {
		err = dump.Remove()
		if err != nil {
			return err
		}
	}

	return nil
}

// DumpSQLError dumps a database error
func (f *Function) DumpSQLError(err error, message string, sql string) *Dump {
	d := f.DumpError(err, message)
	d.AddString("sql.txt", sql)
	return d
}

// DumpError dumps a database error
func (f *Function) DumpError(err error, message string) *Dump {

	d := f.Dump(message)

	d.AddString("error.txt", fmt.Sprintf("%T\n\n%s", err, err.Error()))

	var title string
	var filename string
	var data []byte
	var err3 error

	if err2, ok := err.(*pq.Error); ok {
		title = "pg.error"
		data, err3 = json.Marshal(err2)
	} else if err2, ok := err.(*pgconn.PgError); ok {
		title = "pgconn.error"
		data, err3 = json.Marshal(err2)
	} else {
		title = "error"
		data, err3 = json.Marshal(err)
	}

	if err3 != nil {
		fmt.Println("could not marshal error: " + err3.Error())
	} else {
		filename = filepath.Join(d.Directory, title+".json")
		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			f.Errorf("could not write error to dump\n")
		}
	}

	return d
}
