#!/usr/bin/expect
spawn scp ./ravenor root@ravenor:.
expect "password:"
send "root\r"
# expect "ravenor"
expect eof
exit 0
