import argparse
import logging
import math
import cv2
import dlib
import numpy as np


class Detector():
    def __init__(self, datafile='shape_predictor_68_face_landmarks.dat', verbose=False):
        logging.basicConfig(format='%(asctime)s %(levelname)s %(message)s')
        self.logger = logging.getLogger(__name__)
        if verbose:
            self.logger.setLevel(logging.INFO)
        self.detector = dlib.get_frontal_face_detector()
        self.predictor = dlib.shape_predictor(datafile)
        self.angles = [-48, -36, -24, -12, 0, 12, 24, 36, 48]

    def detect(self, img):
        # Create a large image that does not protrude by rotation
        h, w, c = img.shape
        hypot = math.ceil(math.hypot(h, w))
        hoffset = round((hypot-h)/2)
        woffset = round((hypot-w)/2)
        padded = np.zeros((hypot, hypot, c), np.uint8)
        padded[hoffset:hoffset+h, woffset:woffset+w, :] = img

        # Attempt detection by rotating at multiple angles
        results = []
        for angle in self.angles:
            rotated = self._rotate(padded, angle)
            dets, scores, indices = self.detector.run(rotated, 0, 0.0)
            self.logger.info(f'{angle:3d}: {dets}, {scores}, {indices}')
            if len(dets) == 1:
                results.append([dets[0], scores[0], angle, rotated])
        if len(results) == 0:
            self.logger.info('there are no detected faces')
            return

        # Choose the best angle by scores, and then adjust the angle using the eyes coordinates
        results.sort(key=lambda x: x[1], reverse=True)
        det, _, angle, rotated = results[0]
        shape = self.predictor(rotated, det)
        eyel, eyer = self._eye_center(shape)
        d = eyer - eyel
        angle += math.degrees(math.atan2(d[1], d[0]))
        self.logger.info(f'angle: {angle:.5f}')

        # Detect face and shapes from adjusted angle
        adjusted = self._rotate(padded, angle)
        dets = self.detector(adjusted)
        if len(dets) != 1:
            self.logger.info('faces are not detected in the rotated image')
            return
        shape = self.predictor(adjusted, dets[0])

        # Create a large mirrored image to rotate and crop
        margin = math.ceil(hypot * (math.sqrt(2) - 1.0) / 2)
        mirrored = np.pad(
            img,
            ((hoffset + margin, hypot - h - hoffset + margin),
             (woffset + margin, hypot - w - woffset + margin),
             (0, 0)), mode='symmetric')
        rotated = self._rotate(mirrored, angle)[margin:margin+hypot, margin:margin+hypot, :]

        # Calculate the center position and cropping size
        e0, e1 = self._eye_center(shape)
        m0 = np.array([shape.part(48).x, shape.part(48).y])
        m1 = np.array([shape.part(54).x, shape.part(54).y])
        x = e1 - e0
        y = (e0 + e1) / 2 - (m0 + m1) / 2
        c = (e0 + e1) / 2 + y * 0.1
        s = max(np.linalg.norm(x) * 4.0, np.linalg.norm(y) * 3.6)

        xoffset = int(np.rint(c[0] - s/2))
        yoffset = int(np.rint(c[1] - s/2))
        if xoffset < 0 or yoffset < 0 or xoffset + s >= hypot or yoffset + s >= hypot:
            self.logger.info('cropping area has exceeded the image area')
            return
        size = int(np.rint(s))
        cropped = rotated[yoffset:yoffset+size, xoffset:xoffset+size, :]

        # Attempt detection on the cropped image
        dets = self.detector(cropped)
        if len(dets) != 1:
            self.logger.info('faces are not detected in the cropped image')
            return
        shape = self.predictor(cropped, dets[0])

        return {
            'image': cropped,
            'parts': [(point.x, point.y) for point in shape.parts()],
            'angle': angle,
            'size': size,
        }

    def _rotate(self, img, angle):
        h, w, _ = img.shape
        mat = cv2.getRotationMatrix2D((w/2, h/2), angle, 1.0)
        return cv2.warpAffine(img, mat, (w, h), cv2.INTER_LANCZOS4)

    def _eye_center(self, shape):
        eyel, eyer = np.array([0, 0]), np.array([0, 0])
        for i in range(36, 42):
            eyel[0] += shape.part(i).x
            eyel[1] += shape.part(i).y
        for i in range(42, 48):
            eyer[0] += shape.part(i).x
            eyer[1] += shape.part(i).y
        return eyel / 6, eyer / 6


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('image_file')
    parser.add_argument('-v', '--verbose', action='store_true')
    args = parser.parse_args()

    result = Detector(verbose=args.verbose).detect(cv2.imread(args.image_file))
    if result is None:
        print('detection failed.')
        exit(0)

    img = result['image']
    for part in result['parts']:
        cv2.drawMarker(img, part, (255, 255, 0))
    cv2.imshow('image', img)
    cv2.waitKey(0)
    cv2.destroyAllWindows()
