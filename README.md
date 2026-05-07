# go-logtail

A lightweight terminal utility for monitoring and tailing log files.

## Usage

```bash
go-logtail [OPTIONS] FILE...

OPTIONS
    --debug    enable debug logging
```



## Purpose

I run several long-lived applications that generate large volumes of logs. I wanted to create a lightweight and extensible utility for monitoring those applications in dedicated terminals and surfacing formatted, customizable alerts when abnormal behavior occurs.

This project also serves as an opportunity to improve my Go skills while learning to write more idiomatic Go.

## Concepts Practiced

This project explores core Go concepts, including:

- concurrency
- channels
- goroutines
- Go standard libraries including `bufio`, `flag`, `fmt`, `io`, `os`, and `strings`
- terminal formatting and ANSI escape sequences