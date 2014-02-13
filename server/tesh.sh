COUNTER=1
while [ "$COUNTER" -lt 5000 ]; do
    ./server & 
    echo "$COUNTER connected.\n"
    COUNTER=$(($COUNTER+1))
    r=`expr $COUNTER % 100`
    if [ $r -eq 0 ];then
        sleep 1
    fi
done
