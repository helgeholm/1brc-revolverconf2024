set -e
wat=$(cat cmd/wat)
echo "ready"
while true; do
    inotifywait -qq -e modify cmd/
    wat2=$(cat cmd/wat)
    if [[ "$wat" != "$wat2" ]]; then
        wat=$wat2
        echo "RUN: $wat"
        echo -ne "⚙️  compiling\r"
        if (go build cmd/program.go >/dev/null 2>&1); then
            echo -ne "⏱️  ...        \r"
            run1=$(/bin/time --format="%es" ./program 2>&1 >/dev/null)
            echo -ne "⏱️  ${run1} ...\r"
            run2=$(/bin/time --format="%es" ./program 2>&1 >/dev/null)
            echo -ne "⏱️  ${run1} ${run2} ...\r"
            run3=$(/bin/time --format="%es" ./program 2>&1 >/dev/null)
            echo -ne "✅ ${run1} ${run2} ${run3}\r"
        else
            echo -ne "❌ ERROR       \r"
            sleep 1
        fi
        echo
    fi
done
