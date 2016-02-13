#!/usr/bin/expect
spawn ssh root@ravenor ./ravenor
expect "password:"
send "root\r"
expect "$"
exit 0
