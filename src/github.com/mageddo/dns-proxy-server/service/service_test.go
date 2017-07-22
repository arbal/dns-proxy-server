package service

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/mageddo/log"
)

func TestSetupFor_NormalModeSuccess(t *testing.T) {
	ctx := log.GetContext()

	const SERVICE_FILE = "/tmp/serviceFile"
	sc := NewService(ctx)
	err := sc.SetupFor(SERVICE_FILE, &Script{"ls"})
	if err != nil {
		t.Error(err)
	}

	bytes, err := ioutil.ReadFile(SERVICE_FILE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	assert.Equal(t, EXPECTED_SERVICE_FILE, string(bytes))

}

const EXPECTED_SERVICE_FILE = `#!/bin/sh
### BEGIN INIT INFO
# Provides:          dns-proxy-server
# Required-Start:    $local_fs $network $named $time $syslog
# Required-Stop:     $local_fs $network $named $time $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Description:       DNS PROXY SERVER
### END INIT INFO

SCRIPT=ls
RUNAS=root

PIDFILE=/var/run/dns-proxy-server.pid
LOGFILE=/var/log/dns-proxy-server.log

start() {
  if [ -f /var/run/$PIDNAME ] && kill -0 $(cat /var/run/$PIDNAME); then
    echo 'Service already running' >&2
    return 1
  fi
  echo 'Starting service…' >&2
  local CMD="$SCRIPT &> \"$LOGFILE\" & echo \$!"
  su -c "$CMD" $RUNAS > "$PIDFILE"
  echo 'Service started' >&2
}

stop() {
  if [ ! -f "$PIDFILE" ] || ! kill -0 $(cat "$PIDFILE"); then
    echo 'Service not running' >&2
    return 1
  fi
  echo 'Stopping service…' >&2
  kill -15 $(cat "$PIDFILE") && rm -f "$PIDFILE"
  echo 'Service stopped' >&2
}

uninstall() {
  echo -n "Are you really sure you want to uninstall this service? That cannot be undone. [yes|No] "
  local SURE
  read SURE
  if [ "$SURE" = "yes" ]; then
    stop
    rm -f "$PIDFILE"
    echo "Notice: log file is not be removed: '$LOGFILE'" >&2
    update-rc.d -f dns-proxy-server remove
    rm -fv "$0"
  fi
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  uninstall)
    uninstall
    ;;
  retart)
    stop
    start
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|uninstall}"
esac
`
