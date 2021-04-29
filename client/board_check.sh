#!/bin/bash

# disk
## ram
DESCRIPTION_RAM="ddr ram"

## emmc
DESCRIPTION_EMMC="emmc"

# rtc
CMD_RTC_HWCLOCK="hwclock"
CMD_RTC_SYSCLOCK="date"
DESCRIPTION_RTC="rtc"

# variable
return_value=0
local_network_status=0
abnormal_items="abnormal"

function printmsg() {
    if [ $1 == "error" ]; then
        echo -e "\033[31m [ Error ] $2 $3 \033[0m"
    elif [ $1 == "warn" ]; then
        echo -e "\033[33m [ Warn ] $2 $3 \033[0m"
    elif [ $1 == "title" ]; then
        echo -e "\033[36m >>>>>>>>>> $2  <<<<<<<<<<\033[0m"
    else
        echo -e "\033[32m [ Log ] $2 $3 \033[0m"
    fi
}

function check_ip_ping() {
    if [ $2 != "0.0.0.0" ]; then
        ping -I eth0 -c 2 $2 > /dev/null
        if [ $? == 0 ]; then
            printmsg "log" "ip-ping core-$1: okay"
        else
            printmsg "error" "ip-ping core-$1: failed"
            return_value=1
            local_network_status=1
            return
        fi
    else
        printmsg "error" "couldn't find ip of core-$1"
        return_value=1
        local_network_status=1
    fi
}

function check_local_network() {
    printmsg "title" "Check network of sub cores"
    printmsg "log" "description: check local network"
    printmsg "log" "area: local"

    i2c-cmd -r -m 0 0 > /tmp/local_network.log
    if [ $? == 0 ]; then
        printmsg "log" "i2c-cmd: get ip address of sub core okay"
        cat /tmp/local_network.log
    else
        printmsg "error" "i2c-cmd: get ip address of sub core failed"
        return_value=1
        return
    fi

    core_1_ip=`cat /tmp/local_network.log | grep core1 | awk '{print $2}' | sed 's/.$//'`
    core_2_ip=`cat /tmp/local_network.log | grep core2 | awk '{print $2}' | sed 's/.$//'`
    core_3_ip=`cat /tmp/local_network.log | grep core3 | awk '{print $2}' | sed 's/.$//'`
    core_4_ip=`cat /tmp/local_network.log | grep core4 | awk '{print $2}' | sed 's/.$//'`
    core_5_ip=`cat /tmp/local_network.log | grep core5 | awk '{print $2}' | sed 's/.$//'`
    core_6_ip=`cat /tmp/local_network.log | grep core6 | awk '{print $2}' | sed 's/.$//'`
    core_7_ip=`cat /tmp/local_network.log | grep core7 | awk '{print $2}' | sed 's/.$//'`

    # core-1
    check_ip_ping 1 $core_1_ip

    # core-2
    check_ip_ping 2 $core_2_ip

    # core-3
    check_ip_ping 3 $core_3_ip

    # core-4
    check_ip_ping 4 $core_4_ip

    # core-5
    check_ip_ping 5 $core_5_ip

    # core-6
    check_ip_ping 6 $core_6_ip

    # core-7
    check_ip_ping 7 $core_7_ip

    if [ $local_network_status == 0 ]; then
        printmsg "log" "OK"
    fi
}

function check_ram() {
    printmsg "title" "Check RAM"
    printmsg "log" "description: $DESCRIPTION_RAM"

    mem_free=`cat /proc/meminfo | grep MemFree | awk '{print $2}'`
    mem_total=`cat /proc/meminfo | grep MemTotal | awk '{print $2}'`
    mem_use=`expr $mem_total - $mem_free`
    mem_usage=`expr $mem_use \* 100 / $mem_total`

    printmsg "log" "Mem capacity:" "${mem_total}KB"
    printmsg "log" "Mem usage:" "${mem_usage}%"

    if [[ $mem_total == NULL ]]; then
        abnormal_items="$abnormal_items RAM"
        printmsg "error" "system ram abnormal"
        return_value=1
    else
        printmsg "log" "OK"
    fi
}

function check_emmc() {
    printmsg "title" "Check EMMC"
    printmsg "log" "description: $DESCRIPTION_EMMC"

    #disk_total=`df -h | grep /dev/root | awk '{print $2}'`
    #disk_usage=`df -h | grep /dev/root | awk '{print $5}'`
    disk_total=`cat /proc/partitions | grep mmcblk2p8 | awk '{print $3}'`

    printmsg "log" "Disk capacity:" "${disk_total}KB"
    #printmsg "log" "Disk usage:" "$disk_usage"

    if [[ $disk_total  == NULL ]]; then
        abnormal_items="$abnormal_items EMMC"
        printmsg "error" "system emmc abnormal"
        return_value=1
    else
        printmsg "log" "OK"
    fi
}

function check_rtc() {
    printmsg "title" "Check rtc"
    printmsg "log" "description: $DESCRIPTION_RTC"

    hwclock=`$CMD_RTC_HWCLOCK` > /dev/null
    if [ $? == 0 ]; then
        printmsg "log" "hwclock: $hwclock"
    else
        abnormal_items="$abnormal_items RTC"
        printmsg "error" "no such device"
        return_value=1
        return
    fi

    sysclock=`$CMD_RTC_SYSCLOCK` > /dev/null
    if [ $? == 0 ]; then
        printmsg "log" "sysclock: $sysclock"
        printmsg "log" "OK"
    fi
}

check_ram
check_emmc
check_local_network
check_rtc

# show abnormal items on oled
#echo $abnormal_items
#/usr/bin/oled_cmd $abnormal_items

#echo "voltage_now" $voltage_now >> /home/admin/.board_status

exit $return_value

