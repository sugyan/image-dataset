import argparse
import csv
import os
import tensorflow as tf
from model import cnn, IMAGE_SIZE


def predict(data_dir, weights_path, labels_file):
    labels = []
    with open(labels_file, 'r') as fp:
        labels = [line.strip() for line in fp.readlines()]

    model = tf.keras.Sequential([
        cnn(trainable=False),
        tf.keras.layers.Dense(
            len(labels),
            trainable=False,
            activation='softmax')
    ])
    model.build([None, *IMAGE_SIZE, 3])
    model.summary()
    model.load_weights(weights_path)

    with open('results.tsv', 'w') as fp:
        writer = csv.writer(fp, delimiter='\t')
        for root, dirs, files in os.walk(os.path.join(data_dir, 'validation')):
            if not files:
                continue
            class_name = os.path.basename(root)
            if class_name not in labels:
                continue
            label = labels.index(class_name)
            for filename in files:
                image = tf.io.decode_jpeg(tf.io.gfile.GFile(os.path.join(root, filename), 'rb').read())
                images = tf.expand_dims(tf.image.convert_image_dtype(image, dtype=tf.float32), axis=0)
                result = int(model.predict(images).argmax())
                writer.writerow([
                    os.path.abspath(os.path.join(root, filename)),
                    label,
                    result
                ])


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--data_dir', default=os.path.join(os.path.dirname(__file__), '..', '..', 'images'))
    parser.add_argument('--weights_path', required=True)
    parser.add_argument('--labels_file', default=os.path.join(os.path.dirname(__file__), 'labels.txt'))
    args = parser.parse_args()
    predict(args.data_dir, args.weights_path, args.labels_file)
