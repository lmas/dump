#!/usr/bin/env python
import random

class FractalMap:
	def __init__(self):
		# map size; must be a power of 2
		self.size = 64
		self.map = []
		for i in xrange(self.size * self.size):
			self.map.append(' ')

		# initial pattern, from where the fractal generation is based on
		self.set(0, 0, '~')
		self.set(self.size/2, 0, '~')
		self.set(0, self.size/2, '~')
		self.set(self.size/2, self.size/2, '*')

	def get(self, x, y):
		return self.map[x + self.size * y]

	def set(self, x, y, tile):
		self.map[x + self.size * y] = tile

	def show(self):
		buff = ''
		for y in xrange(self.size):
			line = ''
			for x in xrange(self.size):
				line += str(self.get(x, y))
			buff += line + '\n'
		print buff

	def fractal(self, step):
		# step = detail square width
		for y in xrange(0, self.size, step):
			for x in xrange(0, self.size, step):
				# The inner loop calculates (cx,cy)
				# this is the point from which to copy map colour

				# add random offsets
				cx = x + ((random.random() < 0.5) * step)
				cy = y + ((random.random() < 0.5) * step)

				# truncate to nearest multiple of step*2
				# since step*2 is the previous detail level calculated
				cx = (cx / (step * 2)) * step * 2
				cy = (cy / (step * 2)) * step * 2

				# wrap around to reference world as torus
				# also guarantees getCell() and setCell() are within bounds
				cx = cx % self.size
				cy = cy % self.size

				# copy from (cx,cy) to (x,y)
				self.set(x, y, self.get(cx, cy))

		# recursively calculate finer detail levels
		if (step > 1): self.fractal(step/2)

###############################################################################

class Fractal:
	def __init__(self):
		# map size; must be a power of 2
		self.sx = 64
		self.sy = 64
		self.map = []
		for y in xrange(self.sy):
			self.map.append(['o' for x in xrange(self.sx)])

		# initial pattern, from where the fractal generation is based on
		self.set(0, 0, '~')
		self.set(self.sx/2, 0, '~')
		self.set(0, self.sy/2, '~')
		self.set(self.sx/2, self.sy/2, '*')
		#~for y in xrange(self.sy):
			#~for x in xrange(self.sx):
				#~self.set(x, y, '#')

	def get(self, x, y):
		return self.map[y][x]

	def set(self, x, y, tile):
		self.map[y][x] = tile

	def show(self):
		buff = ''
		for y in xrange(self.sy):
			line = ''
			for x in xrange(self.sx):
				line += str(self.get(x, y))
			buff += line + '\n'
		print buff

	def fractal(self, stepx, stepy):
		# step = detail square width
		for y in xrange(0, self.sy, stepy):
			for x in xrange(0, self.sx, stepx):
				# The inner loop calculates (cx,cy)
				# this is the point from which to copy map colour

				# add random offsets
				cx = x + ((random.random() < 0.5) * stepx)
				cy = y + ((random.random() < 0.5) * stepy)

				# truncate to nearest multiple of step*2
				# since step*2 is the previous detail level calculated
				cx = (cx / (stepx * 2)) * stepx * 2
				cy = (cy / (stepy * 2)) * stepy * 2

				# wrap around to reference world as torus
				# also guarantees getCell() and setCell() are within bounds
				cx = cx % self.sx
				cy = cy % self.sy

				# copy from (cx,cy) to (x,y)
				self.set(x, y, self.get(cx, cy))

		# recursively calculate finer detail levels
		if (stepx > 1) and (stepy > 1): self.fractal(stepx/2, stepy/2)

if __name__ == '__main__':
	test = Fractal()
	test.fractal(16, 16)
	test.show()
