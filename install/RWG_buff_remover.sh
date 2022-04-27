#!/bin/sh

buf_f=/tmp/RWG_app_buffer/buffer.buf
log_f=/tmp/RWG_app_buffer/log.log

while :; do

  sleep 12h

  buf_sz=$(wc -c $buf_f | awk '{print $1}')
  if [[ -$buf_sz -ge 10000000 ]]; then
    rm $buf_f
  fi

  log_sz=$(wc -c $log_f | awk '{print $1}')
  if [[ -$log_sz -ge 1000000 ]]; then
    > $log_f
  fi

done
