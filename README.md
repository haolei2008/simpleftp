# 简单FTP

静态数据指定路径FTP服务器

## 典型用法

默认端口:`2121`，默认用户:`coinv`，默认密码:`coinv1985`
```
simpleftp -root=/mnt/udisk/  -log.dir=/mnt/udisk/log/
```

## 帮助
```bash
simpleftp -h

Usage of simpleftp.exe:
  -ftp-addr string
        ftp bind address (default ":2121")
  -log.alsotostderr
        log to standard error as well as files
  -log.backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log.dir string
        If non-empty, write log files to os temp directory
  -log.flushInterval value
        interval of flush buffer to log file (default 5s)
  -log.stderrthreshold value
        logs at or above this threshold go to stderr
  -log.tostderr
        log to standard error instead of files
  -log.v value
        log level for V logs
  -log.vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
  -pass string
        Password for ftp server login (default "coinv1985")
  -root string
        Root directory to serve, if not exist will create (default ".")
  -user string
        Username for ftp server login (default "coinv")
  -verbose
        show ftp log
  -version
        print application version

```