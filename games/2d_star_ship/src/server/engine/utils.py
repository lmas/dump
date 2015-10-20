
import math


def sign(val):
    if val > 0:
        return 1
    elif val < 0:
        return -1
    return 0


class Vector(object):
    def __init__(self, x=0.0, y=0.0):
        self.x, self.y = x, y

    def round(self, digits=0):
        return Vector(round(self.x, digits), round(self.y, digits))

    def len(self):
        return math.sqrt(self.x * self.x + self.y * self.y)

    def distance(self, other):
        tmp = self - other
        return tmp.len()

    def __check_type(self, other, ignore_int=False):
        # Yep, it's pretty nasty...
        if isinstance(other, Vector):
            return other.x, other.y
        elif isinstance(other, tuple) or isinstance(other, list):
            if (isinstance(other[0], int) or isinstance(other[0], float)) and (
                isinstance(other[1], int) or isinstance(other[1], float)):
                return other[0], other[1]
            else:
                raise TypeError(
                    "unsupported operand type(s) inside '{}'".format(type(other).__name__)
                )
        elif not ignore_int and (isinstance(other, int) or isinstance(other, float)):
            return other, other
        else:
            raise TypeError(
                "unsupported operand type(s): 'vector' and '{}'".format(type(other).__name__)
            )

    def __str__(self):
        return '{}, {}'.format(self.x, self.y)

    def __repr__(self):
        return 'Vector({}, {})'.format(self.x, self.y)

    def __hash__(self):
        return hash((self.x, self.y))

    def __nonzero__(self):
        if self.x == 0 and self.y == 0:
            return False
        else:
            return True

    def __eq__(self, other):
        x, y = self.__check_type(other, True)
        if self.x == x and self.y == y:
            return True
        else:
            return False

    def __ne__(self, other):
        x, y = self.__check_type(other, True)
        if self.x != x or self.y != y:
            return True
        else:
            return False

    def __add__(self, other):
        x, y = self.__check_type(other)
        return Vector(self.x + x, self.y + y)

    def __iadd__(self, other):
        self = self.__add__(other)
        return self

    def __sub__(self, other):
        x, y = self.__check_type(other)
        return Vector(self.x - x, self.y - y)

    def __isub__(self, other):
        self = self.__sub__(other)
        return self

    def __mul__(self, other):
        x, y = self.__check_type(other)
        return Vector(self.x * x, self.y * y)

    def __imul__(self, other):
        self = self.__mul__(other)
        return self

    def __div__(self, other):
        x, y = self.__check_type(other)
        return Vector(self.x / x, self.y / y)

    def __idiv__(self, other):
        self = self.__div__(other)
        return self

