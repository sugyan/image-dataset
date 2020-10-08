import argparse
import dlib
import pandas as pd


def detect_single_face(detector, image):
    for adjust_threshold in [x / 20.0 for x in range(0, -10, -1)]:
        for upsample_num_times in range(0, 3):
            detections, scores, indices = detector.run(
                image, upsample_num_times, adjust_threshold
            )
            if len(detections) == 1:
                return detections, scores, indices
    return None, None, None


def calc_face_detection(df, predictor_path):
    detector = dlib.get_frontal_face_detector()
    predictor = dlib.shape_predictor(predictor_path)

    for row in df[df["face_score"].isnull()].itertuples():
        print(row.Index)
        image = dlib.load_rgb_image(row.Index)
        detections, scores, indices = detect_single_face(detector, image)
        if detections is None:
            print("failed to detect the face")
            continue
        rect = detections[0]
        df.loc[row.Index, "face_score"] = scores[0]
        df.loc[row.Index, "face_left"] = rect.left()
        df.loc[row.Index, "face_top"] = rect.top()
        df.loc[row.Index, "face_right"] = rect.right()
        df.loc[row.Index, "face_bottom"] = rect.bottom()
        shape = predictor(image, rect)
        for i in range(shape.num_parts):
            df.loc[row.Index, f"parts{i:02d}_x"] = shape.part(i).x
            df.loc[row.Index, f"parts{i:02d}_y"] = shape.part(i).y
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
