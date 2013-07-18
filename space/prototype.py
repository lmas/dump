#!/usr/bin/env python

import random
import math

################################################################################

def vector_add(v1, v2):
    return (v1[0] + v2[0], v1[1] + v2[1])

def vector_sub(v1, v2):
    return (v1[0] - v2[0], v1[1] - v2[1])

def vector_mul(v1, v2):
	return (v1[0] * v2[0], v1[1] * v2[1])

def vector_div(v1, v2):
	return (v1[0] / v2[0], v1[1] / v2[1])

def vector_str(vector):
    return ','.join((str(vector[0]), str(vector[1])))

def vector_distance(v1, v2):
    x, y = vector_sub(v1, v2)
    return math.sqrt(x * x + y * y)

################################################################################

tile_grid = {}
world_seed = 0

def gen_world():
    size = 20
    for y in xrange(-size, size + 1):
        for x in xrange(-size, size + 1):
            seed = '%i:%i,%i' % (world_seed, x, y)
            random.seed(seed)
            if random.random() > 0.95:
                tile_grid[(x,y)] = '*'

def load_world():
	with open('map', 'r') as f:
		l = f.readline().strip()
		while l:
			pos = l.split(' ')
			pos = (int(pos[0]), int(pos[1]))
			# this block should be done in the galaxy generator instead
			#delta = (3,3)
			#pos = vector_mul(pos, delta) # scale up
			#delta = (random.randint(-1, 1), random.randint(-1, 1))
			#pos = vector_add(pos, delta) # fuzzy the new up-scaled coords
			# end block
			pos = vector_mul(pos, (10,10))
			tile_grid[pos] = '*'
			l = f.readline().strip()

def draw_map(vector, distance):
    start = vector_sub(vector, (distance, distance))
    stop = vector_add(vector, (distance+1, distance+1))
    grid = []
    for y in xrange(start[1], stop[1]):
        line = []
        for x in xrange(start[0], stop[0]):
            try:
                tmp = tile_grid[(x,y)]
            except KeyError:
                tmp = ' '
            line.append(tmp)
        grid.append(line)
    print 'tiles:', len(tile_grid)
    print 'coords:', vector_str(vector)
    grid[distance][distance] = '@'
    for line in grid:
        print ''.join(line)
    #with open('map.txt', 'w') as f:
     #   for line in grid:
      #      f.write(''.join(line))
       #     f.write('\n')

def draw_radar(vector, distance):
    start = vector_sub(vector, (distance, distance))
    stop = vector_add(vector, (distance+1, distance+1))
    grid = []
    for y in xrange(start[1], stop[1]):
        line = []
        for x in xrange(start[0], stop[0]):
            tmp = tile_grid.get((x,y), ' ')
            if tmp == ' ' and (x,y) in obj_grid:
                tmp = '.'
            line.append(tmp)
        grid.append(line)
    grid[distance][distance] = '@'
    print 'objs:', len(objects)
    print obj_grid
    for line in grid:
        print ''.join(line)

def draw_long(vector, distance):
    scale = 10
    tdistance = distance * scale
    start = vector_sub(vector, (tdistance, tdistance))
    stop = vector_add(vector, (tdistance+1, tdistance+1))
    grid = []
    dic = {}
    for y in xrange(start[1], stop[1]):
        line = []
        for x in xrange(start[0], stop[0]):
            tmp = tile_grid.get((x,y), '')
            if tmp != '':
                pos = vector_div((x,y), (scale, scale))
                dic[pos] = tmp
            #line.append(tmp)
        #grid.append(line)
    #grid[distance][distance] = '@'
    dic[vector_div(vector, (scale,scale))] = '@'
    print dic
    grid = [[' ' for y in xrange(tdistance+1)] for x in xrange(tdistance+1)]
    for k,v in dic.iteritems():
    	x, y = k
    	#print k, v, x, y
    	grid[x][y] = v
    #print grid
    #print 'objs:', len(objects)
    #print obj_grid
    for line in grid:
        print ''.join(line)

################################################################################

objects = {}
obj_grid = {}
obj_ids = 0

def new_obj():
    global obj_ids
    id = obj_ids
    obj_ids += 1
    objects[id] = {'char':'?', 'pos':(0,0)}
    return id

def move(id, vector, delta):
    rem(id, vector)
    pos = vector_add(vector, delta)
    objects[id]['pos'] = pos
    add(id, pos)

def add(id, vector):
    if vector in obj_grid:
        obj_grid[vector].append(id)
    else:
        obj_grid[vector] = [id]

def rem(id, vector):
    if len(obj_grid[vector]) > 1:
        obj_grid[vector].remove(id)
    else:
        del obj_grid[vector]

def update_objs():
    for id in xrange(1, obj_ids):
        vector = objects[id]['pos']
        delta = (random.choice((-1, 0, 1)), random.choice((-1, 0, 1)))
        move(id, vector, delta)

################################################################################

player = new_obj()
#gen_world()
load_world()
print 'stars:', len(tile_grid)
add(player, objects[player]['pos'])
running = True

while running:
    update_objs()
    data = raw_input('>').strip().lower()
    try:
        cmd, arg = data.split(' ', 1)
    except ValueError:
        cmd, arg = data, ''
        
    if cmd == 'm': # map shows the world
        try:
            distance = int(arg)
        except ValueError:
            distance = 10
        draw_map(objects[player]['pos'], distance)

    elif cmd == 'l': # long range star map
    	try:
    		distance = int(arg)
    	except ValueError:
    		distance = 10
    	draw_long(objects[player]['pos'], distance)

    elif cmd == 'r': # radar shows objects
        try:
            distance = int(arg)
        except ValueError:
            distance = 10
        draw_radar(objects[player]['pos'], distance)

    elif cmd == 't': # create a target
        target = new_obj()
        objects[target]['pos'] = objects[player]['pos']
        add(target, objects[target]['pos'])
        print 'target'

    elif cmd == 'd': # distance to target
        try:
            id = int(arg)
        except ValueError:
            pass
        else:
            pos = objects[player]['pos']
            target = objects[id]['pos']
            dis = vector_distance(pos, target)
            print 'player:', pos
            print 'target:', target
            print 'distance:', dis

    elif cmd == 'j': # jump to relative vector
        try:
            x, y = arg.split(',', 1)
            x, y = int(x), int(y)
        except ValueError:
            pass
        else:
            print 'Jump'
            move(player, objects[player]['pos'], (x, y))

    # simple vector movement
    elif cmd == 'n':
        move(player, objects[player]['pos'], (0, -1))
    elif cmd == 's':
        move(player, objects[player]['pos'], (0, 1))
    elif cmd == 'w':
        move(player, objects[player]['pos'], (-1, 0))
    elif cmd == 'e':
        move(player, objects[player]['pos'], (1, 0))
    elif cmd == 'q': # quit
        running = False
    else:
        print 'error:', cmd, ',', arg

