import argparse
import imagehash
import pandas as pd
from PIL import Image


def calc_phash(df):
    for row in df[df["phash"] == ""].itertuples():
        image = Image.open(row.Index)
        phash = imagehash.phash(image)
        print(row.Index, phash)
        df.loc[row.Index, "phash"] = str(phash)

    return df


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--data_file", default="data.h5", help="Name of data file")
    args = parser.parse_args()

    df = pd.read_hdf(args.data_file)
    df = calc_phash(df)
    df.to_hdf(args.data_file, key="images")
