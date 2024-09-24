#! /usr/bin/bash

echo "Launching Write to TUN experiment..."

N_ROUTINES="1 2 4 8 16 32"
PAYLOAD_SIZES="1500 15000 30000 60000"
DURATION=5
nruns=3

res_path=results/write

sudo ip l del dev tun0

mkdir -p $res_path
make

for run in `seq 1 $nruns`; do
    for n_routine in $N_ROUTINES; do
        for payload_size in $PAYLOAD_SIZES; do
            # exp_code=${payload_size}_${n_routine}_${run}
            exp_code=${payload_size}_${n_routine}_${run}
            echo "Eval: type=read payload=$payload_size n_threads=$n_routine run=$run"

            sudo ./proxy-tun write $payload_size $n_routine > $res_path/write_$exp_code.csv &
            sudo ./setup_tun.sh
            sleep 6
            # nping --udp 192.168.0.1 -c 0 --data-length $payload_size -e tun0 --rate 10000000 -N 2> /dev/null 1> /dev/null &
            sar -A -o $res_path/write_sar_$exp_code.data 1 $DURATION > /dev/null &

            sleep $DURATION

            sudo ip l del dev tun0
            sudo pkill nping
        done
    done
done

