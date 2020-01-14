import argparse
import os
import tensorflow as tf
from model import cnn, IMAGE_SIZE

BATCH_SIZE = 32


def train(data_dir, weights_dir, labels_file):
    os.makedirs(weights_dir, exist_ok=True)
    labels = []
    with open(labels_file, 'r') as fp:
        labels = [line.strip() for line in fp.readlines()]

    model = tf.keras.Sequential([
        cnn(),
        tf.keras.layers.Dropout(rate=0.1),
        tf.keras.layers.Dense(
            len(labels),
            activation='softmax',
            kernel_regularizer=tf.keras.regularizers.l2(1e-4)),
    ])
    model.build([None, *IMAGE_SIZE, 3])
    model.summary()
    model.compile(
        optimizer=tf.keras.optimizers.RMSprop(),
        loss=tf.keras.losses.CategoricalCrossentropy(),
        metrics=[tf.keras.metrics.CategoricalAccuracy()])

    training_datagen = tf.keras.preprocessing.image.ImageDataGenerator(
        rotation_range=2,
        width_shift_range=2,
        height_shift_range=2,
        brightness_range=(0.8, 1.2),
        channel_shift_range=0.2,
        zoom_range=0.02,
        rescale=1./255)
    training_data = training_datagen.flow_from_directory(
        os.path.join(data_dir, 'training'),
        target_size=(299, 299),
        classes=labels,
        batch_size=BATCH_SIZE)
    validation_datagen = tf.keras.preprocessing.image.ImageDataGenerator(
        rescale=1./255)
    validation_data = validation_datagen.flow_from_directory(
        os.path.join(data_dir, 'validation'),
        target_size=(299, 299),
        classes=labels,
        batch_size=BATCH_SIZE)

    tf.keras.backend.clear_session()
    history = model.fit(
        training_data,
        epochs=100,
        validation_data=validation_data,
        callbacks=[
            tf.keras.callbacks.TensorBoard(),
            tf.keras.callbacks.ModelCheckpoint(
                os.path.join(weights_dir, 'finetuning_weights-{epoch:02d}.h5'),
                save_weights_only=True),
        ])
    print(history.history)
    model.trainable = False
    model.save('finetuning_classifier.h5')


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--data_dir', default=os.path.join(os.path.dirname(__file__), '..', '..', 'images'))
    parser.add_argument('--weights_dir', default=os.path.join(os.path.dirname(__file__), 'weights'))
    parser.add_argument('--labels_file', default=os.path.join(os.path.dirname(__file__), 'labels.txt'))
    args = parser.parse_args()
    train(args.data_dir, args.weights_dir, args.labels_file)
