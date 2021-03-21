#for gathering data about the OS, the platform, current time and traceback
import os, platform, time, traceback

class logger:
	def __init__(self, name = 'TempPythonLog'):
		''' '''
		self.__logname = str(name)
		self.__filename = str(name) + '.log'
		self.__fulltimeformat = '%a, %d %b %Y %H:%M:%S %Z'
		self.__timestamp = '%d/%m/%y-%H:%M:%S'
		
		i = (80 - len(' START OF ' + str(name) + ' ')) / 2
		
		temp = ':'*i + ' START OF ' + str(name) + ' ' + ':'*i + '\n'
		temp += '='*80 + '\n'
		temp += 'FILE: ' + self.__filename + '\n'
		temp += 'TIME: ' + time.strftime(self.__fulltimeformat) + '\n'
		temp += 'PLATFORM: ' + platform.platform(True, True) + '\n'
		temp += 'PYTHON_VERSION: ' + platform.python_version() + '\n'
		temp += 'PYTHON_BUILD: ' + platform.python_build()[0] + ' (' + platform.python_build()[1] + ')\n'
		temp += '='*80 + '\n'
		
		self.__logdata(temp, 'w')

	def info(self, text):
		''' '''
		head = time.strftime(self.__timestamp) + '-INFO: '
		self.__logdata(str(head + text + '\n'))

	def debug(self, text):
		''' '''
		head = time.strftime(self.__timestamp) + '-DEBUG: '
		self.__logdata(str(head + text + '\n'))
		
	def warning(self, text):
		''' '''
		head = time.strftime(self.__timestamp) + '-WARNING: '
		self.__logdata(str(head + text + '\n'))
		
	def error(self, text):
		''' '''
		head = time.strftime(self.__timestamp) + '-ERROR: '
		self.__logdata(str(head + text + '\n'))
		
	def fatal(self, text):
		''' '''
		head = time.strftime(self.__timestamp) + '-FATAL: '
		self.__logdata(str(head + text + '\n'))

################################################################################

	def latest_error(self):
		''' '''
		data = traceback.format_exc() + '-'*80
		head = time.strftime(self.__timestamp) + '-ERROR: An error was caught by Python:\n'
		self.__logdata(str(head + data + '\n'))		

################################################################################

	def __logdata(self, data, mode = 'a'):
		file = open(self.__filename, mode)
		file.write(data)
		file.close()
