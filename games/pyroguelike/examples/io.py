'''A custom I/O wrapper based on curses.'''
import curses, textwrap

class iowrapper:
	def __init__(self, func, *args, **kwds):
		'''Initiates the I/O Wrapper and various curses stuff.		
			NOTE: 
			most of this code in this function was simply copied from the 
			curses.wrapper module, for easier modification and customization,
			and was then modified for my own purpose!
			'''
		try:
			# Initialize curses
			self.__stdscr = curses.initscr()
			# Turn off echoing of keys, enter cbreak mode,
			# where no buffering is performed on keyboard input,
			#and finally turn of the blinking cursor
			curses.noecho()
			curses.cbreak()
			curses.curs_set(0)

			# In keypad mode, escape sequences for special keys
			# (like the cursor keys) will be interpreted and
			# a special value like curses.KEY_LEFT will be returned
			self.__stdscr.keypad(1)

			# Start color, too.  Harmless if the terminal doesn't have
			# color; user can test with has_color() later on.  The try/catch
			# works around a minor bit of over-conscientiousness in the curses
			# module -- the error return from C start_color() is ignorable.
			try:
				curses.start_color()
				bg = curses.COLOR_BLACK
				curses.init_pair(1, curses.COLOR_BLACK, bg)
				curses.init_pair(2, curses.COLOR_BLUE, bg)
				curses.init_pair(3, curses.COLOR_GREEN, bg)
				curses.init_pair(4, curses.COLOR_CYAN, bg)
				curses.init_pair(5, curses.COLOR_RED, bg)
				curses.init_pair(6, curses.COLOR_MAGENTA, bg)
				curses.init_pair(7, curses.COLOR_YELLOW, bg)
				curses.init_pair(8, curses.COLOR_WHITE, bg)
				self.colorsEnabled = True
			except:
				self.colorsEnabled = False
			
			self.__buffer = ''
			self.__width = self.__stdscr.getmaxyx()[1]
			self.__height = self.__stdscr.getmaxyx()[0]
			
			#instead of passing the stdscr to the funct, we pass the class
			#itself so the caller can use us
			return func(self, *args, **kwds)
		finally:
			# Set everything back to normal
			self.__stdscr.keypad(0)
			curses.echo()
			curses.nocbreak()
			curses.curs_set(1)
			curses.endwin()

################################################################################
#general stuff related to the screen itself, 
#ie. clearing the screen and printing a string/character
#colors being used:
#1=black, 2=blue, 3=green, 4=cyan, 5=red, 6=magenta, 7=yellow, 8=white
#modifiers:
#1=none, 2=dim, 3=bold, 4=standout

	def clearScreen(self):
		'''Clears the screen from everything.'''
		self.__stdscr.erase()

	def showCh(self, x, y, char, color = 8, modifier = 1):
		'''Prints a character to the screen. If colors is enabled, the color and
			modifier(ie. bold, standout) of the character can be changed.'''
		if self.colorsEnabled:
			if not self.__checkAttributes(color, modifier):
				color = 8
				modifier = 1
			attr = curses.color_pair(color) | self.__returnMod(modifier)
			self.__stdscr.addch(y, x, char, attr)
		else:
			self.__stdscr.addch(y, x, char)
		
	def showStr(self, x, y, msg, color = 8, modifier = 1):
		'''Prints a string to the screen. If colors is enabled, the color and
			modifier(ie. bold, standout) of the text can be changed.'''
		if self.colorsEnabled:
			if not self.__checkAttributes(color, modifier):
				color = 8
				modifier = 1
			attr = curses.color_pair(color) | self.__returnMod(modifier)
			self.__stdscr.addstr(y, x, msg, attr)
		else:
			self.__stdscr.addstr(y, x, msg)

	def __checkAttributes(self, color, modifier):
		'''Check if color and modifier has valid values.'''
		if (0 < color < 9) and (0 < modifier < 5):
			return True
		return False

	def __returnMod(self, modifier):
		'''Returns the curses value of a certain modifier, ie. curses.A_BOLD'''
		if modifier == 2:
			return curses.A_DIM
		elif modifier == 3:
			return curses.A_BOLD
		elif modifier == 4:
			return curses.A_STANDOUT
		else:
			return curses.A_NORMAL
################################################################################
#stuff related to keys and key handling,
#ie. get a new key from the user or checking if a key equals something

	def getKey(self):
		'''Returns a new key from curses.'''
		return self.__stdscr.getch()

	def checkKey(self, key, curseskey):
		'''Check if a key is equal to a curses key, ie. key == curses.KEY_END'''
		return (key == curses.curseskey)

	def keyName(self, key):
		'''Returns the name of a key.'''
		return curses.keyname(key)

################################################################################
#stuff related to messages and the messaging system
#ie. adding a new msg to the buffer or pop something out of it

	def addMsg(self, msg):
		'''Add a new message to the buffer.'''
		temp = msg.strip()
		if len(temp) > 0:
			self.__buffer += temp

	def showMsg(self, color = 8, modifier = 1):
		'''Shows the buffered messages. Returns true when there's more to show,
			false when the buffer is empty'''
		self.__stdscr.move(self.__height-3, 0)
		self.__stdscr.clrtobot()
		
		buff = textwrap.wrap(self.__buffer, self.__width-1)
		self.clearMsg()
		if len(buff) > 3:
			#show the first two lines
			line = buff.pop(0)
			self.showStr(0, self.__height-3, line, color, modifier)
			line = buff.pop(0)
			self.showStr(0, self.__height-2, line, color, modifier)
			
			#show the third line, "minimized" to fit with the extra more msg
			more = '[More]'
			temp = textwrap.wrap(buff.pop(0), (self.__width-1)-len(more))
			line = temp[0] + more
			self.showStr(0, self.__height-1, line, color, modifier)
			
			#check if we have a leftover from the last line
			line = ''
			if len(temp) > 1:
				line = temp[1]
			
			#add leftovers to the buffer
			self.__buffer = line
			for line in buff:
				self.__buffer += line
			return True

		else:
			y = self.__height - 3
			for line in buff:
				self.showStr(0, y, line, color, modifier)
				y += 1
			return False	

	def clearMsg(self):
		'''Clear the whole messaging buffer.'''
		self.__buffer = ''
