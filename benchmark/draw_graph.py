import subprocess
import matplotlib.pyplot as plt

PARALLEL_THREADS_NUM = [2, 4, 6, 8, 12]
LINES = ["small", "mixture", "big"]
PATTERNS = ["balance", "steal", ""]
RUN_TIMES = 5

if __name__ == '__main__':
    seq_avg_times = {line: {} for line in LINES}
    balance_avg_times = {line: {} for line in LINES}
    steal_avg_times = {line: {} for line in LINES}

    # traverse to get benchmark result
    for pattern in PATTERNS:
        print("Running: " + pattern)
        for line in LINES:
            if pattern != "":
                # parallel pattern
                for thread_count in PARALLEL_THREADS_NUM:
                    times_p = 0.0
                    cmd_p = "go run ../editor/editor.go {} {} {}".format(line, pattern, str(thread_count))
                    for i in range(RUN_TIMES):
                        output = float(subprocess.check_output(cmd_p, shell=True).strip())
                        times_p += output
                        print("current time is: " + str(output))

                    avg_time_p = times_p/RUN_TIMES
                    print("FINISH -- " + cmd_p + "; AVG TIME: " + str(avg_time_p))
                    if pattern == "balance":
                        balance_avg_times[line][thread_count] = avg_time_p
                    elif pattern == "steal":
                        steal_avg_times[line][thread_count] = avg_time_p
            else: #sequential
                times_s = 0.0
                cmd_s = "go run ../editor/editor.go {}".format(line)
                for i in range(RUN_TIMES):
                    output = subprocess.check_output(cmd_s, shell=True)
                    times_s += float(output.strip())
                avg_time_s = times_s/RUN_TIMES
                print(cmd_s + "; Avg time is: " + str(avg_time_s))
                seq_avg_times[line] = avg_time_s

    # calculating each speedup and plot
    for pattern in ["balance", "steal"]:
        if pattern == "balance":
            avg_times = balance_avg_times
        elif pattern == "steal":
            avg_times = steal_avg_times
        for line in LINES:
            seq_avg_time = seq_avg_times[line]
            speedups = []
            for thread_count in PARALLEL_THREADS_NUM:
                speedup = seq_avg_time / avg_times[line][thread_count]
                speedups.append(speedup)
            plt.plot(PARALLEL_THREADS_NUM, speedups, label = line)
        plt.xlabel("Number of Threads")
        plt.ylabel("Speedup")
        plt.title(pattern + "-speedup Graph")
        plt.legend(loc='best')
        plt.grid()
        plt.savefig('speedup-' + pattern + '.png')
        plt.show()
