import argparse
import cv2
import numpy as np
import pandas as pd


def calc_face_features(df):
    if "face_blur" not in df:
        df["face_blur"] = np.nan

    for row in df[df["face_score"].notna() & df["face_blur"].isnull()].itertuples():
        img = cv2.imread(row.Index, cv2.IMREAD_COLOR)
        face = img[
            round(row.face_top) : round(row.face_bottom),
            round(row.face_left) : round(row.face_right),
            :,
        ]
        gray = cv2.cvtColor(face, cv2.COLOR_BGR2GRAY)
        variance = cv2.Laplacian(gray, cv2.CV_64F).var()
        print(row.Index, variance)
        df.loc[row.Index, "face_blur"] = variance

    return df


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--data_file", default="data.h5", help="Name of data file")
    args = parser.parse_args()

    df = pd.read_hdf(args.data_file)
    df = calc_face_features(df)
    df.to_hdf(args.data_file, key="images")
