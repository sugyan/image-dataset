import tensorflow_hub as hub

IMAGE_SIZE = (299, 299)


def cnn(trainable=True):
    return hub.KerasLayer("https://tfhub.dev/google/imagenet/inception_v3/feature_vector/4",
                          trainable=trainable, arguments=dict(batch_norm_momentum=0.997))
