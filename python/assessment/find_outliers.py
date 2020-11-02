import argparse
import numpy as np
import pandas as pd


def run(df, threshold):
    df.loc[:, "eyel_x"] = (
        df.loc[:, [f"parts{i:02d}_x" for i in range(36, 42)]].sum(axis=1) / 6
    )
    df.loc[:, "eyel_y"] = (
        df.loc[:, [f"parts{i:02d}_y" for i in range(36, 42)]].sum(axis=1) / 6
    )
    df.loc[:, "eyer_x"] = (
        df.loc[:, [f"parts{i:02d}_x" for i in range(42, 48)]].sum(axis=1) / 6
    )
    df.loc[:, "eyer_y"] = (
        df.loc[:, [f"parts{i:02d}_y" for i in range(42, 48)]].sum(axis=1) / 6
    )
    df.loc[:, "angle"] = np.arctan2(
        df["eyer_y"] - df["eyel_y"], df["eyer_x"] - df["eyel_x"]
    )

    parts_x = [f"parts{i:02d}_x" for i in range(0, 68)]
    parts_y = [f"parts{i:02d}_y" for i in range(0, 68)]
    df.loc[:, "face_w"] = df.loc[:, "face_right"] - df.loc[:, "face_left"]
    df.loc[:, "face_h"] = df.loc[:, "face_bottom"] - df.loc[:, "face_top"]
    df.loc[:, "xmin"] = df.loc[:, parts_x].min(axis=1)
    df.loc[:, "xmax"] = df.loc[:, parts_x].max(axis=1)
    df.loc[:, "xmean"] = df.loc[:, parts_x].mean(axis=1)
    df.loc[:, "ymin"] = df.loc[:, parts_y].min(axis=1)
    df.loc[:, "ymax"] = df.loc[:, parts_y].max(axis=1)
    df.loc[:, "ymean"] = df.loc[:, parts_y].mean(axis=1)

    ids = set()
    targets = (
        "angle",
        "face_w",
        "face_h",
        "xmin",
        "xmax",
        "xmean",
        "ymin",
        "ymax",
        "ymean",
    )
    for target in targets:
        x = df.loc[:, target]
        t = (x - x.mean()) / x.std()
        print(f"{target} ({x.mean():.8f}, {x.std():.8f})")
        for idx in t[t.abs() > threshold].sort_values(key=np.abs).index:
            print(f"{idx}: {df.loc[idx, target]:.8f} ({t[idx]:.8f})")
            ids.add(idx)

    print("outliers:")
    for idx in sorted(ids):
        print(idx)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--data_file", default="data.h5", help="Name of data file")
    parser.add_argument("--threshold", default=5.0, type=float)
    args = parser.parse_args()

    df = pd.read_hdf(args.data_file)
    run(df, args.threshold)
