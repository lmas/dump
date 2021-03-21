#!/usr/bin/env python2

import libardrone

import gamepad


drone = libardrone.ARDrone(True)
drone.reset()

def set_speed(speed):
	speed = max(0.1, min(1.0, speed))
	print 'speed:', speed
	drone.speed = speed

def button(id, state):
	print 'button', id, state
	if state == 0:
		drone.hover()
	else:
		if id == 6:
			drone.reset()
		elif id == 7:
			drone.takeoff()
		elif id == 8:
			drone.land()
		elif id == 5:
			speed = drone.speed + 0.1
			set_speed(speed)
		elif id == 4:
			speed = drone.speed - 0.1
			set_speed(speed)
	return True

def axis(id, state):
	print 'axis', id, state
	if state == 0:
		drone.hover()
	else:
		if id == 1:
			if state > 127:
				drone.move_forward()
			else:
				drone.move_backward()
		elif id == 0:
			if state > 127:
				drone.move_left()
			else:
				drone.move_right()
		elif id == 3:
			if state > 127:
				drone.move_up()
			else:
				drone.move_down()
		elif id == 2:
			if state > 127:
				drone.turn_left()
			else:
				drone.turn_right()
	return True

def main():
	gp = gamepad.Gamepad(button, axis)
	running = True
	
	print 'running...'
	try:
		while running:
			running = gp.update()
	except:
		drone.reset()

	print 'shutting down...'
	drone.halt()
	print 'exit'
