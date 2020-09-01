import argparse
import os
import pandas as pd
from collections import defaultdict


def find_duplicates(df):
    results = defaultdict(list)
    s = df.get("phash").value_counts()
    for phash in s[s > 1].index:
        for row in df[df.phash == phash].itertuples():
            results[row.phash].append(os.path.basename(row.Index))
    return results


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--data_file", default="data.h5", help="Name of data file")
    args = parser.parse_args()

    df = pd.read_hdf(args.data_file)
    results = find_duplicates(df)
    for phash, index in results.items():
        print(f"[{phash}]")
        for idx in index:
            print(idx)
