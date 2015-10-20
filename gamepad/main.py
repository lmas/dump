#!/usr/bin/env python2

import libardrone

import gamepad


drone = libardrone.ARDrone(True)
drone.reset()

def update_battery():
	bat = drone.navdata.get(0, dict()).get('battery', 0)
	if bat != old_bat:
		old_bat = bat
		print 'battery: %i%%' % bat

def set_speed(speed):
	tmp = max(0.1, min(1.0, drone.speed + speed))
	print 'speed:', tmp
	drone.speed = tmp

def button(id, state):
	print 'button', id, state
	if state == 0:
		drone.hover()
	else:
		if id == 7:
			drone.reset()
		elif id == 8:
			drone.takeoff()
		elif id == 9:
			drone.land()
		elif id == 5:
			set_speed(+0.1)
		elif id == 4:
			set_speed(-0.1)
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
	old_bat = 0
	running = True
	
	print 'running...'
	try:
		while running:
			update_battery()
			running = gp.update()
	except:
		drone.reset()

	print 'shutting down...'
	drone.halt()
	print 'exit'
