import io

class main:	
	def loop(self, wrapper):
		console = wrapper
		console.clearScreen()
		
		askmore = False
		while 1:
			k = console.getKey()
			console.clearScreen()
			
			if askmore and (k == 13 or k == 10):
				pass
			else:
				console.clearMsg()
				msg = 'pressed: ' + str(k) + ' (' + console.keyName(k) + ')'
				console.addMsg(msg)
				if console.keyName(k) == 'a':
					msg = " ="*(160*2) + 'END!'
					console.addMsg(msg)

			askmore = console.showMsg(5)
			
			if k == 27:
				break

if __name__ == '__main__':
	#io.iowrapper(main)
	test = main()
	io.iowrapper(test.loop)
