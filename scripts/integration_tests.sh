#!/bin/bash

# Color schema
RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[34m"
RESET="\033[0m"

# Run build.sh
echo -e "${YELLOW}-- > Executing build script... < --${RESET}"
./build.sh &
BUILD_PID=$!
wait $BUILD_PID

# Check OS
OS=$(uname -s)
case "$OS" in
    Linux*)
        OS_TYPE="" 
        ;;
    Darwin*)
        OS_TYPE=""
        ;;
    CYGWIN*|MINGW32*|MSYS*|MINGW*)
        OS_TYPE=".exe"
        ;;
    *)
        echo "Unsupported OS"
        exit 1
        ;;
esac

# Run server command
PORT=$((RANDOM % 1001 + 3000))
echo ""
echo -e "${YELLOW}-- > Running server... < --${RESET}"
../build/sumologic_server${OS_TYPE} server -p ${PORT} -a "0.0.0.0" -m 3 > /dev/null 2>&1 &
SERVER_PID=$!
echo -e "${GREEN}-- > [x] The server is up and listening on port ${PORT} < --${RESET}"

sleep 3

# Run client command
sendServerSuccessfullyRequest() {
    local time=$1
    local thread_num=$2
    OUTPUT=$(
        ../build/sumologic_server${OS_TYPE} client -p ${PORT} -a 0.0.0.0 \
        --script "../build/sumologic_server${OS_TYPE}" --script "await" --script "-t" --script "${time}" -t 6000
    )
    
    if echo "$OUTPUT" | grep -q "Await finished"; then
        echo -e "${GREEN}-- > [x] The server handled the request successfully for thread: ${thread_num} with ${time} ms${RESET}"
    else
        echo -e "${RED}-- > [ ] The server could not handle the request for thread: ${thread_num}${RESET}"
    fi
}

echo ""
echo -e "${YELLOW}-- > Running client... < --${RESET}"
echo -e "${BLUE}* Running Parallel Requests - less than the limit of connection to avoid the semaphore... *${RESET}"
PIDS=()
times=(500 100 1000)
completed_threads=()

for i in "${!times[@]}"; do
    sendServerSuccessfullyRequest "${times[i]}" "$((i + 1))" &
    PIDS+=($!)
done


for pid in "${PIDS[@]}"; do
    wait "${PIDS[i]}";
    completed_threads+=("$((i + 1))")
done

echo ""
echo -e "${BLUE}* Running Parallel Requests - more than the limit of connection to test semaphore... *${RESET}"
PIDS=()
times=(1800 1900 2000 200 300 400)
for i in "${!times[@]}"; do
    sendServerSuccessfullyRequest "${times[i]}" "$((i + 1))" &
    PIDS+=($!)
done

for pid in "${PIDS[@]}"; do
    wait "$pid"
done

echo ""
echo -e "${BLUE}* Running Single Requests... *${RESET}"
OUTPUT=$(../build/sumologic_server${OS_TYPE} client -p ${PORT} -a 0.0.0.0 --script "sleep" --script "1" -t 2000)
if echo "$OUTPUT" | grep -q "ExitCode:0"; then
    echo -e "${GREEN}-- > [x] The server could handle general commands such as sleep${RESET}"
else
    echo -e "${RED}-- > [ ] The server could not handle general commands such as sleep${RESET}"
fi

OUTPUT=$(../build/sumologic_server${OS_TYPE} client -p ${PORT} -a 0.0.0.0 --script "echo" --script "sl_integration_test" -t 2000)
if echo "$OUTPUT" | grep -q "Output:sl_integration_test"; then
    echo -e "${GREEN}-- > [x] The server could handle general commands such as echo and save it in output${RESET}"
else
    echo -e "${RED}-- > [ ] The server could not handle general commands such as echo${RESET}"
fi

OUTPUT=$(../build/sumologic_server${OS_TYPE} client -p ${PORT} -a 0.0.0.0 --script "../build/sumologic_server${OS_TYPE}" --script "await" --script "-t" --script "10000" -t 100)
if echo "$OUTPUT" | grep -q "timeout exceeded"; then
    echo -e "${GREEN}-- > [x] The server could handle a time out request successfully${RESET}"
else
    echo -e "${RED}-- > [ ] The server could not handle a time out request successfully${RESET}"
fi

OUTPUT=$(../build/sumologic_server${OS_TYPE} client -p ${PORT} -a 0.0.0.0 --script "../build/invalid" --script "await" --script "-t" --script "10000" -t 20000)
if echo "$OUTPUT" | grep -q "file not found"; then
    echo -e "${GREEN}-- > [x] The server could handle a not found command${RESET}"
else
    echo -e "${RED}-- > [ ] he server could not handle a not found command${RESET}"
fi

echo ""
echo -e "${YELLOW}-- > Killing Server... < --${RESET}"
kill $SERVER_PID
sleep 2
echo -e "${GREEN}-- > [x] Server killed < --${RESET}"

sleep 2
OUTPUT=$(../build/sumologic_server${OS_TYPE} client -p ${PORT} -a 0.0.0.0 --script "../build/sumologic_server${OS_TYPE}" --script "await" --script "-t" --script "10000" -t 100)
if echo "$OUTPUT" | grep -q "actively refused it"; then
    echo -e "${GREEN}-- > [x] The server is not accepting connections anymore < --${RESET}"
else
    echo -e "${RED}-- > [ ] The server is still listening on port ${PORT} < --${RESET}"
fi

echo ""


echo "Press any key to quit..."
read -n 1 