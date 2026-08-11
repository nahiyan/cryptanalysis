file=$1
exit_status=$(tail -n 10 $file | grep "^c exit" | awk '{print $3}')
echo $exit_status
process_time=$(tail -n 50 $file | grep "c total process time since initialization" | awk '{print $7}')
echo $process_time
time_limit=$(head -n 50 $file | grep "setting time limit" | awk '{print $6}')
echo $time_limit
