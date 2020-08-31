import os
import argparse

import pandas as pd


def run(images_dir, data_file):
    s_old = set()
    s_new = set()
    # retrieve or create dataframe
    df = pd.DataFrame([])
    if os.path.exists(data_file):
        df = pd.read_hdf(data_file, "images")
        s_old.update(df.index)
    # search image files
    for filename in os.listdir(images_dir):
        if not filename.endswith(".jpg"):
            continue
        filepath = os.path.abspath(os.path.join(images_dir, filename))
        s_new.add(filepath)

    df = df.drop(s_old - s_new)
    df = df.append(pd.DataFrame(index=pd.Index(s_new - s_old)))
    print(df)
    df.to_hdf(data_file, key="images")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("images_dir", help="Path to images directory")
    parser.add_argument("--data_file", default="data.h5", help="Name of data file")
    args = parser.parse_args()

    run(args.images_dir, args.data_file)
