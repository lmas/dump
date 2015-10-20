#!/usr/bin/env python2

EVENT_UNKNOWN = -1
EVENT_BUTTON = 1
EVENT_AXIS = 2

BUTTON_UP = 0
BUTTON_DOWN = 1

class Gamepad(object):
    def __init__(self, button_callback, axis_callback):
	self.button_callback = button_callback
	self.axis_callback = axis_callback
        self.pipe = open('/dev/input/js0', 'r')

    def _read_buffer(self):
        buffer = [ord(char) for char in self.pipe.read(8)]
        return buffer

    def _parse_buffer(self, buffer):
	event_id = buffer[6]
	input_id = buffer[7]
	button_state = buffer[4]
	axis_state = buffer[5]

        if event_id == EVENT_BUTTON:
	    return self.button_callback(input_id, button_state)

        elif event_id == EVENT_AXIS:
	    return self.axis_callback(input_id, axis_state)

    def update(self):
	buffer = self._read_buffer()
	self._parse_buffer(buffer)

################################################################################

def debug_button(id, state):
	print('button #%i: %i' % (id, state))

def debug_axis(id, state):
	print('axis #%i: %i' % (id, state))

def main():
    gp = Gamepad(debug_button, debug_axis)
    while 1:
	gp.update()

if __name__ == '__main__':
	main()
