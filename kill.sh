#!/usr/bin/expect
spawn ssh root@ravenor killall ravenor
expect "password:"
send "root\r"
expect eof
exit 0
