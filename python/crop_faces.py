import argparse
import csv
import json
import os
import time
import cv2
from detect import Detector


def run(data_dir):
    detector = Detector()
    with open('data.tsv', 'r') as fp:
        r = csv.reader(fp, delimiter='\t')
        next(r)
        for row in r:
            photo_url = row[3]
            basename = os.path.basename(photo_url)

            path0 = f'{ord(basename[0]):02x}'
            path1 = f'{ord(basename[1]):02x}'
            path2 = f'{ord(basename[2]):02x}'
            filepath = os.path.join(data_dir, 'images', path0, path1, path2, basename)
            print(f'processing {filepath} ...')
            try:
                result = detector.detect(cv2.imread(filepath))
                if result is None:
                    print('detection failed.')
                    continue

                outdir = os.path.join(data_dir, 'results', path0, path1, path2)
                os.makedirs(outdir, exist_ok=True)
                img = result['image']
                del result['image']
                result['meta'] = {
                    'face_id': row[0],
                    'photo_id': row[1],
                    'source_url': row[2],
                    'photo_url': row[3],
                    'posted_at': row[4],
                    'label_id': row[5],
                    'label_name': row[6],
                }
                name = os.path.splitext(basename)[0]
                cv2.imwrite(os.path.join(outdir, f'{name}.png'), img)
                with open(os.path.join(outdir, f'{name}.json'), 'w') as fp:
                    json.dump(result, fp, ensure_ascii=False)
            except:
                print(f'error!!: {filepath}')

            time.sleep(1)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--data_dir', default=os.path.expanduser('~/data'))
    args = parser.parse_args()

    print(os.path.abspath(args.data_dir))
