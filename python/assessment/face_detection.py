import argparse
import dlib
import pandas as pd


def calc_face_detection(df, predictor_path):
    parts_names = [f"parts_{i:02d}" for i in range(68)]
    df_faces = pd.DataFrame(
        index=df.index,
        columns=pd.MultiIndex.from_product(
            [["face"], ["score", "left", "top", "right", "bottom"]]
        ),
    )
    df_parts = pd.DataFrame(
        index=df.index,
        columns=pd.MultiIndex.from_product(
            [parts_names, ["x", "y"]]
        ),
    )

    detector = dlib.get_frontal_face_detector()
    predictor = dlib.shape_predictor(predictor_path)
    for row in df.itertuples():
        image = dlib.load_rgb_image(row.Index)
        detections, scores, indices = detector.run(image, 0, 0.0)
        if len(detections) != 1:
            continue

        shape = predictor(image, detections[0])
        print(row.Index, detections, scores, shape.num_parts)
        rect = detections[0]
        df_faces.at[row.Index] = (
            scores[0],
            rect.left(),
            rect.top(),
            rect.right(),
            rect.bottom(),
        )
        for i in range(shape.num_parts):
            df_parts.loc[row.Index][parts_names[i]] = (
                shape.part(i).x,
                shape.part(i).y,
            )
    df = pd.concat([df, df_faces, df_parts], axis=1)
    return df


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--data_file", default="data.h5", help="Name of data file")
    parser.add_argument(
        "--predictor_path",
        default="shape_predictor_68_face_landmarks.dat",
        help="Path to trained model file",
    )
    args = parser.parse_args()

    df = pd.read_hdf(args.data_file)
    df = calc_face_detection(df, args.predictor_path)
    df.to_hdf(args.data_file, key="images")
