package main

import (
	"flag"
	"fmt"
	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
	"github.com/winxxp/hlog"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

var (
	// Version 版本号，编译时自动指定
	// go build -ldflags "-X main.Version=1.1.1 -s -w"
	version     = "1.1.1"
	commit      = ""
	date        = ""
	root        = flag.String("root", ".", "Root directory to serve, if not exist will create")
	ftpAddress  = flag.String("ftp-addr", ":2121", "ftp bind address")
	user        = flag.String("user", "coinv", "Username for ftp server login")
	pass        = flag.String("pass", "coinv1985", "Password for ftp server login")
	showVersion = flag.Bool("version", false, "print application version")
	verbose     = flag.Bool("verbose", false, "show ftp log")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("simpleftp %s build:%s id:%s\n", version, date, commit)
		os.Exit(0)
	}

	if *root == "" {
		*root, _ = os.Getwd()
	}
	RootDirectory, _ := filepath.Abs(*root)
	if _, err := os.Stat(RootDirectory); err != err {
		if os.IsNotExist(err) {
			err = os.MkdirAll(RootDirectory, os.ModePerm)
		}
		if err != nil {
			hlog.WithError(err).Error("init root directory")
		}
	}

	NewFTPServer(RootDirectory, *ftpAddress).Start()
}

type FTPServer struct {
	*server.Server
	opts *server.ServerOpts

	// 文件初始路径
	RootDirectory string
}

func NewFTPServer(root string, addr string) *FTPServer {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		hlog.Fatal(err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		hlog.Fatal(err)
	}

	ftpSrv := &FTPServer{
		opts: &server.ServerOpts{
			Factory: &filedriver.FileDriverFactory{
				RootPath: root,
				Perm:     server.NewSimplePerm("user", "group"),
			},
			Port:     port,
			Hostname: host,
			Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
			Logger: func() server.Logger {
				if *verbose {
					return &FTPLog{}
				}
				return &server.DiscardLogger{}
			}(),
		},
	}
	ftpSrv.Server = server.NewServer(ftpSrv.opts)

	return ftpSrv
}

func (s *FTPServer) Start() {
	hlog.WithFields(hlog.Fields{
		"addr":     net.JoinHostPort(s.Hostname, strconv.Itoa(s.Port)),
		"username": *user,
		"pass":     *pass,
	}).Info("ftp server start")

	err := s.ListenAndServe()
	hlog.WithResult(err).Log("ftp server quit")
}

type FTPLog struct {
}

func (*FTPLog) Print(sessionId string, message interface{}) {
	hlog.WithIDString(sessionId).Info(message)
}

func (*FTPLog) Printf(sessionId string, format string, v ...interface{}) {
	hlog.WithIDString(sessionId).Infof(format, v...)
}

func (*FTPLog) PrintCommand(sessionId string, command string, params string) {
	hlog.WithIDString(sessionId).Infof("> %s %s", command, params)
}

func (*FTPLog) PrintResponse(sessionId string, code int, message string) {
	hlog.WithIDString(sessionId).Infof("< %d %s", code, message)
}
