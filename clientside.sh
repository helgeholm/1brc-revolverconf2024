set -e
echo -ne "⚙️  syncing\r"
rsync --exclude '*.txt' -a . root@37.27.213.152:1brc-prez
echo -ne "✅ ready  \r"
while true; do
    inotifywait -qme modify --format "👉️" cmd/ | while read x; do
        echo -ne "⚙️  compiling        \r"
        if (go build cmd/program.go >/dev/null 2>&1); then
            echo -ne "⚙️  testing       \r"
	    if ./program measurements_1k.txt 2>/dev/null | sort | diff - correct_enough_for_floats_1k.txt >/dev/null \
	    || ./program measurements_1k.txt 2>/dev/null | sort | diff - correct_1k.txt >/dev/null ; then
        	echo -ne "⚙️  syncing       \r"
        	rsync --exclude '*.txt' -a . root@37.27.213.152:1brc-prez
        	echo -ne "✅ synced \r"
	    else
        	echo -ne "❌ INCORRECT  \r"
	    fi
        else
            echo -ne "❌ COMPILE ERROR  \r"
            sleep 1
        fi
    done
done
