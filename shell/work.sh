LOG=/PPPoker/log/punch_kiwi.log
NOW=$(date +"%Y%m%d%H%M%S")
DIR=$(pwd)
echo "start=$NOW, dir=$DIR"  >>  $LOG
killall punch_kiwi
cd /Users/frankieli/goProjects/src/punch_kiwi
./punch_kiwi
echo `end=$(date +"%Y%m%d%I%M%S")` >> $LOG