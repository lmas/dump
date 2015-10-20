
import logging

settings = dict(
    LOGGING_FORMAT = '%(asctime)s %(levelname)s %(message)s',
    LOGGING_LEVEL = logging.INFO,

    SERVER_ADDRESS = '',
    SERVER_PORT = 12345,
    FPS = 60,

    # Number of tiles / second
    SPEED_MAX = 4,
    # Number of tiles / second to accelerate up to the max speed
    SPEED_ACCELERATION = 0.1,

    STARTMAP = '25678a09580d4ad2a7f382132bf4f1da',
    TILEMAPS = {
        '25678a09580d4ad2a7f382132bf4f1da': './maps/test.map',
    },
)
