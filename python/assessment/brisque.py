import argparse
import pandas as pd
from imquality import brisque
from PIL import Image


def calc_brisque(df):
    for row in df[df["brisque"].isnull()].itertuples():
        image = Image.open(row.Index)
        df.loc[row.Index, "brisque"] = brisque.score(image)

    return df


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--data_file", default="data.h5", help="Name of data file")
    args = parser.parse_args()

    df = pd.read_hdf(args.data_file)
    df = calc_brisque(df)
    df.to_hdf(args.data_file, key="images")
