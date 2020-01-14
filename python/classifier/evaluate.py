import argparse
import os
import tensorflow as tf
from model import cnn, IMAGE_SIZE


def evaluate(data_dir, weights_path, labels_file):
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
    model.compile(
        loss=tf.keras.losses.CategoricalCrossentropy(),
        metrics=[tf.keras.metrics.CategoricalAccuracy()])

    evaluation_datagen = tf.keras.preprocessing.image.ImageDataGenerator(
        rescale=1./255)
    evaluation_data = evaluation_datagen.flow_from_directory(
        os.path.join(data_dir, 'validation'),
        target_size=(299, 299),
        classes=labels,
        batch_size=1)
    result = model.evaluate(evaluation_data)
    print(result)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--data_dir', default=os.path.join(os.path.dirname(__file__), '..', '..', 'images'))
    parser.add_argument('--weights_path', required=True)
    parser.add_argument('--labels_file', default=os.path.join(os.path.dirname(__file__), 'labels.txt'))
    args = parser.parse_args()
    evaluate(args.data_dir, args.weights_path, args.labels_file)
