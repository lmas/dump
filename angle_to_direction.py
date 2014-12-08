
import math

def angle_to_direction(angle):
    '''Convert an angle, in degrees, to a cardinal direction.'''
    bearings = ['N', 'NNE', 'NE', 'ENE', 'E', 'ESE', 'SE', 'SSE',
                'S', 'SSW', 'SW', 'WSW', 'W', 'WNW', 'NW', 'NNW']

    # Normalize the angle (normalized = 0 < angle < 360)
    tmp = angle / 360.0
    normalized_angle = (tmp - math.floor(tmp)) * 360.0

    # Calculate index out of the angle and grab the direction
    tmp = 360.0 / len(bearings)
    index = int(round(normalized_angle / tmp )) % len(bearings)
    return bearings[index]
