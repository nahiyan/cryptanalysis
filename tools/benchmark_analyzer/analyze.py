import os

benchmarks_path = "../../benchmarks/saeed/"


def process_session(path):
    # Dot-matrix
    dm_count = 0
    dm_total_time = 0
    # Counter-chain
    cc_count = 0
    cc_total_time = 0

    benchmark_filepath = os.path.join(path, "benchmark.log")
    f = open(benchmark_filepath, "r")

    lines = f.readlines()
    for line in lines:
        segments = line.split()
        inst_name = ""
        exit_code = 0
        time = 0
        i = 0

        for segment in segments:
            if segment == "instance":
                inst_name = segments[i + 2][:-1]

            if segment == "exit":
                exit_code = int(segments[i + 2])

            if segment == "Time:":
                time = float(segments[i + 1][:-2])

            i += 1

        if exit_code == 10:
            inst_segments = inst_name.split("_")
            if inst_segments[2] == "dot":
                dm_count += 1
                dm_total_time += time
            if inst_segments[2] == "counter":
                cc_count += 1
                cc_total_time += time

    f.close()

    return (cc_count, cc_total_time / cc_count, dm_count, dm_total_time / dm_count)


def main():
    for dir in os.walk(benchmarks_path):
        if dir[0].find("session") != -1:
            cc_count, cc_avg_time, dm_count, dm_avg_time = process_session(
                dir[0])

            print("{}, {}, {:.0f}s, {}, {:.0f}s".format(os.path.basename(
                os.path.normpath(dir[0])), cc_count, cc_avg_time, dm_count, dm_avg_time))


main()
