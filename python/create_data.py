import argparse
import csv
import json
import os
import time
import cv2
import numpy as np
from datetime import datetime
from urllib.request import urlopen
from urllib.error import HTTPError

from detect import Detector


def run(datafile, data_dir):
    def download(url):
        try:
            with urlopen(url) as response:
                if response.status == 200:
                    data = np.asarray(bytearray(response.read()), dtype=np.uint8)
                    return cv2.imdecode(data, cv2.IMREAD_COLOR)
        except HTTPError as e:
            print(f'{url}: {e}')
        except (ConnectionError, OSError) as e:
            print(f'{url}: {e}')
            time.sleep(0.1)
        except Exception as e:
            print(f'{url}: {e}')
            time.sleep(0.5)
        return None

    detector = Detector()
    with open(datafile, 'r') as fp:
        r = csv.reader(fp, delimiter='\t')
        for row in r:
            photo_url = row[1]
            print(photo_url)
            img = download(photo_url)
            if img is None:
                continue
            result = detector.detect(img)
            if result is None:
                continue

            basename = os.path.basename(photo_url)
            path0 = f'{ord(basename[0]):02x}'
            path1 = f'{ord(basename[1]):02x}'
            path2 = f'{ord(basename[2]):02x}'
            outdir = os.path.join(data_dir, path0, path1, path2)
            os.makedirs(outdir, exist_ok=True)

            faceimg = result['image']
            del result['image']
            created_at = datetime.fromtimestamp(time.mktime(time.strptime(row[3], '%a %b %d %H:%M:%S +0000 %Y')))
            result['meta'] = {
                'photo_id': row[0],
                'photo_url': row[1],
                'source_url': row[2],
                'published_at': created_at.isoformat(),
                'label_id': row[4],
                'label_name': row[5],
            }

            name = os.path.splitext(basename)[0]
            cv2.imwrite(os.path.join(outdir, f'{name}.jpg'), faceimg, [cv2.IMWRITE_JPEG_QUALITY, 100])
            with open(os.path.join(outdir, f'{name}.json'), 'w') as fp:
                json.dump(result, fp, ensure_ascii=False)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--data_dir', default=os.path.join(os.path.dirname(__file__), 'data'))
    parser.add_argument('datafile')
    args = parser.parse_args()
    run(args.datafile, args.data_dir)
