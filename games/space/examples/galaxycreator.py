import os, sys
import pygame
import math, random

if not pygame.font: print 'Warning, fonts disabled'
if not pygame.mixer: print 'Warning, sound disabled'

class star:
	def __init__(self, x, y, color = (255,255,255)):
		self.x = x
		self.y = y
		self.color = color

class Game:
	def __init__(self, width, height):
		self.width, self.height = width, height
		pygame.init()
		pygame.display.set_caption("Galaxy Creator")
		self.window = pygame.display.set_mode((width, height))
		self.screen = pygame.display.get_surface()
		self.font = pygame.font.SysFont("Times New Roman",10)
		self.viewx, self.viewy = width/2, height/2
		self.running = True
		self.animate = False
		self.borders = True

		self.radius = 500.0
		self.arms = 6
		self.stars = 20000
		self.starmap = []

		#self.makeGalaxy()
		self.screen.blit(self.font.render("Press spacebar to generate a new galaxy!", True, (255,255,255)), (width/2-50, height/2))
		self.screen.blit(self.font.render("Arms: %i" % (self.arms), True, (255,255,255)), (0,0))
		self.screen.blit(self.font.render("Radius: %i" % (self.radius), True, (255,255,255)), (0,10))
		self.screen.blit(self.font.render("Stars: %i" % (self.stars), True, (255,255,255)), (0,20))
		pygame.display.update()
		#self.update()

	def update(self):
		self.screen.fill((0,0,0))
		for star in self.starmap:
			self.screen.set_at(((self.viewx + star.x), (self.viewy + star.y)), star.color)
			if self.animate: pygame.display.update()
		if self.borders:
			pygame.draw.circle(self.screen, (55,55,55), (self.viewx, self.viewy), int(self.radius), 1)
			pygame.draw.circle(self.screen, (255,0,0), (self.viewx, self.viewy), 1, 0)
		self.screen.blit(self.font.render("Arms: %i" % (self.arms), True, (255,255,255)), (0,0))
		self.screen.blit(self.font.render("Radius: %i" % (self.radius), True, (255,255,255)), (0,10))
		self.screen.blit(self.font.render("Stars: %i" % (self.stars), True, (255,255,255)), (0,20))
		pygame.display.update()

	def loop(self):
		changed = False
		while self.running:
			for event in pygame.event.get():
				if event.type == pygame.QUIT:
					#quit
					self.running = False
				elif event.type == pygame.KEYDOWN:
					if event.key == pygame.K_ESCAPE:
						#quit
						self.running = False
					elif event.key == pygame.K_F1:
						#toggle fullscren
						pygame.display.toggle_fullscreen()
					elif event.key == pygame.K_q:
						#decrease galaxy arms
						self.arms -= 1
						if self.arms < 1: self.arms = 1
						self.update()
					elif event.key == pygame.K_w:
						#increase galaxy arms
						self.arms += 1
						self.update()
					elif event.key == pygame.K_a:
						#decrease radius
						self.radius -= 50
						if self.radius < 1: self.radius = 1
						self.update()
					elif event.key == pygame.K_s:
						#increase radius
						self.radius += 50
						self.update()
					elif event.key == pygame.K_z:
						#decrease stars
						self.stars -= 1000
						if self.stars < 1: self.stars = 1
						self.update()
					elif event.key == pygame.K_x:
						#increase stars
						self.stars += 1000
						self.update()
					elif event.key == pygame.K_SPACE:
						self.makeGalaxy()
						self.update()
					elif event.key == pygame.K_RETURN:
						#reset the viewpos
						self.viewx, self.viewy = self.width/2, self.height/2
						self.update()
					elif event.key == pygame.K_i:
						#toggle visibility of extra "borders"
						if self.borders == True:
							self.borders = False
						else:
							self.borders = True
						self.update()
					elif event.key == pygame.K_p:
						#toggle "animation" of the star drawing
						if self.animate == True:
							self.animate = False
						else:
							self.animate = True
				elif event.type == pygame.MOUSEBUTTONDOWN and not self.animate:
					pos = pygame.mouse.get_pos()
					self.viewx, self.viewy = pos
				elif event.type == pygame.MOUSEBUTTONUP and not self.animate:
					self.update()
		sys.exit(0)

	def hat(self, range):
		area = 4 * math.atan(6.0)
		p = area * random.random()
		return math.tan(p/4) * range/6.0

	def makeGalaxy(self):
		self.starmap = []
		deg2rad = math.pi / 180.0
		armangle = ((360 / self.arms) %360)
		angularspread = 250 / self.arms

		num = 0
		temp = []
		for i in range(int(self.radius)):
			tmp = []
			for j in range(int(self.radius)):
				tmp.append([])
			temp.append(tmp)

		for i in range(0, self.stars):
			if (random.randint(0, 1) & 1): extra = 1.0
			else: extra = -1.0
			#R = random.uniform(0.0, radius)
			#Q = random.uniform(0.0, angularspread) * extra
			R = self.hat(self.radius)
			Q = self.hat(angularspread) * extra
			A = (random.randint(1, self.arms)) * armangle
			K = 1

			x = int(R * math.cos(deg2rad * (A + R * K + Q)))
			y = int(R * math.sin(deg2rad * (A + R * K + Q)))

			use = True
			#for j in self.starmap:
			#	if (j.x == x) and (j.y == y): use = False
			if temp[x][y] == 1: use = False
			if use:
				color = random.randint(1, 255)
				color = (color, color, color)
				self.starmap.append(star(x, y, color))
				temp[x][y] = 1
				num += 1
		print "STARS:", num

if __name__ == "__main__":
	game = Game(1024, 800)
	game.loop()
