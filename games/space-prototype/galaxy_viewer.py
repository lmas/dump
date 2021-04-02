
# Source: http://codentronix.com/2011/04/20/simulation-of-3d-point-rotation-with-python-and-pygame/
import sys, math, pygame, random

class Point3D:
    def __init__(self, x = 0, y = 0, z = 0):
        self.x, self.y, self.z = float(x), float(y), float(z)

    def rotateX(self, angle):
        """ Rotates the point around the X axis by the given angle in degrees. """
        rad = angle * math.pi / 180
        cosa = math.cos(rad)
        sina = math.sin(rad)
        y = self.y * cosa - self.z * sina
        z = self.y * sina + self.z * cosa
        return Point3D(self.x, y, z)

    def rotateY(self, angle):
        """ Rotates the point around the Y axis by the given angle in degrees. """
        rad = angle * math.pi / 180
        cosa = math.cos(rad)
        sina = math.sin(rad)
        z = self.z * cosa - self.x * sina
        x = self.z * sina + self.x * cosa
        return Point3D(x, self.y, z)

    def rotateZ(self, angle):
        """ Rotates the point around the Z axis by the given angle in degrees. """
        rad = angle * math.pi / 180
        cosa = math.cos(rad)
        sina = math.sin(rad)
        x = self.x * cosa - self.y * sina
        y = self.x * sina + self.y * cosa
        return Point3D(x, y, self.z)

    def project(self, win_width, win_height, fov, viewer_distance):
        """ Transforms this 3D point to 2D using a perspective projection. """
        #print(fov, float(viewer_distance) + self.z)
        tmp = (viewer_distance + self.z)
        if tmp == 0: tmp = 0.01
        factor = fov / tmp
        x = self.x * factor + win_width / 2
        y = -self.y * factor + win_height / 2
        return Point3D(x, y, 1)

class Simulation:
    def __init__(self, win_width = 640, win_height = 480):
        pygame.init()
        pygame.key.set_repeat(100, 25)

        self.screen = pygame.display.set_mode((win_width, win_height))
        pygame.display.set_caption("Space viewer")

        self.font_size = 12
        self.font = pygame.font.SysFont("monospace", self.font_size)

        self.clock = pygame.time.Clock()

        #self.vertices = [
            #Point3D(-1,1,-1),
            #Point3D(1,1,-1),
            #Point3D(1,-1,-1),
            #Point3D(-1,-1,-1),
            #Point3D(-1,1,1),
            #Point3D(1,1,1),
            #Point3D(1,-1,1),
            #Point3D(-1,-1,1)
        #]
        self.vertices = []
        with open('map.txt', 'r') as f:
            line = f.readline().strip()
            while line:
                x, y, z = line.split(' ')
                x, y, z = int(x), int(y), int(z)
                tmp = Point3D(x, y, z)
                self.vertices.append(tmp)
                line = f.readline().strip()
        print('Stars loaded:', len(self.vertices))
        #self.vertices = []
        #self.range = 10
        #for i in xrange(2000):
            #tmp = Point3D(random.gauss(0, 5), random.gauss(0, 5), random.gauss(0, 5))
            #self.vertices.append(tmp)

        self.angleX, self.angleY, self.angleZ = 0, 0, 0
        self.mx, self.my, self.zoom = 0, 0, 100
        #self.px, self.py = 0, 0

    def text(self, text, x, y):
        label = self.font.render(text, 1, (255,255,255))
        self.screen.blit(label, (x, y))

    def point(self, x, y, rad):
        pygame.draw.circle(self.screen, (255,255,255), (x, y), rad / self.zoom)

    def run(self):
        while 1:
            w, h = self.screen.get_width(), self.screen.get_height()
            pygame.event.pump()
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    sys.exit()

                elif event.type == pygame.MOUSEBUTTONDOWN:
                    if event.button == 4:
                        self.zoom = max(2, min(100, self.zoom - 1))
                    elif event.button == 5:
                        self.zoom = max(2, min(100, self.zoom + 1))
                    #elif event.button == 1:
                        #x, y = event.pos
                        #x = max(1, min(w, x))
                        #y = max(1, min(h, y))
                        #print(x, y)
                        #x = (float(x) / w) + self.zoom #- self.angleX
                        #x = (x / self.zoom) #- self.angleX
                        #y = (float(y) / h) + self.zoom #- self.angleY
                        #print(x, y)
                        #self.px, self.py = int(x), int(y)
                        #print(self.px, self.py)
                        #print((x*w)/self.zoom, (y*h)/self.zoom)

                elif event.type == pygame.MOUSEMOTION:
                    if event.buttons[0] == 1:
                        self.mx = float(event.rel[1])
                        self.my = float(event.rel[0])
                        self.angleX -= self.mx
                        self.angleY -= self.my

                elif event.type == pygame.KEYDOWN:
                    k = event.key
                    if k == pygame.K_s: self.angleX -= 1
                    elif k == pygame.K_w: self.angleX += 1
                    elif k == pygame.K_d: self.angleY -= 1
                    elif k == pygame.K_a: self.angleY += 1
                    elif k == pygame.K_e: self.angleZ -= 1
                    elif k == pygame.K_q: self.angleZ += 1

            self.clock.tick(50)
            self.screen.fill((0,0,0))

            shown = 0
            for v in self.vertices:
                # Rotate the point around X axis, then around Y axis, and finally around Z axis.
                r = v.rotateX(self.angleX).rotateY(self.angleY).rotateZ(self.angleZ)
                # Transform the point from 3D to 2D
                #print(r.x, r.y, r.z)
                p = r.project(self.screen.get_width(), self.screen.get_height(), 256, self.zoom)
                x, y = int(p.x), int(p.y)
                if x > 1 and x < w and y > 1 and y < h:
                    self.point(x, y, 10)
                    shown += 1

            #if self.px != 0 and self.py != 0:
                #self.point(self.px*w, self.py*h, 15)

            self.text('fps %i' % self.clock.get_fps(), 1, 0 * self.font_size)
            self.text('shown %i' % shown, 1, 1 * self.font_size)
            self.text('mouse %i, %i' % (self.mx, self.my), 1, 2 * self.font_size)
            self.text('zoom %i' % (self.zoom-1), 1, 3 * self.font_size)
            self.text('x %i' % self.angleX, 1, 4 * self.font_size)
            self.text('y %i' % self.angleY, 1, 5 * self.font_size)
            self.text('z %i' % self.angleZ, 1, 6 * self.font_size)

            pygame.display.flip()

if __name__ == "__main__":
    Simulation().run()
